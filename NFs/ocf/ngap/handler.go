package ngap

import (
	"encoding/hex"
	"strconv"

	"github.com/free5gc/aper"
	"github.com/free5gc/nas/nasMessage"
	libngap "github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/ocf/consumer"
	"github.com/free5gc/ocf/context"
	gmm_message "github.com/free5gc/ocf/gmm/message"
	"github.com/free5gc/ocf/logger"
	"github.com/free5gc/ocf/nas"
	ngap_message "github.com/free5gc/ocf/ngap/message"
	"github.com/free5gc/openapi/models"
)

func HandleNGSetupRequest(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var globalRANNodeID *ngapType.GlobalRANNodeID
	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	var cause ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	nGSetupRequest := initiatingMessage.Value.NGSetupRequest
	if nGSetupRequest == nil {
		ran.Log.Error("NGSetupRequest is nil")
		return
	}
	ran.Log.Info("Handle NG Setup request")
	for i := 0; i < len(nGSetupRequest.ProtocolIEs.List); i++ {
		ie := nGSetupRequest.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDGlobalRANNodeID:
			globalRANNodeID = ie.Value.GlobalRANNodeID
			ran.Log.Trace("Decode IE GlobalRANNodeID")
			if globalRANNodeID == nil {
				ran.Log.Error("GlobalRANNodeID is nil")
				return
			}
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			ran.Log.Trace("Decode IE SupportedTAList")
			if supportedTAList == nil {
				ran.Log.Error("SupportedTAList is nil")
				return
			}
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			ran.Log.Trace("Decode IE RANNodeName")
			if rANNodeName == nil {
				ran.Log.Error("RANNodeName is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			ran.Log.Trace("Decode IE DefaultPagingDRX")
			if pagingDRX == nil {
				ran.Log.Error("DefaultPagingDRX is nil")
				return
			}
		}
	}

	ran.SetRanId(globalRANNodeID)
	if rANNodeName != nil {
		ran.Name = rANNodeName.Value
	}
	if pagingDRX != nil {
		ran.Log.Tracef("PagingDRX[%d]", pagingDRX.Value)
	}

	for i := 0; i < len(supportedTAList.List); i++ {
		supportedTAItem := supportedTAList.List[i]
		tac := hex.EncodeToString(supportedTAItem.TAC.Value)
		capOfSupportTai := cap(ran.SupportedTAList)
		for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
			supportedTAI := context.NewSupportedTAI()
			supportedTAI.Tai.Tac = tac
			broadcastPLMNItem := supportedTAItem.BroadcastPLMNList.List[j]
			plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
			supportedTAI.Tai.PlmnId = &plmnId
			capOfSNssaiList := cap(supportedTAI.SNssaiList)
			for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList.List); k++ {
				tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList.List[k]
				if len(supportedTAI.SNssaiList) < capOfSNssaiList {
					supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
				} else {
					break
				}
			}
			ran.Log.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
			if len(ran.SupportedTAList) < capOfSupportTai {
				ran.SupportedTAList = append(ran.SupportedTAList, supportedTAI)
			} else {
				break
			}
		}
	}

	if len(ran.SupportedTAList) == 0 {
		ran.Log.Warn("NG-Setup failure: No supported TA exist in NG-Setup request")
		cause.Present = ngapType.CausePresentMisc
		cause.Misc = &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		}
	} else {
		var found bool
		for i, tai := range ran.SupportedTAList {
			if context.InTaiList(tai.Tai, context.OCF_Self().SupportTaiLists) {
				ran.Log.Tracef("SERVED_TAI_INDEX[%d]", i)
				found = true
				break
			}
		}
		if !found {
			ran.Log.Warn("NG-Setup failure: Cannot find Served TAI in OCF")
			cause.Present = ngapType.CausePresentMisc
			cause.Misc = &ngapType.CauseMisc{
				Value: ngapType.CauseMiscPresentUnknownPLMN,
			}
		}
	}

	if cause.Present == ngapType.CausePresentNothing {
		ngap_message.SendNGSetupResponse(ran)
	} else {
		ngap_message.SendNGSetupFailure(ran, cause)
	}
}

func HandleUplinkNasTransport(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var userLocationInformation *ngapType.UserLocationInformation

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}

	uplinkNasTransport := initiatingMessage.Value.UplinkNASTransport
	if uplinkNasTransport == nil {
		ran.Log.Error("UplinkNasTransport is nil")
		return
	}
	ran.Log.Info("Handle Uplink Nas Transport")

	for i := 0; i < len(uplinkNasTransport.ProtocolIEs.List); i++ {
		ie := uplinkNasTransport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDNASPDU:
			nASPDU = ie.Value.NASPDU
			ran.Log.Trace("Decode IE NasPdu")
			if nASPDU == nil {
				ran.Log.Error("nASPDU is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE UserLocationInformation")
			if userLocationInformation == nil {
				ran.Log.Error("UserLocationInformation is nil")
				return
			}
		}
	}

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	ocfUe := ranUe.OcfUe
	if ocfUe == nil {
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorf(err.Error())
		}
		ran.Log.Errorf("No UE Context of RanUe with RANUENGAPID[%d] OCFUENGAPID[%d] ",
			rANUENGAPID.Value, aMFUENGAPID.Value)
		return
	}

	ranUe.Log.Infof("Uplink NAS Transport (RAN UE NGAP ID: %d)", ranUe.RanUeNgapId)

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	nas.HandleNAS(ranUe, ngapType.ProcedureCodeUplinkNASTransport, nASPDU.Value)
}

func HandleNGReset(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var cause *ngapType.Cause
	var resetType *ngapType.ResetType

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	nGReset := initiatingMessage.Value.NGReset
	if nGReset == nil {
		ran.Log.Error("NGReset is nil")
		return
	}

	ran.Log.Info("Handle NG Reset")

	for _, ie := range nGReset.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
			if cause == nil {
				ran.Log.Error("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDResetType:
			resetType = ie.Value.ResetType
			ran.Log.Trace("Decode IE ResetType")
			if resetType == nil {
				ran.Log.Error("ResetType is nil")
				return
			}
		}
	}

	printAndGetCause(ran, cause)

	switch resetType.Present {
	case ngapType.ResetTypePresentNGInterface:
		ran.Log.Trace("ResetType Present: NG Interface")
		ran.RemoveAllUeInRan()
		ngap_message.SendNGResetAcknowledge(ran, nil, nil)
	case ngapType.ResetTypePresentPartOfNGInterface:
		ran.Log.Trace("ResetType Present: Part of NG Interface")

		partOfNGInterface := resetType.PartOfNGInterface
		if partOfNGInterface == nil {
			ran.Log.Error("PartOfNGInterface is nil")
			return
		}

		var ranUe *context.RanUe

		for _, ueAssociatedLogicalNGConnectionItem := range partOfNGInterface.List {
			if ueAssociatedLogicalNGConnectionItem.OCFUENGAPID != nil {
				ran.Log.Tracef("OcfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.OCFUENGAPID.Value)
				for _, ue := range ran.RanUeList {
					if ue.OcfUeNgapId == ueAssociatedLogicalNGConnectionItem.OCFUENGAPID.Value {
						ranUe = ue
						break
					}
				}
			} else if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
				ran.Log.Tracef("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				ranUe = ran.RanUeFindByRanUeNgapID(ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
			}

			if ranUe == nil {
				ran.Log.Warn("Cannot not find UE Context")
				if ueAssociatedLogicalNGConnectionItem.OCFUENGAPID != nil {
					ran.Log.Warnf("OcfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.OCFUENGAPID.Value)
				}
				if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
					ran.Log.Warnf("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				}
			}

			err := ranUe.Remove()
			if err != nil {
				ran.Log.Error(err.Error())
			}
		}
		ngap_message.SendNGResetAcknowledge(ran, partOfNGInterface, nil)
	default:
		ran.Log.Warnf("Invalid ResetType[%d]", resetType.Present)
	}
}

func HandleNGResetAcknowledge(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var uEAssociatedLogicalNGConnectionList *ngapType.UEAssociatedLogicalNGConnectionList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}
	nGResetAcknowledge := successfulOutcome.Value.NGResetAcknowledge
	if nGResetAcknowledge == nil {
		ran.Log.Error("NGResetAcknowledge is nil")
		return
	}

	ran.Log.Info("Handle NG Reset Acknowledge")

	for _, ie := range nGResetAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDUEAssociatedLogicalNGConnectionList:
			uEAssociatedLogicalNGConnectionList = ie.Value.UEAssociatedLogicalNGConnectionList
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if uEAssociatedLogicalNGConnectionList != nil {
		ran.Log.Tracef("%d UE association(s) has been reset", len(uEAssociatedLogicalNGConnectionList.List))
		for i, item := range uEAssociatedLogicalNGConnectionList.List {
			if item.OCFUENGAPID != nil && item.RANUENGAPID != nil {
				ran.Log.Tracef("%d: OcfUeNgapID[%d] RanUeNgapID[%d]", i+1, item.OCFUENGAPID.Value, item.RANUENGAPID.Value)
			} else if item.OCFUENGAPID != nil {
				ran.Log.Tracef("%d: OcfUeNgapID[%d] RanUeNgapID[-1]", i+1, item.OCFUENGAPID.Value)
			} else if item.RANUENGAPID != nil {
				ran.Log.Tracef("%d: OcfUeNgapID[-1] RanUeNgapID[%d]", i+1, item.RANUENGAPID.Value)
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func HandleUEContextReleaseComplete(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation
	var infoOnRecommendedCellsAndRANNodesForPaging *ngapType.InfoOnRecommendedCellsAndRANNodesForPaging
	var pDUSessionResourceList *ngapType.PDUSessionResourceListCxtRelCpl
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}
	uEContextReleaseComplete := successfulOutcome.Value.UEContextReleaseComplete
	if uEContextReleaseComplete == nil {
		ran.Log.Error("NGResetAcknowledge is nil")
		return
	}

	ran.Log.Info("Handle UE Context Release Complete")

	for _, ie := range uEContextReleaseComplete.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDInfoOnRecommendedCellsAndRANNodesForPaging:
			infoOnRecommendedCellsAndRANNodesForPaging = ie.Value.InfoOnRecommendedCellsAndRANNodesForPaging
			ran.Log.Trace("Decode IE InfoOnRecommendedCellsAndRANNodesForPaging")
			if infoOnRecommendedCellsAndRANNodesForPaging != nil {
				ran.Log.Warn("IE infoOnRecommendedCellsAndRANNodesForPaging is not support")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceListCxtRelCpl:
			pDUSessionResourceList = ie.Value.PDUSessionResourceListCxtRelCpl
			ran.Log.Trace("Decode IE PDUSessionResourceList")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	ranUe := context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No RanUe Context[OcfUeNgapID: %d]", aMFUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}
	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	ocfUe := ranUe.OcfUe
	if ocfUe == nil {
		ran.Log.Infof("Release UE Context : RanUe[OcfUeNgapId: %d]", ranUe.OcfUeNgapId)
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorln(err.Error())
		}
		return
	}
	// TODO: OCF shall, if supported, store it and may use it for subsequent paging
	if infoOnRecommendedCellsAndRANNodesForPaging != nil {
		ocfUe.InfoOnRecommendedCellsAndRanNodesForPaging = new(context.InfoOnRecommendedCellsAndRanNodesForPaging)

		recommendedCells := &ocfUe.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedCells
		for _, item := range infoOnRecommendedCellsAndRANNodesForPaging.RecommendedCellsForPaging.RecommendedCellList.List {
			recommendedCell := context.RecommendedCell{}

			switch item.NGRANCGI.Present {
			case ngapType.NGRANCGIPresentNRCGI:
				recommendedCell.NgRanCGI.Present = context.NgRanCgiPresentNRCGI
				recommendedCell.NgRanCGI.NRCGI = new(models.Ncgi)
				plmnID := ngapConvert.PlmnIdToModels(item.NGRANCGI.NRCGI.PLMNIdentity)
				recommendedCell.NgRanCGI.NRCGI.PlmnId = &plmnID
				recommendedCell.NgRanCGI.NRCGI.NrCellId = ngapConvert.BitStringToHex(&item.NGRANCGI.NRCGI.NRCellIdentity.Value)
			case ngapType.NGRANCGIPresentEUTRACGI:
				recommendedCell.NgRanCGI.Present = context.NgRanCgiPresentEUTRACGI
				recommendedCell.NgRanCGI.EUTRACGI = new(models.Ecgi)
				plmnID := ngapConvert.PlmnIdToModels(item.NGRANCGI.EUTRACGI.PLMNIdentity)
				recommendedCell.NgRanCGI.EUTRACGI.PlmnId = &plmnID
				recommendedCell.NgRanCGI.EUTRACGI.EutraCellId = ngapConvert.BitStringToHex(
					&item.NGRANCGI.EUTRACGI.EUTRACellIdentity.Value)
			}

			if item.TimeStayedInCell != nil {
				recommendedCell.TimeStayedInCell = new(int64)
				*recommendedCell.TimeStayedInCell = *item.TimeStayedInCell
			}

			*recommendedCells = append(*recommendedCells, recommendedCell)
		}

		recommendedRanNodes := &ocfUe.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedRanNodes
		ranNodeList := infoOnRecommendedCellsAndRANNodesForPaging.RecommendRANNodesForPaging.RecommendedRANNodeList.List
		for _, item := range ranNodeList {
			recommendedRanNode := context.RecommendRanNode{}

			switch item.OCFPagingTarget.Present {
			case ngapType.OCFPagingTargetPresentGlobalRANNodeID:
				recommendedRanNode.Present = context.RecommendRanNodePresentRanNode
				recommendedRanNode.GlobalRanNodeId = new(models.GlobalRanNodeId)
				// TODO: recommendedRanNode.GlobalRanNodeId = ngapConvert.RanIdToModels(item.OCFPagingTarget.GlobalRANNodeID)
			case ngapType.OCFPagingTargetPresentTAI:
				recommendedRanNode.Present = context.RecommendRanNodePresentTAI
				tai := ngapConvert.TaiToModels(*item.OCFPagingTarget.TAI)
				recommendedRanNode.Tai = &tai
			}
			*recommendedRanNodes = append(*recommendedRanNodes, recommendedRanNode)
		}
	}

	// for each pduSessionID invoke Nsmf_PDUSession_UpdateSMContext Request
	var cause context.CauseAll
	if tmp, exist := ocfUe.ReleaseCause[ran.AnType]; exist {
		cause = *tmp
	}
	if ocfUe.State[ran.AnType].Is(context.Registered) {
		ranUe.Log.Info("Rel Ue Context in GMM-Registered")
		if pDUSessionResourceList != nil {
			for _, pduSessionReourceItem := range pDUSessionResourceList.List {
				pduSessionID := int32(pduSessionReourceItem.PDUSessionID.Value)
				smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				}
				response, _, _, err := consumer.SendUpdateSmContextDeactivateUpCnxState(ocfUe, smContext, cause)
				if err != nil {
					ran.Log.Errorf("Send Update SmContextDeactivate UpCnxState Error[%s]", err.Error())
				} else if response == nil {
					ran.Log.Errorln("Send Update SmContextDeactivate UpCnxState Error")
				}
			}
		}
	}

	// Remove UE N2 Connection
	ocfUe.ReleaseCause[ran.AnType] = nil
	switch ranUe.ReleaseAction {
	case context.UeContextN2NormalRelease:
		ran.Log.Infof("Release UE[%s] Context : N2 Connection Release", ocfUe.Supi)
		// ocfUe.DetachRanUe(ran.AnType)
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorln(err.Error())
		}
	case context.UeContextReleaseUeContext:
		ran.Log.Infof("Release UE[%s] Context : Release Ue Context", ocfUe.Supi)
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorln(err.Error())
		}
		ocfUe.Remove()
	case context.UeContextReleaseHandover:
		ran.Log.Infof("Release UE[%s] Context : Release for Handover", ocfUe.Supi)
		// TODO: it's a workaround, need to fix it.
		targetRanUe := context.OCF_Self().RanUeFindByOcfUeNgapID(ranUe.TargetUe.OcfUeNgapId)

		context.DetachSourceUeTargetUe(ranUe)
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorln(err.Error())
		}
		ocfUe.AttachRanUe(targetRanUe)
		// Todo: remove indirect tunnel
	default:
		ran.Log.Errorf("Invalid Release Action[%d]", ranUe.ReleaseAction)
	}
}

func HandlePDUSessionResourceReleaseResponse(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceReleasedList *ngapType.PDUSessionResourceReleasedListRelRes
	var userLocationInformation *ngapType.UserLocationInformation
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceReleaseResponse := successfulOutcome.Value.PDUSessionResourceReleaseResponse
	if pDUSessionResourceReleaseResponse == nil {
		ran.Log.Error("PDUSessionResourceReleaseResponse is nil")
		return
	}

	ran.Log.Info("Handle PDU Session Resource Release Response")

	for _, ie := range pDUSessionResourceReleaseResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUENgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDPDUSessionResourceReleasedListRelRes:
			pDUSessionResourceReleasedList = ie.Value.PDUSessionResourceReleasedListRelRes
			ran.Log.Trace("Decode IE PDUSessionResourceReleasedList")
			if pDUSessionResourceReleasedList == nil {
				ran.Log.Error("PDUSessionResourceReleasedList is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	ocfUe := ranUe.OcfUe
	if ocfUe == nil {
		ranUe.Log.Error("ocfUe is nil")
		return
	}
	if pDUSessionResourceReleasedList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceReleaseResponseTransfer to SMF")

		for _, item := range pDUSessionResourceReleasedList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceReleaseResponseTransfer
			smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			_, responseErr, problemDetail, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
				models.N2SmInfoType_PDU_RES_REL_RSP, transfer)
			// TODO: error handling
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Error: %+v", err)
			} else if responseErr != nil && responseErr.JsonData.Error != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Error: %+v",
					responseErr.JsonData.Error.Cause)
			} else if problemDetail != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Failed: %+v", problemDetail)
			}
		}
	}
}

func HandleUERadioCapabilityCheckResponse(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var iMSVoiceSupportIndicator *ngapType.IMSVoiceSupportIndicator
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var ranUe *context.RanUe

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}

	uERadioCapabilityCheckResponse := successfulOutcome.Value.UERadioCapabilityCheckResponse
	if uERadioCapabilityCheckResponse == nil {
		ran.Log.Error("UERadioCapabilityCheckResponse is nil")
		return
	}
	ran.Log.Info("Handle UE Radio Capability Check Response")

	for i := 0; i < len(uERadioCapabilityCheckResponse.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityCheckResponse.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDIMSVoiceSupportIndicator:
			iMSVoiceSupportIndicator = ie.Value.IMSVoiceSupportIndicator
			ran.Log.Trace("Decode IE IMSVoiceSupportIndicator")
			if iMSVoiceSupportIndicator == nil {
				ran.Log.Error("iMSVoiceSupportIndicator is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	// TODO: handle iMSVoiceSupportIndicator

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func HandleLocationReportingFailureIndication(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var ranUe *context.RanUe

	var cause *ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	locationReportingFailureIndication := initiatingMessage.Value.LocationReportingFailureIndication
	if locationReportingFailureIndication == nil {
		ran.Log.Error("LocationReportingFailureIndication is nil")
		return
	}

	ran.Log.Info("Handle Location Reporting Failure Indication")

	for i := 0; i < len(locationReportingFailureIndication.ProtocolIEs.List); i++ {
		ie := locationReportingFailureIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
			if cause == nil {
				ran.Log.Error("Cause is nil")
				return
			}
		}
	}

	printAndGetCause(ran, cause)

	ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
}

func HandleInitialUEMessage(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	ocfSelf := context.OCF_Self()

	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var userLocationInformation *ngapType.UserLocationInformation
	var rRCEstablishmentCause *ngapType.RRCEstablishmentCause
	var fiveGSTMSI *ngapType.FiveGSTMSI
	// var aMFSetID *ngapType.OCFSetID
	var uEContextRequest *ngapType.UEContextRequest
	// var allowedNSSAI *ngapType.AllowedNSSAI

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	initialUEMessage := initiatingMessage.Value.InitialUEMessage
	if initialUEMessage == nil {
		ran.Log.Error("InitialUEMessage is nil")
		return
	}

	ran.Log.Info("Handle Initial UE Message")

	for _, ie := range initialUEMessage.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDNASPDU: // reject
			nASPDU = ie.Value.NASPDU
			ran.Log.Trace("Decode IE NasPdu")
			if nASPDU == nil {
				ran.Log.Error("NasPdu is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDNASPDU,
					ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // reject
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE UserLocationInformation")
			if userLocationInformation == nil {
				ran.Log.Error("UserLocationInformation is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDUserLocationInformation, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRRCEstablishmentCause: // ignore
			rRCEstablishmentCause = ie.Value.RRCEstablishmentCause
			ran.Log.Trace("Decode IE RRCEstablishmentCause")
		case ngapType.ProtocolIEIDFiveGSTMSI: // optional, reject
			fiveGSTMSI = ie.Value.FiveGSTMSI
			ran.Log.Trace("Decode IE 5G-S-TMSI")
		case ngapType.ProtocolIEIDOCFSetID: // optional, ignore
			// aMFSetID = ie.Value.OCFSetID
			ran.Log.Trace("Decode IE OcfSetID")
		case ngapType.ProtocolIEIDUEContextRequest: // optional, ignore
			uEContextRequest = ie.Value.UEContextRequest
			ran.Log.Trace("Decode IE UEContextRequest")
		case ngapType.ProtocolIEIDAllowedNSSAI: // optional, reject
			// allowedNSSAI = ie.Value.AllowedNSSAI
			ran.Log.Trace("Decode IE Allowed NSSAI")
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		ran.Log.Trace("Has missing reject IE(s)")

		procedureCode := ngapType.ProcedureCodeInitialUEMessage
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentIgnore
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality,
			&iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
	}

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe != nil && ranUe.OcfUe == nil {
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorln(err.Error())
		}
		ranUe = nil
	}
	if ranUe == nil {
		var err error
		ranUe, err = ran.NewRanUe(rANUENGAPID.Value)
		if err != nil {
			ran.Log.Errorf("NewRanUe Error: %+v", err)
		}
		ran.Log.Debugf("New RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)

		if fiveGSTMSI != nil {
			ranUe.Log.Debug("Receive 5G-S-TMSI")

			servedGuami := ocfSelf.ServedGuamiList[0]

			// <5G-S-TMSI> := <OCF Set ID><OCF Pointer><5G-TMSI>
			// GUAMI := <MCC><MNC><OCF Region ID><OCF Set ID><OCF Pointer>
			// 5G-GUTI := <GUAMI><5G-TMSI>
			tmpReginID, _, _ := ngapConvert.OcfIdToNgap(servedGuami.OcfId)
			ocfID := ngapConvert.OcfIdToModels(tmpReginID, fiveGSTMSI.OCFSetID.Value, fiveGSTMSI.OCFPointer.Value)

			tmsi := hex.EncodeToString(fiveGSTMSI.FiveGTMSI.Value)

			guti := servedGuami.PlmnId.Mcc + servedGuami.PlmnId.Mnc + ocfID + tmsi

			// TODO: invoke Nocf_Communication_UEContextTransfer if serving OCF has changed since
			// last Registration Request procedure
			// Described in TS 23.502 4.2.2.2.2 step 4 (without UDSF deployment)

			if ocfUe, ok := ocfSelf.OcfUeFindByGuti(guti); !ok {
				ranUe.Log.Warnf("Unknown UE [GUTI: %s]", guti)
			} else {
				ranUe.Log.Tracef("find OcfUe [GUTI: %s]", guti)

				if ocfUe.CmConnect(ran.AnType) {
					ranUe.Log.Debug("Implicit Deregistration")
					ranUe.Log.Tracef("RanUeNgapID[%d]", ocfUe.RanUe[ran.AnType].RanUeNgapId)
					ocfUe.DetachRanUe(ran.AnType)
				}
				// TODO: stop Implicit Deregistration timer
				ranUe.Log.Debugf("OcfUe Attach RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)
				ocfUe.AttachRanUe(ranUe)
			}
		}
	} else {
		ranUe.OcfUe.AttachRanUe(ranUe)
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if rRCEstablishmentCause != nil {
		ranUe.Log.Tracef("[Initial UE Message] RRC Establishment Cause[%d]", rRCEstablishmentCause.Value)
		ranUe.RRCEstablishmentCause = strconv.Itoa(int(rRCEstablishmentCause.Value))
	}

	if uEContextRequest != nil {
		ran.Log.Debug("Trigger initial Context Setup procedure")
		ranUe.UeContextRequest = true
		// TODO: Trigger Initial Context Setup procedure
	} else {
		ranUe.UeContextRequest = false
	}

	// TS 23.502 4.2.2.2.3 step 6a Nnrf_NFDiscovery_Request (NF type, OCF Set)
	// if aMFSetID != nil {
	// TODO: This is a rerouted message
	// TS 38.413: OCF shall, if supported, use the IE as described in TS 23.502
	// }

	// ng-ran propagate allowedNssai in the rerouted initial ue message (TS 38.413 8.6.5)
	// TS 23.502 4.2.2.2.3 step 4a Nnssf_NSSelection_Get
	// if allowedNSSAI != nil {
	// TODO: OCF should use it as defined in TS 23.502
	// }

	pdu, err := libngap.Encoder(*message)
	if err != nil {
		ran.Log.Errorf("libngap Encoder Error: %+v", err)
	}
	ranUe.InitialUEMessage = pdu
	nas.HandleNAS(ranUe, ngapType.ProcedureCodeInitialUEMessage, nASPDU.Value)
}

func HandlePDUSessionResourceSetupResponse(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceSetupResponseList *ngapType.PDUSessionResourceSetupListSURes
	var pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListSURes
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *context.RanUe

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	if pDUSessionResourceSetupResponse == nil {
		ran.Log.Error("PDUSessionResourceSetupResponse is nil")
		return
	}

	ran.Log.Info("Handle PDU Session Resource Setup Response")

	for _, ie := range pDUSessionResourceSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes: // ignore
			pDUSessionResourceSetupResponseList = ie.Value.PDUSessionResourceSetupListSURes
			ran.Log.Trace("Decode IE PDUSessionResourceSetupListSURes")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListSURes: // ignore
			pDUSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListSURes
			ran.Log.Trace("Decode IE PDUSessionResourceFailedToSetupListSURes")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if rANUENGAPID != nil {
		ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			ran.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			ran.Log.Warnf("No UE Context[OcfUeNgapID: %d]", aMFUENGAPID.Value)
			return
		}
	}

	if ranUe != nil {
		ranUe.Log.Tracef("OcfUeNgapID[%d] RanUeNgapID[%d]", ranUe.OcfUeNgapId, ranUe.RanUeNgapId)
		ocfUe := ranUe.OcfUe
		if ocfUe == nil {
			ranUe.Log.Error("ocfUe is nil")
			return
		}

		if pDUSessionResourceSetupResponseList != nil {
			ranUe.Log.Trace("Send PDUSessionResourceSetupResponseTransfer to SMF")

			for _, item := range pDUSessionResourceSetupResponseList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceSetupResponseTransfer
				smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				}
				_, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
					models.N2SmInfoType_PDU_RES_SETUP_RSP, transfer)
				if err != nil {
					ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupResponseTransfer] Error: %+v", err)
				}
				// RAN initiated QoS Flow Mobility in subclause 5.2.2.3.7
				// if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
				// } else if response == nil {
				// TODO: error handling
				// }
			}
		}

		if pDUSessionResourceFailedToSetupList != nil {
			ranUe.Log.Trace("Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

			for _, item := range pDUSessionResourceFailedToSetupList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer
				smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				}
				_, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
					models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
				if err != nil {
					ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error: %+v", err)
				}

				// if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
				// } else if response == nil {
				// TODO: error handling
				// }
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func HandlePDUSessionResourceModifyResponse(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pduSessionResourceModifyResponseList *ngapType.PDUSessionResourceModifyListModRes
	var pduSessionResourceFailedToModifyList *ngapType.PDUSessionResourceFailedToModifyListModRes
	var userLocationInformation *ngapType.UserLocationInformation
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *context.RanUe

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceModifyResponse := successfulOutcome.Value.PDUSessionResourceModifyResponse
	if pDUSessionResourceModifyResponse == nil {
		ran.Log.Error("PDUSessionResourceModifyResponse is nil")
		return
	}

	ran.Log.Info("Handle PDU Session Resource Modify Response")

	for _, ie := range pDUSessionResourceModifyResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModRes: // ignore
			pduSessionResourceModifyResponseList = ie.Value.PDUSessionResourceModifyListModRes
			ran.Log.Trace("Decode IE PDUSessionResourceModifyListModRes")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToModifyListModRes: // ignore
			pduSessionResourceFailedToModifyList = ie.Value.PDUSessionResourceFailedToModifyListModRes
			ran.Log.Trace("Decode IE PDUSessionResourceFailedToModifyListModRes")
		case ngapType.ProtocolIEIDUserLocationInformation: // optional, ignore
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if rANUENGAPID != nil {
		ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			ran.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			ran.Log.Warnf("No UE Context[OcfUeNgapID: %d]", aMFUENGAPID.Value)
			return
		}
	}

	if ranUe != nil {
		ranUe.Log.Tracef("OcfUeNgapID[%d] RanUeNgapID[%d]", ranUe.OcfUeNgapId, ranUe.RanUeNgapId)
		ocfUe := ranUe.OcfUe
		if ocfUe == nil {
			ranUe.Log.Error("ocfUe is nil")
			return
		}

		if pduSessionResourceModifyResponseList != nil {
			ranUe.Log.Trace("Send PDUSessionResourceModifyResponseTransfer to SMF")

			for _, item := range pduSessionResourceModifyResponseList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceModifyResponseTransfer
				smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				}
				_, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
					models.N2SmInfoType_PDU_RES_MOD_RSP, transfer)
				if err != nil {
					ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceModifyResponseTransfer] Error: %+v", err)
				}
				// if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
				// } else if response == nil {
				// TODO: error handling
				// }
			}
		}

		if pduSessionResourceFailedToModifyList != nil {
			ranUe.Log.Trace("Send PDUSessionResourceModifyUnsuccessfulTransfer to SMF")

			for _, item := range pduSessionResourceFailedToModifyList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceModifyUnsuccessfulTransfer
				smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				}
				// response, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, pduSessionID,
				_, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
					models.N2SmInfoType_PDU_RES_MOD_FAIL, transfer)
				if err != nil {
					ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceModifyUnsuccessfulTransfer] Error: %+v", err)
				}
				// if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
				// } else if response == nil {
				// TODO: error handling
				// }
			}
		}

		if userLocationInformation != nil {
			ranUe.UpdateLocation(userLocationInformation)
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func HandlePDUSessionResourceNotify(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceNotifyList *ngapType.PDUSessionResourceNotifyList
	var pDUSessionResourceReleasedListNot *ngapType.PDUSessionResourceReleasedListNot
	var userLocationInformation *ngapType.UserLocationInformation

	var ranUe *context.RanUe

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	PDUSessionResourceNotify := initiatingMessage.Value.PDUSessionResourceNotify
	if PDUSessionResourceNotify == nil {
		ran.Log.Error("PDUSessionResourceNotify is nil")
		return
	}

	for _, ie := range PDUSessionResourceNotify.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID // reject
			ran.Log.Trace("Decode IE OcfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID // reject
			ran.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceNotifyList: // reject
			pDUSessionResourceNotifyList = ie.Value.PDUSessionResourceNotifyList
			ran.Log.Trace("Decode IE pDUSessionResourceNotifyList")
			if pDUSessionResourceNotifyList == nil {
				ran.Log.Error("pDUSessionResourceNotifyList is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceReleasedListNot: // ignore
			pDUSessionResourceReleasedListNot = ie.Value.PDUSessionResourceReleasedListNot
			ran.Log.Trace("Decode IE PDUSessionResourceReleasedListNot")
			if pDUSessionResourceReleasedListNot == nil {
				ran.Log.Error("PDUSessionResourceReleasedListNot is nil")
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // optional, ignore
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE userLocationInformation")
			if userLocationInformation == nil {
				ran.Log.Warn("userLocationInformation is nil [optional]")
			}
		}
	}

	ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	}

	ranUe = context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Warnf("No UE Context[OcfUeNgapID: %d]", aMFUENGAPID.Value)
		return
	}

	ranUe.Log.Tracef("OcfUeNgapID[%d] RanUeNgapID[%d]", ranUe.OcfUeNgapId, ranUe.RanUeNgapId)
	ocfUe := ranUe.OcfUe
	if ocfUe == nil {
		ranUe.Log.Error("ocfUe is nil")
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	ranUe.Log.Trace("Send PDUSessionResourceNotifyTransfer to SMF")

	for _, item := range pDUSessionResourceNotifyList.List {
		pduSessionID := int32(item.PDUSessionID.Value)
		transfer := item.PDUSessionResourceNotifyTransfer
		smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
		if !ok {
			ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
		}
		response, errResponse, problemDetail, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
			models.N2SmInfoType_PDU_RES_NTY, transfer)
		if err != nil {
			ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceNotifyTransfer] Error: %+v", err)
		}

		if response != nil {
			responseData := response.JsonData
			n2Info := response.BinaryDataN1SmMessage
			n1Msg := response.BinaryDataN2SmInformation
			if n2Info != nil {
				switch responseData.N2SmInfoType {
				case models.N2SmInfoType_PDU_RES_MOD_REQ:
					ranUe.Log.Debugln("OCF Transfer NGAP PDU Resource Modify Req from SMF")
					var nasPdu []byte
					if n1Msg != nil {
						pduSessionId := uint8(pduSessionID)
						nasPdu, err =
							gmm_message.BuildDLNASTransport(ocfUe, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, pduSessionId, nil, nil, 0)
						if err != nil {
							ranUe.Log.Warnf("GMM Message build DL NAS Transport filaed: %v", err)
						}
					}
					list := ngapType.PDUSessionResourceModifyListModReq{}
					ngap_message.AppendPDUSessionResourceModifyListModReq(&list, pduSessionID, nasPdu, n2Info)
					ngap_message.SendPDUSessionResourceModifyRequest(ranUe, list)
				}
			}
		} else if errResponse != nil {
			errJSON := errResponse.JsonData
			n1Msg := errResponse.BinaryDataN2SmInformation
			ranUe.Log.Warnf("PDU Session Modification is rejected by SMF[pduSessionId:%d], Error[%s]\n",
				pduSessionID, errJSON.Error.Cause)
			if n1Msg != nil {
				gmm_message.SendDLNASTransport(
					ranUe, nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, pduSessionID, 0, nil, 0)
			}
			// TODO: handle n2 info transfer
		} else if err != nil {
			return
		} else {
			// TODO: error handling
			ranUe.Log.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
			return
		}
	}

	if pDUSessionResourceReleasedListNot != nil {
		ranUe.Log.Trace("Send PDUSessionResourceNotifyReleasedTransfer to SMF")
		for _, item := range pDUSessionResourceReleasedListNot.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceNotifyReleasedTransfer
			smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			response, errResponse, problemDetail, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
				models.N2SmInfoType_PDU_RES_NTY_REL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceNotifyReleasedTransfer] Error: %+v", err)
			}
			if response != nil {
				responseData := response.JsonData
				n2Info := response.BinaryDataN1SmMessage
				n1Msg := response.BinaryDataN2SmInformation
				if n2Info != nil {
					switch responseData.N2SmInfoType {
					case models.N2SmInfoType_PDU_RES_REL_CMD:
						ranUe.Log.Debugln("OCF Transfer NGAP PDU Session Resource Rel Co from SMF")
						var nasPdu []byte
						if n1Msg != nil {
							pduSessionId := uint8(pduSessionID)
							nasPdu, err = gmm_message.BuildDLNASTransport(
								ocfUe, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, pduSessionId, nil, nil, 0)
							if err != nil {
								ranUe.Log.Warnf("GMM Message build DL NAS Transport filaed: %v", err)
							}
						}
						list := ngapType.PDUSessionResourceToReleaseListRelCmd{}
						ngap_message.AppendPDUSessionResourceToReleaseListRelCmd(&list, pduSessionID, n2Info)
						ngap_message.SendPDUSessionResourceReleaseCommand(ranUe, nasPdu, list)
					}
				}
			} else if errResponse != nil {
				errJSON := errResponse.JsonData
				n1Msg := errResponse.BinaryDataN2SmInformation
				ranUe.Log.Warnf("PDU Session Release is rejected by SMF[pduSessionId:%d], Error[%s]\n",
					pduSessionID, errJSON.Error.Cause)
				if n1Msg != nil {
					gmm_message.SendDLNASTransport(
						ranUe, nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, pduSessionID, 0, nil, 0)
				}
			} else if err != nil {
				return
			} else {
				// TODO: error handling
				ranUe.Log.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
				return
			}
		}
	}
}

func HandlePDUSessionResourceModifyIndication(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pduSessionResourceModifyIndicationList *ngapType.PDUSessionResourceModifyListModInd

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var ranUe *context.RanUe

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // reject
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		cause := ngapType.Cause{
			Present: ngapType.CausePresentProtocol,
			Protocol: &ngapType.CauseProtocol{
				Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorReject,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}
	pDUSessionResourceModifyIndication := initiatingMessage.Value.PDUSessionResourceModifyIndication
	if pDUSessionResourceModifyIndication == nil {
		ran.Log.Error("PDUSessionResourceModifyIndication is nil")
		cause := ngapType.Cause{
			Present: ngapType.CausePresentProtocol,
			Protocol: &ngapType.CauseProtocol{
				Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorReject,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	ran.Log.Info("Handle PDU Session Resource Modify Indication")

	for _, ie := range pDUSessionResourceModifyIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // reject
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDOCFUENGAPID, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModInd: // reject
			pduSessionResourceModifyIndicationList = ie.Value.PDUSessionResourceModifyListModInd
			ran.Log.Trace("Decode IE PDUSessionResourceModifyListModInd")
			if pduSessionResourceModifyIndicationList == nil {
				ran.Log.Error("PDUSessionResourceModifyListModInd is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDPDUSessionResourceModifyListModInd, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		ran.Log.Error("Has missing reject IE(s)")

		procedureCode := ngapType.ProcedureCodePDUSessionResourceModifyIndication
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality,
			&iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	ran.Log.Tracef("UE Context OcfUeNgapID[%d] RanUeNgapID[%d]", ranUe.OcfUeNgapId, ranUe.RanUeNgapId)

	ocfUe := ranUe.OcfUe
	if ocfUe == nil {
		ran.Log.Error("OcfUe is nil")
		return
	}

	pduSessionResourceModifyListModCfm := ngapType.PDUSessionResourceModifyListModCfm{}
	pduSessionResourceFailedToModifyListModCfm := ngapType.PDUSessionResourceFailedToModifyListModCfm{}

	ran.Log.Trace("Send PDUSessionResourceModifyIndicationTransfer to SMF")
	for _, item := range pduSessionResourceModifyIndicationList.List {
		pduSessionID := int32(item.PDUSessionID.Value)
		transfer := item.PDUSessionResourceModifyIndicationTransfer
		smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
		if !ok {
			ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
		}
		response, errResponse, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
			models.N2SmInfoType_PDU_RES_MOD_IND, transfer)
		if err != nil {
			ran.Log.Errorf("SendUpdateSmContextN2Info Error:\n%s", err.Error())
		}

		if response != nil && response.BinaryDataN2SmInformation != nil {
			ngap_message.AppendPDUSessionResourceModifyListModCfm(&pduSessionResourceModifyListModCfm, int64(pduSessionID),
				response.BinaryDataN2SmInformation)
		}
		if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
			ngap_message.AppendPDUSessionResourceFailedToModifyListModCfm(&pduSessionResourceFailedToModifyListModCfm,
				int64(pduSessionID), errResponse.BinaryDataN2SmInformation)
		}
	}

	ngap_message.SendPDUSessionResourceModifyConfirm(ranUe, pduSessionResourceModifyListModCfm,
		pduSessionResourceFailedToModifyListModCfm, nil)
}

func HandleInitialContextSetupResponse(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceSetupResponseList *ngapType.PDUSessionResourceSetupListCxtRes
	var pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListCxtRes
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}
	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	if initialContextSetupResponse == nil {
		ran.Log.Error("InitialContextSetupResponse is nil")
		return
	}

	ran.Log.Info("Handle Initial Context Setup Response")

	for _, ie := range initialContextSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Warn("OcfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtRes:
			pDUSessionResourceSetupResponseList = ie.Value.PDUSessionResourceSetupListCxtRes
			ran.Log.Trace("Decode IE PDUSessionResourceSetupResponseList")
			if pDUSessionResourceSetupResponseList == nil {
				ran.Log.Warn("PDUSessionResourceSetupResponseList is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtRes:
			pDUSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListCxtRes
			ran.Log.Trace("Decode IE PDUSessionResourceFailedToSetupList")
			if pDUSessionResourceFailedToSetupList == nil {
				ran.Log.Warn("PDUSessionResourceFailedToSetupList is nil")
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE Criticality Diagnostics")
			if criticalityDiagnostics == nil {
				ran.Log.Warn("Criticality Diagnostics is nil")
			}
		}
	}

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	ocfUe := ranUe.OcfUe
	if ocfUe == nil {
		ran.Log.Error("ocfUe is nil")
		return
	}

	ran.Log.Tracef("RanUeNgapID[%d] OcfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.OcfUeNgapId)

	if pDUSessionResourceSetupResponseList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceSetupResponseTransfer to SMF")

		for _, item := range pDUSessionResourceSetupResponseList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupResponseTransfer
			smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			// response, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, pduSessionID,
			_, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_RSP, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupResponseTransfer] Error: %+v", err)
			}
			// RAN initiated QoS Flow Mobility in subclause 5.2.2.3.7
			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if pDUSessionResourceFailedToSetupList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

		for _, item := range pDUSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer
			smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			// response, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, pduSessionID,
			_, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error: %+v", err)
			}

			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if ranUe.Ran.AnType == models.AccessType_NON_3_GPP_ACCESS {
		ngap_message.SendDownlinkNasTransport(ranUe, ocfUe.RegistrationAcceptForNon3GPPAccess, nil)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func HandleInitialContextSetupFailure(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListCxtFail
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		ran.Log.Error("UnsuccessfulOutcome is nil")
		return
	}
	initialContextSetupFailure := unsuccessfulOutcome.Value.InitialContextSetupFailure
	if initialContextSetupFailure == nil {
		ran.Log.Error("InitialContextSetupFailure is nil")
		return
	}

	ran.Log.Info("Handle Initial Context Setup Failure")

	for _, ie := range initialContextSetupFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Warn("OcfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtFail:
			pDUSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListCxtFail
			ran.Log.Trace("Decode IE PDUSessionResourceFailedToSetupList")
			if pDUSessionResourceFailedToSetupList == nil {
				ran.Log.Warn("PDUSessionResourceFailedToSetupList is nil")
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
			if cause == nil {
				ran.Log.Warn("Cause is nil")
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE Criticality Diagnostics")
			if criticalityDiagnostics == nil {
				ran.Log.Warn("CriticalityDiagnostics is nil")
			}
		}
	}

	printAndGetCause(ran, cause)

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	ocfUe := ranUe.OcfUe
	if ocfUe == nil {
		ran.Log.Error("ocfUe is nil")
		return
	}

	if pDUSessionResourceFailedToSetupList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

		for _, item := range pDUSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer
			smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			_, _, _, err := consumer.SendUpdateSmContextN2Info(ocfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error: %+v", err)
			}

			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}
}

func HandleUEContextReleaseRequest(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceList *ngapType.PDUSessionResourceListCxtRelReq
	var cause *ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	uEContextReleaseRequest := initiatingMessage.Value.UEContextReleaseRequest
	if uEContextReleaseRequest == nil {
		ran.Log.Error("UEContextReleaseRequest is nil")
		return
	}

	ran.Log.Info("UE Context Release Request")

	for _, ie := range uEContextReleaseRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDPDUSessionResourceListCxtRelReq:
			pDUSessionResourceList = ie.Value.PDUSessionResourceListCxtRelReq
			ran.Log.Trace("Decode IE Pdu Session Resource List")
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
			if cause == nil {
				ran.Log.Warn("Cause is nil")
			}
		}
	}

	ranUe := context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No RanUe Context[OcfUeNgapID: %d]", aMFUENGAPID.Value)
		cause = &ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, cause, nil)
		return
	}

	ran.Log.Tracef("RanUeNgapID[%d] OcfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.OcfUeNgapId)

	causeGroup := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentUnspecified
	if cause != nil {
		causeGroup, causeValue = printAndGetCause(ran, cause)
	}

	ocfUe := ranUe.OcfUe
	if ocfUe != nil {
		causeAll := context.CauseAll{
			NgapCause: &models.NgApCause{
				Group: int32(causeGroup),
				Value: int32(causeValue),
			},
		}
		if ocfUe.State[ran.AnType].Is(context.Registered) {
			ranUe.Log.Info("Ue Context in GMM-Registered")
			if pDUSessionResourceList != nil {
				for _, pduSessionReourceItem := range pDUSessionResourceList.List {
					pduSessionID := int32(pduSessionReourceItem.PDUSessionID.Value)
					smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
					if !ok {
						ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
					}
					response, _, _, err := consumer.SendUpdateSmContextDeactivateUpCnxState(ocfUe, smContext, causeAll)
					if err != nil {
						ranUe.Log.Errorf("Send Update SmContextDeactivate UpCnxState Error[%s]", err.Error())
					} else if response == nil {
						ranUe.Log.Errorln("Send Update SmContextDeactivate UpCnxState Error")
					}
				}
			}
		} else {
			ranUe.Log.Info("Ue Context in Non GMM-Registered")
			ocfUe.SmContextList.Range(func(key, value interface{}) bool {
				smContext := value.(*context.SmContext)
				detail, err := consumer.SendReleaseSmContextRequest(ocfUe, smContext, &causeAll, "", nil)
				if err != nil {
					ranUe.Log.Errorf("Send ReleaseSmContextRequest Error[%s]", err.Error())
				} else if detail != nil {
					ranUe.Log.Errorf("Send ReleaseSmContextRequeste Error[%s]", detail.Cause)
				}
				return true
			})
			ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextReleaseUeContext, causeGroup, causeValue)
			return
		}
	}
	ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextN2NormalRelease, causeGroup, causeValue)
}

func HandleUEContextModificationResponse(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rRCState *ngapType.RRCState
	var userLocationInformation *ngapType.UserLocationInformation
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *context.RanUe

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}
	uEContextModificationResponse := successfulOutcome.Value.UEContextModificationResponse
	if uEContextModificationResponse == nil {
		ran.Log.Error("UEContextModificationResponse is nil")
		return
	}

	ran.Log.Info("Handle UE Context Modification Response")

	for _, ie := range uEContextModificationResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Warn("OcfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRRCState: // optional, ignore
			rRCState = ie.Value.RRCState
			ran.Log.Trace("Decode IE RRCState")
		case ngapType.ProtocolIEIDUserLocationInformation: // optional, ignore
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if rANUENGAPID != nil {
		ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			ran.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			ran.Log.Warnf("No UE Context[OcfUeNgapID: %d]", aMFUENGAPID.Value)
			return
		}
	}

	if ranUe != nil {
		ranUe.Log.Tracef("OcfUeNgapID[%d] RanUeNgapID[%d]", ranUe.OcfUeNgapId, ranUe.RanUeNgapId)

		if rRCState != nil {
			switch rRCState.Value {
			case ngapType.RRCStatePresentInactive:
				ranUe.Log.Trace("UE RRC State: Inactive")
			case ngapType.RRCStatePresentConnected:
				ranUe.Log.Trace("UE RRC State: Connected")
			}
		}

		if userLocationInformation != nil {
			ranUe.UpdateLocation(userLocationInformation)
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func HandleUEContextModificationFailure(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *context.RanUe

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		ran.Log.Error("UnsuccessfulOutcome is nil")
		return
	}
	uEContextModificationFailure := unsuccessfulOutcome.Value.UEContextModificationFailure
	if uEContextModificationFailure == nil {
		ran.Log.Error("UEContextModificationFailure is nil")
		return
	}

	ran.Log.Info("Handle UE Context Modification Failure")

	for _, ie := range uEContextModificationFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Warn("OcfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDCause: // ignore
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
			if cause == nil {
				ran.Log.Warn("Cause is nil")
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if rANUENGAPID != nil {
		ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			ran.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			ran.Log.Warnf("No UE Context[OcfUeNgapID: %d]", aMFUENGAPID.Value)
		}
	}

	if ranUe != nil {
		ran.Log.Tracef("OcfUeNgapID[%d] RanUeNgapID[%d]", ranUe.OcfUeNgapId, ranUe.RanUeNgapId)
	}

	if cause != nil {
		printAndGetCause(ran, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func HandleRRCInactiveTransitionReport(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rRCState *ngapType.RRCState
	var userLocationInformation *ngapType.UserLocationInformation

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}

	rRCInactiveTransitionReport := initiatingMessage.Value.RRCInactiveTransitionReport
	if rRCInactiveTransitionReport == nil {
		ran.Log.Error("RRCInactiveTransitionReport is nil")
		return
	}
	ran.Log.Info("Handle RRC Inactive Transition Report")

	for i := 0; i < len(rRCInactiveTransitionReport.ProtocolIEs.List); i++ {
		ie := rRCInactiveTransitionReport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // reject
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRRCState: // ignore
			rRCState = ie.Value.RRCState
			ran.Log.Trace("Decode IE RRCState")
			if rRCState == nil {
				ran.Log.Error("RRCState is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // ignore
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE UserLocationInformation")
			if userLocationInformation == nil {
				ran.Log.Error("UserLocationInformation is nil")
				return
			}
		}
	}

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	} else {
		ran.Log.Tracef("RANUENGAPID[%d] OCFUENGAPID[%d]", ranUe.RanUeNgapId, ranUe.OcfUeNgapId)

		if rRCState != nil {
			switch rRCState.Value {
			case ngapType.RRCStatePresentInactive:
				ran.Log.Trace("UE RRC State: Inactive")
			case ngapType.RRCStatePresentConnected:
				ran.Log.Trace("UE RRC State: Connected")
			}
		}
		ranUe.UpdateLocation(userLocationInformation)
	}
}

func HandleHandoverNotify(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	HandoverNotify := initiatingMessage.Value.HandoverNotify
	if HandoverNotify == nil {
		ran.Log.Error("HandoverNotify is nil")
		return
	}

	ran.Log.Info("Handle Handover notification")

	for i := 0; i < len(HandoverNotify.ProtocolIEs.List); i++ {
		ie := HandoverNotify.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OCFUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RANUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE userLocationInformation")
			if userLocationInformation == nil {
				ran.Log.Error("userLocationInformation is nil")
				return
			}
		}
	}

	targetUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)

	if targetUe == nil {
		ran.Log.Errorf("No RanUe Context[OcfUeNgapID: %d]", aMFUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	if userLocationInformation != nil {
		targetUe.UpdateLocation(userLocationInformation)
	}
	ocfUe := targetUe.OcfUe
	if ocfUe == nil {
		ran.Log.Error("OcfUe is nil")
		return
	}
	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: Send to S-OCF
		// Desciibed in (23.502 4.9.1.3.3) [conditional] 6a.Nocf_Communication_N2InfoNotify.
		ran.Log.Error("N2 Handover between OCF has not been implemented yet")
	} else {
		ran.Log.Info("Handle Handover notification Finshed ")
		for _, pduSessionid := range targetUe.SuccessPduSessionId {
			smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionid)
			if !ok {
				sourceUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionid)
			}
			_, _, _, err := consumer.SendUpdateSmContextN2HandoverComplete(ocfUe, smContext, "", nil)
			if err != nil {
				ran.Log.Errorf("Send UpdateSmContextN2HandoverComplete Error[%s]", err.Error())
			}
		}
		ocfUe.AttachRanUe(targetUe)

		ngap_message.SendUEContextReleaseCommand(sourceUe, context.UeContextReleaseHandover, ngapType.CausePresentNas,
			ngapType.CauseNasPresentNormalRelease)
	}

	// TODO: The UE initiates Mobility Registration Update procedure as described in clause 4.2.2.2.2.
}

// TS 23.502 4.9.1
func HandlePathSwitchRequest(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var rANUENGAPID *ngapType.RANUENGAPID
	var sourceOCFUENGAPID *ngapType.OCFUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation
	var uESecurityCapabilities *ngapType.UESecurityCapabilities
	var pduSessionResourceToBeSwitchedInDLList *ngapType.PDUSessionResourceToBeSwitchedDLList
	var pduSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListPSReq

	var ranUe *context.RanUe

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	pathSwitchRequest := initiatingMessage.Value.PathSwitchRequest
	if pathSwitchRequest == nil {
		ran.Log.Error("PathSwitchRequest is nil")
		return
	}

	ran.Log.Info("Handle Path Switch Request")

	for _, ie := range pathSwitchRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDSourceOCFUENGAPID: // reject
			sourceOCFUENGAPID = ie.Value.SourceOCFUENGAPID
			ran.Log.Trace("Decode IE SourceOcfUeNgapID")
			if sourceOCFUENGAPID == nil {
				ran.Log.Error("SourceOcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // ignore
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDUESecurityCapabilities: // ignore
			uESecurityCapabilities = ie.Value.UESecurityCapabilities
			ran.Log.Trace("Decode IE UESecurityCapabilities")
		case ngapType.ProtocolIEIDPDUSessionResourceToBeSwitchedDLList: // reject
			pduSessionResourceToBeSwitchedInDLList = ie.Value.PDUSessionResourceToBeSwitchedDLList
			ran.Log.Trace("Decode IE PDUSessionResourceToBeSwitchedDLList")
			if pduSessionResourceToBeSwitchedInDLList == nil {
				ran.Log.Error("PDUSessionResourceToBeSwitchedDLList is nil")
				return
			}
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListPSReq: // ignore
			pduSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListPSReq
			ran.Log.Trace("Decode IE PDUSessionResourceFailedToSetupListPSReq")
		}
	}

	if sourceOCFUENGAPID == nil {
		ran.Log.Error("SourceOcfUeNgapID is nil")
		return
	}
	ranUe = context.OCF_Self().RanUeFindByOcfUeNgapID(sourceOCFUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("Cannot find UE from sourceAMfUeNgapID[%d]", sourceOCFUENGAPID.Value)
		ngap_message.SendPathSwitchRequestFailure(ran, sourceOCFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	ran.Log.Tracef("OcfUeNgapID[%d] RanUeNgapID[%d]", ranUe.OcfUeNgapId, ranUe.RanUeNgapId)

	ocfUe := ranUe.OcfUe
	if ocfUe == nil {
		ranUe.Log.Error("OcfUe is nil")
		ngap_message.SendPathSwitchRequestFailure(ran, sourceOCFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	if ocfUe.SecurityContextIsValid() {
		// Update NH
		ocfUe.UpdateNH()
	} else {
		ranUe.Log.Errorf("No Security Context : SUPI[%s]", ocfUe.Supi)
		ngap_message.SendPathSwitchRequestFailure(ran, sourceOCFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	if uESecurityCapabilities != nil {
		ocfUe.UESecurityCapability.SetEA1_128_5G(uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes[0] & 0x80)
		ocfUe.UESecurityCapability.SetEA2_128_5G(uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes[0] & 0x40)
		ocfUe.UESecurityCapability.SetEA3_128_5G(uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes[0] & 0x20)
		ocfUe.UESecurityCapability.SetIA1_128_5G(uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes[0] & 0x80)
		ocfUe.UESecurityCapability.SetIA2_128_5G(uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes[0] & 0x40)
		ocfUe.UESecurityCapability.SetIA3_128_5G(uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes[0] & 0x20)
		// not support any E-UTRA algorithms
	}

	if rANUENGAPID != nil {
		ranUe.RanUeNgapId = rANUENGAPID.Value
	}

	ranUe.UpdateLocation(userLocationInformation)

	var pduSessionResourceSwitchedList ngapType.PDUSessionResourceSwitchedList
	var pduSessionResourceReleasedListPSAck ngapType.PDUSessionResourceReleasedListPSAck
	var pduSessionResourceReleasedListPSFail ngapType.PDUSessionResourceReleasedListPSFail

	if pduSessionResourceToBeSwitchedInDLList != nil {
		for _, item := range pduSessionResourceToBeSwitchedInDLList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PathSwitchRequestTransfer
			smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			response, errResponse, _, err := consumer.SendUpdateSmContextXnHandover(ocfUe, smContext,
				models.N2SmInfoType_PATH_SWITCH_REQ, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextXnHandover[PathSwitchRequestTransfer] Error:\n%s", err.Error())
			}
			if response != nil && response.BinaryDataN2SmInformation != nil {
				pduSessionResourceSwitchedItem := ngapType.PDUSessionResourceSwitchedItem{}
				pduSessionResourceSwitchedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceSwitchedItem.PathSwitchRequestAcknowledgeTransfer = response.BinaryDataN2SmInformation
				pduSessionResourceSwitchedList.List = append(pduSessionResourceSwitchedList.List, pduSessionResourceSwitchedItem)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSFail{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSFail.List = append(pduSessionResourceReleasedListPSFail.List,
					pduSessionResourceReleasedItem)
			}
		}
	}

	if pduSessionResourceFailedToSetupList != nil {
		for _, item := range pduSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PathSwitchRequestSetupFailedTransfer
			smContext, ok := ocfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			response, errResponse, _, err := consumer.SendUpdateSmContextXnHandoverFailed(ocfUe, smContext,
				models.N2SmInfoType_PATH_SWITCH_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextXnHandoverFailed[PathSwitchRequestSetupFailedTransfer] Error: %+v", err)
			}
			if response != nil && response.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSAck{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = response.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSAck.List = append(pduSessionResourceReleasedListPSAck.List,
					pduSessionResourceReleasedItem)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSFail{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSFail.List = append(pduSessionResourceReleasedListPSFail.List,
					pduSessionResourceReleasedItem)
			}
		}
	}

	// TS 23.502 4.9.1.2.2 step 7: send ack to Target NG-RAN. If none of the requested PDU Sessions have been switched
	// successfully, the OCF shall send an N2 Path Switch Request Failure message to the Target NG-RAN
	if len(pduSessionResourceSwitchedList.List) > 0 {
		// TODO: set newSecurityContextIndicator to true if there is a new security context
		err := ranUe.SwitchToRan(ran, rANUENGAPID.Value)
		if err != nil {
			ranUe.Log.Error(err.Error())
			return
		}
		ngap_message.SendPathSwitchRequestAcknowledge(ranUe, pduSessionResourceSwitchedList,
			pduSessionResourceReleasedListPSAck, false, nil, nil, nil)
	} else if len(pduSessionResourceReleasedListPSFail.List) > 0 {
		ngap_message.SendPathSwitchRequestFailure(ran, sourceOCFUENGAPID.Value, rANUENGAPID.Value,
			&pduSessionResourceReleasedListPSFail, nil)
	} else {
		ngap_message.SendPathSwitchRequestFailure(ran, sourceOCFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
	}
}

func HandleHandoverRequestAcknowledge(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceAdmittedList *ngapType.PDUSessionResourceAdmittedList
	var pDUSessionResourceFailedToSetupListHOAck *ngapType.PDUSessionResourceFailedToSetupListHOAck
	var targetToSourceTransparentContainer *ngapType.TargetToSourceTransparentContainer
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}
	handoverRequestAcknowledge := successfulOutcome.Value.HandoverRequestAcknowledge // reject
	if handoverRequestAcknowledge == nil {
		ran.Log.Error("HandoverRequestAcknowledge is nil")
		return
	}

	ran.Log.Info("Handle Handover Request Acknowledge")

	for _, ie := range handoverRequestAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceAdmittedList: // ignore
			pDUSessionResourceAdmittedList = ie.Value.PDUSessionResourceAdmittedList
			ran.Log.Trace("Decode IE PduSessionResourceAdmittedList")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListHOAck: // ignore
			pDUSessionResourceFailedToSetupListHOAck = ie.Value.PDUSessionResourceFailedToSetupListHOAck
			ran.Log.Trace("Decode IE PduSessionResourceFailedToSetupListHOAck")
		case ngapType.ProtocolIEIDTargetToSourceTransparentContainer: // reject
			targetToSourceTransparentContainer = ie.Value.TargetToSourceTransparentContainer
			ran.Log.Trace("Decode IE TargetToSourceTransparentContainer")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}
	if targetToSourceTransparentContainer == nil {
		ran.Log.Error("TargetToSourceTransparentContainer is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
			ngapType.ProtocolIEIDTargetToSourceTransparentContainer, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if len(iesCriticalityDiagnostics.List) > 0 {
		ran.Log.Error("Has missing reject IE(s)")

		procedureCode := ngapType.ProcedureCodeHandoverResourceAllocation
		triggeringMessage := ngapType.TriggeringMessagePresentSuccessfulOutcome
		procedureCriticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage,
			&procedureCriticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	targetUe := context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)
	if targetUe == nil {
		ran.Log.Errorf("No UE Context[OCFUENGAPID: %d]", aMFUENGAPID.Value)
		return
	}

	if rANUENGAPID != nil {
		targetUe.RanUeNgapId = rANUENGAPID.Value
	}
	ran.Log.Debugf("Target Ue RanUeNgapID[%d] OcfUeNgapID[%d]", targetUe.RanUeNgapId, targetUe.OcfUeNgapId)

	ocfUe := targetUe.OcfUe
	if ocfUe == nil {
		targetUe.Log.Error("ocfUe is nil")
		return
	}

	var pduSessionResourceHandoverList ngapType.PDUSessionResourceHandoverList
	var pduSessionResourceToReleaseList ngapType.PDUSessionResourceToReleaseListHOCmd

	// describe in 23.502 4.9.1.3.2 step11
	if pDUSessionResourceAdmittedList != nil {
		for _, item := range pDUSessionResourceAdmittedList.List {
			pduSessionID := item.PDUSessionID.Value
			transfer := item.HandoverRequestAcknowledgeTransfer
			pduSessionId := int32(pduSessionID)
			if smContext, exist := ocfUe.SmContextFindByPDUSessionID(pduSessionId); exist {
				response, errResponse, problemDetails, err := consumer.SendUpdateSmContextN2HandoverPrepared(ocfUe,
					smContext, models.N2SmInfoType_HANDOVER_REQ_ACK, transfer)
				if err != nil {
					targetUe.Log.Errorf("Send HandoverRequestAcknowledgeTransfer error: %v", err)
				}
				if problemDetails != nil {
					targetUe.Log.Warnf("ProblemDetails[status: %d, Cause: %s]", problemDetails.Status, problemDetails.Cause)
				}
				if response != nil && response.BinaryDataN2SmInformation != nil {
					handoverItem := ngapType.PDUSessionResourceHandoverItem{}
					handoverItem.PDUSessionID = item.PDUSessionID
					handoverItem.HandoverCommandTransfer = response.BinaryDataN2SmInformation
					pduSessionResourceHandoverList.List = append(pduSessionResourceHandoverList.List, handoverItem)
					targetUe.SuccessPduSessionId = append(targetUe.SuccessPduSessionId, pduSessionId)
				}
				if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
					releaseItem := ngapType.PDUSessionResourceToReleaseItemHOCmd{}
					releaseItem.PDUSessionID = item.PDUSessionID
					releaseItem.HandoverPreparationUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
					pduSessionResourceToReleaseList.List = append(pduSessionResourceToReleaseList.List, releaseItem)
				}
			}
		}
	}

	if pDUSessionResourceFailedToSetupListHOAck != nil {
		for _, item := range pDUSessionResourceFailedToSetupListHOAck.List {
			pduSessionID := item.PDUSessionID.Value
			transfer := item.HandoverResourceAllocationUnsuccessfulTransfer
			pduSessionId := int32(pduSessionID)
			if smContext, exist := ocfUe.SmContextFindByPDUSessionID(pduSessionId); exist {
				_, _, problemDetails, err := consumer.SendUpdateSmContextN2HandoverPrepared(ocfUe, smContext,
					models.N2SmInfoType_HANDOVER_RES_ALLOC_FAIL, transfer)
				if err != nil {
					targetUe.Log.Errorf("Send HandoverResourceAllocationUnsuccessfulTransfer error: %v", err)
				}
				if problemDetails != nil {
					targetUe.Log.Warnf("ProblemDetails[status: %d, Cause: %s]", problemDetails.Status, problemDetails.Cause)
				}
			}
		}
	}

	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: Send Nocf_Communication_CreateUEContext Response to S-OCF
		ran.Log.Error("handover between different Ue has not been implement yet")
	} else {
		ran.Log.Tracef("Source: RanUeNgapID[%d] OcfUeNgapID[%d]", sourceUe.RanUeNgapId, sourceUe.OcfUeNgapId)
		ran.Log.Tracef("Target: RanUeNgapID[%d] OcfUeNgapID[%d]", targetUe.RanUeNgapId, targetUe.OcfUeNgapId)
		if len(pduSessionResourceHandoverList.List) == 0 {
			targetUe.Log.Info("Handle Handover Preparation Failure [HoFailure In Target5GC NgranNode Or TargetSystem]")
			cause := &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem,
				},
			}
			ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
			return
		}
		ngap_message.SendHandoverCommand(sourceUe, pduSessionResourceHandoverList, pduSessionResourceToReleaseList,
			*targetToSourceTransparentContainer, nil)
	}
}

func HandleHandoverFailure(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var cause *ngapType.Cause
	var targetUe *context.RanUe
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}

	unsuccessfulOutcome := message.UnsuccessfulOutcome // reject
	if unsuccessfulOutcome == nil {
		ran.Log.Error("Unsuccessful Message is nil")
		return
	}

	handoverFailure := unsuccessfulOutcome.Value.HandoverFailure
	if handoverFailure == nil {
		ran.Log.Error("HandoverFailure is nil")
		return
	}

	for _, ie := range handoverFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
		case ngapType.ProtocolIEIDCause: // ignore
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	causePresent := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	if cause != nil {
		causePresent, causeValue = printAndGetCause(ran, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	targetUe = context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)

	if targetUe == nil {
		ran.Log.Errorf("No UE Context[OcfUENGAPID: %d]", aMFUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: handle N2 Handover between OCF
		ran.Log.Error("N2 Handover between OCF has not been implemented yet")
	} else {
		ocfUe := targetUe.OcfUe
		if ocfUe != nil {
			ocfUe.SmContextList.Range(func(key, value interface{}) bool {
				pduSessionID := key.(int32)
				smContext := value.(*context.SmContext)
				causeAll := context.CauseAll{
					NgapCause: &models.NgApCause{
						Group: int32(causePresent),
						Value: int32(causeValue),
					},
				}
				_, _, _, err := consumer.SendUpdateSmContextN2HandoverCanceled(ocfUe, smContext, causeAll)
				if err != nil {
					ran.Log.Errorf("Send UpdateSmContextN2HandoverCanceled Error for PduSessionId[%d]", pduSessionID)
				}
				return true
			})
		}
		ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, criticalityDiagnostics)
	}

	ngap_message.SendUEContextReleaseCommand(targetUe, context.UeContextReleaseHandover, causePresent, causeValue)
}

func HandleHandoverRequired(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var handoverType *ngapType.HandoverType
	var cause *ngapType.Cause
	var targetID *ngapType.TargetID
	var pDUSessionResourceListHORqd *ngapType.PDUSessionResourceListHORqd
	var sourceToTargetTransparentContainer *ngapType.SourceToTargetTransparentContainer
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	HandoverRequired := initiatingMessage.Value.HandoverRequired
	if HandoverRequired == nil {
		ran.Log.Error("HandoverRequired is nil")
		return
	}

	ran.Log.Info("Handle HandoverRequired\n")
	for i := 0; i < len(HandoverRequired.ProtocolIEs.List); i++ {
		ie := HandoverRequired.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID // reject
			ran.Log.Trace("Decode IE OcfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDHandoverType: // reject
			handoverType = ie.Value.HandoverType
			ran.Log.Trace("Decode IE HandoverType")
		case ngapType.ProtocolIEIDCause: // ignore
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
		case ngapType.ProtocolIEIDTargetID: // reject
			targetID = ie.Value.TargetID
			ran.Log.Trace("Decode IE TargetID")
		case ngapType.ProtocolIEIDPDUSessionResourceListHORqd: // reject
			pDUSessionResourceListHORqd = ie.Value.PDUSessionResourceListHORqd
			ran.Log.Trace("Decode IE PDUSessionResourceListHORqd")
		case ngapType.ProtocolIEIDSourceToTargetTransparentContainer: // reject
			sourceToTargetTransparentContainer = ie.Value.SourceToTargetTransparentContainer
			ran.Log.Trace("Decode IE SourceToTargetTransparentContainer")
		}
	}

	if aMFUENGAPID == nil {
		ran.Log.Error("OcfUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDOCFUENGAPID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if rANUENGAPID == nil {
		ran.Log.Error("RanUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}

	if handoverType == nil {
		ran.Log.Error("handoverType is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDHandoverType,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if targetID == nil {
		ran.Log.Error("targetID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDTargetID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if pDUSessionResourceListHORqd == nil {
		ran.Log.Error("pDUSessionResourceListHORqd is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
			ngapType.ProtocolIEIDPDUSessionResourceListHORqd, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if sourceToTargetTransparentContainer == nil {
		ran.Log.Error("sourceToTargetTransparentContainer is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
			ngapType.ProtocolIEIDSourceToTargetTransparentContainer, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procedureCode := ngapType.ProcedureCodeHandoverPreparation
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage,
			&procedureCriticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	sourceUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if sourceUe == nil {
		ran.Log.Errorf("Cannot find UE for RAN_UE_NGAP_ID[%d] ", rANUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}
	ocfUe := sourceUe.OcfUe
	if ocfUe == nil {
		ran.Log.Error("Cannot find ocfUE from sourceUE")
		return
	}

	if targetID.Present != ngapType.TargetIDPresentTargetRANNodeID {
		ran.Log.Errorf("targetID type[%d] is not supported", targetID.Present)
		return
	}
	ocfUe.OnGoing[sourceUe.Ran.AnType].Procedure = context.OnGoingProcedureN2Handover
	if !ocfUe.SecurityContextIsValid() {
		sourceUe.Log.Info("Handle Handover Preparation Failure [Authentication Failure]")
		cause = &ngapType.Cause{
			Present: ngapType.CausePresentNas,
			Nas: &ngapType.CauseNas{
				Value: ngapType.CauseNasPresentAuthenticationFailure,
			},
		}
		ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
		return
	}
	aMFSelf := context.OCF_Self()
	targetRanNodeId := ngapConvert.RanIdToModels(targetID.TargetRANNodeID.GlobalRANNodeID)
	targetRan, ok := aMFSelf.OcfRanFindByRanID(targetRanNodeId)
	if !ok {
		// handover between different OCF
		sourceUe.Log.Warnf("Handover required : cannot find target Ran Node Id[%+v] in this OCF", targetRanNodeId)
		sourceUe.Log.Error("Handover between different OCF has not been implemented yet")
		return
		// TODO: Send to T-OCF
		// Described in (23.502 4.9.1.3.2) step 3.Nocf_Communication_CreateUEContext Request
	} else {
		// Handover in same OCF
		sourceUe.HandOverType.Value = handoverType.Value
		tai := ngapConvert.TaiToModels(targetID.TargetRANNodeID.SelectedTAI)
		targetId := models.NgRanTargetId{
			RanNodeId: &targetRanNodeId,
			Tai:       &tai,
		}
		var pduSessionReqList ngapType.PDUSessionResourceSetupListHOReq
		for _, pDUSessionResourceHoItem := range pDUSessionResourceListHORqd.List {
			pduSessionId := int32(pDUSessionResourceHoItem.PDUSessionID.Value)
			if smContext, exist := ocfUe.SmContextFindByPDUSessionID(pduSessionId); exist {
				response, _, _, err := consumer.SendUpdateSmContextN2HandoverPreparing(ocfUe, smContext,
					models.N2SmInfoType_HANDOVER_REQUIRED, pDUSessionResourceHoItem.HandoverRequiredTransfer, "", &targetId)
				if err != nil {
					sourceUe.Log.Errorf("consumer.SendUpdateSmContextN2HandoverPreparing Error: %+v", err)
				}
				if response == nil {
					sourceUe.Log.Errorf("SendUpdateSmContextN2HandoverPreparing Error for PduSessionId[%d]", pduSessionId)
					continue
				} else if response.BinaryDataN2SmInformation != nil {
					ngap_message.AppendPDUSessionResourceSetupListHOReq(&pduSessionReqList, pduSessionId,
						smContext.Snssai(), response.BinaryDataN2SmInformation)
				}
			}
		}
		if len(pduSessionReqList.List) == 0 {
			sourceUe.Log.Info("Handle Handover Preparation Failure [HoFailure In Target5GC NgranNode Or TargetSystem]")
			cause = &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem,
				},
			}
			ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
			return
		}
		// Update NH
		ocfUe.UpdateNH()
		ngap_message.SendHandoverRequest(sourceUe, targetRan, *cause, pduSessionReqList,
			*sourceToTargetTransparentContainer, false)
	}
}

func HandleHandoverCancel(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	HandoverCancel := initiatingMessage.Value.HandoverCancel
	if HandoverCancel == nil {
		ran.Log.Error("Handover Cancel is nil")
		return
	}

	ran.Log.Info("Handle Handover Cancel")
	for i := 0; i < len(HandoverCancel.ProtocolIEs.List); i++ {
		ie := HandoverCancel.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OCFUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RANUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
			if cause == nil {
				ran.Log.Error(cause, "cause is nil")
				return
			}
		}
	}

	sourceUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if sourceUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	if sourceUe.OcfUeNgapId != aMFUENGAPID.Value {
		ran.Log.Warnf("Conflict OCF_UE_NGAP_ID : %d != %d", sourceUe.OcfUeNgapId, aMFUENGAPID.Value)
	}
	ran.Log.Tracef("Source: RAN_UE_NGAP_ID[%d] OCF_UE_NGAP_ID[%d]", sourceUe.RanUeNgapId, sourceUe.OcfUeNgapId)

	causePresent := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	if cause != nil {
		causePresent, causeValue = printAndGetCause(ran, cause)
	}
	targetUe := sourceUe.TargetUe
	if targetUe == nil {
		// Described in (23.502 4.11.1.2.3) step 2
		// Todo : send to T-OCF invoke Nocf_UeContextReleaseRequest(targetUe)
		ran.Log.Error("N2 Handover between OCF has not been implemented yet")
	} else {
		ran.Log.Tracef("Target : RAN_UE_NGAP_ID[%d] OCF_UE_NGAP_ID[%d]", targetUe.RanUeNgapId, targetUe.OcfUeNgapId)
		ocfUe := sourceUe.OcfUe
		if ocfUe != nil {
			ocfUe.SmContextList.Range(func(key, value interface{}) bool {
				pduSessionID := key.(int32)
				smContext := value.(*context.SmContext)
				causeAll := context.CauseAll{
					NgapCause: &models.NgApCause{
						Group: int32(causePresent),
						Value: int32(causeValue),
					},
				}
				_, _, _, err := consumer.SendUpdateSmContextN2HandoverCanceled(ocfUe, smContext, causeAll)
				if err != nil {
					sourceUe.Log.Errorf("Send UpdateSmContextN2HandoverCanceled Error for PduSessionId[%d]", pduSessionID)
				}
				return true
			})
		}
		ngap_message.SendUEContextReleaseCommand(targetUe, context.UeContextReleaseHandover, causePresent, causeValue)
		ngap_message.SendHandoverCancelAcknowledge(sourceUe, nil)
	}
}

func HandleUplinkRanStatusTransfer(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rANStatusTransferTransparentContainer *ngapType.RANStatusTransferTransparentContainer
	var ranUe *context.RanUe

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	uplinkRanStatusTransfer := initiatingMessage.Value.UplinkRANStatusTransfer
	if uplinkRanStatusTransfer == nil {
		ran.Log.Error("UplinkRanStatusTransfer is nil")
		return
	}

	ran.Log.Info("Handle Uplink Ran Status Transfer")

	for _, ie := range uplinkRanStatusTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // reject
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANStatusTransferTransparentContainer: // reject
			rANStatusTransferTransparentContainer = ie.Value.RANStatusTransferTransparentContainer
			ran.Log.Trace("Decode IE RANStatusTransferTransparentContainer")
			if rANStatusTransferTransparentContainer == nil {
				ran.Log.Error("RANStatusTransferTransparentContainer is nil")
			}
		}
	}

	ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	ranUe.Log.Tracef("UE Context OcfUeNgapID[%d] RanUeNgapID[%d]", ranUe.OcfUeNgapId, ranUe.RanUeNgapId)

	ocfUe := ranUe.OcfUe
	if ocfUe == nil {
		ranUe.Log.Error("OcfUe is nil")
		return
	}
	// send to T-OCF using N1N2MessageTransfer (R16)
}

func HandleNasNonDeliveryIndication(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var cause *ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	nASNonDeliveryIndication := initiatingMessage.Value.NASNonDeliveryIndication
	if nASNonDeliveryIndication == nil {
		ran.Log.Error("NASNonDeliveryIndication is nil")
		return
	}

	ran.Log.Info("Handle Nas Non Delivery Indication")

	for _, ie := range nASNonDeliveryIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDNASPDU:
			nASPDU = ie.Value.NASPDU
			if nASPDU == nil {
				ran.Log.Error("NasPdu is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			if cause == nil {
				ran.Log.Error("Cause is nil")
				return
			}
		}
	}

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	ran.Log.Tracef("RanUeNgapID[%d] OcfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.OcfUeNgapId)

	printAndGetCause(ran, cause)

	nas.HandleNAS(ranUe, ngapType.ProcedureCodeNASNonDeliveryIndication, nASPDU.Value)
}

func HandleRanConfigurationUpdate(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	var cause ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}

	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	rANConfigurationUpdate := initiatingMessage.Value.RANConfigurationUpdate
	if rANConfigurationUpdate == nil {
		ran.Log.Error("RAN Configuration is nil")
		return
	}
	ran.Log.Info("Handle Ran Configuration Update")
	for i := 0; i < len(rANConfigurationUpdate.ProtocolIEs.List); i++ {
		ie := rANConfigurationUpdate.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			if rANNodeName == nil {
				ran.Log.Error("RAN Node Name is nil")
				return
			}
			ran.Log.Tracef("Decode IE RANNodeName = [%s]", rANNodeName.Value)
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			ran.Log.Trace("Decode IE SupportedTAList")
			if supportedTAList == nil {
				ran.Log.Error("Supported TA List is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			if pagingDRX == nil {
				ran.Log.Error("PagingDRX is nil")
				return
			}
			ran.Log.Tracef("Decode IE PagingDRX = [%d]", pagingDRX.Value)
		}
	}

	for i := 0; i < len(supportedTAList.List); i++ {
		supportedTAItem := supportedTAList.List[i]
		tac := hex.EncodeToString(supportedTAItem.TAC.Value)
		capOfSupportTai := cap(ran.SupportedTAList)
		for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
			supportedTAI := context.NewSupportedTAI()
			supportedTAI.Tai.Tac = tac
			broadcastPLMNItem := supportedTAItem.BroadcastPLMNList.List[j]
			plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
			supportedTAI.Tai.PlmnId = &plmnId
			capOfSNssaiList := cap(supportedTAI.SNssaiList)
			for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList.List); k++ {
				tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList.List[k]
				if len(supportedTAI.SNssaiList) < capOfSNssaiList {
					supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
				} else {
					break
				}
			}
			ran.Log.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
			if len(ran.SupportedTAList) < capOfSupportTai {
				ran.SupportedTAList = append(ran.SupportedTAList, supportedTAI)
			} else {
				break
			}
		}
	}

	if len(ran.SupportedTAList) == 0 {
		ran.Log.Warn("RanConfigurationUpdate failure: No supported TA exist in RanConfigurationUpdate")
		cause.Present = ngapType.CausePresentMisc
		cause.Misc = &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		}
	} else {
		var found bool
		for i, tai := range ran.SupportedTAList {
			if context.InTaiList(tai.Tai, context.OCF_Self().SupportTaiLists) {
				ran.Log.Tracef("SERVED_TAI_INDEX[%d]", i)
				found = true
				break
			}
		}
		if !found {
			ran.Log.Warn("RanConfigurationUpdate failure: Cannot find Served TAI in OCF")
			cause.Present = ngapType.CausePresentMisc
			cause.Misc = &ngapType.CauseMisc{
				Value: ngapType.CauseMiscPresentUnknownPLMN,
			}
		}
	}

	if cause.Present == ngapType.CausePresentNothing {
		ran.Log.Info("Handle RanConfigurationUpdateAcknowledge")
		ngap_message.SendRanConfigurationUpdateAcknowledge(ran, nil)
	} else {
		ran.Log.Info("Handle RanConfigurationUpdateAcknowledgeFailure")
		ngap_message.SendRanConfigurationUpdateFailure(ran, cause, nil)
	}
}

func HandleUplinkRanConfigurationTransfer(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var sONConfigurationTransferUL *ngapType.SONConfigurationTransfer

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	uplinkRANConfigurationTransfer := initiatingMessage.Value.UplinkRANConfigurationTransfer
	if uplinkRANConfigurationTransfer == nil {
		ran.Log.Error("ErrorIndication is nil")
		return
	}

	for _, ie := range uplinkRANConfigurationTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDSONConfigurationTransferUL: // optional, ignore
			sONConfigurationTransferUL = ie.Value.SONConfigurationTransferUL
			ran.Log.Trace("Decode IE SONConfigurationTransferUL")
			if sONConfigurationTransferUL == nil {
				ran.Log.Warn("sONConfigurationTransferUL is nil")
			}
		}
	}

	if sONConfigurationTransferUL != nil {
		targetRanNodeID := ngapConvert.RanIdToModels(sONConfigurationTransferUL.TargetRANNodeID.GlobalRANNodeID)

		if targetRanNodeID.GNbId.GNBValue != "" {
			ran.Log.Tracef("targerRanID [%s]", targetRanNodeID.GNbId.GNBValue)
		}

		aMFSelf := context.OCF_Self()

		targetRan, ok := aMFSelf.OcfRanFindByRanID(targetRanNodeID)
		if !ok {
			ran.Log.Warn("targetRan is nil")
		}

		ngap_message.SendDownlinkRanConfigurationTransfer(targetRan, sONConfigurationTransferUL)
	}
}

func HandleUplinkUEAssociatedNRPPATransport(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var routingID *ngapType.RoutingID
	var nRPPaPDU *ngapType.NRPPaPDU

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	uplinkUEAssociatedNRPPaTransport := initiatingMessage.Value.UplinkUEAssociatedNRPPaTransport
	if uplinkUEAssociatedNRPPaTransport == nil {
		ran.Log.Error("uplinkUEAssociatedNRPPaTransport is nil")
		return
	}

	ran.Log.Info("Handle Uplink UE Associated NRPPA Transpor")

	for _, ie := range uplinkUEAssociatedNRPPaTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // reject
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE aMFUENGAPID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE rANUENGAPID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRoutingID: // reject
			routingID = ie.Value.RoutingID
			ran.Log.Trace("Decode IE routingID")
			if routingID == nil {
				ran.Log.Error("routingID is nil")
				return
			}
		case ngapType.ProtocolIEIDNRPPaPDU: // reject
			nRPPaPDU = ie.Value.NRPPaPDU
			ran.Log.Trace("Decode IE nRPPaPDU")
			if nRPPaPDU == nil {
				ran.Log.Error("nRPPaPDU is nil")
				return
			}
		}
	}

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	ran.Log.Tracef("RanUeNgapId[%d] OcfUeNgapId[%d]", ranUe.RanUeNgapId, ranUe.OcfUeNgapId)

	ranUe.RoutingID = hex.EncodeToString(routingID.Value)

	// TODO: Forward NRPPaPDU to LMF
}

func HandleUplinkNonUEAssociatedNRPPATransport(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var routingID *ngapType.RoutingID
	var nRPPaPDU *ngapType.NRPPaPDU

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	uplinkNonUEAssociatedNRPPATransport := initiatingMessage.Value.UplinkNonUEAssociatedNRPPaTransport
	if uplinkNonUEAssociatedNRPPATransport == nil {
		ran.Log.Error("Uplink Non UE Associated NRPPA Transport is nil")
		return
	}

	ran.Log.Info("Handle Uplink Non UE Associated NRPPA Transport")

	for i := 0; i < len(uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List); i++ {
		ie := uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRoutingID:
			routingID = ie.Value.RoutingID
			ran.Log.Trace("Decode IE RoutingID")

		case ngapType.ProtocolIEIDNRPPaPDU:
			nRPPaPDU = ie.Value.NRPPaPDU
			ran.Log.Trace("Decode IE NRPPaPDU")
		}
	}

	if routingID == nil {
		ran.Log.Error("RoutingID is nil")
		return
	}
	// Forward routingID to LMF
	// Described in (23.502 4.13.5.6)

	if nRPPaPDU == nil {
		ran.Log.Error("NRPPaPDU is nil")
		return
	}
	// TODO: Forward NRPPaPDU to LMF
}

func HandleLocationReport(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation
	var uEPresenceInAreaOfInterestList *ngapType.UEPresenceInAreaOfInterestList
	var locationReportingRequestType *ngapType.LocationReportingRequestType

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	locationReport := initiatingMessage.Value.LocationReport
	if locationReport == nil {
		ran.Log.Error("LocationReport is nil")
		return
	}

	ran.Log.Info("Handle Location Report")
	for _, ie := range locationReport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // reject
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // ignore
			userLocationInformation = ie.Value.UserLocationInformation
			ran.Log.Trace("Decode IE userLocationInformation")
			if userLocationInformation == nil {
				ran.Log.Warn("userLocationInformation is nil")
			}
		case ngapType.ProtocolIEIDUEPresenceInAreaOfInterestList: // optional, ignore
			uEPresenceInAreaOfInterestList = ie.Value.UEPresenceInAreaOfInterestList
			ran.Log.Trace("Decode IE uEPresenceInAreaOfInterestList")
			if uEPresenceInAreaOfInterestList == nil {
				ran.Log.Warn("uEPresenceInAreaOfInterestList is nil [optional]")
			}
		case ngapType.ProtocolIEIDLocationReportingRequestType: // ignore
			locationReportingRequestType = ie.Value.LocationReportingRequestType
			ran.Log.Trace("Decode IE LocationReportingRequestType")
			if locationReportingRequestType == nil {
				ran.Log.Warn("LocationReportingRequestType is nil")
			}
		}
	}

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	ranUe.UpdateLocation(userLocationInformation)

	ranUe.Log.Tracef("Report Area[%d]", locationReportingRequestType.ReportArea.Value)

	switch locationReportingRequestType.EventType.Value {
	case ngapType.EventTypePresentDirect:
		ranUe.Log.Trace("To report directly")

	case ngapType.EventTypePresentChangeOfServeCell:
		ranUe.Log.Trace("To report upon change of serving cell")

	case ngapType.EventTypePresentUePresenceInAreaOfInterest:
		ranUe.Log.Trace("To report UE presence in the area of interest")
		for _, uEPresenceInAreaOfInterestItem := range uEPresenceInAreaOfInterestList.List {
			uEPresence := uEPresenceInAreaOfInterestItem.UEPresence.Value
			referenceID := uEPresenceInAreaOfInterestItem.LocationReportingReferenceID.Value

			for _, AOIitem := range locationReportingRequestType.AreaOfInterestList.List {
				if referenceID == AOIitem.LocationReportingReferenceID.Value {
					ran.Log.Tracef("uEPresence[%d], presence AOI ReferenceID[%d]", uEPresence, referenceID)
				}
			}
		}

	case ngapType.EventTypePresentStopChangeOfServeCell:
		ranUe.Log.Trace("To stop reporting at change of serving cell")
		ngap_message.SendLocationReportingControl(ranUe, nil, 0, locationReportingRequestType.EventType)
		// TODO: Clear location report

	case ngapType.EventTypePresentStopUePresenceInAreaOfInterest:
		ranUe.Log.Trace("To stop reporting UE presence in the area of interest")
		ranUe.Log.Tracef("ReferenceID To Be Cancelled[%d]",
			locationReportingRequestType.LocationReportingReferenceIDToBeCancelled.Value)
		// TODO: Clear location report

	case ngapType.EventTypePresentCancelLocationReportingForTheUe:
		ranUe.Log.Trace("To cancel location reporting for the UE")
		// TODO: Clear location report
	}
}

func HandleUERadioCapabilityInfoIndication(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	var uERadioCapability *ngapType.UERadioCapability
	var uERadioCapabilityForPaging *ngapType.UERadioCapabilityForPaging

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	uERadioCapabilityInfoIndication := initiatingMessage.Value.UERadioCapabilityInfoIndication
	if uERadioCapabilityInfoIndication == nil {
		ran.Log.Error("UERadioCapabilityInfoIndication is nil")
		return
	}

	ran.Log.Info("Handle UE Radio Capability Info Indication")

	for i := 0; i < len(uERadioCapabilityInfoIndication.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityInfoIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDUERadioCapability:
			uERadioCapability = ie.Value.UERadioCapability
			ran.Log.Trace("Decode IE UERadioCapability")
			if uERadioCapability == nil {
				ran.Log.Error("UERadioCapability is nil")
				return
			}
		case ngapType.ProtocolIEIDUERadioCapabilityForPaging:
			uERadioCapabilityForPaging = ie.Value.UERadioCapabilityForPaging
			ran.Log.Trace("Decode IE UERadioCapabilityForPaging")
			if uERadioCapabilityForPaging == nil {
				ran.Log.Error("UERadioCapabilityForPaging is nil")
				return
			}
		}
	}

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	ran.Log.Tracef("RanUeNgapID[%d] OcfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.OcfUeNgapId)
	ocfUe := ranUe.OcfUe

	if ocfUe == nil {
		ranUe.Log.Errorln("ocfUe is nil")
		return
	}
	if uERadioCapability != nil {
		ocfUe.UeRadioCapability = hex.EncodeToString(uERadioCapability.Value)
	}
	if uERadioCapabilityForPaging != nil {
		ocfUe.UeRadioCapabilityForPaging = &context.UERadioCapabilityForPaging{}
		if uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR != nil {
			ocfUe.UeRadioCapabilityForPaging.NR = hex.EncodeToString(
				uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR.Value)
		}
		if uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA != nil {
			ocfUe.UeRadioCapabilityForPaging.EUTRA = hex.EncodeToString(
				uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA.Value)
		}
	}

	// TS 38.413 8.14.1.2/TS 23.502 4.2.8a step5/TS 23.501, clause 5.4.4.1.
	// send its most up to date UE Radio Capability information to the RAN in the N2 REQUEST message.
}

func HandleOCFconfigurationUpdateFailure(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		ran.Log.Error("Unsuccessful Message is nil")
		return
	}

	OCFconfigurationUpdateFailure := unsuccessfulOutcome.Value.OCFConfigurationUpdateFailure
	if OCFconfigurationUpdateFailure == nil {
		ran.Log.Error("OCFConfigurationUpdateFailure is nil")
		return
	}

	ran.Log.Info("Handle OCF Confioguration Update Failure")

	for _, ie := range OCFconfigurationUpdateFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
			if cause == nil {
				ran.Log.Error("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	//	TODO: Time To Wait

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func HandleOCFconfigurationUpdateAcknowledge(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFTNLAssociationSetupList *ngapType.OCFTNLAssociationSetupList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var aMFTNLAssociationFailedToSetupList *ngapType.TNLAssociationList
	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		ran.Log.Error("SuccessfulOutcome is nil")
		return
	}
	aMFConfigurationUpdateAcknowledge := successfulOutcome.Value.OCFConfigurationUpdateAcknowledge
	if aMFConfigurationUpdateAcknowledge == nil {
		ran.Log.Error("OCFConfigurationUpdateAcknowledge is nil")
		return
	}

	ran.Log.Info("Handle OCF Configuration Update Acknowledge")

	for i := 0; i < len(aMFConfigurationUpdateAcknowledge.ProtocolIEs.List); i++ {
		ie := aMFConfigurationUpdateAcknowledge.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFTNLAssociationSetupList:
			aMFTNLAssociationSetupList = ie.Value.OCFTNLAssociationSetupList
			ran.Log.Trace("Decode IE OCFTNLAssociationSetupList")
			if aMFTNLAssociationSetupList == nil {
				ran.Log.Error("OCFTNLAssociationSetupList is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE Criticality Diagnostics")

		case ngapType.ProtocolIEIDOCFTNLAssociationFailedToSetupList:
			aMFTNLAssociationFailedToSetupList = ie.Value.OCFTNLAssociationFailedToSetupList
			ran.Log.Trace("Decode IE OCFTNLAssociationFailedToSetupList")
			if aMFTNLAssociationFailedToSetupList == nil {
				ran.Log.Error("OCFTNLAssociationFailedToSetupList is nil")
				return
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func HandleErrorIndication(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	errorIndication := initiatingMessage.Value.ErrorIndication
	if errorIndication == nil {
		ran.Log.Error("ErrorIndication is nil")
		return
	}

	for _, ie := range errorIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID:
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
			if aMFUENGAPID == nil {
				ran.Log.Error("OcfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				ran.Log.Error("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			ran.Log.Trace("Decode IE Cause")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			ran.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if cause == nil && criticalityDiagnostics == nil {
		ran.Log.Error("[ErrorIndication] both Cause IE and CriticalityDiagnostics IE are nil, should have at least one")
		return
	}

	if cause != nil {
		printAndGetCause(ran, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	// TODO: handle error based on cause/criticalityDiagnostics
}

func HandleCellTrafficTrace(ran *context.OcfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.OCFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nGRANTraceID *ngapType.NGRANTraceID
	var nGRANCGI *ngapType.NGRANCGI
	var traceCollectionEntityIPAddress *ngapType.TransportLayerAddress

	var ranUe *context.RanUe

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		ran.Log.Error("InitiatingMessage is nil")
		return
	}
	cellTrafficTrace := initiatingMessage.Value.CellTrafficTrace
	if cellTrafficTrace == nil {
		ran.Log.Error("CellTrafficTrace is nil")
		return
	}

	ran.Log.Info("Handle Cell Traffic Trace")

	for _, ie := range cellTrafficTrace.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDOCFUENGAPID: // reject
			aMFUENGAPID = ie.Value.OCFUENGAPID
			ran.Log.Trace("Decode IE OcfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			ran.Log.Trace("Decode IE RanUeNgapID")

		case ngapType.ProtocolIEIDNGRANTraceID: // ignore
			nGRANTraceID = ie.Value.NGRANTraceID
			ran.Log.Trace("Decode IE NGRANTraceID")
		case ngapType.ProtocolIEIDNGRANCGI: // ignore
			nGRANCGI = ie.Value.NGRANCGI
			ran.Log.Trace("Decode IE NGRANCGI")
		case ngapType.ProtocolIEIDTraceCollectionEntityIPAddress: // ignore
			traceCollectionEntityIPAddress = ie.Value.TraceCollectionEntityIPAddress
			ran.Log.Trace("Decode IE TraceCollectionEntityIPAddress")
		}
	}
	if aMFUENGAPID == nil {
		ran.Log.Error("OcfUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDOCFUENGAPID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if rANUENGAPID == nil {
		ran.Log.Error("RanUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procedureCode := ngapType.ProcedureCodeCellTrafficTrace
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentIgnore
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality,
			&iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	if aMFUENGAPID != nil {
		ranUe = context.OCF_Self().RanUeFindByOcfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			ran.Log.Errorf("No UE Context[OcfUeNgapID: %d]", aMFUENGAPID.Value)
			cause := ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
				},
			}
			ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
			return
		}
	}

	ran.Log.Debugf("UE: OcfUeNgapID[%d], RanUeNgapID[%d]", ranUe.OcfUeNgapId, ranUe.RanUeNgapId)

	ranUe.Trsr = hex.EncodeToString(nGRANTraceID.Value[6:])

	ranUe.Log.Tracef("TRSR[%s]", ranUe.Trsr)

	switch nGRANCGI.Present {
	case ngapType.NGRANCGIPresentNRCGI:
		plmnID := ngapConvert.PlmnIdToModels(nGRANCGI.NRCGI.PLMNIdentity)
		cellID := ngapConvert.BitStringToHex(&nGRANCGI.NRCGI.NRCellIdentity.Value)
		ranUe.Log.Debugf("NRCGI[plmn: %s, cellID: %s]", plmnID, cellID)
	case ngapType.NGRANCGIPresentEUTRACGI:
		plmnID := ngapConvert.PlmnIdToModels(nGRANCGI.EUTRACGI.PLMNIdentity)
		cellID := ngapConvert.BitStringToHex(&nGRANCGI.EUTRACGI.EUTRACellIdentity.Value)
		ranUe.Log.Debugf("EUTRACGI[plmn: %s, cellID: %s]", plmnID, cellID)
	}

	tceIpv4, tceIpv6 := ngapConvert.IPAddressToString(*traceCollectionEntityIPAddress)
	if tceIpv4 != "" {
		ranUe.Log.Debugf("TCE IP Address[v4: %s]", tceIpv4)
	}
	if tceIpv6 != "" {
		ranUe.Log.Debugf("TCE IP Address[v6: %s]", tceIpv6)
	}

	// TODO: TS 32.422 4.2.2.10
	// When OCF receives this new NG signalling message containing the Trace Recording Session Reference (TRSR)
	// and Trace Reference (TR), the OCF shall look up the SUPI/IMEI(SV) of the given call from its database and
	// shall send the SUPI/IMEI(SV) numbers together with the Trace Recording Session Reference and Trace Reference
	// to the Trace Collection Entity.
}

func printAndGetCause(ran *context.OcfRan, cause *ngapType.Cause) (present int, value aper.Enumerated) {
	present = cause.Present
	switch cause.Present {
	case ngapType.CausePresentRadioNetwork:
		ran.Log.Warnf("Cause RadioNetwork[%d]", cause.RadioNetwork.Value)
		value = cause.RadioNetwork.Value
	case ngapType.CausePresentTransport:
		ran.Log.Warnf("Cause Transport[%d]", cause.Transport.Value)
		value = cause.Transport.Value
	case ngapType.CausePresentProtocol:
		ran.Log.Warnf("Cause Protocol[%d]", cause.Protocol.Value)
		value = cause.Protocol.Value
	case ngapType.CausePresentNas:
		ran.Log.Warnf("Cause Nas[%d]", cause.Nas.Value)
		value = cause.Nas.Value
	case ngapType.CausePresentMisc:
		ran.Log.Warnf("Cause Misc[%d]", cause.Misc.Value)
		value = cause.Misc.Value
	default:
		ran.Log.Errorf("Invalid Cause group[%d]", cause.Present)
	}
	return
}

func printCriticalityDiagnostics(ran *context.OcfRan, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {
	ran.Log.Trace("Criticality Diagnostics")

	if criticalityDiagnostics.ProcedureCriticality != nil {
		switch criticalityDiagnostics.ProcedureCriticality.Value {
		case ngapType.CriticalityPresentReject:
			ran.Log.Trace("Procedure Criticality: Reject")
		case ngapType.CriticalityPresentIgnore:
			ran.Log.Trace("Procedure Criticality: Ignore")
		case ngapType.CriticalityPresentNotify:
			ran.Log.Trace("Procedure Criticality: Notify")
		}
	}

	if criticalityDiagnostics.IEsCriticalityDiagnostics != nil {
		for _, ieCriticalityDiagnostics := range criticalityDiagnostics.IEsCriticalityDiagnostics.List {
			ran.Log.Tracef("IE ID: %d", ieCriticalityDiagnostics.IEID.Value)

			switch ieCriticalityDiagnostics.IECriticality.Value {
			case ngapType.CriticalityPresentReject:
				ran.Log.Trace("Criticality Reject")
			case ngapType.CriticalityPresentNotify:
				ran.Log.Trace("Criticality Notify")
			}

			switch ieCriticalityDiagnostics.TypeOfError.Value {
			case ngapType.TypeOfErrorPresentNotUnderstood:
				ran.Log.Trace("Type of error: Not understood")
			case ngapType.TypeOfErrorPresentMissing:
				ran.Log.Trace("Type of error: Missing")
			}
		}
	}
}

func buildCriticalityDiagnostics(
	procedureCode *int64,
	triggeringMessage *aper.Enumerated,
	procedureCriticality *aper.Enumerated,
	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList) (
	criticalityDiagnostics ngapType.CriticalityDiagnostics) {
	if procedureCode != nil {
		criticalityDiagnostics.ProcedureCode = new(ngapType.ProcedureCode)
		criticalityDiagnostics.ProcedureCode.Value = *procedureCode
	}

	if triggeringMessage != nil {
		criticalityDiagnostics.TriggeringMessage = new(ngapType.TriggeringMessage)
		criticalityDiagnostics.TriggeringMessage.Value = *triggeringMessage
	}

	if procedureCriticality != nil {
		criticalityDiagnostics.ProcedureCriticality = new(ngapType.Criticality)
		criticalityDiagnostics.ProcedureCriticality.Value = *procedureCriticality
	}

	if iesCriticalityDiagnostics != nil {
		criticalityDiagnostics.IEsCriticalityDiagnostics = iesCriticalityDiagnostics
	}

	return criticalityDiagnostics
}

func buildCriticalityDiagnosticsIEItem(ieCriticality aper.Enumerated, ieID int64, typeOfErr aper.Enumerated) (
	item ngapType.CriticalityDiagnosticsIEItem) {
	item = ngapType.CriticalityDiagnosticsIEItem{
		IECriticality: ngapType.Criticality{
			Value: ieCriticality,
		},
		IEID: ngapType.ProtocolIEID{
			Value: ieID,
		},
		TypeOfError: ngapType.TypeOfError{
			Value: typeOfErr,
		},
	}

	return item
}

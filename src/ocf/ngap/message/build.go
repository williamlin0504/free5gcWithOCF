package message

import (
	"encoding/binary"
	"encoding/hex"

	"free5gcWithOCF/lib/aper"
	"free5gcWithOCF/lib/ngap"
	"free5gcWithOCF/lib/ngap/ngapConvert"
	"free5gcWithOCF/lib/ngap/ngapType"
	"free5gcWithOCF/src/ocf/context"
	"free5gcWithOCF/src/ocf/logger"
	"free5gcWithOCF/src/ocf/util"
)

func BuildNGSetupRequest() ([]byte, error) {

	ocfSelf := context.OCFSelf()
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeNGSetup
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentNGSetupRequest
	initiatingMessage.Value.NGSetupRequest = new(ngapType.NGSetupRequest)

	nGSetupRequest := initiatingMessage.Value.NGSetupRequest
	nGSetupRequestIEs := &nGSetupRequest.ProtocolIEs

	// GlobalRANNodeID
	ie := ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDGlobalRANNodeID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentGlobalRANNodeID
	ie.Value.GlobalRANNodeID = new(ngapType.GlobalRANNodeID)

	globalRANNodeID := ie.Value.GlobalRANNodeID
	globalRANNodeID.Present = ngapType.GlobalRANNodeIDPresentGlobalOCFID
	globalRANNodeID.GlobalOCFID = new(ngapType.GlobalOCFID)

	globalOCFID := globalRANNodeID.GlobalOCFID
	globalOCFID.PLMNIdentity = util.PlmnIdToNgap(ocfSelf.NFInfo.GlobalOCFID.PLMNID)
	globalOCFID.OCFID.Present = ngapType.OCFIDPresentOCFID
	globalOCFID.OCFID.OCFID = util.OcfIdToNgap(ocfSelf.NFInfo.GlobalOCFID.OCFID)
	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	// RANNodeName
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANNodeName
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentRANNodeName
	ie.Value.RANNodeName = new(ngapType.RANNodeName)

	rANNodeName := ie.Value.RANNodeName
	rANNodeName.Value = ocfSelf.NFInfo.RanNodeName
	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)
	// SupportedTAList
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSupportedTAList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentSupportedTAList
	ie.Value.SupportedTAList = new(ngapType.SupportedTAList)

	supportedTAList := ie.Value.SupportedTAList

	for _, supportedTAItemLocal := range ocfSelf.NFInfo.SupportedTAList {
		// SupportedTAItem in SupportedTAList
		supportedTAItem := ngapType.SupportedTAItem{}
		var err error
		supportedTAItem.TAC.Value, err = hex.DecodeString(supportedTAItemLocal.TAC)
		if err != nil {
			logger.NgapLog.Errorf("DecodeString error: %+v", err)
		}

		broadcastPLMNList := &supportedTAItem.BroadcastPLMNList

		for _, broadcastPLMNListLocal := range supportedTAItemLocal.BroadcastPLMNList {
			// BroadcastPLMNItem in BroadcastPLMNList
			broadcastPLMNItem := ngapType.BroadcastPLMNItem{}
			broadcastPLMNItem.PLMNIdentity = util.PlmnIdToNgap(broadcastPLMNListLocal.PLMNID)

			sliceSupportList := &broadcastPLMNItem.TAISliceSupportList

			for _, sliceSupportItemLocal := range broadcastPLMNListLocal.TAISliceSupportList {
				// SliceSupportItem in SliceSupportList
				sliceSupportItem := ngapType.SliceSupportItem{}
				sliceSupportItem.SNSSAI.SST.Value, err = hex.DecodeString(sliceSupportItemLocal.SNSSAI.SST)
				if err != nil {
					logger.NgapLog.Errorf("DecodeString error: %+v", err)
				}

				if sliceSupportItemLocal.SNSSAI.SD != "" {
					sliceSupportItem.SNSSAI.SD = new(ngapType.SD)
					sliceSupportItem.SNSSAI.SD.Value, err = hex.DecodeString(sliceSupportItemLocal.SNSSAI.SD)
					if err != nil {
						logger.NgapLog.Errorf("DecodeString error: %+v", err)
					}
				}

				sliceSupportList.List = append(sliceSupportList.List, sliceSupportItem)
			}

			broadcastPLMNList.List = append(broadcastPLMNList.List, broadcastPLMNItem)
		}

		supportedTAList.List = append(supportedTAList.List, supportedTAItem)
	}

	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	/*
		* The reason PagingDRX ie was commented is that in TS23.501
		* PagingDRX was mentioned to be used only for 3GPP access.
		* However, the question that if the paging function for OCF
		* is needed requires verification.

			// PagingDRX
			ie = ngapType.NGSetupRequestIEs{}
			ie.Id.Value = ngapType.ProtocolIEIDDefaultPagingDRX
			ie.Criticality.Value = ngapType.CriticalityPresentIgnore
			ie.Value.Present = ngapType.NGSetupRequestIEsPresentDefaultPagingDRX
			ie.Value.DefaultPagingDRX = new(ngapType.PagingDRX)

			pagingDRX := ie.Value.DefaultPagingDRX
			pagingDRX.Value = ngapType.PagingDRXPresentV128
			nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)
	*/

	return ngap.Encoder(pdu)
}

func BuildNGReset(
	ngCause ngapType.Cause,
	partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeNGReset
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentNGReset
	initiatingMessage.Value.NGReset = new(ngapType.NGReset)

	nGReset := initiatingMessage.Value.NGReset
	nGResetIEs := &nGReset.ProtocolIEs
	// Cause
	{
		ie := ngapType.NGResetIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCause
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.NGResetIEsPresentCause
		ie.Value.Cause = new(ngapType.Cause)

		cause := ie.Value.Cause
		*cause = ngCause

		nGResetIEs.List = append(nGResetIEs.List, ie)
	}
	// ResetType
	{
		ie := ngapType.NGResetIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDResetType
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.NGResetIEsPresentResetType
		ie.Value.ResetType = new(ngapType.ResetType)

		resetType := ie.Value.ResetType
		if partOfNGInterface == nil {
			resetType.Present = ngapType.ResetTypePresentNGInterface
			resetType.NGInterface = new(ngapType.ResetAll)
			resetType.NGInterface.Value = ngapType.ResetAllPresentResetAll
		} else {
			resetType.Present = ngapType.ResetTypePresentPartOfNGInterface
			resetType.PartOfNGInterface = new(ngapType.UEAssociatedLogicalNGConnectionList)
			resetType.PartOfNGInterface = partOfNGInterface
		}

		nGResetIEs.List = append(nGResetIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildNGResetAcknowledge(
	partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList,
	diagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeNGReset
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentNGResetAcknowledge
	successfulOutcome.Value.NGResetAcknowledge = new(ngapType.NGResetAcknowledge)

	nGResetAcknowledge := successfulOutcome.Value.NGResetAcknowledge
	nGResetAcknowledgeIEs := &nGResetAcknowledge.ProtocolIEs
	// UEAssociatedLogicalNGConnectionList
	if partOfNGInterface != nil {
		ie := ngapType.NGResetAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUEAssociatedLogicalNGConnectionList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.NGResetAcknowledgeIEsPresentUEAssociatedLogicalNGConnectionList
		ie.Value.UEAssociatedLogicalNGConnectionList = new(ngapType.UEAssociatedLogicalNGConnectionList)

		uEAssociatedLogicalNGConnectionList := ie.Value.UEAssociatedLogicalNGConnectionList
		*uEAssociatedLogicalNGConnectionList = *partOfNGInterface

		nGResetAcknowledgeIEs.List = append(nGResetAcknowledgeIEs.List, ie)
	}
	// CriticalityDiagnostics
	if diagnostics != nil {
		ie := ngapType.NGResetAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.NGResetAcknowledgeIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		criticalityDiagnostics := ie.Value.CriticalityDiagnostics
		*criticalityDiagnostics = *diagnostics

		nGResetAcknowledgeIEs.List = append(nGResetAcknowledgeIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildInitialContextSetupResponse(
	ue *context.OCFUe,
	responseList *ngapType.PDUSessionResourceSetupListCxtRes,
	failedList *ngapType.PDUSessionResourceFailedToSetupListCxtRes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentInitialContextSetupResponse
	successfulOutcome.Value.InitialContextSetupResponse = new(ngapType.InitialContextSetupResponse)

	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	initialContextSetupResponseIEs := &initialContextSetupResponse.ProtocolIEs

	// OCF UE NGAP ID
	ie := ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentAmfUENGAPID
	ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

	AmfUENGAPID := ie.Value.AmfUENGAPID
	AmfUENGAPID.Value = ue.AmfUENGAPID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// PDU Session Resource Setup Response List (optional)
	if responseList != nil && len(responseList.List) > 0 {
		ie = ngapType.InitialContextSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtRes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentPDUSessionResourceSetupListCxtRes
		ie.Value.PDUSessionResourceSetupListCxtRes = responseList
		initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)
	}

	// PDU Session Resource Failed to Setup List (optional)
	if failedList != nil && len(failedList.List) > 0 {
		ie = ngapType.InitialContextSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtRes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentPDUSessionResourceFailedToSetupListCxtRes
		ie.Value.PDUSessionResourceFailedToSetupListCxtRes = failedList
		initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)
	}

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.InitialContextSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildInitialContextSetupFailure(
	ue *context.OCFUe,
	cause ngapType.Cause,
	failedList *ngapType.PDUSessionResourceFailedToSetupListCxtFail,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentInitialContextSetupFailure
	unsuccessfulOutcome.Value.InitialContextSetupFailure = new(ngapType.InitialContextSetupFailure)

	initialContextSetupFailure := unsuccessfulOutcome.Value.InitialContextSetupFailure
	initialContextSetupFailureIEs := &initialContextSetupFailure.ProtocolIEs

	// OCF UE NGAP ID
	ie := ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentAmfUENGAPID
	ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

	AmfUENGAPID := ie.Value.AmfUENGAPID
	AmfUENGAPID.Value = ue.AmfUENGAPID

	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// PDU Session Resource Failed to Setup List
	if failedList != nil && len(failedList.List) > 0 {
		ie = ngapType.InitialContextSetupFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtFail
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentPDUSessionResourceFailedToSetupListCxtFail
		ie.Value.PDUSessionResourceFailedToSetupListCxtFail = failedList
		initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)
	}

	// Cause
	ie = ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentCause
	ie.Value.Cause = &cause
	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.InitialContextSetupFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildUEContextModificationResponse(
	ue *context.OCFUe, criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUEContextModification
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentUEContextModificationResponse
	successfulOutcome.Value.UEContextModificationResponse = new(ngapType.UEContextModificationResponse)

	uEContextModificationResponse := successfulOutcome.Value.UEContextModificationResponse
	uEContextModificationResponseIEs := &uEContextModificationResponse.ProtocolIEs

	// OCF UE NGAP ID
	ie := ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationResponseIEsPresentAmfUENGAPID
	ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

	AmfUENGAPID := ie.Value.AmfUENGAPID
	AmfUENGAPID.Value = ue.AmfUENGAPID

	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	// Criticality Diagnostics (optional)
	ie = ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.CriticalityDiagnostics = criticalityDiagnostics
	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildUEContextModificationFailure(ue *context.OCFUe, cause ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUEContextModification
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentUEContextModificationFailure
	unsuccessfulOutcome.Value.UEContextModificationFailure = new(ngapType.UEContextModificationFailure)

	uEContextModificationFailure := unsuccessfulOutcome.Value.UEContextModificationFailure
	uEContextModificationFailureIEs := &uEContextModificationFailure.ProtocolIEs

	// OCF UE NGAP ID
	ie := ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationFailureIEsPresentAmfUENGAPID
	ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

	AmfUENGAPID := ie.Value.AmfUENGAPID
	AmfUENGAPID.Value = ue.AmfUENGAPID

	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationFailureIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	// Cause
	ie = ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationFailureIEsPresentCause
	ie.Value.Cause = &cause
	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	// Criticality Diagnostics (optional)
	ie = ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.CriticalityDiagnostics = criticalityDiagnostics
	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildUEContextReleaseComplete(ue *context.OCFUe,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUEContextRelease
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentUEContextReleaseComplete
	successfulOutcome.Value.UEContextReleaseComplete = new(ngapType.UEContextReleaseComplete)

	uEContextReleaseComplete := successfulOutcome.Value.UEContextReleaseComplete
	uEContextReleaseCompleteIEs := &uEContextReleaseComplete.ProtocolIEs

	// OCF UE NGAP ID
	ie := ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentAmfUENGAPID
	ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

	AmfUENGAPID := ie.Value.AmfUENGAPID
	AmfUENGAPID.Value = ue.AmfUENGAPID

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// User Location Information (optional)
	ie = ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationOCF
	userLocationInformation.UserLocationInformationOCF = new(ngapType.UserLocationInformationOCF)

	userLocationInfoOCF := userLocationInformation.UserLocationInformationOCF
	userLocationInfoOCF.IPAddress = ngapConvert.IPAddressToNgap(ue.IPAddrv4, ue.IPAddrv6)
	userLocationInfoOCF.PortNumber = ngapConvert.PortNumberToNgap(ue.PortNumber)

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// PDU Session Resource List
	ie = ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceListCxtRelCpl
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentPDUSessionResourceListCxtRelCpl
	ie.Value.PDUSessionResourceListCxtRelCpl = new(ngapType.PDUSessionResourceListCxtRelCpl)

	pDUSessionResourceListCxtRelCpl := ie.Value.PDUSessionResourceListCxtRelCpl

	// PDU Session Resource Item (in PDU Session Resource List)
	for _, pduSession := range ue.PduSessionList {
		pDUSessionResourceItemCxtRelCpl := ngapType.PDUSessionResourceItemCxtRelCpl{}
		pDUSessionResourceItemCxtRelCpl.PDUSessionID.Value = pduSession.Id
		pDUSessionResourceListCxtRelCpl.List =
			append(pDUSessionResourceListCxtRelCpl.List, pDUSessionResourceItemCxtRelCpl)
	}

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.UEContextReleaseCompleteIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildUEContextReleaseRequest(ue *context.OCFUe, cause ngapType.Cause) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUEContextReleaseRequest
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUEContextReleaseRequest
	initiatingMessage.Value.UEContextReleaseRequest = new(ngapType.UEContextReleaseRequest)

	uEContextReleaseRequest := initiatingMessage.Value.UEContextReleaseRequest
	uEContextReleaseRequestIEs := &uEContextReleaseRequest.ProtocolIEs

	// OCF UE NGAP ID
	ie := ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentAmfUENGAPID
	ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

	AmfUENGAPID := ie.Value.AmfUENGAPID
	AmfUENGAPID.Value = ue.AmfUENGAPID

	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	// PDU Session Resource List
	ie = ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceListCxtRelReq
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentPDUSessionResourceListCxtRelReq
	ie.Value.PDUSessionResourceListCxtRelReq = new(ngapType.PDUSessionResourceListCxtRelReq)

	pDUSessionResourceListCxtRelReq := ie.Value.PDUSessionResourceListCxtRelReq

	// PDU Session Resource Item in PDU session Resource List
	for _, pduSession := range ue.PduSessionList {
		pDUSessionResourceItem := ngapType.PDUSessionResourceItemCxtRelReq{}
		pDUSessionResourceItem.PDUSessionID.Value = pduSession.Id
		pDUSessionResourceListCxtRelReq.List =
			append(pDUSessionResourceListCxtRelReq.List, pDUSessionResourceItem)
	}
	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	// Cause
	ie = ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentCause
	ie.Value.Cause = &cause
	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildInitialUEMessage(ue *context.OCFUe, nasPdu []byte,
	allowedNSSAI *ngapType.AllowedNSSAI) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeInitialUEMessage
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentInitialUEMessage
	initiatingMessage.Value.InitialUEMessage = new(ngapType.InitialUEMessage)

	initialUEMessage := initiatingMessage.Value.InitialUEMessage
	initialUEMessageIEs := &initialUEMessage.ProtocolIEs
	// RANUENGAPID
	{
		ie := ngapType.InitialUEMessageIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialUEMessageIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

		rANUENGAPID := ie.Value.RANUENGAPID
		rANUENGAPID.Value = ue.RanUeNgapId

		initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)
	}
	// NASPDU
	{
		ie := ngapType.InitialUEMessageIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNASPDU
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialUEMessageIEsPresentNASPDU
		ie.Value.NASPDU = new(ngapType.NASPDU)

		nASPDU := ie.Value.NASPDU
		nASPDU.Value = nasPdu

		initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)
	}
	// UserLocationInformation
	{
		ie := ngapType.InitialUEMessageIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialUEMessageIEsPresentUserLocationInformation
		ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

		userLocationInformation := ie.Value.UserLocationInformation
		userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationOCF
		userLocationInformation.UserLocationInformationOCF = new(ngapType.UserLocationInformationOCF)

		userLocationInfoOCF := userLocationInformation.UserLocationInformationOCF
		userLocationInfoOCF.IPAddress = ngapConvert.IPAddressToNgap(ue.IPAddrv4, ue.IPAddrv6)
		userLocationInfoOCF.PortNumber = ngapConvert.PortNumberToNgap(ue.PortNumber)

		initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)
	}
	// RRCEstablishmentCause
	{
		ie := ngapType.InitialUEMessageIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRRCEstablishmentCause
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialUEMessageIEsPresentRRCEstablishmentCause
		ie.Value.RRCEstablishmentCause = new(ngapType.RRCEstablishmentCause)

		rRCEstablishmentCause := ie.Value.RRCEstablishmentCause
		rRCEstablishmentCause.Value = aper.Enumerated(ue.RRCEstablishmentCause)
		initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)
	}
	// FiveGSTMSI
	if len(ue.Guti) != 0 {
		ie := ngapType.InitialUEMessageIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDFiveGSTMSI
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialUEMessageIEsPresentFiveGSTMSI
		ie.Value.FiveGSTMSI = new(ngapType.FiveGSTMSI)

		fiveGSTMSI := ie.Value.FiveGSTMSI
		var amfID string
		var tmsi string
		if len(ue.Guti) == 19 {
			amfID = ue.Guti[5:11]
			tmsi = ue.Guti[11:]
		} else {
			amfID = ue.Guti[6:12]
			tmsi = ue.Guti[12:]
		}
		_, amfSetID, amfPointer := ngapConvert.AmfIdToNgap(amfID)

		fiveGSTMSI.AMFSetID.Value = amfSetID
		fiveGSTMSI.AMFPointer.Value = amfPointer
		var err error
		fiveGSTMSI.FiveGTMSI.Value, err = hex.DecodeString(tmsi)
		if err != nil {
			logger.NgapLog.Errorf("DecodeString error: %+v", err)
		}
		initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)
	}
	// OCFSetID
	if len(ue.Guti) != 0 {
		ie := ngapType.InitialUEMessageIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFSetID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialUEMessageIEsPresentAMFSetID
		ie.Value.AMFSetID = new(ngapType.AMFSetID)

		aMFSetID := ie.Value.AMFSetID
		// <MCC><MNC><OCF Region ID><OCF Set ID><OCF Pointer><5G-TMSI>
		// <MCC><MNC> is 3 bytes, <OCF Region ID><OCF Set ID><OCF Pointer> is 3 bytes
		// 1 byte is 2 characters
		var amfID string
		if len(ue.Guti) == 19 { // MNC is 2 char
			amfID = ue.Guti[5:11]
		} else {
			amfID = ue.Guti[6:12]
		}
		_, aMFSetID.Value, _ = ngapConvert.AmfIdToNgap(amfID)

		initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)
	}
	// UEContextRequest
	ie := ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUEContextRequest
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentUEContextRequest
	ie.Value.UEContextRequest = new(ngapType.UEContextRequest)

	ie.Value.UEContextRequest.Value = ngapType.UEContextRequestPresentRequested

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// AllowedNSSAI
	if allowedNSSAI != nil {
		ie := ngapType.InitialUEMessageIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAllowedNSSAI
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialUEMessageIEsPresentAllowedNSSAI
		ie.Value.AllowedNSSAI = new(ngapType.AllowedNSSAI)

		ie.Value.AllowedNSSAI = allowedNSSAI

		initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildUplinkNASTransport(ue *context.OCFUe, nasPdu []byte) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUplinkNASTransport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUplinkNASTransport
	initiatingMessage.Value.UplinkNASTransport = new(ngapType.UplinkNASTransport)

	uplinkNasTransport := initiatingMessage.Value.UplinkNASTransport
	uplinkNasTransportIEs := &uplinkNasTransport.ProtocolIEs

	// OCF UE NGAP ID
	ie := ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentAmfUENGAPID
	ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

	AmfUENGAPID := ie.Value.AmfUENGAPID
	AmfUENGAPID.Value = ue.AmfUENGAPID

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	// NAS-PDU
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNASPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentNASPDU
	ie.Value.NASPDU = new(ngapType.NASPDU)
	nASPDU := ie.Value.NASPDU
	nASPDU.Value = nasPdu
	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	// User Location Information
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationOCF
	userLocationInformation.UserLocationInformationOCF = new(ngapType.UserLocationInformationOCF)
	userLocationInformationOCF := userLocationInformation.UserLocationInformationOCF
	userLocationInformationOCF.IPAddress = ngapConvert.IPAddressToNgap(ue.IPAddrv4, ue.IPAddrv6)
	userLocationInformationOCF.PortNumber = ngapConvert.PortNumberToNgap(ue.PortNumber)

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildNASNonDeliveryIndication(ue *context.OCFUe, nasPdu []byte, cause ngapType.Cause) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeNASNonDeliveryIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentNASNonDeliveryIndication
	initiatingMessage.Value.NASNonDeliveryIndication = new(ngapType.NASNonDeliveryIndication)

	nASNonDeliveryIndication := initiatingMessage.Value.NASNonDeliveryIndication
	nASNonDeliveryIndicationIEs := &nASNonDeliveryIndication.ProtocolIEs
	// AmfUENGAPID
	{
		ie := ngapType.NASNonDeliveryIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.NASNonDeliveryIndicationIEsPresentAmfUENGAPID
		ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

		AmfUENGAPID := ie.Value.AmfUENGAPID
		AmfUENGAPID.Value = ue.AmfUENGAPID

		nASNonDeliveryIndicationIEs.List = append(nASNonDeliveryIndicationIEs.List, ie)
	}
	// RANUENGAPID
	{
		ie := ngapType.NASNonDeliveryIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.NASNonDeliveryIndicationIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

		rANUENGAPID := ie.Value.RANUENGAPID
		rANUENGAPID.Value = ue.RanUeNgapId

		nASNonDeliveryIndicationIEs.List = append(nASNonDeliveryIndicationIEs.List, ie)
	}
	// NASPDU
	{
		ie := ngapType.NASNonDeliveryIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDNASPDU
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.NASNonDeliveryIndicationIEsPresentNASPDU
		ie.Value.NASPDU = new(ngapType.NASPDU)

		nASPDU := ie.Value.NASPDU
		nASPDU.Value = nasPdu
		nASNonDeliveryIndicationIEs.List = append(nASNonDeliveryIndicationIEs.List, ie)
	}
	// Cause
	{
		ie := ngapType.NASNonDeliveryIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCause
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.NASNonDeliveryIndicationIEsPresentCause
		ie.Value.Cause = new(ngapType.Cause)

		ie.Value.Cause = &cause

		nASNonDeliveryIndicationIEs.List = append(nASNonDeliveryIndicationIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildRerouteNASRequest() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceSetupResponse(
	ue *context.OCFUe,
	responseList *ngapType.PDUSessionResourceSetupListSURes,
	failedList *ngapType.PDUSessionResourceFailedToSetupListSURes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceSetupResponse
	successfulOutcome.Value.PDUSessionResourceSetupResponse = new(ngapType.PDUSessionResourceSetupResponse)

	pduSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	pduSessionResourceSetupResponseIEs := &pduSessionResourceSetupResponse.ProtocolIEs

	// OCF UE NGAP ID
	ie := ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentAmfUENGAPID
	ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

	AmfUENGAPID := ie.Value.AmfUENGAPID
	AmfUENGAPID.Value = ue.AmfUENGAPID

	pduSessionResourceSetupResponseIEs.List = append(pduSessionResourceSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	pduSessionResourceSetupResponseIEs.List = append(pduSessionResourceSetupResponseIEs.List, ie)

	// PDU Session Resource Setup Response List (optional)
	if responseList != nil && len(responseList.List) > 0 {
		ie = ngapType.PDUSessionResourceSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceSetupListSURes
		ie.Value.PDUSessionResourceSetupListSURes = responseList
		pduSessionResourceSetupResponseIEs.List = append(pduSessionResourceSetupResponseIEs.List, ie)
	}

	// PDU Session Resource Failed to Setup List (optional)
	if failedList != nil && len(failedList.List) > 0 {
		ie = ngapType.PDUSessionResourceSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListSURes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceFailedToSetupListSURes
		ie.Value.PDUSessionResourceFailedToSetupListSURes = failedList
		pduSessionResourceSetupResponseIEs.List = append(pduSessionResourceSetupResponseIEs.List, ie)
	}

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.PDUSessionResourceSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		pduSessionResourceSetupResponseIEs.List = append(pduSessionResourceSetupResponseIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceModifyResponse(
	ue *context.OCFUe,
	responseList *ngapType.PDUSessionResourceModifyListModRes,
	failedList *ngapType.PDUSessionResourceFailedToModifyListModRes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceModify
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceModifyResponse
	successfulOutcome.Value.PDUSessionResourceModifyResponse = new(ngapType.PDUSessionResourceModifyResponse)

	pduSessionResourceModifyResponse := successfulOutcome.Value.PDUSessionResourceModifyResponse
	pduSessionResourceModifyResponseIEs := &pduSessionResourceModifyResponse.ProtocolIEs

	// OCF UE NGAP ID
	ie := ngapType.PDUSessionResourceModifyResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentAmfUENGAPID
	ie.Value.AmfUENGAPID = &ngapType.AmfUENGAPID{
		Value: ue.AmfUENGAPID,
	}
	pduSessionResourceModifyResponseIEs.List = append(pduSessionResourceModifyResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceModifyResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = &ngapType.RANUENGAPID{
		Value: ue.RanUeNgapId,
	}
	pduSessionResourceModifyResponseIEs.List = append(pduSessionResourceModifyResponseIEs.List, ie)

	// PDU Session Resource Modify Response List (optional)
	if responseList != nil && len(responseList.List) > 0 {
		ie = ngapType.PDUSessionResourceModifyResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceModifyListModRes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentPDUSessionResourceModifyListModRes
		ie.Value.PDUSessionResourceModifyListModRes = responseList
		pduSessionResourceModifyResponseIEs.List = append(pduSessionResourceModifyResponseIEs.List, ie)
	}

	// PDU Session Resource Failed to Modify List (optional)
	if failedList != nil && len(failedList.List) > 0 {
		ie = ngapType.PDUSessionResourceModifyResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToModifyListModRes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentPDUSessionResourceFailedToModifyListModRes
		ie.Value.PDUSessionResourceFailedToModifyListModRes = failedList
		pduSessionResourceModifyResponseIEs.List = append(pduSessionResourceModifyResponseIEs.List, ie)
	}

	// User Location Information
	ie = ngapType.PDUSessionResourceModifyResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceModifyResponseIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationOCF
	userLocationInformation.UserLocationInformationOCF = new(ngapType.UserLocationInformationOCF)

	userLocationInformationOCF := userLocationInformation.UserLocationInformationOCF
	userLocationInformationOCF.IPAddress = ngapConvert.IPAddressToNgap(ue.IPAddrv4, ue.IPAddrv6)
	userLocationInformationOCF.PortNumber = ngapConvert.PortNumberToNgap(ue.PortNumber)

	pduSessionResourceModifyResponseIEs.List = append(pduSessionResourceModifyResponseIEs.List, ie)

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.PDUSessionResourceModifyResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		pduSessionResourceModifyResponseIEs.List = append(pduSessionResourceModifyResponseIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceModifyIndication(
	ue *context.OCFUe,
	modifyList []ngapType.PDUSessionResourceModifyItemModInd) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceModifyIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPDUSessionResourceModifyIndication
	initiatingMessage.Value.PDUSessionResourceModifyIndication = new(ngapType.PDUSessionResourceModifyIndication)

	pDUSessionResourceModifyIndication := initiatingMessage.Value.PDUSessionResourceModifyIndication
	pDUSessionResourceModifyIndicationIEs := &pDUSessionResourceModifyIndication.ProtocolIEs
	// AmfUENGAPID
	{
		ie := ngapType.PDUSessionResourceModifyIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.PDUSessionResourceModifyIndicationIEsPresentAmfUENGAPID
		ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

		AmfUENGAPID := ie.Value.AmfUENGAPID
		AmfUENGAPID.Value = ue.AmfUENGAPID

		pDUSessionResourceModifyIndicationIEs.List = append(pDUSessionResourceModifyIndicationIEs.List, ie)
	}
	// RANUENGAPID
	{
		ie := ngapType.PDUSessionResourceModifyIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.PDUSessionResourceModifyIndicationIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

		rANUENGAPID := ie.Value.RANUENGAPID
		rANUENGAPID.Value = ue.RanUeNgapId

		pDUSessionResourceModifyIndicationIEs.List = append(pDUSessionResourceModifyIndicationIEs.List, ie)
	}
	// PDUSessionResourceModifyListModInd
	{
		ie := ngapType.PDUSessionResourceModifyIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceModifyListModInd
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.PDUSessionResourceModifyIndicationIEsPresentPDUSessionResourceModifyListModInd
		ie.Value.PDUSessionResourceModifyListModInd = new(ngapType.PDUSessionResourceModifyListModInd)

		pDUSessionResourceModifyListModInd := ie.Value.PDUSessionResourceModifyListModInd
		pDUSessionResourceModifyListModInd.List = modifyList

		pDUSessionResourceModifyIndicationIEs.List = append(pDUSessionResourceModifyIndicationIEs.List, ie)
	}

	return ngap.Encoder(pdu)

}

func BuildPDUSessionResourceNotify(
	ue *context.OCFUe,
	notiList *ngapType.PDUSessionResourceNotifyList,
	relList *ngapType.PDUSessionResourceReleasedListNot) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceNotify
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPDUSessionResourceNotify
	initiatingMessage.Value.PDUSessionResourceNotify = new(ngapType.PDUSessionResourceNotify)

	pDUSessionResourceNotify := initiatingMessage.Value.PDUSessionResourceNotify
	pDUSessionResourceNotifyIEs := &pDUSessionResourceNotify.ProtocolIEs
	// AmfUENGAPID
	{
		ie := ngapType.PDUSessionResourceNotifyIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentAmfUENGAPID
		ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

		AmfUENGAPID := ie.Value.AmfUENGAPID
		AmfUENGAPID.Value = ue.AmfUENGAPID

		pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)
	}
	// RANUENGAPID
	{
		ie := ngapType.PDUSessionResourceNotifyIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

		rANUENGAPID := ie.Value.RANUENGAPID
		rANUENGAPID.Value = ue.RanUeNgapId

		pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)
	}
	// PDUSessionResourceNotifyList
	if notiList != nil {
		ie := ngapType.PDUSessionResourceNotifyIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceNotifyList
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentPDUSessionResourceNotifyList
		ie.Value.PDUSessionResourceNotifyList = new(ngapType.PDUSessionResourceNotifyList)

		pDUSessionResourceNotifyList := ie.Value.PDUSessionResourceNotifyList
		*pDUSessionResourceNotifyList = *notiList

		pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)
	}
	// PDUSessionResourceReleasedListNot
	if relList != nil {
		ie := ngapType.PDUSessionResourceNotifyIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceReleasedListNot
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentPDUSessionResourceReleasedListNot
		ie.Value.PDUSessionResourceReleasedListNot = new(ngapType.PDUSessionResourceReleasedListNot)

		pDUSessionResourceReleasedListNot := ie.Value.PDUSessionResourceReleasedListNot
		*pDUSessionResourceReleasedListNot = *relList

		pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)
	}
	// UserLocationInformation
	if (ue.IPAddrv4 != "" || ue.IPAddrv6 != "") && ue.PortNumber != 0 {
		ie := ngapType.PDUSessionResourceNotifyIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceNotifyIEsPresentUserLocationInformation
		ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

		userLocationInformation := ie.Value.UserLocationInformation
		*userLocationInformation = ngapType.UserLocationInformation{
			Present: ngapType.UserLocationInformationPresentUserLocationInformationOCF,
			UserLocationInformationOCF: &ngapType.UserLocationInformationOCF{
				IPAddress:  ngapConvert.IPAddressToNgap(ue.IPAddrv4, ue.IPAddrv6),
				PortNumber: ngapConvert.PortNumberToNgap(ue.PortNumber),
			},
		}

		pDUSessionResourceNotifyIEs.List = append(pDUSessionResourceNotifyIEs.List, ie)
	}

	return ngap.Encoder(pdu)

}

func BuildPDUSessionResourceReleaseResponse(
	ue *context.OCFUe,
	relList ngapType.PDUSessionResourceReleasedListRelRes,
	diagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceRelease
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceReleaseResponse
	successfulOutcome.Value.PDUSessionResourceReleaseResponse = new(ngapType.PDUSessionResourceReleaseResponse)

	pDUSessionResourceReleaseResponse := successfulOutcome.Value.PDUSessionResourceReleaseResponse
	pDUSessionResourceReleaseResponseIEs := &pDUSessionResourceReleaseResponse.ProtocolIEs
	// AmfUENGAPID
	{
		ie := ngapType.PDUSessionResourceReleaseResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentAmfUENGAPID
		ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

		AmfUENGAPID := ie.Value.AmfUENGAPID
		AmfUENGAPID.Value = ue.AmfUENGAPID

		pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)
	}
	// RANUENGAPID
	{
		ie := ngapType.PDUSessionResourceReleaseResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

		rANUENGAPID := ie.Value.RANUENGAPID
		rANUENGAPID.Value = ue.RanUeNgapId

		pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)
	}
	// PDUSessionResourceReleasedListRelRes
	{
		ie := ngapType.PDUSessionResourceReleaseResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceReleasedListRelRes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentPDUSessionResourceReleasedListRelRes
		ie.Value.PDUSessionResourceReleasedListRelRes = new(ngapType.PDUSessionResourceReleasedListRelRes)

		pDUSessionResourceReleasedListRelRes := ie.Value.PDUSessionResourceReleasedListRelRes
		*pDUSessionResourceReleasedListRelRes = relList

		pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)
	}
	// UserLocationInformation
	if (ue.IPAddrv4 != "" || ue.IPAddrv6 != "") && ue.PortNumber != 0 {
		ie := ngapType.PDUSessionResourceReleaseResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentUserLocationInformation
		ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

		userLocationInformation := ie.Value.UserLocationInformation
		*userLocationInformation = ngapType.UserLocationInformation{
			Present: ngapType.UserLocationInformationPresentUserLocationInformationOCF,
			UserLocationInformationOCF: &ngapType.UserLocationInformationOCF{
				IPAddress:  ngapConvert.IPAddressToNgap(ue.IPAddrv4, ue.IPAddrv6),
				PortNumber: ngapConvert.PortNumberToNgap(ue.PortNumber),
			},
		}

		pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)
	}
	// CriticalityDiagnostics
	if diagnostics != nil {
		ie := ngapType.PDUSessionResourceReleaseResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		criticalityDiagnostics := ie.Value.CriticalityDiagnostics
		*criticalityDiagnostics = *diagnostics

		pDUSessionResourceReleaseResponseIEs.List = append(pDUSessionResourceReleaseResponseIEs.List, ie)
	}

	return ngap.Encoder(pdu)

}

func BuildErrorIndication(
	AmfUENGAPID *int64,
	ranUENGAPID *int64,
	cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeErrorIndication
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentErrorIndication
	initiatingMessage.Value.ErrorIndication = new(ngapType.ErrorIndication)

	errorIndication := initiatingMessage.Value.ErrorIndication
	errorIndicationIEs := &errorIndication.ProtocolIEs

	if AmfUENGAPID != nil && ranUENGAPID != nil {

		// OCF UE NGAP ID
		ie := ngapType.ErrorIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.ErrorIndicationIEsPresentAmfUENGAPID
		ie.Value.AmfUENGAPID = &ngapType.AmfUENGAPID{Value: *AmfUENGAPID}
		errorIndicationIEs.List = append(errorIndicationIEs.List, ie)

		// RAN UE NGAP ID
		ie = ngapType.ErrorIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.ErrorIndicationIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = &ngapType.RANUENGAPID{Value: *ranUENGAPID}
		errorIndicationIEs.List = append(errorIndicationIEs.List, ie)

	}

	// Cause
	if cause != nil {
		ie := ngapType.ErrorIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCause
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.ErrorIndicationIEsPresentCause
		ie.Value.Cause = cause
		errorIndicationIEs.List = append(errorIndicationIEs.List, ie)
	}

	// Criticality Diagnostics
	if criticalityDiagnostics != nil {
		ie := ngapType.ErrorIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.ErrorIndicationIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		errorIndicationIEs.List = append(errorIndicationIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildUERadioCapabilityInfoIndication() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildUERadioCapabilityCheckResponse(
	ue *context.OCFUe,
	diagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUERadioCapabilityCheck
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentUERadioCapabilityCheckResponse
	successfulOutcome.Value.UERadioCapabilityCheckResponse = new(ngapType.UERadioCapabilityCheckResponse)

	uERadioCapabilityCheckResponse := successfulOutcome.Value.UERadioCapabilityCheckResponse
	uERadioCapabilityCheckResponseIEs := &uERadioCapabilityCheckResponse.ProtocolIEs
	// AmfUENGAPID
	{
		ie := ngapType.UERadioCapabilityCheckResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAmfUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.UERadioCapabilityCheckResponseIEsPresentAmfUENGAPID
		ie.Value.AmfUENGAPID = new(ngapType.AmfUENGAPID)

		AmfUENGAPID := ie.Value.AmfUENGAPID
		AmfUENGAPID.Value = ue.AmfUENGAPID
		uERadioCapabilityCheckResponseIEs.List = append(uERadioCapabilityCheckResponseIEs.List, ie)
	}
	// RANUENGAPID
	{
		ie := ngapType.UERadioCapabilityCheckResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.UERadioCapabilityCheckResponseIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

		rANUENGAPID := ie.Value.RANUENGAPID
		rANUENGAPID.Value = ue.RanUeNgapId
		uERadioCapabilityCheckResponseIEs.List = append(uERadioCapabilityCheckResponseIEs.List, ie)
	}
	// IMSVoiceSupportIndicator
	{
		ie := ngapType.UERadioCapabilityCheckResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDIMSVoiceSupportIndicator
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.UERadioCapabilityCheckResponseIEsPresentIMSVoiceSupportIndicator
		ie.Value.IMSVoiceSupportIndicator = new(ngapType.IMSVoiceSupportIndicator)

		iMSVoiceSupportIndicator := ie.Value.IMSVoiceSupportIndicator
		iMSVoiceSupportIndicator.Value = aper.Enumerated(ue.IMSVoiceSupported)
		uERadioCapabilityCheckResponseIEs.List = append(uERadioCapabilityCheckResponseIEs.List, ie)
	}
	// CriticalityDiagnostics
	if diagnostics != nil {
		ie := ngapType.UERadioCapabilityCheckResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.UERadioCapabilityCheckResponseIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		criticalityDiagnostics := ie.Value.CriticalityDiagnostics
		*criticalityDiagnostics = *diagnostics

		uERadioCapabilityCheckResponseIEs.List = append(uERadioCapabilityCheckResponseIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildAMFConfigurationUpdateAcknowledge(
	setupList *ngapType.AMFTNLAssociationSetupList,
	failList *ngapType.TNLAssociationList,
	diagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeAMFConfigurationUpdate
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentAMFConfigurationUpdateAcknowledge
	successfulOutcome.Value.AMFConfigurationUpdateAcknowledge = new(ngapType.AMFConfigurationUpdateAcknowledge)

	aMFConfigurationUpdateAcknowledge := successfulOutcome.Value.AMFConfigurationUpdateAcknowledge
	aMFConfigurationUpdateAcknowledgeIEs := &aMFConfigurationUpdateAcknowledge.ProtocolIEs
	// AMFTNLAssociationSetupList
	if setupList != nil {
		ie := ngapType.AMFConfigurationUpdateAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationSetupList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateAcknowledgeIEsPresentAMFTNLAssociationSetupList
		ie.Value.AMFTNLAssociationSetupList = new(ngapType.AMFTNLAssociationSetupList)

		aMFTNLAssociationSetupList := ie.Value.AMFTNLAssociationSetupList
		*aMFTNLAssociationSetupList = *setupList

		aMFConfigurationUpdateAcknowledgeIEs.List = append(aMFConfigurationUpdateAcknowledgeIEs.List, ie)
	}
	// OCFTNLAssociationFailedToSetupList
	if failList != nil {
		ie := ngapType.AMFConfigurationUpdateAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFTNLAssociationFailedToSetupList
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateAcknowledgeIEsPresentAMFTNLAssociationFailedToSetupList
		ie.Value.AMFTNLAssociationFailedToSetupList = new(ngapType.TNLAssociationList)

		aMFTNLAssociationFailedToSetupList := ie.Value.AMFTNLAssociationFailedToSetupList
		*aMFTNLAssociationFailedToSetupList = *failList

		aMFConfigurationUpdateAcknowledgeIEs.List = append(aMFConfigurationUpdateAcknowledgeIEs.List, ie)
	}
	// CriticalityDiagnostics
	if diagnostics != nil {
		ie := ngapType.AMFConfigurationUpdateAcknowledgeIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateAcknowledgeIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		criticalityDiagnostics := ie.Value.CriticalityDiagnostics
		*criticalityDiagnostics = *diagnostics

		aMFConfigurationUpdateAcknowledgeIEs.List = append(aMFConfigurationUpdateAcknowledgeIEs.List, ie)
	}

	return ngap.Encoder(pdu)

}

func BuildAMFConfigurationUpdateFailure(
	ngCause ngapType.Cause,
	time *ngapType.TimeToWait,
	diagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeAMFConfigurationUpdate
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentAMFConfigurationUpdateFailure
	unsuccessfulOutcome.Value.AMFConfigurationUpdateFailure = new(ngapType.AMFConfigurationUpdateFailure)

	aMFConfigurationUpdateFailure := unsuccessfulOutcome.Value.AMFConfigurationUpdateFailure
	aMFConfigurationUpdateFailureIEs := &aMFConfigurationUpdateFailure.ProtocolIEs
	// Cause
	{
		ie := ngapType.AMFConfigurationUpdateFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCause
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateFailureIEsPresentCause
		ie.Value.Cause = new(ngapType.Cause)

		cause := ie.Value.Cause
		*cause = ngCause

		aMFConfigurationUpdateFailureIEs.List = append(aMFConfigurationUpdateFailureIEs.List, ie)
	}
	// TimeToWait
	if time != nil {
		ie := ngapType.AMFConfigurationUpdateFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDTimeToWait
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateFailureIEsPresentTimeToWait
		ie.Value.TimeToWait = new(ngapType.TimeToWait)

		timeToWait := ie.Value.TimeToWait
		*timeToWait = *time

		aMFConfigurationUpdateFailureIEs.List = append(aMFConfigurationUpdateFailureIEs.List, ie)
	}
	// CriticalityDiagnostics
	if diagnostics != nil {
		ie := ngapType.AMFConfigurationUpdateFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.AMFConfigurationUpdateFailureIEsPresentCriticalityDiagnostics
		ie.Value.CriticalityDiagnostics = new(ngapType.CriticalityDiagnostics)

		criticalityDiagnostics := ie.Value.CriticalityDiagnostics
		*criticalityDiagnostics = *diagnostics

		aMFConfigurationUpdateFailureIEs.List = append(aMFConfigurationUpdateFailureIEs.List, ie)
	}

	return ngap.Encoder(pdu)

}

func BuildRANConfigurationUpdate() ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeRANConfigurationUpdate
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentRANConfigurationUpdate
	initiatingMessage.Value.RANConfigurationUpdate = new(ngapType.RANConfigurationUpdate)

	rANConfigurationUpdate := initiatingMessage.Value.RANConfigurationUpdate
	rANConfigurationUpdateIEs := &rANConfigurationUpdate.ProtocolIEs

	ocfSelf := context.OCFSelf()

	// RANNodeName
	if ocfSelf.NFInfo.RanNodeName != "" {
		ie := ngapType.RANConfigurationUpdateIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANNodeName
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.RANConfigurationUpdateIEsPresentRANNodeName
		ie.Value.RANNodeName = new(ngapType.RANNodeName)

		rANNodeName := ie.Value.RANNodeName
		rANNodeName.Value = ocfSelf.NFInfo.RanNodeName

		rANConfigurationUpdateIEs.List = append(rANConfigurationUpdateIEs.List, ie)
	}
	// SupportedTAList
	if len(ocfSelf.NFInfo.SupportedTAList) > 0 {
		ie := ngapType.RANConfigurationUpdateIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDSupportedTAList
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.RANConfigurationUpdateIEsPresentSupportedTAList
		ie.Value.SupportedTAList = new(ngapType.SupportedTAList)

		supportedTAList := ie.Value.SupportedTAList

		for _, supportedTAItemLocal := range ocfSelf.NFInfo.SupportedTAList {
			// SupportedTAItem in SupportedTAList
			supportedTAItem := ngapType.SupportedTAItem{}
			var err error
			supportedTAItem.TAC.Value, err = hex.DecodeString(supportedTAItemLocal.TAC)
			if err != nil {
				logger.NgapLog.Errorf("DecodeString error: %+v", err)
			}

			broadcastPLMNList := &supportedTAItem.BroadcastPLMNList

			for _, broadcastPLMNListLocal := range supportedTAItemLocal.BroadcastPLMNList {
				// BroadcastPLMNItem in BroadcastPLMNList
				broadcastPLMNItem := ngapType.BroadcastPLMNItem{}
				broadcastPLMNItem.PLMNIdentity = util.PlmnIdToNgap(broadcastPLMNListLocal.PLMNID)

				sliceSupportList := &broadcastPLMNItem.TAISliceSupportList

				for _, sliceSupportItemLocal := range broadcastPLMNListLocal.TAISliceSupportList {
					// SliceSupportItem in SliceSupportList
					sliceSupportItem := ngapType.SliceSupportItem{}
					sliceSupportItem.SNSSAI.SST.Value, err = hex.DecodeString(sliceSupportItemLocal.SNSSAI.SST)
					if err != nil {
						logger.NgapLog.Errorf("DecodeString error: %+v", err)
					}

					if sliceSupportItemLocal.SNSSAI.SD != "" {
						sliceSupportItem.SNSSAI.SD = new(ngapType.SD)
						sliceSupportItem.SNSSAI.SD.Value, err = hex.DecodeString(sliceSupportItemLocal.SNSSAI.SD)
						if err != nil {
							logger.NgapLog.Errorf("DecodeString error: %+v", err)
						}
					}

					sliceSupportList.List = append(sliceSupportList.List, sliceSupportItem)
				}

				broadcastPLMNList.List = append(broadcastPLMNList.List, broadcastPLMNItem)
			}

			supportedTAList.List = append(supportedTAList.List, supportedTAItem)
		}

		rANConfigurationUpdateIEs.List = append(rANConfigurationUpdateIEs.List, ie)
	}
	// DefaultPagingDRX
	// {
	// 	ie := ngapType.RANConfigurationUpdateIEs{}
	// 	ie.Id.Value = ngapType.ProtocolIEIDDefaultPagingDRX
	// 	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	// 	ie.Value.Present = ngapType.RANConfigurationUpdateIEsPresentDefaultPagingDRX
	// 	ie.Value.DefaultPagingDRX = new(ngapType.PagingDRX)

	// 	defaultPagingDRX := ie.Value.DefaultPagingDRX

	// 	rANConfigurationUpdateIEs.List = append(rANConfigurationUpdateIEs.List, ie)
	// }

	return ngap.Encoder(pdu)

}

func BuildUplinkRANConfigurationTransfer() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildUplinkRANStatusTransfer() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildLocationReportingFailureIndication() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildLocationReport() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildRRCInactiveTransitionReport() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceSetupResponseTransfer(pduSession *context.PDUSession) ([]byte, error) {
	// OCF context
	ocfSelf := context.OCFSelf()

	transfer := ngapType.PDUSessionResourceSetupResponseTransfer{}

	// TODO: use tunnel info allocated by ocf
	// QOS Flow Per TNL Information
	qosFlowPerTNLInformation := &transfer.QosFlowPerTNLInformation

	// UP transport layer information - UE(RAN) side
	qosFlowPerTNLInformation.UPTransportLayerInformation.Present =
		ngapType.UPTransportLayerInformationPresentGTPTunnel
	qosFlowPerTNLInformation.UPTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)

	gtpTunnel := qosFlowPerTNLInformation.UPTransportLayerInformation.GTPTunnel
	teid := make([]byte, 4)
	binary.BigEndian.PutUint32(teid, pduSession.GTPConnection.IncomingTEID)
	gtpTunnel.GTPTEID.Value = teid
	gtpTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(ocfSelf.GTPBindAddress, "")

	// Associated Qos Flow List
	for _, qfi := range pduSession.QFIList {
		associatedQosFlowItem := ngapType.AssociatedQosFlowItem{
			QosFlowIdentifier: ngapType.QosFlowIdentifier{
				Value: int64(qfi),
			},
		}
		qosFlowPerTNLInformation.AssociatedQosFlowList.List =
			append(qosFlowPerTNLInformation.AssociatedQosFlowList.List, associatedQosFlowItem)
	}

	return aper.MarshalWithParams(transfer, "valueExt")
}

func BuildPDUSessionResourceSetupUnsuccessfulTransfer(
	cause ngapType.Cause, criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	transfer := ngapType.PDUSessionResourceSetupUnsuccessfulTransfer{}

	// Cause
	transfer.Cause = cause

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		transfer.CriticalityDiagnostics = criticalityDiagnostics
	}

	return aper.MarshalWithParams(transfer, "valueExt")
}

func BuildPDUSessionResourceModifyResponseTransfer(
	ulNGUUPTNLInformation *ngapType.UPTransportLayerInformation,
	dlNGUUPTNLInformation *ngapType.UPTransportLayerInformation,
	responseList *ngapType.QosFlowAddOrModifyResponseList,
	failedList *ngapType.QosFlowList) ([]byte, error) {

	transfer := ngapType.PDUSessionResourceModifyResponseTransfer{}

	// ulNGUUPTNLInformation store user plane tunnel information of
	// 5GC's endpoint (optional)
	if ulNGUUPTNLInformation != nil {
		transfer.ULNGUUPTNLInformation = ulNGUUPTNLInformation
	}

	// dlNGUUPTNLInformation store user plane tunnel information of
	// ran's endpoint (optional)
	if dlNGUUPTNLInformation != nil {
		transfer.DLNGUUPTNLInformation = dlNGUUPTNLInformation
	}

	if responseList != nil && len(responseList.List) != 0 {
		transfer.QosFlowAddOrModifyResponseList = responseList
	}

	// Additional Qos Flow per TNL Information (optional)

	// Qos Flow Failed to Add or Modify List (optional)
	if failedList != nil && len(failedList.List) != 0 {
		transfer.QosFlowFailedToAddOrModifyList = failedList
	}

	return aper.MarshalWithParams(transfer, "valueExt")
}

func BuildPDUSessionResourceModifyUnsuccessfulTransfer(cause ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	transfer := ngapType.PDUSessionResourceModifyUnsuccessfulTransfer{}

	// Cause
	transfer.Cause = cause

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		transfer.CriticalityDiagnostics = criticalityDiagnostics
	}

	return aper.MarshalWithParams(transfer, "valueExt")
}

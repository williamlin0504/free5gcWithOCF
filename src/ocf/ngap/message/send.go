package message

import (
	"free5gcWithOCF/lib/aper"
	"free5gcWithOCF/lib/ngap/ngapType"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/src/ocf/context"
	"free5gcWithOCF/src/ocf/logger"
	"free5gcWithOCF/src/ocf/producer/callback"
	"free5gcWithOCF/src/ocf/util"
	"time"

	"github.com/sirupsen/logrus"
)

var ngaplog *logrus.Entry

func init() {
	ngaplog = logger.NgapLog
}

func SendToRan(ran *context.OcfRan, packet []byte) {

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	if len(packet) == 0 {
		ngaplog.Error("packet len is 0")
		return
	}

	ngaplog.Debugf("[NGAP] Send To Ran [IP: %s]", ran.Conn.RemoteAddr().String())

	if n, err := ran.Conn.Write(packet); err != nil {
		ngaplog.Errorf("Send error: %+v", err)
		return
	} else {
		ngaplog.Debugf("Write %d bytes", n)
	}
}

func SendToRanUe(ue *context.RanUe, packet []byte) {

	var ran *context.OcfRan

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if ran = ue.Ran; ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	if ue.OcfUe == nil {
		ngaplog.Warn("OcfUe is nil")
	}

	SendToRan(ran, packet)
}

func NasSendToRan(ue *context.OcfUe, accessType models.AccessType, packet []byte) {

	if ue == nil {
		ngaplog.Error("OcfUe is nil")
		return
	}

	ranUe := ue.RanUe[accessType]
	if ranUe == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	SendToRanUe(ranUe, packet)
}

func SendNGSetupResponse(ran *context.OcfRan) {

	ngaplog.Info("[OCF] Send NG-Setup response")

	pkt, err := BuildNGSetupResponse()
	if err != nil {
		ngaplog.Errorf("Build NGSetupResponse failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendNGSetupFailure(ran *context.OcfRan, cause ngapType.Cause) {

	ngaplog.Info("[OCF] Send NG-Setup failure")

	if cause.Present == ngapType.CausePresentNothing {
		ngaplog.Errorf("Cause present is nil")
		return
	}

	pkt, err := BuildNGSetupFailure(cause)
	if err != nil {
		ngaplog.Errorf("Build NGSetupFailure failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

// partOfNGInterface: if reset type is "reset all", set it to nil TS 38.413 9.2.6.11
func SendNGReset(ran *context.OcfRan, cause ngapType.Cause,
	partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList) {

	ngaplog.Info("[OCF] Send NG Reset")

	pkt, err := BuildNGReset(cause, partOfNGInterface)
	if err != nil {
		ngaplog.Errorf("Build NGReset failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendNGResetAcknowledge(ran *context.OcfRan, partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send NG Reset Acknowledge")

	if partOfNGInterface != nil && len(partOfNGInterface.List) == 0 {
		ngaplog.Error("length of partOfNGInterface is 0")
		return
	}

	pkt, err := BuildNGResetAcknowledge(partOfNGInterface, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build NGResetAcknowledge failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendDownlinkNasTransport(ue *context.RanUe, nasPdu []byte,
	mobilityRestrictionList *ngapType.MobilityRestrictionList) {

	ngaplog.Info("[OCF] Send Downlink Nas Transport")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(nasPdu) == 0 {
		ngaplog.Errorf("[Send DownlinkNasTransport] Error: nasPdu is nil")
	}

	pkt, err := BuildDownlinkNasTransport(ue, nasPdu, mobilityRestrictionList)
	if err != nil {
		ngaplog.Errorf("Build DownlinkNasTransport failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendPDUSessionResourceReleaseCommand(ue *context.RanUe, nasPdu []byte,
	pduSessionResourceReleasedList ngapType.PDUSessionResourceToReleaseListRelCmd) {

	ngaplog.Info("[OCF] Send PDU Session Resource Release Command")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildPDUSessionResourceReleaseCommand(ue, nasPdu, pduSessionResourceReleasedList)
	if err != nil {
		ngaplog.Errorf("Build PDUSessionResourceReleaseCommand failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendUEContextReleaseCommand(ue *context.RanUe, action context.RelAction, causePresent int, cause aper.Enumerated) {

	ngaplog.Info("[OCF] Send UE Context Release Command")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildUEContextReleaseCommand(ue, causePresent, cause)
	if err != nil {
		ngaplog.Errorf("Build UEContextReleaseCommand failed : %s", err.Error())
		return
	}
	ue.ReleaseAction = action
	if ue.OcfUe != nil && ue.Ran != nil {
		ue.OcfUe.ReleaseCause[ue.Ran.AnType] = &context.CauseAll{
			NgapCause: &models.NgApCause{
				Group: int32(causePresent),
				Value: int32(cause),
			},
		}
	}
	SendToRanUe(ue, pkt)
}

func SendErrorIndication(ran *context.OcfRan, ocfUeNgapId, ranUeNgapId *int64, cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send Error Indication")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildErrorIndication(ocfUeNgapId, ranUeNgapId, cause, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build ErrorIndication failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendUERadioCapabilityCheckRequest(ue *context.RanUe) {

	ngaplog.Info("[OCF] Send UE Radio Capability Check Request")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildUERadioCapabilityCheckRequest(ue)
	if err != nil {
		ngaplog.Errorf("Build UERadioCapabilityCheckRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendHandoverCancelAcknowledge(ue *context.RanUe, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send Handover Cancel Acknowledge")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildHandoverCancelAcknowledge(ue, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build HandoverCancelAcknowledge failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// nasPDU: from nas layer
// pduSessionResourceSetupRequestList: provided by OCF, and transfer data is from SMF
func SendPDUSessionResourceSetupRequest(ue *context.RanUe, nasPdu []byte,
	pduSessionResourceSetupRequestList ngapType.PDUSessionResourceSetupListSUReq) {

	ngaplog.Info("[OCF] Send PDU Session Resource Setup Request")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(pduSessionResourceSetupRequestList.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPDUSessionResourceSetupRequest(ue, nasPdu, pduSessionResourceSetupRequestList)
	if err != nil {
		ngaplog.Errorf("Build PDUSessionResourceSetupRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// pduSessionResourceModifyConfirmList: provided by OCF, and transfer data is return from SMF
// pduSessionResourceFailedToModifyList: provided by OCF, and transfer data is return from SMF
func SendPDUSessionResourceModifyConfirm(
	ue *context.RanUe,
	pduSessionResourceModifyConfirmList ngapType.PDUSessionResourceModifyListModCfm,
	pduSessionResourceFailedToModifyList ngapType.PDUSessionResourceFailedToModifyListModCfm,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send PDU Session Resource Modify Confirm")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(pduSessionResourceModifyConfirmList.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	if len(pduSessionResourceFailedToModifyList.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPDUSessionResourceModifyConfirm(ue, pduSessionResourceModifyConfirmList,
		pduSessionResourceFailedToModifyList, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build PDUSessionResourceModifyConfirm failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// pduSessionResourceModifyRequestList: from SMF
func SendPDUSessionResourceModifyRequest(ue *context.RanUe,
	pduSessionResourceModifyRequestList ngapType.PDUSessionResourceModifyListModReq) {

	ngaplog.Info("[OCF] Send PDU Session Resource Modify Request")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(pduSessionResourceModifyRequestList.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPDUSessionResourceModifyRequest(ue, pduSessionResourceModifyRequestList)
	if err != nil {
		ngaplog.Errorf("Build PDUSessionResourceModifyRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendInitialContextSetupRequest(
	ocfUe *context.OcfUe,
	anType models.AccessType,
	nasPdu []byte,
	pduSessionResourceSetupRequestList *ngapType.PDUSessionResourceSetupListCxtReq,
	rrcInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest,
	coreNetworkAssistanceInfo *ngapType.CoreNetworkAssistanceInformation,
	emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator) {

	ngaplog.Info("[OCF] Send Initial Context Setup Request")

	if ocfUe == nil {
		ngaplog.Error("OcfUe is nil")
		return
	}

	if pduSessionResourceSetupRequestList != nil {
		if len(pduSessionResourceSetupRequestList.List) > context.MaxNumOfPDUSessions {
			ngaplog.Error("Pdu List out of range")
			return
		}
	}

	pkt, err := BuildInitialContextSetupRequest(ocfUe, anType, nasPdu, pduSessionResourceSetupRequestList,
		rrcInactiveTransitionReportRequest, coreNetworkAssistanceInfo, emergencyFallbackIndicator)
	if err != nil {
		ngaplog.Errorf("Build InitialContextSetupRequest failed : %s", err.Error())
		return
	}
	ocfUe.RanUe[anType].SentInitialContextSetupRequest = true
	NasSendToRan(ocfUe, anType, pkt)
}

func SendUEContextModificationRequest(
	ocfUe *context.OcfUe,
	anType models.AccessType,
	oldOcfUeNgapID *int64,
	rrcInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest,
	coreNetworkAssistanceInfo *ngapType.CoreNetworkAssistanceInformation,
	mobilityRestrictionList *ngapType.MobilityRestrictionList,
	emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator) {

	ngaplog.Info("[OCF] Send UE Context Modification Request")

	if ocfUe == nil {
		ngaplog.Error("OcfUe is nil")
		return
	}

	pkt, err := BuildUEContextModificationRequest(ocfUe, anType, oldOcfUeNgapID, rrcInactiveTransitionReportRequest,
		coreNetworkAssistanceInfo, mobilityRestrictionList, emergencyFallbackIndicator)
	if err != nil {
		ngaplog.Errorf("Build UEContextModificationRequest failed : %s", err.Error())
		return
	}
	NasSendToRan(ocfUe, anType, pkt)
}

// pduSessionResourceHandoverList: provided by ocf and transfer is return from smf
// pduSessionResourceToReleaseList: provided by ocf and transfer is return from smf
// criticalityDiagnostics = criticalityDiagonstics IE in receiver node's error indication
// when received node can't comprehend the IE or missing IE
func SendHandoverCommand(
	sourceUe *context.RanUe,
	pduSessionResourceHandoverList ngapType.PDUSessionResourceHandoverList,
	pduSessionResourceToReleaseList ngapType.PDUSessionResourceToReleaseListHOCmd,
	container ngapType.TargetToSourceTransparentContainer,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send Handover Command")

	if sourceUe == nil {
		ngaplog.Error("SourceUe is nil")
		return
	}

	if len(pduSessionResourceHandoverList.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	if len(pduSessionResourceToReleaseList.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildHandoverCommand(sourceUe, pduSessionResourceHandoverList, pduSessionResourceToReleaseList,
		container, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build HandoverCommand failed : %s", err.Error())
		return
	}
	SendToRanUe(sourceUe, pkt)
}

// cause = initiate the Handover Cancel procedure with the appropriate value for the Cause IE.
// criticalityDiagnostics = criticalityDiagonstics IE in receiver node's error indication
// when received node can't comprehend the IE or missing IE
func SendHandoverPreparationFailure(sourceUe *context.RanUe, cause ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send Handover Preparation Failure")

	if sourceUe == nil {
		ngaplog.Error("SourceUe is nil")
		return
	}
	ocfUe := sourceUe.OcfUe
	if ocfUe == nil {
		ngaplog.Error("ocfUe is nil")
		return
	}
	ocfUe.OnGoing[sourceUe.Ran.AnType].Procedure = context.OnGoingProcedureNothing
	pkt, err := BuildHandoverPreparationFailure(sourceUe, cause, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build HandoverPreparationFailure failed : %s", err.Error())
		return
	}
	SendToRanUe(sourceUe, pkt)
}

/*The PGW-C+SMF (V-SMF in the case of home-routed roaming scenario only) sends
a Nsmf_PDUSession_CreateSMContext Response(N2 SM Information (PDU Session ID, cause code)) to the OCF.*/
// Cause is from SMF
// pduSessionResourceSetupList provided by OCF, and the transfer data is from SMF
// sourceToTargetTransparentContainer is received from S-RAN
// nsci: new security context indicator, if ocfUe has updated security context, set nsci to true, otherwise set to false
// N2 handover in same OCF
func SendHandoverRequest(sourceUe *context.RanUe, targetRan *context.OcfRan, cause ngapType.Cause,
	pduSessionResourceSetupListHOReq ngapType.PDUSessionResourceSetupListHOReq,
	sourceToTargetTransparentContainer ngapType.SourceToTargetTransparentContainer, nsci bool) {

	ngaplog.Info("[OCF] Send Handover Request")

	if sourceUe == nil {
		ngaplog.Error("sourceUe is nil")
		return
	}
	ocfUe := sourceUe.OcfUe
	if ocfUe == nil {
		ngaplog.Error("ocfUe is nil")
		return
	}
	if targetRan == nil {
		ngaplog.Error("targetRan is nil")
		return
	}

	if sourceUe.TargetUe != nil {
		ngaplog.Error("Handover Required Duplicated")
		return
	}

	if len(pduSessionResourceSetupListHOReq.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	if len(sourceToTargetTransparentContainer.Value) == 0 {
		ngaplog.Error("Source To Target TransparentContainer is nil")
		return
	}

	var targetUe *context.RanUe
	if targetUeTmp, err := targetRan.NewRanUe(context.RanUeNgapIdUnspecified); err != nil {
		ngaplog.Errorf("Create target UE error: %+v", err)
	} else {
		targetUe = targetUeTmp
	}

	ngaplog.Tracef("Source : OCF_UE_NGAP_ID[%d], RAN_UE_NGAP_ID[%d]", sourceUe.OcfUeNgapId, sourceUe.RanUeNgapId)
	ngaplog.Tracef("Target : OCF_UE_NGAP_ID[%d], RAN_UE_NGAP_ID[Unknown]", targetUe.OcfUeNgapId)
	context.AttachSourceUeTargetUe(sourceUe, targetUe)

	pkt, err := BuildHandoverRequest(targetUe, cause, pduSessionResourceSetupListHOReq,
		sourceToTargetTransparentContainer, nsci)
	if err != nil {
		ngaplog.Errorf("Build HandoverRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(targetUe, pkt)
}

// pduSessionResourceSwitchedList: provided by OCF, and the transfer data is from SMF
// pduSessionResourceReleasedList: provided by OCF, and the transfer data is from SMF
// newSecurityContextIndicator: if OCF has activated a new 5G NAS security context, set it to true,
// otherwise set to false
// coreNetworkAssistanceInformation: provided by OCF, based on collection of UE behaviour statistics
// and/or other available
// information about the expected UE behaviour. TS 23.501 5.4.6, 5.4.6.2
// rrcInactiveTransitionReportRequest: configured by ocf
// criticalityDiagnostics: from received node when received not comprehended IE or missing IE
func SendPathSwitchRequestAcknowledge(
	ue *context.RanUe,
	pduSessionResourceSwitchedList ngapType.PDUSessionResourceSwitchedList,
	pduSessionResourceReleasedList ngapType.PDUSessionResourceReleasedListPSAck,
	newSecurityContextIndicator bool,
	coreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation,
	rrcInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send Path Switch Request Acknowledge")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(pduSessionResourceSwitchedList.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	if len(pduSessionResourceReleasedList.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPathSwitchRequestAcknowledge(ue, pduSessionResourceSwitchedList, pduSessionResourceReleasedList,
		newSecurityContextIndicator, coreNetworkAssistanceInformation, rrcInactiveTransitionReportRequest,
		criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build PathSwitchRequestAcknowledge failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// pduSessionResourceReleasedList: provided by OCF, and the transfer data is from SMF
// criticalityDiagnostics: from received node when received not comprehended IE or missing IE
func SendPathSwitchRequestFailure(
	ran *context.OcfRan,
	ocfUeNgapId,
	ranUeNgapId int64,
	pduSessionResourceReleasedList *ngapType.PDUSessionResourceReleasedListPSFail,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send Path Switch Request Failure")

	if pduSessionResourceReleasedList != nil && len(pduSessionResourceReleasedList.List) > context.MaxNumOfPDUSessions {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildPathSwitchRequestFailure(ocfUeNgapId, ranUeNgapId, pduSessionResourceReleasedList,
		criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build PathSwitchRequestFailure failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

//RanStatusTransferTransparentContainer from Uplink Ran Configuration Transfer
func SendDownlinkRanStatusTransfer(ue *context.RanUe, container ngapType.RANStatusTransferTransparentContainer) {

	ngaplog.Info("[OCF] Send Downlink Ran Status Transfer")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(container.DRBsSubjectToStatusTransferList.List) > context.MaxNumOfDRBs {
		ngaplog.Error("Pdu List out of range")
		return
	}

	pkt, err := BuildDownlinkRanStatusTransfer(ue, container)
	if err != nil {
		ngaplog.Errorf("Build DownlinkRanStatusTransfer failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// anType indicate ocfUe send this msg for which accessType
// Paging Priority: is included only if the OCF receives an Nocf_Communication_N1N2MessageTransfer message
// with an ARP value associated with
// priority services (e.g., MPS, MCS), as configured by the operator. (TS 23.502 4.2.3.3, TS 23.501 5.22.3)
// pagingOriginNon3GPP: TS 23.502 4.2.3.3 step 4b: If the UE is simultaneously registered over 3GPP and non-3GPP
// accesses in the same PLMN,
// the UE is in CM-IDLE state in both 3GPP access and non-3GPP access, and the PDU Session ID in step 3a
// is associated with non-3GPP access, the OCF sends a Paging message with associated access "non-3GPP" to
// NG-RAN node(s) via 3GPP access.
// more paging policy with 3gpp/non-3gpp access is described in TS 23.501 5.6.8
func SendPaging(ue *context.OcfUe, ngapBuf []byte) {

	// var pagingPriority *ngapType.PagingPriority
	if ue == nil {
		ngaplog.Error("OcfUe is nil")
		return
	}

	// if ppi != nil {
	// pagingPriority = new(ngapType.PagingPriority)
	// pagingPriority.Value = aper.Enumerated(*ppi)
	// }
	// pkt, err := BuildPaging(ue, pagingPriority, pagingOriginNon3GPP)
	// if err != nil {
	// 	ngaplog.Errorf("Build Paging failed : %s", err.Error())
	// }
	taiList := ue.RegistrationArea[models.AccessType__3_GPP_ACCESS]
	context.OCF_Self().OcfRanPool.Range(func(key, value interface{}) bool {
		ran := value.(*context.OcfRan)
		for _, item := range ran.SupportedTAList {
			if context.InTaiList(item.Tai, taiList) {
				ngaplog.Infof("[OCF] Send Paging to TAI(%+v, Tac:%+v) for Ue[%s]",
					item.Tai.PlmnId, item.Tai.Tac, ue.Supi)
				SendToRan(ran, ngapBuf)
				break
			}
		}
		return true
	})

	ue.T3513RetryTimes = 0
	ue.T3513 = time.AfterFunc(context.TimeT3513, func() {
		ue.T3513RetryTimes++
		if ue.T3513RetryTimes > context.MaxT3513RetryTimes {
			logger.GmmLog.Warnf("UE[%s] T3513 expires %d times, abort paging procedure", ue.Supi, ue.T3513RetryTimes)
			if ue.OnGoing[models.AccessType__3_GPP_ACCESS].Procedure != context.OnGoingProcedureN2Handover {
				callback.SendN1N2TransferFailureNotification(ue, models.N1N2MessageTransferCause_UE_NOT_RESPONDING)
			}
			util.StopT3513(ue)
		} else {
			logger.NgapLog.Warnf("[NGAP] T3513 expires, retransmit Paging (UE: [%s], retry: %d)",
				ue.Supi, ue.T3513RetryTimes)
			context.OCF_Self().OcfRanPool.Range(func(key, value interface{}) bool {
				ran := value.(*context.OcfRan)
				for _, item := range ran.SupportedTAList {
					if context.InTaiList(item.Tai, taiList) {
						SendToRan(ran, ngapBuf)
						break
					}
				}
				return true
			})
			ue.T3513.Reset(context.TimeT3513)
		}
	})
}

// TS 23.502 4.2.2.2.3
// anType: indicate ocfUe send this msg for which accessType
// ocfUeNgapID: initial OCF get it from target OCF
// ngapMessage: initial UE Message to reroute
// allowedNSSAI: provided by OCF, and OCF get it from NSSF (4.2.2.2.3 step 4b)
func SendRerouteNasRequest(ue *context.OcfUe, anType models.AccessType, ocfUeNgapID *int64, ngapMessage []byte,
	allowedNSSAI *ngapType.AllowedNSSAI) {

	ngaplog.Info("[OCF] Send Reroute Nas Request")

	if ue == nil {
		ngaplog.Error("OcfUe is nil")
		return
	}

	if len(ngapMessage) == 0 {
		ngaplog.Error("Ngap Message is nil")
		return
	}

	pkt, err := BuildRerouteNasRequest(ue, anType, ocfUeNgapID, ngapMessage, allowedNSSAI)
	if err != nil {
		ngaplog.Errorf("Build RerouteNasRequest failed : %s", err.Error())
		return
	}
	NasSendToRan(ue, anType, pkt)
}

// criticality ->from received node when received node can't comprehend the IE or missing IE
func SendRanConfigurationUpdateAcknowledge(
	ran *context.OcfRan, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send Ran Configuration Update Acknowledge")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildRanConfigurationUpdateAcknowledge(criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build RanConfigurationUpdateAcknowledge failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

// criticality ->from received node when received node can't comprehend the IE or missing IE
// If the OCF cannot accept the update,
// it shall respond with a RAN CONFIGURATION UPDATE FAILURE message and appropriate cause value.
func SendRanConfigurationUpdateFailure(ran *context.OcfRan, cause ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngaplog.Info("[OCF] Send Ran Configuration Update Failure")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildRanConfigurationUpdateFailure(cause, criticalityDiagnostics)
	if err != nil {
		ngaplog.Errorf("Build RanConfigurationUpdateFailure failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

//An OCF shall be able to instruct other peer CP NFs, subscribed to receive such a notification,
//that it will be unavailable on this OCF and its corresponding target OCF(s).
//If CP NF does not subscribe to receive OCF unavailable notification, the CP NF may attempt
//forwarding the transaction towards the old OCF and detect that the OCF is unavailable. When
//it detects unavailable, it marks the OCF and its associated GUAMI(s) as unavailable.
//Defined in 23.501 5.21.2.2.2
func SendOCFStatusIndication(ran *context.OcfRan, unavailableGUAMIList ngapType.UnavailableGUAMIList) {

	ngaplog.Info("[OCF] Send OCF Status Indication")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	if len(unavailableGUAMIList.List) > context.MaxNumOfServedGuamiList {
		ngaplog.Error("GUAMI List out of range")
		return
	}

	pkt, err := BuildOCFStatusIndication(unavailableGUAMIList)
	if err != nil {
		ngaplog.Errorf("Build OCFStatusIndication failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

// TS 23.501 5.19.5.2
// ocfOverloadResponse: the required behaviour of NG-RAN, provided by OCF
// ocfTrafficLoadReductionIndication(int 1~99): indicates the percentage of the type, set to 0 if does not need this ie
// of traffic relative to the instantaneous incoming rate at the NG-RAN node, provided by OCF
// overloadStartNSSAIList: overload slices, provide by OCF
func SendOverloadStart(
	ran *context.OcfRan,
	ocfOverloadResponse *ngapType.OverloadResponse,
	ocfTrafficLoadReductionIndication int64,
	overloadStartNSSAIList *ngapType.OverloadStartNSSAIList) {

	ngaplog.Info("[OCF] Send Overload Start")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	if ocfTrafficLoadReductionIndication != 0 &&
		(ocfTrafficLoadReductionIndication < 1 || ocfTrafficLoadReductionIndication > 99) {
		ngaplog.Error("OcfTrafficLoadReductionIndication out of range (should be 1 ~ 99)")
		return
	}

	if overloadStartNSSAIList != nil && len(overloadStartNSSAIList.List) > context.MaxNumOfSlice {
		ngaplog.Error("NSSAI List out of range")
		return
	}

	pkt, err := BuildOverloadStart(ocfOverloadResponse, ocfTrafficLoadReductionIndication, overloadStartNSSAIList)
	if err != nil {
		ngaplog.Errorf("Build OverloadStart failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

func SendOverloadStop(ran *context.OcfRan) {

	ngaplog.Info("[OCF] Send Overload Stop")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildOverloadStop()
	if err != nil {
		ngaplog.Errorf("Build OverloadStop failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

// SONConfigurationTransfer = sONConfigurationTransfer from uplink Ran Configuration Transfer
func SendDownlinkRanConfigurationTransfer(ran *context.OcfRan, transfer *ngapType.SONConfigurationTransfer) {

	ngaplog.Info("[OCF] Send Downlink Ran Configuration Transfer")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildDownlinkRanConfigurationTransfer(transfer)
	if err != nil {
		ngaplog.Errorf("Build DownlinkRanConfigurationTransfer failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

//NRPPa PDU is by pass
//NRPPa PDU is from LMF define in 4.13.5.6
func SendDownlinkNonUEAssociatedNRPPATransport(ue *context.RanUe, nRPPaPDU ngapType.NRPPaPDU) {

	ngaplog.Info("[OCF] Send Downlink Non UE Associated NRPPA Transport")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(nRPPaPDU.Value) == 0 {
		ngaplog.Error("length of NRPPA-PDU is 0")
		return
	}

	pkt, err := BuildDownlinkNonUEAssociatedNRPPATransport(ue, nRPPaPDU)
	if err != nil {
		ngaplog.Errorf("Build DownlinkNonUEAssociatedNRPPATransport failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendDeactivateTrace(ocfUe *context.OcfUe, anType models.AccessType) {

	ngaplog.Info("[OCF] Send Deactivate Trace")

	if ocfUe == nil {
		ngaplog.Error("OcfUe is nil")
		return
	}

	ranUe := ocfUe.RanUe[anType]
	if ranUe == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildDeactivateTrace(ocfUe, anType)
	if err != nil {
		ngaplog.Errorf("Build DeactivateTrace failed : %s", err.Error())
		return
	}
	SendToRanUe(ranUe, pkt)
}

// AOI List is from SMF
// The SMF may subscribe to the UE mobility event notification from the OCF
// (e.g. location reporting, UE moving into or out of Area Of Interest) TS 23.502 4.3.2.2.1 Step.17
// The Location Reporting Control message shall identify the UE for which reports are requested and may include
// Reporting Type, Location Reporting Level, Area Of Interest and Request Reference ID
// TS 23.502 4.10 LocationReportingProcedure
// The OCF may request the NG-RAN location reporting with event reporting type (e.g. UE location or UE presence
// in Area of Interest), reporting mode and its related parameters (e.g. number of reporting) TS 23.501 5.4.7
// Location Reference ID To Be Cancelled IE shall be present if the Event Type IE is set to "Stop UE presence
// in the area of interest". otherwise set it to 0
func SendLocationReportingControl(
	ue *context.RanUe,
	AOIList *ngapType.AreaOfInterestList,
	LocationReportingReferenceIDToBeCancelled int64,
	eventType ngapType.EventType) {

	ngaplog.Info("[OCF] Send Location Reporting Control")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if AOIList != nil && len(AOIList.List) > context.MaxNumOfAOI {
		ngaplog.Error("AOI List out of range")
		return
	}

	if eventType.Value == ngapType.EventTypePresentStopUePresenceInAreaOfInterest {
		if LocationReportingReferenceIDToBeCancelled < 1 || LocationReportingReferenceIDToBeCancelled > 64 {
			ngaplog.Error("LocationReportingReferenceIDToBeCancelled out of range (should be 1 ~ 64)")
			return
		}
	}

	pkt, err := BuildLocationReportingControl(ue, AOIList, LocationReportingReferenceIDToBeCancelled, eventType)
	if err != nil {
		ngaplog.Errorf("Build LocationReportingControl failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

func SendUETNLABindingReleaseRequest(ue *context.RanUe) {

	ngaplog.Info("[OCF] Send UE TNLA Binging Release Request")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	pkt, err := BuildUETNLABindingReleaseRequest(ue)
	if err != nil {
		ngaplog.Errorf("Build UETNLABindingReleaseRequest failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

// Weight Factor associated with each of the TNL association within the OCF
func SendOCFConfigurationUpdate(ran *context.OcfRan, usage ngapType.TNLAssociationUsage,
	weightfactor ngapType.TNLAddressWeightFactor) {

	ngaplog.Info("[OCF] Send OCF Configuration Update")

	if ran == nil {
		ngaplog.Error("Ran is nil")
		return
	}

	pkt, err := BuildOCFConfigurationUpdate(usage, weightfactor)
	if err != nil {
		ngaplog.Errorf("Build OCFConfigurationUpdate failed : %s", err.Error())
		return
	}
	SendToRan(ran, pkt)
}

//NRPPa PDU is a pdu from LMF to RAN defined in TS 23.502 4.13.5.5 step 3
//NRPPa PDU is by pass
func SendDownlinkUEAssociatedNRPPaTransport(ue *context.RanUe, nRPPaPDU ngapType.NRPPaPDU) {

	ngaplog.Info("[OCF] Send Downlink UE Associated NRPPa Transport")

	if ue == nil {
		ngaplog.Error("RanUe is nil")
		return
	}

	if len(nRPPaPDU.Value) == 0 {
		ngaplog.Error("length of NRPPA-PDU is 0")
		return
	}

	pkt, err := BuildDownlinkUEAssociatedNRPPaTransport(ue, nRPPaPDU)
	if err != nil {
		ngaplog.Errorf("Build DownlinkUEAssociatedNRPPaTransport failed : %s", err.Error())
		return
	}
	SendToRanUe(ue, pkt)
}

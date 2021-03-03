package message

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/ocf/context"
	"free5gc/src/ocf/logger"
	ngap_message "free5gc/src/ocf/ngap/message"
	"free5gc/src/ocf/producer/callback"
	"free5gc/src/ocf/util"
	"time"
)

// backOffTimerUint = 7 means backoffTimer is null
func SendDLNASTransport(ue *context.RanUe, payloadContainerType uint8, nasPdu []byte,
	pduSessionId int32, cause uint8, backOffTimerUint *uint8, backOffTimer uint8) {

	logger.GmmLog.Info("[NAS] Send DL NAS Transport")
	var causePtr *uint8
	if cause != 0 {
		causePtr = &cause
	}
	nasMsg, err := BuildDLNASTransport(ue.OcfUe, payloadContainerType, nasPdu,
		uint8(pduSessionId), causePtr, backOffTimerUint, backOffTimer)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
}

func SendNotification(ue *context.RanUe, nasMsg []byte) {

	logger.GmmLog.Info("[NAS] Send Notification")

	amfUe := ue.OcfUe
	if amfUe == nil {
		logger.GmmLog.Error("OcfUe is nil")
		return
	}

	amfUe.T3565 = time.AfterFunc(context.TimeT3565, func() {
		amfUe.T3565RetryTimes++
		if amfUe.T3565RetryTimes > context.MaxT3565RetryTimes {
			logger.GmmLog.Warnf("UE[%s] T3565 Expires %d times, abort notification procedure",
				amfUe.Supi, amfUe.T3565RetryTimes)
			if amfUe.OnGoing[models.AccessType__3_GPP_ACCESS].Procedure != context.OnGoingProcedureN2Handover {
				callback.SendN1N2TransferFailureNotification(amfUe, models.N1N2MessageTransferCause_UE_NOT_RESPONDING)
			}
			util.StopT3565(amfUe)
		} else {
			logger.GmmLog.Warnf("[NAS] T3565 expires, retransmit Notification (retry: %d)", amfUe.T3565RetryTimes)
			ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
			amfUe.T3565.Reset(context.TimeT3565)
		}
	})
}

func SendIdentityRequest(ue *context.RanUe, typeOfIdentity uint8) {

	logger.GmmLog.Info("[NAS] Send Identity Request")

	nasMsg, err := BuildIdentityRequest(typeOfIdentity)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
}

func SendAuthenticationRequest(ue *context.RanUe) {

	amfUe := ue.OcfUe
	if amfUe == nil {
		logger.GmmLog.Error("OcfUe is nil")
		return
	}

	logger.GmmLog.Infof("[NAS] Send Authentication Request")

	if amfUe.AuthenticationCtx == nil {
		logger.GmmLog.Error("Authentication Context of UE is nil")
		return
	}

	nasMsg, err := BuildAuthenticationRequest(amfUe)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)

	amfUe.T3560RetryTimes = 0
	amfUe.T3560 = time.AfterFunc(context.TimeT3560, func() {
		amfUe.T3560RetryTimes++
		if amfUe.T3560RetryTimes > context.MaxT3560RetryTimes {
			logger.GmmLog.Warnf("T3560 Expires %d times, abort authentication procedure & ongoing 5GMM procedure",
				amfUe.T3560RetryTimes)
			util.StopT3560(amfUe)
			amfUe.Remove()
		} else {
			logger.GmmLog.Warnf("[NAS] T3560 expires, retransmit Authentication Request (retry: %d)", amfUe.T3560RetryTimes)
			ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
			amfUe.T3560.Reset(context.TimeT3560)
		}
	})
}

func SendServiceAccept(ue *context.RanUe, pDUSessionStatus *[16]bool, reactivationResult *[16]bool,
	errPduSessionId, errCause []uint8) {

	logger.GmmLog.Info("[NAS] Send Service Accept")

	nasMsg, err := BuildServiceAccept(ue.OcfUe, pDUSessionStatus, reactivationResult, errPduSessionId, errCause)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
}

func SendConfigurationUpdateCommand(amfUe *context.OcfUe, accessType models.AccessType,
	networkSlicingIndication *nasType.NetworkSlicingIndication) {

	logger.GmmLog.Info("[NAS] Configuration Update Command")

	nasMsg, err := BuildConfigurationUpdateCommand(amfUe, accessType, networkSlicingIndication)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	mobilityRestrictionList := ngap_message.BuildIEMobilityRestrictionList(amfUe)
	ngap_message.SendDownlinkNasTransport(amfUe.RanUe[accessType], nasMsg, &mobilityRestrictionList)
}

func SendAuthenticationReject(ue *context.RanUe, eapMsg string) {

	logger.GmmLog.Info("[NAS] Send Authentication Reject")

	nasMsg, err := BuildAuthenticationReject(ue.OcfUe, eapMsg)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
}

func SendAuthenticationResult(ue *context.RanUe, eapSuccess bool, eapMsg string) {

	logger.GmmLog.Info("[NAS] Send Authentication Result")

	if ue.OcfUe == nil {
		logger.GmmLog.Errorf("OcfUe is nil")
		return
	}

	nasMsg, err := BuildAuthenticationResult(ue.OcfUe, eapSuccess, eapMsg)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
}
func SendServiceReject(ue *context.RanUe, pDUSessionStatus *[16]bool, cause uint8) {

	logger.GmmLog.Info("[NAS] Send Service Reject")

	nasMsg, err := BuildServiceReject(pDUSessionStatus, cause)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
}

// T3502: This IE may be included to indicate a value for timer T3502 during the initial registration
// eapMessage: if the REGISTRATION REJECT message is used to convey EAP-failure message
func SendRegistrationReject(ue *context.RanUe, cause5GMM uint8, eapMessage string) {

	logger.GmmLog.Info("[NAS] Send Registration Reject")

	nasMsg, err := BuildRegistrationReject(ue.OcfUe, cause5GMM, eapMessage)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
}

// eapSuccess: only used when authType is EAP-AKA', set the value to false if authType is not EAP-AKA'
// eapMessage: only used when authType is EAP-AKA', set the value to "" if authType is not EAP-AKA'
func SendSecurityModeCommand(ue *context.RanUe, eapSuccess bool, eapMessage string) {

	logger.GmmLog.Info("[NAS] Send Security Mode Command")

	nasMsg, err := BuildSecurityModeCommand(ue.OcfUe, eapSuccess, eapMessage)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)

	amfUe := ue.OcfUe

	amfUe.T3560RetryTimes = 0
	amfUe.T3560 = time.AfterFunc(context.TimeT3560, func() {
		amfUe.T3560RetryTimes++
		if amfUe.T3560RetryTimes > context.MaxT3560RetryTimes {
			logger.GmmLog.Warnf("T3560 Expires %d times, abort security mode control procedure", amfUe.T3560RetryTimes)
			util.StopT3560(amfUe)
			amfUe.Remove()
		} else {
			logger.GmmLog.Warnf("[NAS] T3560 expires, retransmit Security Mode Command (retry: %d)", amfUe.T3560RetryTimes)
			ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
			amfUe.T3560.Reset(context.TimeT3560)
		}
	})
}

func SendDeregistrationRequest(ue *context.RanUe, accessType uint8, reRegistrationRequired bool, cause5GMM uint8) {

	logger.GmmLog.Info("[NAS] Send Deregistration Request")

	nasMsg, err := BuildDeregistrationRequest(ue, accessType, reRegistrationRequired, cause5GMM)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)

	amfUe := ue.OcfUe

	amfUe.T3522RetryTimes = 0
	amfUe.T3522 = time.AfterFunc(context.TimeT3522, func() {
		amfUe.T3522RetryTimes++
		if amfUe.T3522RetryTimes > context.MaxT3522RetryTimes {
			logger.GmmLog.Warnf("T3522 Expires %d times, abort deregistration procedure", amfUe.T3522RetryTimes)
			if accessType == nasMessage.AccessType3GPP {
				logger.GmmLog.Warnln("UE accessType3GPP transfer to Deregistered state")
				amfUe.State[models.AccessType__3_GPP_ACCESS].Set(context.Deregistered)
			} else if accessType == nasMessage.AccessTypeNon3GPP {
				logger.GmmLog.Warnln("UE accessTypeNon3GPP transfer to Deregistered state")
				amfUe.State[models.AccessType_NON_3_GPP_ACCESS].Set(context.Deregistered)
			} else {
				logger.GmmLog.Warnln("UE accessType3GPP transfer to Deregistered state")
				amfUe.State[models.AccessType__3_GPP_ACCESS].Set(context.Deregistered)
				logger.GmmLog.Warnln("UE accessTypeNon3GPP transfer to Deregistered state")
				amfUe.State[models.AccessType_NON_3_GPP_ACCESS].Set(context.Deregistered)
			}
			util.StopT3522(amfUe)
		} else {
			logger.GmmLog.Warnf("[NAS] T3522 expires, retransmit Deregistration Request (retry: %d)", amfUe.T3522RetryTimes)
			ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
			amfUe.T3522.Reset(context.TimeT3522)
		}
	})
}

func SendDeregistrationAccept(ue *context.RanUe) {

	logger.GmmLog.Info("[NAS] Send Deregistration Accept")

	nasMsg, err := BuildDeregistrationAccept()
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
}

func SendRegistrationAccept(
	ue *context.OcfUe,
	anType models.AccessType,
	pDUSessionStatus *[16]bool,
	reactivationResult *[16]bool,
	errPduSessionId, errCause []uint8,
	pduSessionResourceSetupList *ngapType.PDUSessionResourceSetupListCxtReq) {

	logger.GmmLog.Info("[NAS] Send Registration Accept")

	nasMsg, err := BuildRegistrationAccept(ue, anType, pDUSessionStatus, reactivationResult, errPduSessionId, errCause)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}

	if ue.RanUe[anType].UeContextRequest {
		ngap_message.SendInitialContextSetupRequest(ue, anType, nasMsg, pduSessionResourceSetupList, nil, nil, nil)
	} else {
		ngap_message.SendDownlinkNasTransport(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasMsg, nil)
	}

	ue.T3550RetryTimes = 0
	ue.T3550 = time.AfterFunc(context.TimeT3550, func() {
		ue.T3550RetryTimes++
		if ue.T3550RetryTimes > context.MaxT3550RetryTimes {
			logger.GmmLog.Warnf("T3550 Expires %d times, abort retransmission of Registration Accept", ue.T3550RetryTimes)
			// TS 24.501 5.5.1.2.8 case c, 5.5.1.3.8 case c
			ue.State[anType].Set(context.Registered)
			ue.ClearRegistrationRequestData(anType)
			util.StopT3550(ue)
		} else {
			logger.GmmLog.Warnf("[NAS] T3550 expires, retransmit Registration Accept (retry: %d)", ue.T3550RetryTimes)
			ngap_message.SendDownlinkNasTransport(ue.RanUe[anType], nasMsg, nil)
			ue.T3550.Reset(context.TimeT3550)
		}
	})
}

func SendStatus5GMM(ue *context.RanUe, cause uint8) {

	logger.GmmLog.Info("[NAS] Send Status 5GMM")

	nasMsg, err := BuildStatus5GMM(cause)
	if err != nil {
		logger.GmmLog.Error(err.Error())
		return
	}
	ngap_message.SendDownlinkNasTransport(ue, nasMsg, nil)
}
package consumer

import (
	"context"
	"free5gcWithOCF/lib/openapi"
	"free5gcWithOCF/lib/openapi/Nsmf_PDUSession"
	"free5gcWithOCF/lib/openapi/models"
	ocf_context "free5gcWithOCF/src/ocf/context"
	"strconv"
)

type UpdateSmContextPresent string

const (
	UpdateSmContextPresentActivateUpCnxState        UpdateSmContextPresent = "Activate_User_Plane_Connectivity"
	UpdateSmContextPresentDeactivateUpCnxState      UpdateSmContextPresent = "Dectivate_User_Plane_Connectivity"
	UpdateSmContextPresentChangeAccessType          UpdateSmContextPresent = "Change_AccessType"
	UpdateSmContextPresentXnHandover                UpdateSmContextPresent = "Xn_Handover"
	UpdateSmContextPresentXnHandoverFailed          UpdateSmContextPresent = "Xn_Handover_Failed"
	UpdateSmContextPresentN2HandoverPreparing       UpdateSmContextPresent = "N2_Handover_Preparing"
	UpdateSmContextPresentN2HandoverPrepared        UpdateSmContextPresent = "N2_Handover_Prepared"
	UpdateSmContextPresentN2HandoverComplete        UpdateSmContextPresent = "N2_Handover_Complete"
	UpdateSmContextPresentN2HandoverCanceled        UpdateSmContextPresent = "N2_Handover_Canceled"
	UpdateSmContextPresentHandoverBetweenAccessType UpdateSmContextPresent = "Handover_Between_AccessType"
	UpdateSmContextPresentHandoverBetweenOCF        UpdateSmContextPresent = "Handover_Between_OCF"
	UpdateSmContextPresentOnlyN2SmInfo              UpdateSmContextPresent = "N2SmInfo"
)

type updateSmContextRequsetParam struct {
	accessType         models.AccessType
	cause              ocf_context.CauseAll
	n2SmType           models.N2SmInfoType
	anTypeCanBeChanged bool
}
type updateSmContextRequsetHandoverParam struct {
	accessType models.AccessType
	targetId   *models.NgRanTargetId
	guami      *models.Guami
	ocfid      string
	cause      ocf_context.CauseAll
	n2SmType   models.N2SmInfoType
	n1SmMsg    bool
	activation bool
}

func SendCreateSmContextRequest(
	ue *ocf_context.OcfUe, smfUri string, nasPdu []byte, smContextCreateData models.SmContextCreateData) (
	response *models.PostSmContextsResponse, smContextRef string, errorResponse *models.PostSmContextsErrorResponse,
	problemDetail *models.ProblemDetails, err1 error) {
	configuration := Nsmf_PDUSession.NewConfiguration()
	configuration.SetBasePath(smfUri)

	client := Nsmf_PDUSession.NewAPIClient(configuration)

	var postSmContextsRequest models.PostSmContextsRequest
	postSmContextsRequest.JsonData = &smContextCreateData
	postSmContextsRequest.BinaryDataN1SmMessage = nasPdu

	postSmContextReponse, httpResponse, err :=
		client.SMContextsCollectionApi.PostSmContexts(context.Background(), postSmContextsRequest)
	if err == nil {
		response = &postSmContextReponse
		smContextRef = httpResponse.Header.Get("Location")
	} else if httpResponse != nil {
		if httpResponse.Status != err.Error() {
			err1 = err
			return
		}
		switch httpResponse.StatusCode {
		case 400, 403, 404, 500, 503, 504:
			errResponse := err.(openapi.GenericOpenAPIError).Model().(models.PostSmContextsErrorResponse)
			errorResponse = &errResponse
		case 411, 413, 415, 429:
			problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
			problemDetail = &problem
		}
	} else {
		err1 = openapi.ReportError("server no response")
	}
	return response, smContextRef, errorResponse, problemDetail, err1
}
func BuildCreateSmContextRequest(ue *ocf_context.OcfUe, pduSessionContext models.PduSessionContext,
	requestType models.RequestType) (smContextCreateData models.SmContextCreateData) {
	context := ocf_context.OCF_Self()
	smContextCreateData.Supi = ue.Supi
	smContextCreateData.UnauthenticatedSupi = ue.UnauthenticatedSupi
	smContextCreateData.Pei = ue.Pei
	smContextCreateData.Gpsi = ue.Gpsi
	smContextCreateData.PduSessionId = pduSessionContext.PduSessionId
	smContextCreateData.SNssai = pduSessionContext.SNssai
	smContextCreateData.Dnn = pduSessionContext.Dnn
	smContextCreateData.ServingNfId = context.NfId
	smContextCreateData.Guami = &context.ServedGuamiList[0]
	smContextCreateData.ServingNetwork = context.ServedGuamiList[0].PlmnId
	if requestType == models.RequestType_EXISTING_PDU_SESSION ||
		requestType == models.RequestType_EXISTING_EMERGENCY_PDU_SESSION {
		smContextCreateData.RequestType = requestType
	}
	smContextCreateData.N1SmMsg = new(models.RefToBinaryData)
	smContextCreateData.N1SmMsg.ContentId = "n1SmMsg"
	smContextCreateData.AnType = pduSessionContext.AccessType
	if ue.RatType != "" {
		smContextCreateData.RatType = ue.RatType
	}
	// TODO: location is used in roaming scenerio
	// if ue.Location != nil {
	// 	smContextCreateData.UeLocation = ue.Location
	// }
	smContextCreateData.UeTimeZone = ue.TimeZone
	smContextCreateData.SmContextStatusUri = context.GetIPv4Uri() + "/nocf-callback/v1/smContextStatus/" +
		ue.Guti + "/" + strconv.Itoa(int(pduSessionContext.PduSessionId))

	return smContextCreateData
}

// Upadate SmContext Request
// servingNfId, smContextStatusUri, guami, servingNetwork -> ocf change
// anType -> anType change
// ratType -> ratType change
// presenceInLadn -> Service Request , Xn handover, N2 handover and dnn is a ladn
// ueLocation -> the user location has changed or the user plane of the PDU session is deactivated
// upCnxState -> request the activation or the deactivation of the user plane connection of the PDU session
// hoState -> the preparation, execution or cancellation of a handover of the PDU session
// toBeSwitch -> Xn Handover to request to switch the PDU session to a new downlink N3 tunnel endpoint
// failedToBeSwitch -> indicate that the PDU session failed to be setup in the target RAN
// targetId, targetServingNfId(preparation with OCF change) -> N2 handover
// release -> duplicated PDU Session Id in subclause 5.2.2.3.11, slice not available in subclause 5.2.2.3.12
// ngApCause -> e.g. the NGAP cause for requesting to deactivate the user plane connection of the PDU session.
// 5gMmCauseValue -> OCF received a 5GMM cause code from the UE e.g 5GMM Status message in response to
// a Downlink NAS Transport message carrying 5GSM payload
// anTypeCanBeChanged

func SendUpdateSmContextActivateUpCnxState(ue *ocf_context.OcfUe, pduSessionId int32, accessType models.AccessType) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetParam{
		accessType: accessType,
	}
	updateData := BuildUpdateSmContextRequset(ue, UpdateSmContextPresentActivateUpCnxState, pduSessionId, param)
	return SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, nil)
}

func SendUpdateSmContextDeactivateUpCnxState(ue *ocf_context.OcfUe, pduSessionId int32, cause ocf_context.CauseAll) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetParam{
		cause: cause,
	}
	updateData := BuildUpdateSmContextRequset(ue, UpdateSmContextPresentDeactivateUpCnxState, pduSessionId, param)
	return SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, nil)
}
func SendUpdateSmContextChangeAccessType(ue *ocf_context.OcfUe, pduSessionId int32, anTypeCanBeChanged bool) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetParam{
		anTypeCanBeChanged: anTypeCanBeChanged,
	}
	updateData := BuildUpdateSmContextRequset(ue, UpdateSmContextPresentChangeAccessType, pduSessionId, param)
	return SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, nil)
}

func SendUpdateSmContextN2Info(
	ue *ocf_context.OcfUe, pduSessionId int32, n2SmType models.N2SmInfoType, N2SmInfo []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetParam{
		n2SmType: n2SmType,
	}
	updateData := BuildUpdateSmContextRequset(ue, UpdateSmContextPresentOnlyN2SmInfo, pduSessionId, param)
	return SendUpdateSmContextRequest(
		ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, N2SmInfo)
}

func SendUpdateSmContextXnHandover(
	ue *ocf_context.OcfUe, pduSessionId int32, n2SmType models.N2SmInfoType, N2SmInfo []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetHandoverParam{
		n2SmType: n2SmType,
	}
	updateData := BuildUpdateSmContextRequsetHandover(ue, UpdateSmContextPresentXnHandover, pduSessionId, param)
	return SendUpdateSmContextRequest(
		ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, N2SmInfo)
}

func SendUpdateSmContextXnHandoverFailed(
	ue *ocf_context.OcfUe, pduSessionId int32, n2SmType models.N2SmInfoType, N2SmInfo []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetHandoverParam{
		n2SmType: n2SmType,
	}
	updateData :=
		BuildUpdateSmContextRequsetHandover(ue, UpdateSmContextPresentXnHandoverFailed, pduSessionId, param)
	return SendUpdateSmContextRequest(
		ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, N2SmInfo)
}

func SendUpdateSmContextN2HandoverPreparing(ue *ocf_context.OcfUe, pduSessionId int32, n2SmType models.N2SmInfoType,
	N2SmInfo []byte, ocfid string, targetId *models.NgRanTargetId) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetHandoverParam{
		targetId: targetId,
		ocfid:    ocfid,
		n2SmType: n2SmType,
	}
	updateData := BuildUpdateSmContextRequsetHandover(ue, UpdateSmContextPresentN2HandoverPreparing, pduSessionId, param)
	return SendUpdateSmContextRequest(
		ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, N2SmInfo)
}
func SendUpdateSmContextN2HandoverPrepared(
	ue *ocf_context.OcfUe, pduSessionId int32, n2SmType models.N2SmInfoType, N2SmInfo []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetHandoverParam{
		n2SmType: n2SmType,
	}
	updateData := BuildUpdateSmContextRequsetHandover(ue, UpdateSmContextPresentN2HandoverPrepared, pduSessionId, param)
	return SendUpdateSmContextRequest(
		ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, N2SmInfo)
}

func SendUpdateSmContextN2HandoverComplete(
	ue *ocf_context.OcfUe, pduSessionId int32, ocfid string, guami *models.Guami) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetHandoverParam{
		guami: guami,
		ocfid: ocfid,
	}
	updateData := BuildUpdateSmContextRequsetHandover(ue, UpdateSmContextPresentN2HandoverComplete, pduSessionId, param)
	return SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, nil)
}
func SendUpdateSmContextN2HandoverCanceled(ue *ocf_context.OcfUe, pduSessionId int32, cause ocf_context.CauseAll) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetHandoverParam{
		cause: cause,
	}
	updateData := BuildUpdateSmContextRequsetHandover(ue, UpdateSmContextPresentN2HandoverCanceled, pduSessionId, param)
	return SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, nil)
}

func SendUpdateSmContextHandoverBetweenAccessType(
	ue *ocf_context.OcfUe, pduSessionId int32, targetAccessType models.AccessType, N1SmMsg []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	isN1SmMsg := false
	if N1SmMsg != nil {
		isN1SmMsg = true
	}
	param := updateSmContextRequsetHandoverParam{
		accessType: targetAccessType,
		n1SmMsg:    isN1SmMsg,
	}
	updateData :=
		BuildUpdateSmContextRequsetHandover(ue, UpdateSmContextPresentHandoverBetweenAccessType, pduSessionId, param)
	return SendUpdateSmContextRequest(
		ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, N1SmMsg, nil)
}

func SendUpdateSmContextHandoverBetweenOCF(
	ue *ocf_context.OcfUe, pduSessionId int32, ocfid string, guami *models.Guami, activate bool) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		return nil, nil, nil, openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
	}
	param := updateSmContextRequsetHandoverParam{
		guami:      guami,
		ocfid:      ocfid,
		activation: activate,
	}
	updateData := BuildUpdateSmContextRequsetHandover(ue, UpdateSmContextPresentHandoverBetweenOCF, pduSessionId, param)
	return SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, nil)
}

func SendUpdateSmContextRequest(ue *ocf_context.OcfUe, smfUri, smContextRef string,
	updateData models.SmContextUpdateData, n1Msg []byte, n2Info []byte) (
	response *models.UpdateSmContextResponse, errorResponse *models.UpdateSmContextErrorResponse,
	problemDetail *models.ProblemDetails, err1 error) {
	configuration := Nsmf_PDUSession.NewConfiguration()
	configuration.SetBasePath(smfUri)
	client := Nsmf_PDUSession.NewAPIClient(configuration)

	var updateSmContextRequest models.UpdateSmContextRequest
	updateSmContextRequest.JsonData = &updateData
	updateSmContextRequest.BinaryDataN1SmMessage = n1Msg
	updateSmContextRequest.BinaryDataN2SmInformation = n2Info
	updateSmContextReponse, httpResponse, err :=
		client.IndividualSMContextApi.UpdateSmContext(context.Background(), smContextRef, updateSmContextRequest)
	if err == nil {
		response = &updateSmContextReponse
	} else if httpResponse != nil {
		if httpResponse.Status != err.Error() {
			err1 = err
			return
		}
		switch httpResponse.StatusCode {
		case 400, 403, 404, 500, 503:
			errResponse := err.(openapi.GenericOpenAPIError).Model().(models.UpdateSmContextErrorResponse)
			errorResponse = &errResponse
		case 411, 413, 415, 429:
			problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
			problemDetail = &problem
		}
	} else {
		err1 = openapi.ReportError("server no response")
	}
	return response, errorResponse, problemDetail, err1

}

func BuildUpdateSmContextRequset(
	ue *ocf_context.OcfUe, present UpdateSmContextPresent, pduSessionId int32, param updateSmContextRequsetParam) (
	updateData models.SmContextUpdateData) {
	smContext := ue.SmContextList[pduSessionId]
	context := ocf_context.OCF_Self()
	switch present {
	case UpdateSmContextPresentActivateUpCnxState:
		updateData.UpCnxState = models.UpCnxState_ACTIVATING
		if !ocf_context.CompareUserLocation(ue.Location, smContext.UserLocation) {
			updateData.UeLocation = &ue.Location
		}
		if param.accessType != "" && smContext.PduSessionContext.AccessType != param.accessType {
			updateData.AnType = param.accessType
		}
		if ladn, ok := context.LadnPool[smContext.PduSessionContext.Dnn]; ok {
			if ocf_context.InTaiList(ue.Tai, ladn.TaiLists) {
				updateData.PresenceInLadn = models.PresenceState_IN_AREA
			}
		}
	case UpdateSmContextPresentDeactivateUpCnxState:
		updateData.UpCnxState = models.UpCnxState_DEACTIVATED
		updateData.UeLocation = &ue.Location
		cause := param.cause
		if cause.Cause != nil {
			updateData.Cause = *cause.Cause
		}
		if cause.NgapCause != nil {
			updateData.NgApCause = cause.NgapCause
		}
		if cause.Var5GmmCause != nil {
			updateData.Var5gMmCauseValue = *cause.Var5GmmCause
		}
	case UpdateSmContextPresentChangeAccessType:
		updateData.AnTypeCanBeChanged = param.anTypeCanBeChanged
	case UpdateSmContextPresentOnlyN2SmInfo:
		updateData.N2SmInfoType = param.n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = "N2SmInfo"
		updateData.UeLocation = &ue.Location
	}
	return updateData
}

func BuildUpdateSmContextRequsetHandover(
	ue *ocf_context.OcfUe, present UpdateSmContextPresent, pduSessionId int32, param updateSmContextRequsetHandoverParam) (
	updateData models.SmContextUpdateData) {
	smContext := ue.SmContextList[pduSessionId]
	context := ocf_context.OCF_Self()
	if param.n2SmType != "" {
		updateData.N2SmInfoType = param.n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = "N2SmInfo"
	}
	switch present {
	case UpdateSmContextPresentXnHandover:
		updateData.ToBeSwitched = true
		updateData.UeLocation = &ue.Location
		if ladn, ok := context.LadnPool[smContext.PduSessionContext.Dnn]; ok {
			if ocf_context.InTaiList(ue.Tai, ladn.TaiLists) {
				updateData.PresenceInLadn = models.PresenceState_IN_AREA
			} else {
				updateData.PresenceInLadn = models.PresenceState_OUT_OF_AREA
			}
		}
	case UpdateSmContextPresentXnHandoverFailed:
		updateData.FailedToBeSwitched = true
	case UpdateSmContextPresentN2HandoverPreparing:
		updateData.HoState = models.HoState_PREPARING
		updateData.TargetId = param.targetId
		// ocf changed in same plmn
		if param.ocfid != "" {
			updateData.TargetServingNfId = param.ocfid
		}
	case UpdateSmContextPresentN2HandoverPrepared:
		updateData.HoState = models.HoState_PREPARED
	case UpdateSmContextPresentN2HandoverComplete:
		updateData.HoState = models.HoState_COMPLETED
		if param.ocfid != "" {
			updateData.ServingNfId = param.ocfid
			updateData.ServingNetwork = param.guami.PlmnId
			updateData.Guami = param.guami
		}
		if ladn, ok := context.LadnPool[smContext.PduSessionContext.Dnn]; ok {
			if ocf_context.InTaiList(ue.Tai, ladn.TaiLists) {
				updateData.PresenceInLadn = models.PresenceState_IN_AREA
			} else {
				updateData.PresenceInLadn = models.PresenceState_OUT_OF_AREA
			}
		}
	case UpdateSmContextPresentN2HandoverCanceled:
		updateData.HoState = models.HoState_CANCELLED
		cause := param.cause
		if cause.Cause != nil {
			updateData.Cause = *cause.Cause
		}
		if cause.NgapCause != nil {
			updateData.NgApCause = cause.NgapCause
		}
		if cause.Var5GmmCause != nil {
			updateData.Var5gMmCauseValue = *cause.Var5GmmCause
		}
	case UpdateSmContextPresentHandoverBetweenAccessType:
		updateData.AnType = param.accessType
		if param.n1SmMsg {
			updateData.N1SmMsg = new(models.RefToBinaryData)
			updateData.N1SmMsg.ContentId = "N1Msg"
		}
	case UpdateSmContextPresentHandoverBetweenOCF:
		updateData.ServingNfId = param.ocfid
		updateData.ServingNetwork = param.guami.PlmnId
		updateData.Guami = param.guami
		if param.activation {
			updateData.UpCnxState = models.UpCnxState_ACTIVATING
			if !ocf_context.CompareUserLocation(ue.Location, smContext.UserLocation) {
				updateData.UeLocation = &ue.Location
			}
			if param.accessType != "" && smContext.PduSessionContext.AccessType != param.accessType {
				updateData.AnType = param.accessType
			}
			if ladn, ok := context.LadnPool[smContext.PduSessionContext.Dnn]; ok {
				if ocf_context.InTaiList(ue.Tai, ladn.TaiLists) {
					updateData.PresenceInLadn = models.PresenceState_IN_AREA
				}
			}
		}
	}
	return updateData
}

// Release SmContext Request

func SendReleaseSmContextRequest(ue *ocf_context.OcfUe, pduSessionId int32,
	smContextReleaseData models.SmContextReleaseData) (detail *models.ProblemDetails, err error) {
	smContext, ok := ue.SmContextList[pduSessionId]
	if !ok {
		err = openapi.ReportError("[OCF] pduSessionId : %d is not in Ue", pduSessionId)
		return
	}

	configuration := Nsmf_PDUSession.NewConfiguration()
	configuration.SetBasePath(smContext.SmfUri)
	client := Nsmf_PDUSession.NewAPIClient(configuration)

	var releaseSmContextRequest models.ReleaseSmContextRequest
	releaseSmContextRequest.JsonData = &smContextReleaseData

	response, err1 := client.IndividualSMContextApi.ReleaseSmContext(
		context.Background(), smContext.PduSessionContext.SmContextRef, releaseSmContextRequest)
	if err1 == nil {
		delete(ue.SmContextList, pduSessionId)
	} else if response != nil && response.Status == err1.Error() {
		problem := err1.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		detail = &problem
	} else {
		err = err1
	}
	return
}
func BuildReleaseSmContextRequest(
	ue *ocf_context.OcfUe, cause *ocf_context.CauseAll, n2SmInfoType models.N2SmInfoType, n2Info []byte) (
	releaseData models.SmContextReleaseData) {
	if cause != nil {
		if cause.Cause != nil {
			releaseData.Cause = *cause.Cause
		}
		if cause.NgapCause != nil {
			releaseData.NgApCause = cause.NgapCause
		}
		if cause.Var5GmmCause != nil {
			releaseData.Var5gMmCauseValue = *cause.Var5GmmCause
		}
	}
	if ue.TimeZone != "" {
		releaseData.UeTimeZone = ue.TimeZone
	}
	if n2Info != nil {
		releaseData.N2SmInfoType = n2SmInfoType
		releaseData.N2SmInfo = &models.RefToBinaryData{
			ContentId: "n2SmInfo",
		}
	}
	// TODO: other param(ueLocation...)
	return
}

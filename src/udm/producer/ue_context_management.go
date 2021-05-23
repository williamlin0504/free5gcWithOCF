package producer

import (
	"context"
	"fmt"
	" free5gc/lib/http_wrapper"
	" free5gcenapi"
	" free5gcenapi/Nudr_DataRepository"
	Nudr " free5gcenapi/Nudr_DataRepository"
	" free5gcenapi/models"
	" free5gcm/consumer"
	udm_context " free5gcm/context"
	" free5gcm/factory"
	" free5gcm/logger"
	" free5gcm/producer/callback"
	"net/http"
	"strconv"
	"strings"

	"github.com/antihax/optional"
)

func createUDMClientToUDR(id string, nonUe bool) *Nudr_DataRepository.APIClient {
	var addr string
	if !nonUe {
		addr = getUdrUri(id)
	}
	if addr == "" {
		// dafault
		if !nonUe {
			logger.Handlelog.Warnf("Use default UDR Uri bacause ID[%s] does not match any UDR", id)
		}
		config := factory.UdmConfig
		udrclient := config.Configuration.Udrclient
		addr = fmt.Sprintf("%s://%s:%d", udrclient.Scheme, udrclient.Ipv4Addr, udrclient.Port)
	}
	cfg := Nudr.NewConfiguration()
	cfg.SetBasePath(addr)
	clientAPI := Nudr.NewAPIClient(cfg)
	return clientAPI
}

func getUdrUri(id string) string {
	// supi
	if strings.Contains(id, "imsi") || strings.Contains(id, "nai") {
		ue, ok := udm_context.UDM_Self().UdmUeFindBySupi(id)
		if ok {
			if ue.UdrUri == "" {
				ue.UdrUri = consumer.SendNFIntancesUDR(id, consumer.NFDiscoveryToUDRParamSupi)
			}
			return ue.UdrUri
		} else {
			ue = udm_context.UDM_Self().NewUdmUe(id)
			ue.UdrUri = consumer.SendNFIntancesUDR(id, consumer.NFDiscoveryToUDRParamSupi)
			return ue.UdrUri
		}
	} else if strings.Contains(id, "pei") {
		var udrUri string
		udm_context.UDM_Self().UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if ue.Amf3GppAccessRegistration != nil && ue.Amf3GppAccessRegistration.Pei == id {
				if ue.UdrUri == "" {
					ue.UdrUri = consumer.SendNFIntancesUDR(ue.Supi, consumer.NFDiscoveryToUDRParamSupi)
				}
				udrUri = ue.UdrUri
				return false
			} else if ue.AmfNon3GppAccessRegistration != nil && ue.AmfNon3GppAccessRegistration.Pei == id {
				if ue.UdrUri == "" {
					ue.UdrUri = consumer.SendNFIntancesUDR(ue.Supi, consumer.NFDiscoveryToUDRParamSupi)
				}
				udrUri = ue.UdrUri
				return false
			}
			return true
		})
		return udrUri
	} else if strings.Contains(id, "extgroupid") {
		// extra group id
		return consumer.SendNFIntancesUDR(id, consumer.NFDiscoveryToUDRParamExtGroupId)
	} else if strings.Contains(id, "msisdn") || strings.Contains(id, "extid") {
		// gpsi
		return consumer.SendNFIntancesUDR(id, consumer.NFDiscoveryToUDRParamGpsi)
	}
	return ""
}

func HandleGetAmf3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle HandleGetAmf3gppAccessRequest")

	// step 2: retrieve request
	ueID := request.Params["ueId"]
	supportedFeatures := request.Query.Get("supported-features")

	// step 3: handle the message
	response, problemDetails := GetAmf3gppAccessProcedure(ueID, supportedFeatures)

	// step 4: process the return value from step 3
	if response != nil {
		// status code is based on SPEC, and option headers
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func GetAmf3gppAccessProcedure(ueID string, supportedFeatures string) (
	response *models.Amf3GppAccessRegistration, problemDetails *models.ProblemDetails) {
	var queryAmfContext3gppParamOpts Nudr_DataRepository.QueryAmfContext3gppParamOpts
	queryAmfContext3gppParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)

	clientAPI := createUDMClientToUDR(ueID, false)
	amf3GppAccessRegistration, resp, err := clientAPI.AMF3GPPAccessRegistrationDocumentApi.
		QueryAmfContext3gpp(context.Background(), ueID, &queryAmfContext3gppParamOpts)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return nil, problemDetails
	}
	return &amf3GppAccessRegistration, nil
}

func HandleGetAmfNon3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infoln("Handle GetAmfNon3gppAccessRequest")

	// step 2: retrieve request
	ueId := request.Params["ueId"]
	supportedFeatures := request.Query.Get("supported-features")

	var queryAmfContextNon3gppParamOpts Nudr_DataRepository.QueryAmfContextNon3gppParamOpts
	queryAmfContextNon3gppParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)
	// step 3: handle the message
	response, problemDetails := GetAmfNon3gppAccessProcedure(queryAmfContextNon3gppParamOpts, ueId)

	// step 4: process the return value from step 3
	if response != nil {
		// status code is based on SPEC, and option headers
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)

}

func GetAmfNon3gppAccessProcedure(queryAmfContextNon3gppParamOpts Nudr_DataRepository.
	QueryAmfContextNon3gppParamOpts, ueID string) (response *models.AmfNon3GppAccessRegistration,
	problemDetails *models.ProblemDetails) {
	clientAPI := createUDMClientToUDR(ueID, false)
	amfNon3GppAccessRegistration, resp, err := clientAPI.AMFNon3GPPAccessRegistrationDocumentApi.
		QueryAmfContextNon3gpp(context.Background(), ueID, &queryAmfContextNon3gppParamOpts)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return nil, problemDetails
	}

	return &amfNon3GppAccessRegistration, nil
}

func HandleRegistrationAmf3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle RegistrationAmf3gppAccess")

	// step 2: retrieve request
	registerRequest := request.Body.(models.Amf3GppAccessRegistration)
	ueID := request.Params["ueId"]
	logger.UecmLog.Info("UEID: ", ueID)

	// step 3: handle the message
	header, response, problemDetails := RegistrationAmf3gppAccessProcedure(registerRequest, ueID)

	// step 4: process the return value from step 3
	if response != nil {
		// status code is based on SPEC, and option headers
		return http_wrapper.NewResponse(http.StatusCreated, header, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

// TS 29.503 5.3.2.2.2
func RegistrationAmf3gppAccessProcedure(registerRequest models.Amf3GppAccessRegistration, ueID string) (
	header http.Header, response *models.Amf3GppAccessRegistration, problemDetails *models.ProblemDetails) {
	// TODO: EPS interworking with N26 is not supported yet in this stage
	var oldAmf3GppAccessRegContext *models.Amf3GppAccessRegistration
	if udm_context.UDM_Self().UdmAmf3gppRegContextExists(ueID) {
		ue, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		oldAmf3GppAccessRegContext = ue.Amf3GppAccessRegistration
	}

	udm_context.UDM_Self().CreateAmf3gppRegContext(ueID, registerRequest)

	clientAPI := createUDMClientToUDR(ueID, false)
	var createAmfContext3gppParamOpts Nudr_DataRepository.CreateAmfContext3gppParamOpts
	optInterface := optional.NewInterface(registerRequest)
	createAmfContext3gppParamOpts.Amf3GppAccessRegistration = optInterface
	resp, err := clientAPI.AMF3GPPAccessRegistrationDocumentApi.CreateAmfContext3gpp(context.Background(),
		ueID, &createAmfContext3gppParamOpts)
	if err != nil {
		logger.UecmLog.Errorln("CreateAmfContext3gpp error : ", err)
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return nil, nil, problemDetails
	}

	// TS 23.502 4.2.2.2.2 14d: UDM initiate a Nudm_UECM_DeregistrationNotification to the old AMF
	// corresponding to the same (e.g. 3GPP) access, if one exists
	if oldAmf3GppAccessRegContext != nil {
		deregistData := models.DeregistrationData{
			DeregReason: models.DeregistrationReason_SUBSCRIPTION_WITHDRAWN,
			AccessType:  models.AccessType__3_GPP_ACCESS,
		}
		callback.SendOnDeregistrationNotification(ueID, oldAmf3GppAccessRegContext.DeregCallbackUri,
			deregistData) // Deregistration Notify Triggered

		return nil, nil, nil
	} else {
		header = make(http.Header)
		udmUe, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		header.Set("Location", udmUe.GetLocationURI(udm_context.LocationUriAmf3GppAccessRegistration))
		return header, &registerRequest, nil
	}
}

// TS 29.503 5.3.2.2.3
func HandleRegisterAmfNon3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle RegisterAmfNon3gppAccessRequest")

	// step 2: retrieve request
	registerRequest := request.Body.(models.AmfNon3GppAccessRegistration)
	ueID := request.Params["ueId"]

	// step 3: handle the message
	header, response, problemDetails := RegisterAmfNon3gppAccessProcedure(registerRequest, ueID)

	// step 4: process the return value from step 3
	if response != nil {
		// status code is based on SPEC, and option headers
		return http_wrapper.NewResponse(http.StatusCreated, header, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func RegisterAmfNon3gppAccessProcedure(registerRequest models.AmfNon3GppAccessRegistration, ueID string) (
	header http.Header, response *models.AmfNon3GppAccessRegistration, problemDetails *models.ProblemDetails) {
	var oldAmfNon3GppAccessRegContext *models.AmfNon3GppAccessRegistration
	if udm_context.UDM_Self().UdmAmfNon3gppRegContextExists(ueID) {
		ue, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		oldAmfNon3GppAccessRegContext = ue.AmfNon3GppAccessRegistration
	}

	udm_context.UDM_Self().CreateAmfNon3gppRegContext(ueID, registerRequest)

	clientAPI := createUDMClientToUDR(ueID, false)
	var createAmfContextNon3gppParamOpts Nudr_DataRepository.CreateAmfContextNon3gppParamOpts
	optInterface := optional.NewInterface(registerRequest)
	createAmfContextNon3gppParamOpts.AmfNon3GppAccessRegistration = optInterface
	resp, err := clientAPI.AMFNon3GPPAccessRegistrationDocumentApi.CreateAmfContextNon3gpp(
		context.Background(), ueID, &createAmfContextNon3gppParamOpts)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return nil, nil, problemDetails
	}

	// TS 23.502 4.2.2.2.2 14d: UDM initiate a Nudm_UECM_DeregistrationNotification to the old AMF
	// corresponding to the same (e.g. 3GPP) access, if one exists
	if oldAmfNon3GppAccessRegContext != nil {
		deregistData := models.DeregistrationData{
			DeregReason: models.DeregistrationReason_SUBSCRIPTION_WITHDRAWN,
			AccessType:  models.AccessType_NON_3_GPP_ACCESS,
		}
		callback.SendOnDeregistrationNotification(ueID, oldAmfNon3GppAccessRegContext.DeregCallbackUri,
			deregistData) // Deregistration Notify Triggered

		return nil, nil, nil
	} else {
		header = make(http.Header)
		udmUe, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		header.Set("Location", udmUe.GetLocationURI(udm_context.LocationUriAmfNon3GppAccessRegistration))
		return header, &registerRequest, nil
	}
}

// TODO: ueID may be SUPI or GPSI, but this function did not handle this condition
func HandleUpdateAmf3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle UpdateAmf3gppAccessRequest")

	// step 2: retrieve request
	amf3GppAccessRegistrationModification := request.Body.(models.Amf3GppAccessRegistrationModification)
	ueID := request.Params["ueId"]

	// step 3: handle the message
	problemDetails := UpdateAmf3gppAccessProcedure(amf3GppAccessRegistrationModification, ueID)

	// step 4: process the return value from step 3
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func UpdateAmf3gppAccessProcedure(request models.Amf3GppAccessRegistrationModification, ueID string) (
	problemDetails *models.ProblemDetails) {
	var patchItemReqArray []models.PatchItem
	currentContext := udm_context.UDM_Self().GetAmf3gppRegContext(ueID)
	if currentContext == nil {
		logger.UecmLog.Errorln("[UpdateAmf3gppAccess] Empty Amf3gppRegContext")
		problemDetails = &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return problemDetails
	}

	if request.Guami != nil {
		udmUe, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		if udmUe.SameAsStoredGUAMI3gpp(*request.Guami) { // deregistration
			logger.UecmLog.Infoln("UpdateAmf3gppAccess - deregistration")
			request.PurgeFlag = true
		} else {
			logger.UecmLog.Errorln("INVALID_GUAMI")
			problemDetails = &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "INVALID_GUAMI",
			}
			return problemDetails
		}

		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "Guami"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = *request.Guami
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.PurgeFlag {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "PurgeFlag"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.PurgeFlag
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.Pei != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "Pei"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.Pei
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.ImsVoPs != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "ImsVoPs"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.ImsVoPs
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.BackupAmfInfo != nil {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "BackupAmfInfo"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.BackupAmfInfo
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.AMF3GPPAccessRegistrationDocumentApi.AmfContext3gpp(context.Background(), ueID,
		patchItemReqArray)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}

		return problemDetails
	}

	return nil
}

// TODO: ueID may be SUPI or GPSI, but this function did not handle this condition
func HandleUpdateAmfNon3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle UpdateAmfNon3gppAccessRequest")

	// step 2: retrieve request
	requestMSG := request.Body.(models.AmfNon3GppAccessRegistrationModification)
	ueID := request.Params["ueId"]

	// step 3: handle the message
	problemDetails := UpdateAmfNon3gppAccessProcedure(requestMSG, ueID)

	// step 4: process the return value from step 3
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func UpdateAmfNon3gppAccessProcedure(request models.AmfNon3GppAccessRegistrationModification, ueID string) (
	problemDetails *models.ProblemDetails) {
	var patchItemReqArray []models.PatchItem
	currentContext := udm_context.UDM_Self().GetAmfNon3gppRegContext(ueID)
	if currentContext == nil {
		logger.UecmLog.Errorln("[UpdateAmfNon3gppAccess] Empty AmfNon3gppRegContext")
		problemDetails = &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return problemDetails
	}

	if request.Guami != nil {
		udmUe, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		if udmUe.SameAsStoredGUAMINon3gpp(*request.Guami) { // deregistration
			logger.UecmLog.Infoln("UpdateAmfNon3gppAccess - deregistration")
			request.PurgeFlag = true
		} else {
			logger.UecmLog.Errorln("INVALID_GUAMI")
			problemDetails = &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "INVALID_GUAMI",
			}
		}

		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "Guami"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = *request.Guami
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.PurgeFlag {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "PurgeFlag"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.PurgeFlag
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.Pei != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "Pei"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.Pei
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.ImsVoPs != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "ImsVoPs"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.ImsVoPs
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.BackupAmfInfo != nil {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "BackupAmfInfo"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.BackupAmfInfo
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.AMFNon3GPPAccessRegistrationDocumentApi.AmfContextNon3gpp(context.Background(),
		ueID, patchItemReqArray)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return problemDetails
	}
	return nil
}

func HandleDeregistrationSmfRegistrations(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle DeregistrationSmfRegistrations")

	// step 2: retrieve request
	ueID := request.Params["ueId"]
	pduSessionID := request.Params["pduSessionId"]

	// step 3: handle the message
	problemDetails := DeregistrationSmfRegistrationsProcedure(ueID, pduSessionID)

	// step 4: process the return value from step 3
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func DeregistrationSmfRegistrationsProcedure(ueID string, pduSessionID string) (problemDetails *models.ProblemDetails) {
	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.SMFRegistrationDocumentApi.DeleteSmfContext(context.Background(), ueID, pduSessionID)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return problemDetails
	}
	return nil
}

// SmfRegistrations
func HandleRegistrationSmfRegistrationsRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle RegistrationSmfRegistrations")

	// step 2: retrieve request
	registerRequest := request.Body.(models.SmfRegistration)
	ueID := request.Params["ueId"]
	pduSessionID := request.Params["pduSessionId"]

	// step 3: handle the message
	header, response, problemDetails := RegistrationSmfRegistrationsProcedure(&registerRequest, ueID, pduSessionID)

	// step 4: process the return value from step 3
	if response != nil {
		// status code is based on SPEC, and option headers
		return http_wrapper.NewResponse(http.StatusCreated, header, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		//all nil
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

// SmfRegistrationsProcedure
func RegistrationSmfRegistrationsProcedure(request *models.SmfRegistration, ueID string, pduSessionID string) (
	header http.Header, response *models.SmfRegistration, problemDetails *models.ProblemDetails) {
	contextExisted := false
	udm_context.UDM_Self().CreateSmfRegContext(ueID, pduSessionID)
	if !udm_context.UDM_Self().UdmSmfRegContextNotExists(ueID) {
		contextExisted = true
	}

	pduID64, err := strconv.ParseInt(pduSessionID, 10, 32)
	if err != nil {
		logger.UecmLog.Errorln(err.Error())
	}
	pduID32 := int32(pduID64)

	var createSmfContextNon3gppParamOpts Nudr_DataRepository.CreateSmfContextNon3gppParamOpts
	optInterface := optional.NewInterface(request)
	createSmfContextNon3gppParamOpts.SmfRegistration = optInterface

	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.SMFRegistrationDocumentApi.CreateSmfContextNon3gpp(context.Background(), ueID,
		pduID32, &createSmfContextNon3gppParamOpts)
	if err != nil {
		problemDetails.Cause = err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return nil, nil, problemDetails
	}

	if contextExisted {
		return nil, nil, nil
	} else {
		header = make(http.Header)
		udmUe, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		header.Set("Location", udmUe.GetLocationURI(udm_context.LocationUriSmfRegistration))
		return header, request, nil
	}
}

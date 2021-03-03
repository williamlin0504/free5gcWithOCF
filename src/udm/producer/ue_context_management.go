package producer

import (
	"context"
	"fmt"
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/Nudr_DataRepository"
	Nudr "free5gc/lib/openapi/Nudr_DataRepository"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/consumer"
	udm_context "free5gc/src/udm/context"
	"free5gc/src/udm/factory"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/producer/callback"
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
			if ue.Ocf3GppAccessRegistration != nil && ue.Ocf3GppAccessRegistration.Pei == id {
				if ue.UdrUri == "" {
					ue.UdrUri = consumer.SendNFIntancesUDR(ue.Supi, consumer.NFDiscoveryToUDRParamSupi)
				}
				udrUri = ue.UdrUri
				return false
			} else if ue.OcfNon3GppAccessRegistration != nil && ue.OcfNon3GppAccessRegistration.Pei == id {
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

func HandleGetOcf3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle HandleGetOcf3gppAccessRequest")

	// step 2: retrieve request
	ueID := request.Params["ueId"]
	supportedFeatures := request.Query.Get("supported-features")

	// step 3: handle the message
	response, problemDetails := GetOcf3gppAccessProcedure(ueID, supportedFeatures)

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

func GetOcf3gppAccessProcedure(ueID string, supportedFeatures string) (
	response *models.Ocf3GppAccessRegistration, problemDetails *models.ProblemDetails) {
	var queryOcfContext3gppParamOpts Nudr_DataRepository.QueryOcfContext3gppParamOpts
	queryOcfContext3gppParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)

	clientAPI := createUDMClientToUDR(ueID, false)
	ocf3GppAccessRegistration, resp, err := clientAPI.OCF3GPPAccessRegistrationDocumentApi.
		QueryOcfContext3gpp(context.Background(), ueID, &queryOcfContext3gppParamOpts)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return nil, problemDetails
	}
	return &ocf3GppAccessRegistration, nil
}

func HandleGetOcfNon3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infoln("Handle GetOcfNon3gppAccessRequest")

	// step 2: retrieve request
	ueId := request.Params["ueId"]
	supportedFeatures := request.Query.Get("supported-features")

	var queryOcfContextNon3gppParamOpts Nudr_DataRepository.QueryOcfContextNon3gppParamOpts
	queryOcfContextNon3gppParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)
	// step 3: handle the message
	response, problemDetails := GetOcfNon3gppAccessProcedure(queryOcfContextNon3gppParamOpts, ueId)

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

func GetOcfNon3gppAccessProcedure(queryOcfContextNon3gppParamOpts Nudr_DataRepository.
	QueryOcfContextNon3gppParamOpts, ueID string) (response *models.OcfNon3GppAccessRegistration,
	problemDetails *models.ProblemDetails) {
	clientAPI := createUDMClientToUDR(ueID, false)
	ocfNon3GppAccessRegistration, resp, err := clientAPI.OCFNon3GPPAccessRegistrationDocumentApi.
		QueryOcfContextNon3gpp(context.Background(), ueID, &queryOcfContextNon3gppParamOpts)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return nil, problemDetails
	}

	return &ocfNon3GppAccessRegistration, nil
}

func HandleRegistrationOcf3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle RegistrationOcf3gppAccess")

	// step 2: retrieve request
	registerRequest := request.Body.(models.Ocf3GppAccessRegistration)
	ueID := request.Params["ueId"]
	logger.UecmLog.Info("UEID: ", ueID)

	// step 3: handle the message
	header, response, problemDetails := RegistrationOcf3gppAccessProcedure(registerRequest, ueID)

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
func RegistrationOcf3gppAccessProcedure(registerRequest models.Ocf3GppAccessRegistration, ueID string) (
	header http.Header, response *models.Ocf3GppAccessRegistration, problemDetails *models.ProblemDetails) {
	// TODO: EPS interworking with N26 is not supported yet in this stage
	var oldOcf3GppAccessRegContext *models.Ocf3GppAccessRegistration
	if udm_context.UDM_Self().UdmOcf3gppRegContextExists(ueID) {
		ue, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		oldOcf3GppAccessRegContext = ue.Ocf3GppAccessRegistration
	}

	udm_context.UDM_Self().CreateOcf3gppRegContext(ueID, registerRequest)

	clientAPI := createUDMClientToUDR(ueID, false)
	var createOcfContext3gppParamOpts Nudr_DataRepository.CreateOcfContext3gppParamOpts
	optInterface := optional.NewInterface(registerRequest)
	createOcfContext3gppParamOpts.Ocf3GppAccessRegistration = optInterface
	resp, err := clientAPI.OCF3GPPAccessRegistrationDocumentApi.CreateOcfContext3gpp(context.Background(),
		ueID, &createOcfContext3gppParamOpts)
	if err != nil {
		logger.UecmLog.Errorln("CreateOcfContext3gpp error : ", err)
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return nil, nil, problemDetails
	}

	// TS 23.502 4.2.2.2.2 14d: UDM initiate a Nudm_UECM_DeregistrationNotification to the old OCF
	// corresponding to the same (e.g. 3GPP) access, if one exists
	if oldOcf3GppAccessRegContext != nil {
		deregistData := models.DeregistrationData{
			DeregReason: models.DeregistrationReason_SUBSCRIPTION_WITHDRAWN,
			AccessType:  models.AccessType__3_GPP_ACCESS,
		}
		callback.SendOnDeregistrationNotification(ueID, oldOcf3GppAccessRegContext.DeregCallbackUri,
			deregistData) // Deregistration Notify Triggered

		return nil, nil, nil
	} else {
		header = make(http.Header)
		udmUe, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		header.Set("Location", udmUe.GetLocationURI(udm_context.LocationUriOcf3GppAccessRegistration))
		return header, &registerRequest, nil
	}
}

// TS 29.503 5.3.2.2.3
func HandleRegisterOcfNon3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle RegisterOcfNon3gppAccessRequest")

	// step 2: retrieve request
	registerRequest := request.Body.(models.OcfNon3GppAccessRegistration)
	ueID := request.Params["ueId"]

	// step 3: handle the message
	header, response, problemDetails := RegisterOcfNon3gppAccessProcedure(registerRequest, ueID)

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

func RegisterOcfNon3gppAccessProcedure(registerRequest models.OcfNon3GppAccessRegistration, ueID string) (
	header http.Header, response *models.OcfNon3GppAccessRegistration, problemDetails *models.ProblemDetails) {
	var oldOcfNon3GppAccessRegContext *models.OcfNon3GppAccessRegistration
	if udm_context.UDM_Self().UdmOcfNon3gppRegContextExists(ueID) {
		ue, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		oldOcfNon3GppAccessRegContext = ue.OcfNon3GppAccessRegistration
	}

	udm_context.UDM_Self().CreateOcfNon3gppRegContext(ueID, registerRequest)

	clientAPI := createUDMClientToUDR(ueID, false)
	var createOcfContextNon3gppParamOpts Nudr_DataRepository.CreateOcfContextNon3gppParamOpts
	optInterface := optional.NewInterface(registerRequest)
	createOcfContextNon3gppParamOpts.OcfNon3GppAccessRegistration = optInterface
	resp, err := clientAPI.OCFNon3GPPAccessRegistrationDocumentApi.CreateOcfContextNon3gpp(
		context.Background(), ueID, &createOcfContextNon3gppParamOpts)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: int32(resp.StatusCode),
			Cause:  err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails).Cause,
			Detail: err.Error(),
		}
		return nil, nil, problemDetails
	}

	// TS 23.502 4.2.2.2.2 14d: UDM initiate a Nudm_UECM_DeregistrationNotification to the old OCF
	// corresponding to the same (e.g. 3GPP) access, if one exists
	if oldOcfNon3GppAccessRegContext != nil {
		deregistData := models.DeregistrationData{
			DeregReason: models.DeregistrationReason_SUBSCRIPTION_WITHDRAWN,
			AccessType:  models.AccessType_NON_3_GPP_ACCESS,
		}
		callback.SendOnDeregistrationNotification(ueID, oldOcfNon3GppAccessRegContext.DeregCallbackUri,
			deregistData) // Deregistration Notify Triggered

		return nil, nil, nil
	} else {
		header = make(http.Header)
		udmUe, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		header.Set("Location", udmUe.GetLocationURI(udm_context.LocationUriOcfNon3GppAccessRegistration))
		return header, &registerRequest, nil
	}
}

// TODO: ueID may be SUPI or GPSI, but this function did not handle this condition
func HandleUpdateOcf3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle UpdateOcf3gppAccessRequest")

	// step 2: retrieve request
	ocf3GppAccessRegistrationModification := request.Body.(models.Ocf3GppAccessRegistrationModification)
	ueID := request.Params["ueId"]

	// step 3: handle the message
	problemDetails := UpdateOcf3gppAccessProcedure(ocf3GppAccessRegistrationModification, ueID)

	// step 4: process the return value from step 3
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func UpdateOcf3gppAccessProcedure(request models.Ocf3GppAccessRegistrationModification, ueID string) (
	problemDetails *models.ProblemDetails) {
	var patchItemReqArray []models.PatchItem
	currentContext := udm_context.UDM_Self().GetOcf3gppRegContext(ueID)
	if currentContext == nil {
		logger.UecmLog.Errorln("[UpdateOcf3gppAccess] Empty Ocf3gppRegContext")
		problemDetails = &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return problemDetails
	}

	if request.Guami != nil {
		udmUe, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		if udmUe.SameAsStoredGUAMI3gpp(*request.Guami) { // deregistration
			logger.UecmLog.Infoln("UpdateOcf3gppAccess - deregistration")
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

	if request.BackupOcfInfo != nil {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "BackupOcfInfo"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.BackupOcfInfo
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.OCF3GPPAccessRegistrationDocumentApi.OcfContext3gpp(context.Background(), ueID,
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
func HandleUpdateOcfNon3gppAccessRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// step 1: log
	logger.UecmLog.Infof("Handle UpdateOcfNon3gppAccessRequest")

	// step 2: retrieve request
	requestMSG := request.Body.(models.OcfNon3GppAccessRegistrationModification)
	ueID := request.Params["ueId"]

	// step 3: handle the message
	problemDetails := UpdateOcfNon3gppAccessProcedure(requestMSG, ueID)

	// step 4: process the return value from step 3
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func UpdateOcfNon3gppAccessProcedure(request models.OcfNon3GppAccessRegistrationModification, ueID string) (
	problemDetails *models.ProblemDetails) {
	var patchItemReqArray []models.PatchItem
	currentContext := udm_context.UDM_Self().GetOcfNon3gppRegContext(ueID)
	if currentContext == nil {
		logger.UecmLog.Errorln("[UpdateOcfNon3gppAccess] Empty OcfNon3gppRegContext")
		problemDetails = &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return problemDetails
	}

	if request.Guami != nil {
		udmUe, _ := udm_context.UDM_Self().UdmUeFindBySupi(ueID)
		if udmUe.SameAsStoredGUAMINon3gpp(*request.Guami) { // deregistration
			logger.UecmLog.Infoln("UpdateOcfNon3gppAccess - deregistration")
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

	if request.BackupOcfInfo != nil {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "BackupOcfInfo"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.BackupOcfInfo
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.OCFNon3GPPAccessRegistrationDocumentApi.OcfContextNon3gpp(context.Background(),
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

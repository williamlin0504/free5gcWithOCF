package producer

import (
	"context"
	"fmt"
	"free5gcWithOCF/lib/http_wrapper"
	"free5gcWithOCF/lib/openapi"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/src/chf/consumer"
	chf_context "free5gcWithOCF/src/chf/context"
	"free5gcWithOCF/src/chf/logger"
	"free5gcWithOCF/src/chf/util"
	"net/http"
	"reflect"

	"github.com/mohae/deepcopy"
)

func HandleDeletePoliciesPolAssoId(request *http_wrapper.Request) *http_wrapper.Response {
	logger.AMpolicylog.Infof("Handle Policy Association Delete")

	polAssoId := request.Params["polAssoId"]

	problemDetails := DeletePoliciesPolAssoIdProcedure(polAssoId)
	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	} else {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func DeletePoliciesPolAssoIdProcedure(polAssoId string) *models.ProblemDetails {
	ue := chf_context.CHF_Self().CHFUeFindByPolicyId(polAssoId)
	if ue == nil || ue.AMPolicyData[polAssoId] == nil {
		problemDetails := util.GetProblemDetail("polAssoId not found  in CHF", util.CONTEXT_NOT_FOUND)
		return &problemDetails
	}
	delete(ue.AMPolicyData, polAssoId)
	return nil
}

// PoliciesPolAssoIdGet -
func HandleGetPoliciesPolAssoId(request *http_wrapper.Request) *http_wrapper.Response {
	logger.AMpolicylog.Infof("Handle Policy Association Get")

	polAssoId := request.Params["polAssoId"]

	response, problemDetails := GetPoliciesPolAssoIdProcedure(polAssoId)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func GetPoliciesPolAssoIdProcedure(polAssoId string) (*models.PolicyAssociation, *models.ProblemDetails) {
	ue := chf_context.CHF_Self().CHFUeFindByPolicyId(polAssoId)
	if ue == nil || ue.AMPolicyData[polAssoId] == nil {
		problemDetails := util.GetProblemDetail("polAssoId not found  in CHF", util.CONTEXT_NOT_FOUND)
		return nil, &problemDetails
	}
	amPolicyData := ue.AMPolicyData[polAssoId]
	rsp := models.PolicyAssociation{
		SuppFeat: amPolicyData.SuppFeat,
	}
	if amPolicyData.Rfsp != 0 {
		rsp.Rfsp = amPolicyData.Rfsp
	}
	if amPolicyData.ServAreaRes != nil {
		rsp.ServAreaRes = amPolicyData.ServAreaRes
	}
	if amPolicyData.Triggers != nil {
		rsp.Triggers = amPolicyData.Triggers
		for _, trigger := range amPolicyData.Triggers {
			if trigger == models.RequestTrigger_PRA_CH {
				rsp.Pras = amPolicyData.Pras
				break
			}
		}
	}
	return &rsp, nil
}

func HandleUpdatePostPoliciesPolAssoId(request *http_wrapper.Request) *http_wrapper.Response {
	logger.AMpolicylog.Infof("Handle Policy Association Update")

	polAssoId := request.Params["polAssoId"]
	policyAssociationUpdateRequest := request.Body.(models.PolicyAssociationUpdateRequest)

	response, problemDetails := UpdatePostPoliciesPolAssoIdProcedure(polAssoId, policyAssociationUpdateRequest)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func UpdatePostPoliciesPolAssoIdProcedure(polAssoId string,
	policyAssociationUpdateRequest models.PolicyAssociationUpdateRequest) (*models.PolicyUpdate, *models.ProblemDetails) {
	ue := chf_context.CHF_Self().CHFUeFindByPolicyId(polAssoId)
	if ue == nil || ue.AMPolicyData[polAssoId] == nil {
		problemDetails := util.GetProblemDetail("polAssoId not found  in CHF", util.CONTEXT_NOT_FOUND)
		return nil, &problemDetails
	}

	amPolicyData := ue.AMPolicyData[polAssoId]
	var response models.PolicyUpdate
	if policyAssociationUpdateRequest.NotificationUri != "" {
		amPolicyData.NotificationUri = policyAssociationUpdateRequest.NotificationUri
	}
	if policyAssociationUpdateRequest.AltNotifIpv4Addrs != nil {
		amPolicyData.AltNotifIpv4Addrs = policyAssociationUpdateRequest.AltNotifIpv4Addrs
	}
	if policyAssociationUpdateRequest.AltNotifIpv6Addrs != nil {
		amPolicyData.AltNotifIpv6Addrs = policyAssociationUpdateRequest.AltNotifIpv6Addrs
	}
	for _, trigger := range policyAssociationUpdateRequest.Triggers {
		//TODO: Modify the value according to policies
		switch trigger {
		case models.RequestTrigger_LOC_CH:
			//TODO: report to AF subscriber
			if policyAssociationUpdateRequest.UserLoc == nil {
				problemDetail := util.GetProblemDetail("UserLoc are nli", util.ERROR_REQUEST_PARAMETERS)
				logger.AMpolicylog.Warnln(
					"UserLoc doesn't exist in Policy Association Requset Update while Triggers include LOC_CH")
				return nil, &problemDetail
			}
			amPolicyData.UserLoc = policyAssociationUpdateRequest.UserLoc
			logger.AMpolicylog.Infof("Ue[%s] UserLocation %+v", ue.Supi, amPolicyData.UserLoc)
		case models.RequestTrigger_PRA_CH:
			if policyAssociationUpdateRequest.PraStatuses == nil {
				problemDetail := util.GetProblemDetail("PraStatuses are nli", util.ERROR_REQUEST_PARAMETERS)
				logger.AMpolicylog.Warnln("PraStatuses doesn't exist in Policy Association",
					"Requset Update while Triggers include PRA_CH")
				return nil, &problemDetail
			}
			for praId, praInfo := range policyAssociationUpdateRequest.PraStatuses {
				//TODO: report to AF subscriber
				logger.AMpolicylog.Infof("Policy Association Presence Id[%s] change state to %s", praId, praInfo.PresenceState)
			}
		case models.RequestTrigger_SERV_AREA_CH:
			if policyAssociationUpdateRequest.ServAreaRes == nil {
				problemDetail := util.GetProblemDetail("ServAreaRes are nli", util.ERROR_REQUEST_PARAMETERS)
				logger.AMpolicylog.Warnln("ServAreaRes doesn't exist in Policy Association",
					"Requset Update while Triggers include SERV_AREA_CH")
				return nil, &problemDetail
			} else {
				amPolicyData.ServAreaRes = policyAssociationUpdateRequest.ServAreaRes
				response.ServAreaRes = policyAssociationUpdateRequest.ServAreaRes
			}
		case models.RequestTrigger_RFSP_CH:
			if policyAssociationUpdateRequest.Rfsp == 0 {
				problemDetail := util.GetProblemDetail("Rfsp are nli", util.ERROR_REQUEST_PARAMETERS)
				logger.AMpolicylog.Warnln("Rfsp doesn't exist in Policy Association Requset Update while Triggers include RFSP_CH")
				return nil, &problemDetail
			} else {
				amPolicyData.Rfsp = policyAssociationUpdateRequest.Rfsp
				response.Rfsp = policyAssociationUpdateRequest.Rfsp
			}
		}
	}
	//TODO: handle TraceReq
	//TODO: Change Request Trigger Policies if needed
	response.Triggers = amPolicyData.Triggers
	//TODO: Change Policies if needed
	// rsp.Pras
	return &response, nil
}

// Create AM Policy
func HandlePostPolicies(request *http_wrapper.Request) *http_wrapper.Response {
	logger.AMpolicylog.Infof("Handle Policy Association Request")

	polAssoId := request.Params["polAssoId"]
	policyAssociationRequest := request.Body.(models.PolicyAssociationRequest)

	response, locationHeader, problemDetails := PostPoliciesProcedure(polAssoId, policyAssociationRequest)
	headers := http.Header{
		"Location": {locationHeader},
	}
	if response != nil {
		return http_wrapper.NewResponse(http.StatusCreated, headers, response)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PostPoliciesProcedure(polAssoId string,
	policyAssociationRequest models.PolicyAssociationRequest) (*models.PolicyAssociation, string, *models.ProblemDetails) {
	var response models.PolicyAssociation
	chfSelf := chf_context.CHF_Self()
	var ue *chf_context.UeContext
	if val, ok := chfSelf.UePool.Load(policyAssociationRequest.Supi); ok {
		ue = val.(*chf_context.UeContext)
	}
	if ue == nil {
		if newUe, err := chfSelf.NewCHFUe(policyAssociationRequest.Supi); err != nil {
			// supi format dose not match "imsi-..."
			problemDetail := util.GetProblemDetail("Supi Format Error", util.ERROR_REQUEST_PARAMETERS)
			logger.AMpolicylog.Errorln(err.Error())
			return nil, "", &problemDetail
		} else {
			ue = newUe
		}
	}
	udrUri := getUdrUri(ue)
	if udrUri == "" {
		// Can't find any UDR support this Ue
		chfSelf.UePool.Delete(ue.Supi)
		problemDetail := util.GetProblemDetail("Ue is not supported in CHF", util.USER_UNKNOWN)
		logger.AMpolicylog.Errorf("Ue[%s] is not supported in CHF", ue.Supi)
		return nil, "", &problemDetail
	}
	ue.UdrUri = udrUri

	response.Request = deepcopy.Copy(&policyAssociationRequest).(*models.PolicyAssociationRequest)
	assolId := fmt.Sprintf("%s-%d", ue.Supi, ue.PolAssociationIDGenerator)
	amPolicy := ue.AMPolicyData[assolId]

	if amPolicy == nil || amPolicy.AmPolicyData == nil {
		client := util.GetNudrClient(udrUri)
		var response *http.Response
		amData, response, err := client.DefaultApi.PolicyDataUesUeIdAmDataGet(context.Background(), ue.Supi)
		if err != nil || response == nil || response.StatusCode != http.StatusOK {
			problemDetail := util.GetProblemDetail("Can't find UE AM Policy Data in UDR", util.USER_UNKNOWN)
			logger.AMpolicylog.Errorf("Can't find UE[%s] AM Policy Data in UDR", ue.Supi)
			return nil, "", &problemDetail
		}
		if amPolicy == nil {
			amPolicy = ue.NewUeAMPolicyData(assolId, policyAssociationRequest)
		}
		amPolicy.AmPolicyData = &amData
	}

	// TODO: according to CHF Policy to determine ServAreaRes, Rfsp, SuppFeat
	// amPolicy.ServAreaRes =
	// amPolicy.Rfsp =
	var requestSuppFeat openapi.SupportedFeature
	if suppFeat, err := openapi.NewSupportedFeature(policyAssociationRequest.SuppFeat); err != nil {
		logger.AMpolicylog.Warnln(err)
	} else {
		requestSuppFeat = suppFeat
	}
	amPolicy.SuppFeat = chfSelf.ChfSuppFeats[models.
		ServiceName_NCHF_AM_POLICY_CONTROL].NegotiateWith(
		requestSuppFeat).String()
	if amPolicy.Rfsp != 0 {
		response.Rfsp = amPolicy.Rfsp
	}
	response.SuppFeat = amPolicy.SuppFeat
	// TODO: add Reports
	// rsp.Triggers
	// rsp.Pras
	ue.PolAssociationIDGenerator++
	// Create location header for update, delete, get
	locationHeader := util.GetResourceUri(models.ServiceName_NCHF_AM_POLICY_CONTROL, assolId)
	logger.AMpolicylog.Tracef("AMPolicy association Id[%s] Create", assolId)

	if policyAssociationRequest.Guami != nil {
		// if consumer is AMF then subscribe this AMF Status
		for _, statusSubsData := range chfSelf.AMFStatusSubsData {
			for _, guami := range statusSubsData.GuamiList {
				if reflect.DeepEqual(guami, policyAssociationRequest.Guami) {
					amPolicy.AmfStatusChangeSubscription = &statusSubsData
				}
			}
		}
	}
	return &response, locationHeader, nil
}

// Send AM Policy Update to AMF if policy has changed
func SendAMPolicyUpdateNotification(ue *chf_context.UeContext, PolId string, request models.PolicyUpdate) {
	if ue == nil {
		logger.AMpolicylog.Warnln("Policy Update Notification Error[Ue is nil]")
		return
	}
	amPolicyData := ue.AMPolicyData[PolId]
	if amPolicyData == nil {
		logger.AMpolicylog.Warnf("Policy Update Notification Error[Can't find polAssoId[%s] in UE(%s)]", PolId, ue.Supi)
		return
	}
	client := util.GetNchfAMPolicyCallbackClient()
	uri := amPolicyData.NotificationUri
	for uri != "" {

		rsp, err := client.DefaultCallbackApi.PolicyUpdateNotification(context.Background(), uri, request)
		if err != nil {
			if rsp != nil && rsp.StatusCode != http.StatusNoContent {
				logger.AMpolicylog.Warnf("Policy Update Notification Error[%s]", rsp.Status)
			} else {
				logger.AMpolicylog.Warnf("Policy Update Notification Failed[%s]", err.Error())
			}
			return
		} else if rsp == nil {
			logger.AMpolicylog.Warnln("Policy Update Notification Failed[HTTP Response is nil]")
			return
		}
		if rsp.StatusCode == http.StatusTemporaryRedirect {
			// for redirect case, resend the notification to redirect target
			uRI, err := rsp.Location()
			if err != nil {
				logger.AMpolicylog.Warnln("Policy Update Notification Redirect Need Supply URI")
				return
			}
			uri = uRI.String()
			continue
		}

		logger.AMpolicylog.Infoln("Policy Update Notification Success")
		return
	}

}

// Send AM Policy Update to AMF if policy has been terminated
func SendAMPolicyTerminationRequestNotification(ue *chf_context.UeContext,
	PolId string, request models.TerminationNotification) {
	if ue == nil {
		logger.AMpolicylog.Warnln("Policy Assocition Termination Request Notification Error[Ue is nil]")
		return
	}
	amPolicyData := ue.AMPolicyData[PolId]
	if amPolicyData == nil {
		logger.AMpolicylog.Warnf(
			"Policy Assocition Termination Request Notification Error[Can't find polAssoId[%s] in UE(%s)]", PolId, ue.Supi)
		return
	}
	client := util.GetNchfAMPolicyCallbackClient()
	uri := amPolicyData.NotificationUri
	for uri != "" {

		rsp, err := client.DefaultCallbackApi.PolicyAssocitionTerminationRequestNotification(
			context.Background(), uri, request)
		if err != nil {
			if rsp != nil && rsp.StatusCode != http.StatusNoContent {
				logger.AMpolicylog.Warnf("Policy Assocition Termination Request Notification Error[%s]", rsp.Status)
			} else {
				logger.AMpolicylog.Warnf("Policy Assocition Termination Request Notification Failed[%s]", err.Error())
			}
			return
		} else if rsp == nil {
			logger.AMpolicylog.Warnln("Policy Assocition Termination Request Notification Failed[HTTP Response is nil]")
			return
		}
		if rsp.StatusCode == http.StatusTemporaryRedirect {
			// for redirect case, resend the notification to redirect target
			uRI, err := rsp.Location()
			if err != nil {
				logger.AMpolicylog.Warnln("Policy Assocition Termination Request Notification Redirect Need Supply URI")
				return
			}
			uri = uRI.String()
			continue
		}
		return
	}

}

// returns UDR Uri of Ue, if ue.UdrUri dose not exist, query NRF to get supported Udr Uri
func getUdrUri(ue *chf_context.UeContext) string {
	if ue.UdrUri != "" {
		return ue.UdrUri
	}
	return consumer.SendNFIntancesUDR(chf_context.CHF_Self().NrfUri, ue.Supi)
}

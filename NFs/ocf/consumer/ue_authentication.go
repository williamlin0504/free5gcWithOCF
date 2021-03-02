package consumer

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"github.com/antihax/optional"

	"github.com/free5gc/nas/nasType"
	ocf_context "github.com/free5gc/ocf/context"
	"github.com/free5gc/ocf/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nausf_UEAuthentication"
	"github.com/free5gc/openapi/models"
)

func SendUEAuthenticationAuthenticateRequest(ue *ocf_context.OcfUe,
	resynchronizationInfo *models.ResynchronizationInfo) (*models.UeAuthenticationCtx, *models.ProblemDetails, error) {
	configuration := Nausf_UEAuthentication.NewConfiguration()
	configuration.SetBasePath(ue.AusfUri)

	client := Nausf_UEAuthentication.NewAPIClient(configuration)

	ocfSelf := ocf_context.OCF_Self()
	servedGuami := ocfSelf.ServedGuamiList[0]

	var authInfo models.AuthenticationInfo
	authInfo.SupiOrSuci = ue.Suci
	if mnc, err := strconv.Atoi(servedGuami.PlmnId.Mnc); err != nil {
		return nil, nil, err
	} else {
		authInfo.ServingNetworkName = fmt.Sprintf("5G:mnc%03d.mcc%s.3gppnetwork.org", mnc, servedGuami.PlmnId.Mcc)
	}
	if resynchronizationInfo != nil {
		authInfo.ResynchronizationInfo = resynchronizationInfo
	}

	ueAuthenticationCtx, httpResponse, err := client.DefaultApi.UeAuthenticationsPost(context.Background(), authInfo)
	if err == nil {
		return &ueAuthenticationCtx, nil, nil
	} else if httpResponse != nil {
		if httpResponse.Status != err.Error() {
			return nil, nil, err
		}
		problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		return nil, &problem, nil
	} else {
		return nil, nil, openapi.ReportError("server no response")
	}
}

func SendAuth5gAkaConfirmRequest(ue *ocf_context.OcfUe, resStar string) (
	*models.ConfirmationDataResponse, *models.ProblemDetails, error) {
	var ausfUri string
	if confirmUri, err := url.Parse(ue.AuthenticationCtx.Links["link"].Href); err != nil {
		return nil, nil, err
	} else {
		ausfUri = fmt.Sprintf("%s://%s", confirmUri.Scheme, confirmUri.Host)
	}

	configuration := Nausf_UEAuthentication.NewConfiguration()
	configuration.SetBasePath(ausfUri)
	client := Nausf_UEAuthentication.NewAPIClient(configuration)

	confirmData := &Nausf_UEAuthentication.UeAuthenticationsAuthCtxId5gAkaConfirmationPutParamOpts{
		ConfirmationData: optional.NewInterface(models.ConfirmationData{
			ResStar: resStar,
		}),
	}

	confirmResult, httpResponse, err := client.DefaultApi.UeAuthenticationsAuthCtxId5gAkaConfirmationPut(
		context.Background(), ue.Suci, confirmData)
	if err == nil {
		return &confirmResult, nil, nil
	} else if httpResponse != nil {
		if httpResponse.Status != err.Error() {
			return nil, nil, err
		}
		switch httpResponse.StatusCode {
		case 400, 500:
			problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
			return nil, &problem, nil
		}
		return nil, nil, nil
	} else {
		return nil, nil, openapi.ReportError("server no response")
	}
}

func SendEapAuthConfirmRequest(ue *ocf_context.OcfUe, eapMsg nasType.EAPMessage) (
	response *models.EapSession, problemDetails *models.ProblemDetails, err1 error) {
	confirmUri, err := url.Parse(ue.AuthenticationCtx.Links["link"].Href)
	if err != nil {
		logger.ConsumerLog.Errorf("url Parse failed: %+v", err)
	}
	ausfUri := fmt.Sprintf("%s://%s", confirmUri.Scheme, confirmUri.Host)

	configuration := Nausf_UEAuthentication.NewConfiguration()
	configuration.SetBasePath(ausfUri)
	client := Nausf_UEAuthentication.NewAPIClient(configuration)

	eapSessionReq := &Nausf_UEAuthentication.EapAuthMethodParamOpts{
		EapSession: optional.NewInterface(models.EapSession{
			EapPayload: base64.StdEncoding.EncodeToString(eapMsg.GetEAPMessage()),
		}),
	}

	eapSession, httpResponse, err := client.DefaultApi.EapAuthMethod(context.Background(), ue.Suci, eapSessionReq)
	if err == nil {
		response = &eapSession
	} else if httpResponse != nil {
		if httpResponse.Status != err.Error() {
			err1 = err
			return
		}
		switch httpResponse.StatusCode {
		case 400, 500:
			problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
			problemDetails = &problem
		}
	} else {
		err1 = openapi.ReportError("server no response")
	}

	return response, problemDetails, err1
}
package consumer

import (
	"context"
	"fmt"
	"free5gcWithOCF/lib/openapi"
	"free5gcWithOCF/lib/openapi/models"
	chf_context "free5gcWithOCF/src/chf/context"
	"free5gcWithOCF/src/chf/logger"
	"free5gcWithOCF/src/chf/util"
	"strings"
)

func AmfStatusChangeSubscribe(amfInfo chf_context.AMFStatusSubscriptionData) (
	problemDetails *models.ProblemDetails, err error) {
	logger.Consumerlog.Debugf("CHF Subscribe to AMF status[%+v]", amfInfo.AmfUri)
	chfSelf := chf_context.CHF_Self()
	client := util.GetNamfClient(amfInfo.AmfUri)

	subscriptionData := models.SubscriptionData{
		AmfStatusUri: fmt.Sprintf("%s/nchf-callback/v1/amfstatus", chfSelf.GetIPv4Uri()),
		GuamiList:    amfInfo.GuamiList,
	}

	res, httpResp, localErr :=
		client.SubscriptionsCollectionDocumentApi.AMFStatusChangeSubscribe(context.Background(), subscriptionData)
	if localErr == nil {
		locationHeader := httpResp.Header.Get("Location")
		logger.Consumerlog.Debugf("location header: %+v", locationHeader)

		subscriptionId := locationHeader[strings.LastIndex(locationHeader, "/")+1:]
		amfStatusSubsData := chf_context.AMFStatusSubscriptionData{
			AmfUri:       amfInfo.AmfUri,
			AmfStatusUri: res.AmfStatusUri,
			GuamiList:    res.GuamiList,
		}
		chfSelf.AMFStatusSubsData[subscriptionId] = amfStatusSubsData
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			err = localErr
			return
		}
		problem := localErr.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = openapi.ReportError("%s: server no response", amfInfo.AmfUri)
	}
	return problemDetails, err
}

package consumer

import (
	"context"
	"fmt"
	" free5gcWithOCF/lib/openapi"
	" free5gcWithOCF/lib/openapi/models"
	ccf_context " free5gcWithOCF/src/ccf/context"
	" free5gcWithOCF/src/ccf/logger"
	" free5gcWithOCF/src/ccf/util"
	"strings"
)

func AmfStatusChangeSubscribe(amfInfo ccf_context.AMFStatusSubscriptionData) (
	problemDetails *models.ProblemDetails, err error) {
	logger.Consumerlog.Debugf("CCF Subscribe to AMF status[%+v]", amfInfo.AmfUri)
	ccfSelf := ccf_context.CCF_Self()
	client := util.GetNamfClient(amfInfo.AmfUri)

	subscriptionData := models.SubscriptionData{
		AmfStatusUri: fmt.Sprintf("%s/nccf-callback/v1/amfstatus", ccfSelf.GetIPv4Uri()),
		GuamiList:    amfInfo.GuamiList,
	}

	res, httpResp, localErr :=
		client.SubscriptionsCollectionDocumentApi.AMFStatusChangeSubscribe(context.Background(), subscriptionData)
	if localErr == nil {
		locationHeader := httpResp.Header.Get("Location")
		logger.Consumerlog.Debugf("location header: %+v", locationHeader)

		subscriptionId := locationHeader[strings.LastIndex(locationHeader, "/")+1:]
		amfStatusSubsData := ccf_context.AMFStatusSubscriptionData{
			AmfUri:       amfInfo.AmfUri,
			AmfStatusUri: res.AmfStatusUri,
			GuamiList:    res.GuamiList,
		}
		ccfSelf.AMFStatusSubsData[subscriptionId] = amfStatusSubsData
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

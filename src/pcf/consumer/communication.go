package consumer

import (
	"context"
	"fmt"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/models"
	pcf_context "free5gc/src/pcf/context"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/util"
	"strings"
)

func OcfStatusChangeSubscribe(ocfInfo pcf_context.OCFStatusSubscriptionData) (
	problemDetails *models.ProblemDetails, err error) {
	logger.Consumerlog.Debugf("PCF Subscribe to OCF status[%+v]", ocfInfo.OcfUri)
	pcfSelf := pcf_context.PCF_Self()
	client := util.GetNocfClient(ocfInfo.OcfUri)

	subscriptionData := models.SubscriptionData{
		OcfStatusUri: fmt.Sprintf("%s/npcf-callback/v1/ocfstatus", pcfSelf.GetIPv4Uri()),
		GuamiList:    ocfInfo.GuamiList,
	}

	res, httpResp, localErr :=
		client.SubscriptionsCollectionDocumentApi.OCFStatusChangeSubscribe(context.Background(), subscriptionData)
	if localErr == nil {
		locationHeader := httpResp.Header.Get("Location")
		logger.Consumerlog.Debugf("location header: %+v", locationHeader)

		subscriptionId := locationHeader[strings.LastIndex(locationHeader, "/")+1:]
		ocfStatusSubsData := pcf_context.OCFStatusSubscriptionData{
			OcfUri:       ocfInfo.OcfUri,
			OcfStatusUri: res.OcfStatusUri,
			GuamiList:    res.GuamiList,
		}
		pcfSelf.OCFStatusSubsData[subscriptionId] = ocfStatusSubsData
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			err = localErr
			return
		}
		problem := localErr.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = openapi.ReportError("%s: server no response", ocfInfo.OcfUri)
	}
	return problemDetails, err
}

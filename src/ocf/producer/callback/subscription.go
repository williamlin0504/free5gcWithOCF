package callback

import (
	"context"
	"free5gc/lib/openapi/Nocf_Communication"
	"free5gc/lib/openapi/models"
	ocf_context "free5gc/src/ocf/context"
	"free5gc/src/ocf/logger"
	"reflect"
)

func SendOcfStatusChangeNotify(ocfStatus string, guamiList []models.Guami) {
	ocfSelf := ocf_context.OCF_Self()

	ocfSelf.OCFStatusSubscriptions.Range(func(key, value interface{}) bool {
		subscriptionData := value.(models.SubscriptionData)

		configuration := Nocf_Communication.NewConfiguration()
		client := Nocf_Communication.NewAPIClient(configuration)
		ocfStatusNotification := models.OcfStatusChangeNotification{}
		var ocfStatusInfo = models.OcfStatusInfo{}

		for _, guami := range guamiList {
			for _, subGumi := range subscriptionData.GuamiList {
				if reflect.DeepEqual(guami, subGumi) {
					//OCF status is available
					ocfStatusInfo.GuamiList = append(ocfStatusInfo.GuamiList, guami)
				}
			}
		}

		ocfStatusInfo = models.OcfStatusInfo{
			StatusChange:     (models.StatusChange)(ocfStatus),
			TargetOcfRemoval: "",
			TargetOcfFailure: "",
		}

		ocfStatusNotification.OcfStatusInfoList = append(ocfStatusNotification.OcfStatusInfoList, ocfStatusInfo)
		uri := subscriptionData.OcfStatusUri

		logger.ProducerLog.Infof("[OCF] Send Ocf Status Change Notify to %s", uri)
		httpResponse, err := client.OcfStatusChangeCallbackDocumentApiServiceCallbackDocumentApi.
			OcfStatusChangeNotify(context.Background(), uri, ocfStatusNotification)
		if err != nil {
			if httpResponse == nil {
				HttpLog.Errorln(err.Error())
			} else if err.Error() != httpResponse.Status {
				HttpLog.Errorln(err.Error())
			}
		}
		return true
	})
}

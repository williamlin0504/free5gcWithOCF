package callback

import (
	"context"
	"reflect"

	ocf_context "github.com/free5gc/ocf/context"
	"github.com/free5gc/ocf/logger"
	"github.com/free5gc/openapi/Nocf_Communication"
	"github.com/free5gc/openapi/models"
)

func SendOcfStatusChangeNotify(ocfStatus string, guamiList []models.Guami) {
	ocfSelf := ocf_context.OCF_Self()

	ocfSelf.OCFStatusSubscriptions.Range(func(key, value interface{}) bool {
		subscriptionData := value.(models.SubscriptionData)

		configuration := Nocf_Communication.NewConfiguration()
		client := Nocf_Communication.NewAPIClient(configuration)
		ocfStatusNotification := models.OcfStatusChangeNotification{}
		ocfStatusInfo := models.OcfStatusInfo{}

		for _, guami := range guamiList {
			for _, subGumi := range subscriptionData.GuamiList {
				if reflect.DeepEqual(guami, subGumi) {
					// OCF status is available
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

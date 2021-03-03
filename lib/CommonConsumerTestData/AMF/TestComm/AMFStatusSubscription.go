package TestComm

import (
	"free5gc/lib/openapi/models"
)

const (
	OCFStatusSubscription403      = "OCFStatusSubscription403"
	OCFStatusSubscription201      = "OCFStatusSubscription201"
	OCFStatusUnSubscription403    = "OCFStatusUnSubscription403"
	OCFStatusUnSubscription204    = "OCFStatusUnSubscription204"
	OCFStatusSubscriptionModfy403 = "OCFStatusSubscriptionModfy403"
	OCFStatusSubscriptionModfy200 = "OCFStatusSubscriptionModfy200"
)

var ConsumerOCFStatusSubscriptionTable = make(map[string]models.SubscriptionData)

func init() {
	ConsumerOCFStatusSubscriptionTable[OCFStatusSubscription403] = models.SubscriptionData{
		OcfStatusUri: "",
		GuamiList:    nil,
	}

	ConsumerOCFStatusSubscriptionTable[OCFStatusSubscription201] = models.SubscriptionData{
		OcfStatusUri: "https://127.0.0.1:29333/OCFStatusNotify",
		GuamiList: []models.Guami{
			{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				OcfId: "cafe00",
			},
		},
	}
}

var ConsumerOCFStatusUnSubscriptionTable = make(map[string]string)

func init() {
	ConsumerOCFStatusUnSubscriptionTable[OCFStatusUnSubscription403] = "0"
	ConsumerOCFStatusUnSubscriptionTable[OCFStatusUnSubscription204] = "1"
}

var ConsumerOCFStatusChangeSubscribeModfyTable = make(map[string]models.SubscriptionData)

func init() {
	ConsumerOCFStatusChangeSubscribeModfyTable[OCFStatusSubscriptionModfy403] = models.SubscriptionData{
		OcfStatusUri: "",
		GuamiList:    nil,
	}

	ConsumerOCFStatusChangeSubscribeModfyTable[OCFStatusSubscriptionModfy200] = models.SubscriptionData{
		OcfStatusUri: "https://127.0.0.1:29333/OCFStatusNotify/1",
		GuamiList: []models.Guami{
			{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				OcfId: "cafe00",
			},
		},
	}
}

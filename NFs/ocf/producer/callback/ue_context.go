package callback

import (
	"context"
	"fmt"

	ocf_context "github.com/free5gc/ocf/context"
	"github.com/free5gc/openapi/Nocf_Communication"
	"github.com/free5gc/openapi/models"
)

func SendN2InfoNotifyN2Handover(ue *ocf_context.OcfUe, releaseList []int32) error {
	if ue.HandoverNotifyUri == "" {
		return fmt.Errorf("N2 Info Notify N2Handover failed(uri dose not exist)")
	}
	configuration := Nocf_Communication.NewConfiguration()
	client := Nocf_Communication.NewAPIClient(configuration)

	n2InformationNotification := models.N2InformationNotification{
		N2NotifySubscriptionId: ue.Supi,
		ToReleaseSessionList:   releaseList,
		NotifyReason:           models.N2InfoNotifyReason_HANDOVER_COMPLETED,
	}

	_, httpResponse, err := client.N2MessageNotifyCallbackDocumentApiServiceCallbackDocumentApi.
		N2InfoNotify(context.Background(), ue.HandoverNotifyUri, n2InformationNotification)

	if err == nil {
		// TODO: handle Msg
	} else {
		if httpResponse == nil {
			HttpLog.Errorln(err.Error())
		} else if err.Error() != httpResponse.Status {
			HttpLog.Errorln(err.Error())
		}
	}
	return nil
}

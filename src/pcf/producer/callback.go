package producer

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"net/http"
)

func HandleOcfStatusChangeNotify(request *http_wrapper.Request) *http_wrapper.Response {
	logger.CallbackLog.Warnf("[PCF] Handle Ocf Status Change Notify is not implemented.")

	notification := request.Body.(models.OcfStatusChangeNotification)

	OcfStatusChangeNotifyProcedure(notification)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
}

// TODO: handle OCF Status Change Notify
func OcfStatusChangeNotifyProcedure(notification models.OcfStatusChangeNotification) {
	logger.CallbackLog.Debugf("receive OCF status change notification[%+v]", notification)
}

func HandleSmPolicyNotify(request *http_wrapper.Request) *http_wrapper.Response {
	logger.CallbackLog.Warnf("[PCF] Handle Sm Policy Notify is not implemented.")

	notification := request.Body.(models.PolicyDataChangeNotification)
	supi := request.Params["ReqURI"]

	SmPolicyNotifyProcedure(supi, notification)

	return http_wrapper.NewResponse(http.StatusNotImplemented, nil, nil)
}

// TODO: handle SM Policy Notify
func SmPolicyNotifyProcedure(supi string, notification models.PolicyDataChangeNotification) {
}

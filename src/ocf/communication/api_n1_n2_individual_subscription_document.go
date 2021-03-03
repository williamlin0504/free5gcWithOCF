/*
 * Namf_Communication
 *
 * OCF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package communication

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/models"
	"free5gc/src/ocf/logger"
	"free5gc/src/ocf/producer"
	"net/http"

	"github.com/gin-gonic/gin"
)

// N1N2MessageUnSubscribe - Namf_Communication N1N2 Message UnSubscribe (UE Specific) service Operation
func HTTPN1N2MessageUnSubscribe(c *gin.Context) {

	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["ueContextId"] = c.Params.ByName("ueContextId")
	req.Params["subscriptionId"] = c.Params.ByName("subscriptionId")

	rsp := producer.HandleN1N2MessageUnSubscribeRequest(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.CommLog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}
}

/*
 * Nsmf_EventExposure
 *
 * Session Management Event Exposure Service API
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package callback

import (
	" free5gc/lib/http_wrapper"
	" free5gcenapi"
	" free5gcenapi/models"
	" free5gcf/logger"
	" free5gcf/producer"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubscriptionsPost -
func HTTPSmPolicyUpdateNotification(c *gin.Context) {
	var request models.SmPolicyNotification

	reqBody, _ := c.GetRawData()

	err := openapi.Deserialize(&request, reqBody, c.ContentType())
	if err != nil {
		logger.PduSessLog.Errorln("Deserialize request failed")
	}

	reqWrapper := http_wrapper.NewRequest(c.Request, request)
	reqWrapper.Params["smContextRef"] = c.Params.ByName("smContextRef")

	smContextRef := reqWrapper.Params["smContextRef"]
	HTTPResponse := producer.HandleSMPolicyUpdateNotify(smContextRef, reqWrapper.Body.(models.SmPolicyNotification))

	for key, val := range HTTPResponse.Header {
		c.Header(key, val[0])
	}

	resBody, err := openapi.Serialize(HTTPResponse.Body, "application/json")
	c.Writer.Write(resBody)
	c.Status(HTTPResponse.Status)
}

func SmPolicyControlTerminationRequestNotification(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

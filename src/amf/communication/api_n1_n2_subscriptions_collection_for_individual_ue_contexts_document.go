/*
 * Namf_Communication
 *
 * AMF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package communication

import (
	"free5gcWithOCF/lib/http_wrapper"
	"free5gcWithOCF/lib/openapi"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/src/amf/logger"
	"free5gcWithOCF/src/amf/producer"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HTTPN1N2MessageSubscribe(c *gin.Context) {
	var ueN1N2InfoSubscriptionCreateData models.UeN1N2InfoSubscriptionCreateData

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&ueN1N2InfoSubscriptionCreateData, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, ueN1N2InfoSubscriptionCreateData)
	req.Params["ueContextId"] = c.Params.ByName("ueContextId")

	rsp := producer.HandleN1N2MessageSubscirbeRequest(req)

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

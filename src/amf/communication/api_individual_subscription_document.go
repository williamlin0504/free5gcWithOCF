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
	" free5gc/lib/http_wrapper"
	" free5gcenapi"
	" free5gcenapi/models"
	" free5gcf/logger"
	" free5gcf/producer"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AMFStatusChangeSubscribeModify - Namf_Communication AMF Status Change Subscribe Modify service Operation
func HTTPAMFStatusChangeSubscribeModify(c *gin.Context) {
	var subscriptionData models.SubscriptionData

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

	err = openapi.Deserialize(&subscriptionData, requestBody, "application/json")
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

	req := http_wrapper.NewRequest(c.Request, subscriptionData)
	req.Params["subscriptionId"] = c.Params.ByName("subscriptionId")

	rsp := producer.HandleAMFStatusChangeSubscribeModify(req)

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

// AMFStatusChangeUnSubscribe - Namf_Communication AMF Status Change UnSubscribe service Operation
func HTTPAMFStatusChangeUnSubscribe(c *gin.Context) {

	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["subscriptionId"] = c.Params.ByName("subscriptionId")

	rsp := producer.HandleAMFStatusChangeUnSubscribeRequest(req)

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

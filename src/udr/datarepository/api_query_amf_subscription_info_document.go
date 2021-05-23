/*
 * Nudr_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package datarepository

import (
	" free5gc/lib/http_wrapper"
	" free5gcenapi"
	" free5gcenapi/models"
	" free5gcr/logger"
	" free5gcr/producer"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPGetAmfSubscriptionInfo - Retrieve AMF subscription Info
func HTTPGetAmfSubscriptionInfo(c *gin.Context) {
	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Params.ByName("ueId")
	req.Params["subsId"] = c.Params.ByName("subsId")

	rsp := producer.HandleGetAmfSubscriptionInfo(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.DataRepoLog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}

	// req := http_wrapper.NewRequest(c.Request, nil)
	// req.Params["ueId"] = c.Params.ByName("ueId")
	// req.Params["subsId"] = c.Params.ByName("subsId")

	// handlerMsg := message.NewHandlerMessage(message.EventGetAmfSubscriptionInfo, req)
	// message.SendMessage(handlerMsg)

	// rsp := <-handlerMsg.ResponseChan

	// HTTPResponse := rsp.HTTPResponse

	// c.JSON(HTTPResponse.Status, HTTPResponse.Body)
}

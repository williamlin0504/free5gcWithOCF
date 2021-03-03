/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package nssaiavailability

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"free5gcWithOCF/lib/http_wrapper"
	"free5gcWithOCF/lib/openapi"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/src/nssf/logger"
	"free5gcWithOCF/src/nssf/producer"
)

func HTTPNSSAIAvailabilityUnsubscribe(c *gin.Context) {
	// Due to conflict of route matching, 'subscriptions' in the route is replaced with the existing wildcard ':nfId'
	nfID := c.Param("nfId")
	if nfID != "subscriptions" {
		c.JSON(http.StatusNotFound, gin.H{})
		logger.Nssaiavailability.Infof("404 Not Found")
		return
	}

	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["subscriptionId"] = c.Params.ByName("subscriptionId")

	rsp := producer.HandleNSSAIAvailabilityUnsubscribe(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.HandlerLog.Errorln(err)
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
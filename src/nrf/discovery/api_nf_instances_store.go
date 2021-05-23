/*
 * NRF NFDiscovery Service
 *
 * NRF NFDiscovery  Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package discovery

import (
	" free5gcWithOCF/lib/http_wrapper"
	" free5gcWithOCF/lib/openapi"
	" free5gcWithOCF/lib/openapi/models"
	" free5gcWithOCF/src/nrf/logger"
	" free5gcWithOCF/src/nrf/producer"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SearchNFInstances - Search a collection of NF Instances
func HTTPSearchNFInstances(c *gin.Context) {
	// var searchNFInstance context.SearchNFInstances
	// c.BindQuery(&searchNFInstance)
	//logger.DiscoveryLog.Infoln("searchNFInstance: ", searchNFInstance)
	// logger.DiscoveryLog.Infoln("targetNFType: ", searchNFInstance.TargetNFType)
	// logger.DiscoveryLog.Infoln("requesterNFType: ", searchNFInstance.RequesterNFType)
	//logger.DiscoveryLog.Infoln("ChfSupportedPlmn: ", searchNFInstance.ChfSupportedPlmn)

	req := http_wrapper.NewRequest(c.Request, nil)
	req.Query = c.Request.URL.Query()
	httpResponse := producer.HandleNFDiscoveryRequest(req)

	responseBody, err := openapi.Serialize(httpResponse.Body, "application/json")
	if err != nil {
		logger.DiscoveryLog.Warnln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(httpResponse.Status, "application/json", responseBody)
	}

}

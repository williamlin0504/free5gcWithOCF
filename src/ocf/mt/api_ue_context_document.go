/*
 * Namf_MT
 *
 * OCF Mobile Termination Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package mt

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/models"
	"free5gc/src/ocf/logger"
	"free5gc/src/ocf/producer"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProvideDomainSelectionInfo - Namf_MT Provide Domain Selection Info service Operation
func HTTPProvideDomainSelectionInfo(c *gin.Context) {

	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["ueContextId"] = c.Params.ByName("ueContextId")
	infoClassQuery := c.Query("info-class")
	req.Query.Add("info-class", infoClassQuery)
	supportedFeaturesQuery := c.Query("supported-features")
	req.Query.Add("supported-features", supportedFeaturesQuery)

	rsp := producer.HandleProvideDomainSelectionInfoRequest(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.MtLog.Errorln(err)
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
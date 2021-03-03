/*
 * Nudm_UECM
 *
 * Nudm Context Management Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package uecontextmanagement

import (
	// "fmt"
	"free5gcWithOCF/lib/http_wrapper"
	"free5gcWithOCF/lib/openapi"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/src/udm/logger"
	"free5gcWithOCF/src/udm/producer"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAmf3gppAccess - retrieve the AMF registration for 3GPP access information
func HTTPGetAmf3gppAccess(c *gin.Context) {
	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Param("ueId")
	req.Query.Add("supported-features", c.Query("supported-features"))

	rsp := producer.HandleGetAmf3gppAccessRequest(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.UecmLog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}
	return
}
/*
 * Nudm_SDM
 *
 * Nudm Subscriber Data Management Service
 *
 * API version: 2.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package subscriberdatamanagement

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/producer"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetNssai - retrieve a UE's subscribed NSSAI
func HTTPGetNssai(c *gin.Context) {

	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["supi"] = c.Params.ByName("supi")
	req.Query.Set("plmn-id", c.Query("plmn-id"))
	req.Query.Set("supported-features", c.Query("supported-features"))

	rsp := producer.HandleGetNssaiRequest(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.SdmLog.Errorln(err)
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

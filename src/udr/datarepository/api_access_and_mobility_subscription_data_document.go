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

// HTTPQueryAmData - Retrieves the access and mobility subscription data of a UE
func HTTPQueryAmData(c *gin.Context) {

	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Params.ByName("ueId")
	req.Params["servingPlmnId"] = c.Params.ByName("servingPlmnId")

	rsp := producer.HandleQueryAmData(req)

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
}

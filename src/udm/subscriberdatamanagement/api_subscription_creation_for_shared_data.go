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
	" free5gcWithOCF/lib/http_wrapper"
	" free5gcWithOCF/lib/openapi"
	" free5gcWithOCF/lib/openapi/models"
	" free5gcWithOCF/src/udm/logger"
	" free5gcWithOCF/src/udm/producer"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SubscribeToSharedData - subscribe to notifications for shared data
func HTTPSubscribeToSharedData(c *gin.Context) {

	var sharedDataSubsReq models.SdmSubscription
	// step 1: retrieve http request body
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.SdmLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	// step 2: convert requestBody to openapi models
	err = openapi.Deserialize(&sharedDataSubsReq, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.SdmLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, sharedDataSubsReq)
	rsp := producer.HandleSubscribeToSharedDataRequest(req)
	// step 5: response
	for key, val := range rsp.Header { // header response is optional
		c.Header(key, val[0])
	}
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

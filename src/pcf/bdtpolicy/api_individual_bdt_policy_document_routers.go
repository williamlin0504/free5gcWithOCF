/*
 * Nccf_BDTPolicyControl Service API
 *
 * The Nccf_BDTPolicyControl Service is used by an NF service consumer to retrieve background data transfer policies from the ccf and to update the ccf with the background data transfer policy selected by the NF service consumer.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package bdtpolicy

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/models"
	"free5gc/src/ccf/logger"
	"free5gc/src/ccf/producer"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetBDTPolicy - Read an Individual BDT policy
func HTTPGetBDTPolicy(c *gin.Context) {
	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["bdtPolicyId"] = c.Params.ByName("bdtPolicyId")

	rsp := producer.HandleGetBDTPolicyContextRequest(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.Bdtpolicylog.Errorln(err)
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

// UpdateBDTPolicy - Update an Individual BDT policy
func HTTPUpdateBDTPolicy(c *gin.Context) {
	var bdtPolicyDataPatch models.BdtPolicyDataPatch
	// step 1: retrieve http request body
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.Bdtpolicylog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	// step 2: convert requestBody to openapi models
	err = openapi.Deserialize(&bdtPolicyDataPatch, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.Bdtpolicylog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, bdtPolicyDataPatch)
	req.Params["bdtPolicyId"] = c.Params.ByName("bdtPolicyId")

	rsp := producer.HandleUpdateBDTPolicyContextProcedure(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.Bdtpolicylog.Errorln(err)
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

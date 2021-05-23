/*
 * Npcf_AMPolicyControl
 *
 * Access and Mobility Policy Control Service API
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package ampolicy

import (
	" free5gcWithOCF/lib/http_wrapper"
	" free5gcWithOCF/lib/openapi"
	" free5gcWithOCF/lib/openapi/models"
	" free5gcWithOCF/src/pcf/logger"
	" free5gcWithOCF/src/pcf/producer"
	" free5gcWithOCF/src/pcf/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HTTPPoliciesPolAssoIdDelete(c *gin.Context) {
	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["polAssoId"], _ = c.Params.Get("polAssoId")

	rsp := producer.HandleDeletePoliciesPolAssoId(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.AMpolicylog.Errorln(err)
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

// HTTPPoliciesPolAssoIdGet -
func HTTPPoliciesPolAssoIdGet(c *gin.Context) {
	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["polAssoId"], _ = c.Params.Get("polAssoId")

	rsp := producer.HandleGetPoliciesPolAssoId(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.AMpolicylog.Errorln(err)
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

// HTTPPoliciesPolAssoIdUpdatePost -
func HTTPPoliciesPolAssoIdUpdatePost(c *gin.Context) {
	var policyAssociationUpdateRequest models.PolicyAssociationUpdateRequest

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.AMpolicylog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&policyAssociationUpdateRequest, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.AMpolicylog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, policyAssociationUpdateRequest)
	req.Params["polAssoId"], _ = c.Params.Get("polAssoId")

	rsp := producer.HandleUpdatePostPoliciesPolAssoId(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.AMpolicylog.Errorln(err)
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

// HTTPPoliciesPost -
func HTTPPoliciesPost(c *gin.Context) {
	var policyAssociationRequest models.PolicyAssociationRequest

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.AMpolicylog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&policyAssociationRequest, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.AMpolicylog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	if policyAssociationRequest.Supi == "" || policyAssociationRequest.NotificationUri == "" {
		rsp := util.GetProblemDetail("Miss Mandotory IE", util.ERROR_REQUEST_PARAMETERS)
		logger.HandlerLog.Errorln(rsp.Detail)
		c.JSON(int(rsp.Status), rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, policyAssociationRequest)
	req.Params["polAssoId"], _ = c.Params.Get("polAssoId")

	rsp := producer.HandlePostPolicies(req)

	for key, val := range rsp.Header {
		c.Header(key, val[0])
	}

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.AMpolicylog.Errorln(err)
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

/*
 * Nccf_PolicyAuthorization Service API
 *
 * This is the Policy Authorization Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package policyauthorization

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/models"
	"free5gc/src/ccf/logger"
	"free5gc/src/ccf/producer"
	"free5gc/src/ccf/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPDeleteEventsSubsc - deletes the Events Subscription subresource
func HTTPDeleteEventsSubsc(c *gin.Context) {
	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["appSessionId"], _ = c.Params.Get("appSessionId")

	rsp := producer.HandleDeleteEventsSubscContext(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.PolicyAuthorizationlog.Errorln(err)
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

// HTTPUpdateEventsSubsc - creates or modifies an Events Subscription subresource
func HTTPUpdateEventsSubsc(c *gin.Context) {
	var eventsSubscReqData models.EventsSubscReqData

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.PolicyAuthorizationlog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&eventsSubscReqData, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.PolicyAuthorizationlog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	if eventsSubscReqData.Events == nil || eventsSubscReqData.NotifUri == "" {
		problemDetail := util.GetProblemDetail("Errorneous/Missing Mandotory IE", util.ERROR_REQUEST_PARAMETERS)
		logger.PolicyAuthorizationlog.Errorln(problemDetail.Detail)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	req := http_wrapper.NewRequest(c.Request, eventsSubscReqData)
	req.Params["appSessionId"], _ = c.Params.Get("appSessionId")

	rsp := producer.HandleUpdateEventsSubscContext(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.PolicyAuthorizationlog.Errorln(err)
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

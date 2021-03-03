/*
 * Nsmf_PDUSession
 *
 * SMF PDU Session Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package pdusession

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/models"
	"free5gc/src/smf/logger"
	"free5gc/src/smf/producer"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HTTPPostSmContexts - Create SM Context
func HTTPPostSmContexts(c *gin.Context) {
	logger.PduSessLog.Info("Recieve Create SM Context Request")
	var request models.PostSmContextsRequest

	request.JsonData = new(models.SmContextCreateData)

	s := strings.Split(c.GetHeader("Content-Type"), ";")
	var err error
	switch s[0] {
	case "application/json":
		err = c.ShouldBindJSON(request.JsonData)
	case "multipart/related":
		err = c.ShouldBindWith(&request, openapi.MultipartRelatedBinding{})
	}

	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.PduSessLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, request)
	HTTPResponse := producer.HandlePDUSessionSMContextCreate(req.Body.(models.PostSmContextsRequest))
	//Http Response to OCF
	for key, val := range HTTPResponse.Header {
		c.Header(key, val[0])
	}
	switch HTTPResponse.Status {
	case http.StatusCreated,
		http.StatusBadRequest,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		c.Render(HTTPResponse.Status, openapi.MultipartRelatedRender{Data: HTTPResponse.Body})
	default:
		c.JSON(HTTPResponse.Status, HTTPResponse.Body)
	}
}

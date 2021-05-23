/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package producer

import (
	"net/http"

	" free5gc/lib/http_wrapper"
	" free5gcenapi/models"
	" free5gcsf/logger"
)

// HandleNSSAIAvailabilityUnsubscribe - Deletes an already existing NSSAI availability notification subscription
func HandleNSSAIAvailabilityUnsubscribe(request *http_wrapper.Request) *http_wrapper.Response {
	logger.Nssaiavailability.Infof("Handle NSSAIAvailabilityUnsubscribe")

	subscriptionID := request.Params["subscriptionId"]

	problemDetails := NSSAIAvailabilityUnsubscribeProcedure(subscriptionID)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

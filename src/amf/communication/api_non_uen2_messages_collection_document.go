/*
 * Namf_Communication
 *
 * AMF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package communication

import (
	" free5gcWithOCF/src/amf/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NonUeN2MessageTransfer - Namf_Communication Non UE N2 Message Transfer service Operation
func HTTPNonUeN2MessageTransfer(c *gin.Context) {
	logger.CommLog.Warnf("Handle Non Ue N2 Message Transfer is not implemented.")
	c.JSON(http.StatusOK, gin.H{})
}

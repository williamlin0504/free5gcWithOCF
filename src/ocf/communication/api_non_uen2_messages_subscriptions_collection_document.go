/*
 * Nocf_Communication
 *
 * OCF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package communication

import (
	"free5gcWithOCF/src/ocf/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NonUeN2InfoSubscribe - Nocf_Communication Non UE N2 Info Subscribe service Operation
func HTTPNonUeN2InfoSubscribe(c *gin.Context) {
	logger.CommLog.Warnf("Handle Non Ue N2 Info Subscribe is not implemented.")
	c.JSON(http.StatusOK, gin.H{})
}

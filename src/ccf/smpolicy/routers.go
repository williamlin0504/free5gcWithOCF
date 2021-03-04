/*
 * Nchf_SMPolicyControl
 *
 * Session Management Policy Control Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package smpolicy

import (
	"free5gcWithOCF/lib/logger_util"
	"free5gcWithOCF/src/chf/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)
	AddService(router)
	return router
}

func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("/nchf-smpolicycontrol/v1")

	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		case "PUT":
			group.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.HandlerFunc)
		case "PATCH":
			group.PATCH(route.Pattern, route.HandlerFunc)
		}
	}
	return group
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		"GET",
		"/",
		Index,
	},

	{
		"SmPoliciesPost",
		strings.ToUpper("Post"),
		"/sm-policies",
		HTTPSmPoliciesPost,
	},

	{
		"SmPoliciesSmPolicyIdDeletePost",
		strings.ToUpper("Post"),
		"/sm-policies/:smPolicyId/delete",
		HTTPSmPoliciesSmPolicyIdDeletePost,
	},

	{
		"SmPoliciesSmPolicyIdGet",
		strings.ToUpper("Get"),
		"/sm-policies/:smPolicyId",
		HTTPSmPoliciesSmPolicyIDGet,
	},

	{
		"SmPoliciesSmPolicyIdUpdatePost",
		strings.ToUpper("Post"),
		"/sm-policies/:smPolicyId/update",
		HTTPSmPoliciesSmPolicyIdUpdatePost,
	},
}

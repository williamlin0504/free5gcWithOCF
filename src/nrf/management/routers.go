/*
 * NRF NFManagement Service
 *
 * NRF NFManagement Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package management

import (
	"free5gcWithOCF/lib/logger_util"
	"free5gcWithOCF/src/nrf/logger"
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
	group := engine.Group("/nnrf-nfm/v1")

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
		"DeregisterNFInstance",
		strings.ToUpper("Delete"),
		"/nf-instances/:nfInstanceID",
		HTTPDeregisterNFInstance,
	},

	{
		"GetNFInstance",
		strings.ToUpper("Get"),
		"/nf-instances/:nfInstanceID",
		HTTPGetNFInstance,
	},

	{
		"RegisterNFInstance",
		strings.ToUpper("Put"),
		"/nf-instances/:nfInstanceID",
		HTTPRegisterNFInstance,
	},

	{
		"UpdateNFInstance",
		strings.ToUpper("Patch"),
		"/nf-instances/:nfInstanceID",
		HTTPUpdateNFInstance,
	},

	{
		"GetNFInstances",
		strings.ToUpper("Get"),
		"/nf-instances",
		HTTPGetNFInstances,
	},

	{
		"RemoveSubscription",
		strings.ToUpper("Delete"),
		"/subscriptions/:subscriptionID",
		HTTPRemoveSubscription,
	},

	{
		"UpdateSubscription",
		strings.ToUpper("Patch"),
		"/subscriptions/:subscriptionID",
		HTTPUpdateSubscription,
	},

	{
		"CreateSubscription",
		strings.ToUpper("Post"),
		"/subscriptions",
		HTTPCreateSubscription,
	},
}
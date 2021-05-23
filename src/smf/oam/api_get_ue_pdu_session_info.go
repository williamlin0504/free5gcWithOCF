package oam

import (
	" free5gcWithOCF/lib/http_wrapper"
	" free5gcWithOCF/src/smf/producer"
	"github.com/gin-gonic/gin"
)

func HTTPGetUEPDUSessionInfo(c *gin.Context) {

	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["smContextRef"] = c.Params.ByName("smContextRef")

	smContextRef := req.Params["smContextRef"]
	HTTPResponse := producer.HandleOAMGetUEPDUSessionInfo(smContextRef)

	c.JSON(HTTPResponse.Status, HTTPResponse.Body)
}

package main

import (
	"fmt"
	" free5gcWithOCF/lib/http2_util"
	" free5gcWithOCF/lib/logger_util"
	. " free5gcWithOCF/lib/openapi/models"
	" free5gcWithOCF/lib/path_util"
	" free5gcWithOCF/src/nrf/logger"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	NrfLogPath = path_util.Go free5gcPath(" free5gcWithOCF/src/nrf/management/sslkeylog.log")
	NrfPemPath = path_util.Go free5gcPath(" free5gcWithOCF/support/TLS/nrf.pem")
	NrfKeyPath = path_util.Go free5gcPath(" free5gcWithOCF/support/TLS/nrf.key")
)

func main() {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	router.POST("", func(c *gin.Context) {
		/*buf, err := c.GetRawData()
		if err != nil {
			t.Errorf(err.Error())
		}
		// Remove NL line feed, new line character
		//requestBody = string(buf[:len(b uf)-1])*/
		var ND NotificationData

		if err := c.ShouldBindJSON(&ND); err != nil {
			log.Panic(err.Error())
		}
		fmt.Println(ND)
		c.JSON(http.StatusNoContent, gin.H{})
	})

	srv, err := http2_util.NewServer(":30678", NrfLogPath, router)
	if err != nil {
		log.Panic(err.Error())
	}

	err2 := srv.ListenAndServeTLS(NrfPemPath, NrfKeyPath)
	if err2 != nil && err2 != http.ErrServerClosed {
		log.Panic(err2.Error())
	}

}

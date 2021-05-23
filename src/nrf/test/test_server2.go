package main

import (
	"fmt"
	" free5gc/lib/http2_util"
	" free5gcgger_util"
	. " free5gcenapi/models"
	" free5gcth_util"
	" free5gcf/logger"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	NrfLogPath = path_util.Go free5gcPath(" free5gcf/management/sslkeylog.log")
	NrfPemPath = path_util.Go free5gcPath(" free5gct/TLS/nrf.pem")
	NrfKeyPath = path_util.Go free5gcPath(" free5gct/TLS/nrf.key")
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

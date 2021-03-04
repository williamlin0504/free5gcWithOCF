package pdusession

import (
	"free5gcWithOCF/lib/http2_util"
	"free5gcWithOCF/lib/logger_util"
	"free5gcWithOCF/lib/path_util"
	"free5gcWithOCF/src/smf/logger"
	"free5gcWithOCF/src/smf/pfcp"
	"free5gcWithOCF/src/smf/pfcp/udp"
	"log"
	"net/http"
)

func DummyServer() {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	AddService(router)

	go udp.Run(pfcp.Dispatch)

	smfKeyLogPath := path_util.Gofree5gcPath("free5gcWithOCF/smfsslkey.log")
	smfPemPath := path_util.Gofree5gcPath("free5gcWithOCF/support/TLS/smf.pem")
	smfkeyPath := path_util.Gofree5gcPath("free5gcWithOCF/support/TLS/smf.key")

	var server *http.Server
	if srv, err := http2_util.NewServer(":29502", smfKeyLogPath, router); err != nil {
	} else {
		server = srv
	}

	if err := server.ListenAndServeTLS(smfPemPath, smfkeyPath); err != nil {
		log.Fatal(err)
	}

}

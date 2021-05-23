package pdusession

import (
	" free5gc/lib/http2_util"
	" free5gcgger_util"
	" free5gcth_util"
	" free5gcf/logger"
	" free5gcf/pfcp"
	" free5gcf/pfcp/udp"
	"log"
	"net/http"
)

func DummyServer() {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	AddService(router)

	go udp.Run(pfcp.Dispatch)

	smfKeyLogPath := path_util.Go free5gcPath(" free5gckey.log")
	smfPemPath := path_util.Go free5gcPath(" free5gct/TLS/smf.pem")
	smfkeyPath := path_util.Go free5gcPath(" free5gct/TLS/smf.key")

	var server *http.Server
	if srv, err := http2_util.NewServer(":29502", smfKeyLogPath, router); err != nil {
	} else {
		server = srv
	}

	if err := server.ListenAndServeTLS(smfPemPath, smfkeyPath); err != nil {
		log.Fatal(err)
	}

}

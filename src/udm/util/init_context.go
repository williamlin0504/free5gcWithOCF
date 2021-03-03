package util

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/src/udm/context"
	"free5gcWithOCF/src/udm/factory"
	"free5gcWithOCF/src/udm/logger"
)

func InitUDMContext(udmContext *context.UDMContext) {
	config := factory.UdmConfig
	logger.UtilLog.Info("udmconfig Info: Version[", config.Info.Version, "] Description[", config.Info.Description, "]")
	configuration := config.Configuration
	udmContext.NfId = uuid.New().String()
	if configuration.UdmName != "" {
		udmContext.Name = configuration.UdmName
	}
	nrfclient := config.Configuration.Nrfclient
	udmContext.NrfUri = fmt.Sprintf("%s://%s:%d", nrfclient.Scheme, nrfclient.Ipv4Addr, nrfclient.Port)
	sbi := configuration.Sbi
	udmContext.UriScheme = ""
	udmContext.SBIPort = 29503
	udmContext.RegisterIPv4 = "127.0.0.1"
	if sbi != nil {
		if sbi.Scheme != "" {
			udmContext.UriScheme = models.UriScheme(sbi.Scheme)
		}
		if sbi.RegisterIPv4 != "" {
			udmContext.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			udmContext.SBIPort = sbi.Port
		}

		udmContext.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if udmContext.BindingIPv4 != "" {
			logger.UtilLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			udmContext.BindingIPv4 = sbi.BindingIPv4
			if udmContext.BindingIPv4 == "" {
				logger.UtilLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				udmContext.BindingIPv4 = "0.0.0.0"
			}
		}
	}
	if configuration.NrfUri != "" {
		udmContext.NrfUri = configuration.NrfUri
	} else {
		logger.UtilLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
		udmContext.NrfUri = fmt.Sprintf("%s://%s:%d", udmContext.UriScheme, "127.0.0.1", 29510)
	}
	servingNameList := configuration.ServiceNameList

	udmContext.Keys = configuration.Keys

	udmContext.InitNFService(servingNameList, config.Info.Version)
}
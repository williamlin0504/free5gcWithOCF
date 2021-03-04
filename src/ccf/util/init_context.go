package util

import (
	"os"

	"github.com/google/uuid"

	"free5gcWithOCF/lib/openapi"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/src/chf/context"
	"free5gcWithOCF/src/chf/factory"
	"free5gcWithOCF/src/chf/logger"
)

// Init CHF Context from config flie
func InitchfContext(context *context.CHFContext) {
	config := factory.ChfConfig
	logger.UtilLog.Infof("chfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	if configuration.ChfName != "" {
		context.Name = configuration.ChfName
	}
	sbi := configuration.Sbi
	context.NrfUri = configuration.NrfUri
	context.UriScheme = ""
	context.RegisterIPv4 = "127.0.0.1" // default localhost
	context.SBIPort = 29507            // default port
	if sbi != nil {
		if sbi.Scheme != "" {
			context.UriScheme = models.UriScheme(sbi.Scheme)
		}
		if sbi.RegisterIPv4 != "" {
			context.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			context.SBIPort = sbi.Port
		}
		if sbi.Scheme == "https" {
			context.UriScheme = models.UriScheme_HTTPS
		} else {
			context.UriScheme = models.UriScheme_HTTP
		}

		context.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if context.BindingIPv4 != "" {
			logger.UtilLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			context.BindingIPv4 = sbi.BindingIPv4
			if context.BindingIPv4 == "" {
				logger.UtilLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				context.BindingIPv4 = "0.0.0.0"
			}
		}
	}
	serviceList := configuration.ServiceList
	context.InitNFService(serviceList, config.Info.Version)
	context.TimeFormat = configuration.TimeFormat
	context.DefaultBdtRefId = configuration.DefaultBdtRefId
	for _, service := range context.NfService {
		var err error
		context.ChfServiceUris[service.ServiceName] =
			service.ApiPrefix + "/" + string(service.ServiceName) + "/" + (*service.Versions)[0].ApiVersionInUri
		context.ChfSuppFeats[service.ServiceName], err = openapi.NewSupportedFeature(service.SupportedFeatures)
		if err != nil {
			logger.UtilLog.Errorf("openapi NewSupportedFeature error: %+v", err)
		}
	}
}

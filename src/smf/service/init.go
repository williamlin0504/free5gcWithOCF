package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"free5gcWithOCF/lib/http2_util"
	"free5gcWithOCF/lib/logger_util"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/lib/path_util"
	"free5gcWithOCF/src/app"
	"free5gcWithOCF/src/smf/callback"
	"free5gcWithOCF/src/smf/consumer"
	"free5gcWithOCF/src/smf/context"
	"free5gcWithOCF/src/smf/eventexposure"
	"free5gcWithOCF/src/smf/factory"
	"free5gcWithOCF/src/smf/logger"
	"free5gcWithOCF/src/smf/oam"
	"free5gcWithOCF/src/smf/pdusession"
	"free5gcWithOCF/src/smf/pfcp"
	"free5gcWithOCF/src/smf/pfcp/message"
	"free5gcWithOCF/src/smf/pfcp/udp"
	"free5gcWithOCF/src/smf/util"
)

type SMF struct{}

type (
	// Config information.
	Config struct {
		smfcfg    string
		uerouting string
	}
)

var config Config

var smfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gcWithOCFcfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "smfcfg",
		Usage: "config file",
	},
	cli.StringFlag{
		Name:  "uerouting",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*SMF) GetCliCmd() (flags []cli.Flag) {
	return smfCLi
}

func (*SMF) Initialize(c *cli.Context) {

	config = Config{
		smfcfg:    c.String("smfcfg"),
		uerouting: c.String("uerouting"),
	}

	if config.smfcfg != "" {
		factory.InitConfigFactory(config.smfcfg)
	} else {
		DefaultSmfConfigPath := path_util.Gofree5gcWithOCFPath("free5gcWithOCF/config/smfcfg.conf")
		factory.InitConfigFactory(DefaultSmfConfigPath)
	}

	if config.uerouting != "" {
		factory.InitRoutingConfigFactory(config.uerouting)
	} else {
		DefaultUERoutingPath := path_util.Gofree5gcWithOCFPath("free5gcWithOCF/config/uerouting.yaml")
		factory.InitRoutingConfigFactory(DefaultUERoutingPath)
	}

	if app.ContextSelf().Logger.SMF.DebugLevel != "" {
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.SMF.DebugLevel)
		if err != nil {
			initLog.Warnf("Log level [%s] is not valid, set to [info] level", app.ContextSelf().Logger.SMF.DebugLevel)
			logger.SetLogLevel(logrus.InfoLevel)
		} else {
			logger.SetLogLevel(level)
			initLog.Infof("Log level is set to [%s] level", level)
		}
	} else {
		initLog.Infoln("Log level is default set to [info] level")
		logger.SetLogLevel(logrus.InfoLevel)
	}
	logger.SetReportCaller(app.ContextSelf().Logger.SMF.ReportCaller)
}

func (smf *SMF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range smf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (smf *SMF) Start() {
	context.InitSmfContext(&factory.SmfConfig)
	//allocate id for each upf
	context.AllocateUPFID()
	context.InitSMFUERouting(&factory.UERoutingConfig)

	initLog.Infoln("Server started")
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	err := consumer.SendNFRegistration()

	if err != nil {
		retry_err := consumer.RetrySendNFRegistration(10)
		if retry_err != nil {
			logger.InitLog.Errorln(retry_err)
			return
		}
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		smf.Terminate()
		os.Exit(0)
	}()

	oam.AddService(router)
	callback.AddService(router)
	for _, serviceName := range factory.SmfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serviceName) {
		case models.ServiceName_NSMF_PDUSESSION:
			pdusession.AddService(router)
		case models.ServiceName_NSMF_EVENT_EXPOSURE:
			eventexposure.AddService(router)
		}
	}
	udp.Run(pfcp.Dispatch)

	for _, upf := range context.SMF_Self().UserPlaneInformation.UPFs {
		logger.AppLog.Infof("Send PFCP Association Request to UPF[%s]\n", upf.NodeID.NodeIdValue)
		message.SendPfcpAssociationSetupRequest(upf.NodeID)
	}

	time.Sleep(1000 * time.Millisecond)

	HTTPAddr := fmt.Sprintf("%s:%d", context.SMF_Self().BindingIPv4, context.SMF_Self().SBIPort)
	server, err := http2_util.NewServer(HTTPAddr, util.SmfLogPath, router)

	if server == nil {
		initLog.Error("Initialize HTTP server failed:", err)
		return
	}

	if err != nil {
		initLog.Warnln("Initialize HTTP server:", err)
	}

	serverScheme := factory.SmfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.SmfPemPath, util.SmfKeyPath)
	}

	if err != nil {
		initLog.Fatalln("HTTP server setup failed:", err)
	}

}

func (smf *SMF) Terminate() {
	logger.InitLog.Infof("Terminating SMF...")
	// deregister with NRF
	problemDetails, err := consumer.SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("Deregister from NRF successfully")
	}
}

func (smf *SMF) Exec(c *cli.Context) error {
	return nil
}

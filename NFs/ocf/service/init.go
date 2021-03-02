package service

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	aperLogger "github.com/free5gc/aper/logger"
	fsmLogger "github.com/free5gc/fsm/logger"
	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	nasLogger "github.com/free5gc/nas/logger"
	ngapLogger "github.com/free5gc/ngap/logger"
	"github.com/free5gc/ocf/communication"
	"github.com/free5gc/ocf/consumer"
	"github.com/free5gc/ocf/context"
	"github.com/free5gc/ocf/eventexposure"
	"github.com/free5gc/ocf/factory"
	"github.com/free5gc/ocf/httpcallback"
	"github.com/free5gc/ocf/location"
	"github.com/free5gc/ocf/logger"
	"github.com/free5gc/ocf/mt"
	"github.com/free5gc/ocf/ngap"
	ngap_message "github.com/free5gc/ocf/ngap/message"
	ngap_service "github.com/free5gc/ocf/ngap/service"
	"github.com/free5gc/ocf/oam"
	"github.com/free5gc/ocf/producer/callback"
	"github.com/free5gc/ocf/util"
	openApiLogger "github.com/free5gc/openapi/logger"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/path_util"
	pathUtilLogger "github.com/free5gc/path_util/logger"
)

type OCF struct{}

type (
	// Config information.
	Config struct {
		ocfcfg string
	}
)

var config Config

var ocfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "ocfcfg",
		Usage: "ocf config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*OCF) GetCliCmd() (flags []cli.Flag) {
	return ocfCLi
}

func (ocf *OCF) Initialize(c *cli.Context) error {
	config = Config{
		ocfcfg: c.String("ocfcfg"),
	}

	if config.ocfcfg != "" {
		if err := factory.InitConfigFactory(config.ocfcfg); err != nil {
			return err
		}
	} else {
		DefaultOcfConfigPath := path_util.Free5gcPath("free5gc/config/ocfcfg.yaml")
		if err := factory.InitConfigFactory(DefaultOcfConfigPath); err != nil {
			return err
		}
	}

	ocf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (ocf *OCF) setLogLevel() {
	if factory.OcfConfig.Logger == nil {
		initLog.Warnln("OCF config without log level setting!!!")
		return
	}

	if factory.OcfConfig.Logger.OCF != nil {
		if factory.OcfConfig.Logger.OCF.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.OcfConfig.Logger.OCF.DebugLevel); err != nil {
				initLog.Warnf("OCF Log level [%s] is invalid, set to [info] level",
					factory.OcfConfig.Logger.OCF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("OCF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Warnln("OCF Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.OcfConfig.Logger.OCF.ReportCaller)
	}

	if factory.OcfConfig.Logger.NAS != nil {
		if factory.OcfConfig.Logger.NAS.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.OcfConfig.Logger.NAS.DebugLevel); err != nil {
				nasLogger.NasLog.Warnf("NAS Log level [%s] is invalid, set to [info] level",
					factory.OcfConfig.Logger.NAS.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				nasLogger.SetLogLevel(level)
			}
		} else {
			nasLogger.NasLog.Warnln("NAS Log level not set. Default set to [info] level")
			nasLogger.SetLogLevel(logrus.InfoLevel)
		}
		nasLogger.SetReportCaller(factory.OcfConfig.Logger.NAS.ReportCaller)
	}

	if factory.OcfConfig.Logger.NGAP != nil {
		if factory.OcfConfig.Logger.NGAP.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.OcfConfig.Logger.NGAP.DebugLevel); err != nil {
				ngapLogger.NgapLog.Warnf("NGAP Log level [%s] is invalid, set to [info] level",
					factory.OcfConfig.Logger.NGAP.DebugLevel)
				ngapLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				ngapLogger.SetLogLevel(level)
			}
		} else {
			ngapLogger.NgapLog.Warnln("NGAP Log level not set. Default set to [info] level")
			ngapLogger.SetLogLevel(logrus.InfoLevel)
		}
		ngapLogger.SetReportCaller(factory.OcfConfig.Logger.NGAP.ReportCaller)
	}

	if factory.OcfConfig.Logger.FSM != nil {
		if factory.OcfConfig.Logger.FSM.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.OcfConfig.Logger.FSM.DebugLevel); err != nil {
				fsmLogger.FsmLog.Warnf("FSM Log level [%s] is invalid, set to [info] level",
					factory.OcfConfig.Logger.FSM.DebugLevel)
				fsmLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				fsmLogger.SetLogLevel(level)
			}
		} else {
			fsmLogger.FsmLog.Warnln("FSM Log level not set. Default set to [info] level")
			fsmLogger.SetLogLevel(logrus.InfoLevel)
		}
		fsmLogger.SetReportCaller(factory.OcfConfig.Logger.FSM.ReportCaller)
	}

	if factory.OcfConfig.Logger.Aper != nil {
		if factory.OcfConfig.Logger.Aper.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.OcfConfig.Logger.Aper.DebugLevel); err != nil {
				aperLogger.AperLog.Warnf("Aper Log level [%s] is invalid, set to [info] level",
					factory.OcfConfig.Logger.Aper.DebugLevel)
				aperLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				aperLogger.SetLogLevel(level)
			}
		} else {
			aperLogger.AperLog.Warnln("Aper Log level not set. Default set to [info] level")
			aperLogger.SetLogLevel(logrus.InfoLevel)
		}
		aperLogger.SetReportCaller(factory.OcfConfig.Logger.Aper.ReportCaller)
	}

	if factory.OcfConfig.Logger.PathUtil != nil {
		if factory.OcfConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.OcfConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.OcfConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.OcfConfig.Logger.PathUtil.ReportCaller)
	}

	if factory.OcfConfig.Logger.OpenApi != nil {
		if factory.OcfConfig.Logger.OpenApi.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.OcfConfig.Logger.OpenApi.DebugLevel); err != nil {
				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.OcfConfig.Logger.OpenApi.DebugLevel)
				openApiLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				openApiLogger.SetLogLevel(level)
			}
		} else {
			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
			openApiLogger.SetLogLevel(logrus.InfoLevel)
		}
		openApiLogger.SetReportCaller(factory.OcfConfig.Logger.OpenApi.ReportCaller)
	}
}

func (ocf *OCF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range ocf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (ocf *OCF) Start() {
	initLog.Infoln("Server started")

	router := logger_util.NewGinWithLogrus(logger.GinLog)
	router.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host",
			"Token", "X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	httpcallback.AddService(router)
	oam.AddService(router)
	for _, serviceName := range factory.OcfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serviceName) {
		case models.ServiceName_NOCF_COMM:
			communication.AddService(router)
		case models.ServiceName_NOCF_EVTS:
			eventexposure.AddService(router)
		case models.ServiceName_NOCF_MT:
			mt.AddService(router)
		case models.ServiceName_NOCF_LOC:
			location.AddService(router)
		}
	}

	self := context.OCF_Self()
	util.InitOcfContext(self)

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	ngapHandler := ngap_service.NGAPHandler{
		HandleMessage:      ngap.Dispatch,
		HandleNotification: ngap.HandleSCTPNotification,
	}
	ngap_service.Run(self.NgapIpList, 38412, ngapHandler)

	// Register to NRF
	var profile models.NfProfile
	if profileTmp, err := consumer.BuildNFInstance(self); err != nil {
		initLog.Error("Build OCF Profile Error")
	} else {
		profile = profileTmp
	}

	if _, nfId, err := consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile); err != nil {
		initLog.Warnf("Send Register NF Instance failed: %+v", err)
	} else {
		self.NfId = nfId
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		ocf.Terminate()
		os.Exit(0)
	}()

	server, err := http2_util.NewServer(addr, util.OcfLogPath, router)

	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: %+v", err)
	}

	serverScheme := factory.OcfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.OcfPemPath, util.OcfKeyPath)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (ocf *OCF) Exec(c *cli.Context) error {
	// OCF.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("ocfcfg"))
	args := ocf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./ocf", args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	stderr, err := command.StderrPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	go func() {
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		if err = command.Start(); err != nil {
			initLog.Errorf("OCF Start error: %+v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}

// Used in OCF planned removal procedure
func (ocf *OCF) Terminate() {
	logger.InitLog.Infof("Terminating OCF...")
	ocfSelf := context.OCF_Self()

	// TODO: forward registered UE contexts to target OCF in the same OCF set if there is one

	// deregister with NRF
	problemDetails, err := consumer.SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.InitLog.Infof("[OCF] Deregister from NRF successfully")
	}

	// send OCF status indication to ran to notify ran that this OCF will be unavailable
	logger.InitLog.Infof("Send OCF Status Indication to Notify RANs due to OCF terminating")
	unavailableGuamiList := ngap_message.BuildUnavailableGUAMIList(ocfSelf.ServedGuamiList)
	ocfSelf.OcfRanPool.Range(func(key, value interface{}) bool {
		ran := value.(*context.OcfRan)
		ngap_message.SendOCFStatusIndication(ran, unavailableGuamiList)
		return true
	})

	ngap_service.Stop()

	callback.SendOcfStatusChangeNotify((string)(models.StatusChange_UNAVAILABLE), ocfSelf.ServedGuamiList)
	logger.InitLog.Infof("OCF terminated")
}

package service

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/antihax/optional"
	"github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"free5gcWithOCF/lib/http2_util"
	"free5gcWithOCF/lib/logger_util"
	"free5gcWithOCF/lib/openapi/Nnrf_NFDiscovery"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/lib/path_util"
	"free5gcWithOCF/src/app"
	"free5gcWithOCF/src/chf/ampolicy"
	"free5gcWithOCF/src/chf/bdtpolicy"
	"free5gcWithOCF/src/chf/consumer"
	"free5gcWithOCF/src/chf/context"
	"free5gcWithOCF/src/chf/factory"
	"free5gcWithOCF/src/chf/httpcallback"
	"free5gcWithOCF/src/chf/logger"
	"free5gcWithOCF/src/chf/oam"
	"free5gcWithOCF/src/chf/policyauthorization"
	"free5gcWithOCF/src/chf/smpolicy"
	"free5gcWithOCF/src/chf/uepolicy"
	"free5gcWithOCF/src/chf/util"
)

type CHF struct{}

type (
	// Config information.
	Config struct {
		chfcfg string
	}
)

var config Config

var chfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gcWithOCFcfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "chfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*CHF) GetCliCmd() (flags []cli.Flag) {
	return chfCLi
}

func (*CHF) Initialize(c *cli.Context) {

	config = Config{
		chfcfg: c.String("chfcfg"),
	}
	if config.chfcfg != "" {
		factory.InitConfigFactory(config.chfcfg)
	} else {
		DefaultChfConfigPath := path_util.Gofree5gcPath("free5gcWithOCF/config/chfcfg.conf")
		factory.InitConfigFactory(DefaultChfConfigPath)
	}

	if app.ContextSelf().Logger.CHF.DebugLevel != "" {
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.CHF.DebugLevel)
		if err != nil {
			initLog.Warnf("Log level [%s] is not valid, set to [info] level", app.ContextSelf().Logger.CHF.DebugLevel)
			logger.SetLogLevel(logrus.InfoLevel)
		} else {
			logger.SetLogLevel(level)
			initLog.Infof("Log level is set to [%s] level", level)
		}
	} else {
		initLog.Infoln("Log level is default set to [info] level")
		logger.SetLogLevel(logrus.InfoLevel)
	}

	logger.SetReportCaller(app.ContextSelf().Logger.CHF.ReportCaller)
}

func (chf *CHF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range chf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (chf *CHF) Start() {
	initLog.Infoln("Server started")
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	bdtpolicy.AddService(router)
	smpolicy.AddService(router)
	ampolicy.AddService(router)
	uepolicy.AddService(router)
	policyauthorization.AddService(router)
	httpcallback.AddService(router)
	oam.AddService(router)

	router.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "User-Agent",
			"Referrer", "Host", "Token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	self := context.CHF_Self()
	util.InitchfContext(self)

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	profile, err := consumer.BuildNFInstance(self)
	if err != nil {
		initLog.Error("Build CHF Profile Error")
	}
	_, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
	if err != nil {
		initLog.Errorf("CHF register to NRF Error[%s]", err.Error())
	}

	// subscribe to all Amfs' status change
	amfInfos := consumer.SearchAvailableAMFs(self.NrfUri, models.ServiceName_NAMF_COMM)
	for _, amfInfo := range amfInfos {
		guamiList := util.GetNotSubscribedGuamis(amfInfo.GuamiList)
		if len(guamiList) == 0 {
			continue
		}
		var problemDetails *models.ProblemDetails
		problemDetails, err = consumer.AmfStatusChangeSubscribe(amfInfo)
		if problemDetails != nil {
			logger.InitLog.Warnf("AMF status subscribe Failed[%+v]", problemDetails)
		} else if err != nil {
			logger.InitLog.Warnf("AMF status subscribe Error[%+v]", err)
		}
	}

	// TODO: subscribe NRF NFstatus

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NUDR_DR}),
	}
	resp, err := consumer.SendSearchNFInstances(self.NrfUri, models.NfType_UDR, models.NfType_CHF, param)
	for _, nfProfile := range resp.NfInstances {
		udruri := util.SearchNFServiceUri(nfProfile, models.ServiceName_NUDR_DR, models.NfServiceStatus_REGISTERED)
		if udruri != "" {
			self.SetDefaultUdrURI(udruri)
			break
		}
	}
	if err != nil {
		initLog.Errorln(err)
	}
	server, err := http2_util.NewServer(addr, util.CHF_LOG_PATH, router)
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: +%v", err)
	}

	serverScheme := factory.ChfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.CHF_PEM_PATH, util.CHF_KEY_PATH)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (chf *CHF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("chfcfg"))
	args := chf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./chf", args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(4)
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
		fmt.Println("CHF log start")
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		fmt.Println("CHF start")
		if err = command.Start(); err != nil {
			fmt.Printf("command.Start() error: %v", err)
		}
		fmt.Println("CHF end")
		wg.Done()
	}()

	wg.Wait()

	return err
}

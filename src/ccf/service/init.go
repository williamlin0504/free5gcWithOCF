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

	"free5gc/lib/http2_util"
	"free5gc/lib/logger_util"
	"free5gc/lib/openapi/Nnrf_NFDiscovery"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/app"
	"free5gc/src/ccf/ampolicy"
	"free5gc/src/ccf/bdtpolicy"
	"free5gc/src/ccf/consumer"
	"free5gc/src/ccf/context"
	"free5gc/src/ccf/factory"
	"free5gc/src/ccf/httpcallback"
	"free5gc/src/ccf/logger"
	"free5gc/src/ccf/oam"
	"free5gc/src/ccf/policyauthorization"
	"free5gc/src/ccf/smpolicy"
	"free5gc/src/ccf/uepolicy"
	"free5gc/src/ccf/util"
)

type CCF struct{}

type (
	// Config information.
	Config struct {
		ccfcfg string
	}
)

var config Config

var ccfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "ccfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*CCF) GetCliCmd() (flags []cli.Flag) {
	return ccfCLi
}

func (*CCF) Initialize(c *cli.Context) {

	config = Config{
		ccfcfg: c.String("ccfcfg"),
	}
	if config.ccfcfg != "" {
		factory.InitConfigFactory(config.ccfcfg)
	} else {
		DefaultCcfConfigPath := path_util.Gofree5gcPath("free5gc/config/ccfcfg.conf")
		factory.InitConfigFactory(DefaultCcfConfigPath)
	}

	if app.ContextSelf().Logger.CCF.DebugLevel != "" {
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.CCF.DebugLevel)
		if err != nil {
			initLog.Warnf("Log level [%s] is not valid, set to [info] level", app.ContextSelf().Logger.CCF.DebugLevel)
			logger.SetLogLevel(logrus.InfoLevel)
		} else {
			logger.SetLogLevel(level)
			initLog.Infof("Log level is set to [%s] level", level)
		}
	} else {
		initLog.Infoln("Log level is default set to [info] level")
		logger.SetLogLevel(logrus.InfoLevel)
	}

	logger.SetReportCaller(app.ContextSelf().Logger.CCF.ReportCaller)
}

func (ccf *CCF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range ccf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (ccf *CCF) Start() {
	initLog.Infoln("Server started")
	fmt.Fprint(w, "CCF Started Ready to Start Session...")
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

	self := context.CCF_Self()
	util.InitccfContext(self)

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	profile, err := consumer.BuildNFInstance(self)
	if err != nil {
		initLog.Error("Build CCF Profile Error")
	}
	_, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
	if err != nil {
		initLog.Errorf("CCF register to NRF Error[%s]", err.Error())
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
	resp, err := consumer.SendSearchNFInstances(self.NrfUri, models.NfType_UDR, models.NfType_CCF, param)
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
	server, err := http2_util.NewServer(addr, util.CCF_LOG_PATH, router)
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: +%v", err)
	}

	serverScheme := factory.CcfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.CCF_PEM_PATH, util.CCF_KEY_PATH)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (ccf *CCF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("ccfcfg"))
	args := ccf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./ccf", args...)

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
		fmt.Println("CCF log start")
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		fmt.Println("CCF start")
		if err = command.Start(); err != nil {
			fmt.Printf("command.Start() error: %v", err)
		}
		fmt.Println("CCF end")
		wg.Done()
	}()

	wg.Wait()

	return err
}

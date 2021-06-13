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

	"free5gc/src/app"
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
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*OCF) GetCliCmd() (flags []cli.Flag) {
	return ocfCLi
}

func (*OCF) Initialize(c *cli.Context) {

	config = Config{
		ocfcfg: c.String("ocfcfg"),
	}
	if config.ocfcfg != "" {
		factory.InitConfigFactory(config.ocfcfg)
	} else {
		DefaultOcfConfigPath := path_util.Gofree5gcPath("free5gc/config/ocfcfg.conf")
		factory.InitConfigFactory(DefaultOcfConfigPath)
	}

	if app.ContextSelf().Logger.OCF.DebugLevel != "" {
		level, err := logrus.ParseLevel(app.ContextSelf().Logger.OCF.DebugLevel)
		if err != nil {
			initLog.Warnf("Log level [%s] is not valid, set to [info] level", app.ContextSelf().Logger.OCF.DebugLevel)
			logger.SetLogLevel(logrus.InfoLevel)
		} else {
			logger.SetLogLevel(level)
			initLog.Infof("Log level is set to [%s] level", level)
		}
	} else {
		initLog.Infoln("Log level is default set to [info] level")
		logger.SetLogLevel(logrus.InfoLevel)
	}

	logger.SetReportCaller(app.ContextSelf().Logger.OCF.ReportCaller)
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

	self := context.OCF_Self()
	util.InitocfContext(self)

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	profile, err := consumer.BuildNFInstance(self)
	if err != nil {
		initLog.Error("Build OCF Profile Error")
	}
	_, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
	if err != nil {
		initLog.Errorf("OCF register to NRF Error[%s]", err.Error())
	}

	// TODO: subscribe NRF NFstatus

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NUDR_DR}),
	}
	resp, err := consumer.SendSearchNFInstances(self.NrfUri, models.NfType_UDR, models.NfType_OCF, param)
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
	server, err := http2_util.NewServer(addr, util.OCF_LOG_PATH, router)
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: +%v", err)
	}

	serverScheme := factory.OcfConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.OCF_PEM_PATH, util.OCF_KEY_PATH)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (ocf *OCF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("ocfcfg"))
	args := ocf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./ocf", args...)

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
		fmt.Println("OCF log start")
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		fmt.Println("OCF start")
		if err = command.Start(); err != nil {
			fmt.Printf("command.Start() error: %v", err)
		}
		fmt.Println("OCF end")
		wg.Done()
	}()

	wg.Wait()

	return err
}

package service

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"free5gcWithOCF/lib/path_util"
	"free5gcWithOCF/src/app"
	"free5gcWithOCF/src/ccf/factory"
	ike_service "free5gcWithOCF/src/ccf/ike/service"
	"free5gcWithOCF/src/ccf/logger"
	ngap_service "free5gcWithOCF/src/ccf/ngap/service"
	nwucp_service "free5gcWithOCF/src/ccf/nwucp/service"
	nwuup_service "free5gcWithOCF/src/ccf/nwuup/service"
	"free5gcWithOCF/src/ccf/util"
)

type CCF struct{}

type (
	// Config information.
	Config struct {
		ocfcfg string
	}
)

var config Config

var ocfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gcWithOCFcfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "ocfcfg",
		Usage: "ccf config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*CCF) GetCliCmd() (flags []cli.Flag) {
	return ocfCLi
}

func (*CCF) Initialize(c *cli.Context) {

	config = Config{
		ocfcfg: c.String("ocfcfg"),
	}

	if config.ocfcfg != "" {
		factory.InitConfigFactory(config.ocfcfg)
	} else {
		DefaultCcfConfigPath := path_util.Gofree5gcWithOCFPath("free5gcWithOCF/config/ocfcfg.conf")
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

	if !util.InitCCFContext() {
		initLog.Error("Initicating context failed")
		return
	}

	wg := sync.WaitGroup{}

	// NGAP
	if err := ngap_service.Run(); err != nil {
		initLog.Errorf("Start NGAP service failed: %+v", err)
		return
	} else {
		initLog.Info("NGAP service running.")
		wg.Add(1)
	}

	// Relay listeners
	// Control plane
	if err := nwucp_service.Run(); err != nil {
		initLog.Errorf("Listen NWu control plane traffic failed: %+v", err)
	} else {
		initLog.Info("NAS TCP server successfully started.")
		wg.Add(1)
	}
	// User plane
	if err := nwuup_service.Run(); err != nil {
		initLog.Errorf("Listen NWu user plane traffic failed: %+v", err)
		return
	} else {
		initLog.Info("Listening NWu user plane traffic")
		wg.Add(1)
	}

	// IKE
	if err := ike_service.Run(); err != nil {
		initLog.Errorf("Start IKE service failed: %+v", err)
		return
	} else {
		initLog.Info("IKE service running.")
		wg.Add(1)
	}

	initLog.Info("CCF running...")

	wg.Wait()

}

func (ccf *CCF) Exec(c *cli.Context) error {

	//CCF.Initialize(cfgPath, c)

	initLog.Traceln("args:", c.String("ocfcfg"))
	args := ccf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./ccf", args...)

	wg := sync.WaitGroup{}
	wg.Add(3)

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
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
		if errCom := command.Start(); errCom != nil {
			initLog.Errorf("CCF start error: %v", errCom)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}

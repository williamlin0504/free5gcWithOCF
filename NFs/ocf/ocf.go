package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/ocf/logger"
	"github.com/free5gc/ocf/service"
	"github.com/free5gc/version"
)

var OCF = &service.OCF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "ocf"
	appLog.Infoln(app.Name)
	appLog.Infoln("OCF version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -ocfcfg ocf configuration file"
	app.Action = action
	app.Flags = OCF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("OCF Run error: %v", err)
		return
	}
}

func action(c *cli.Context) error {
	if err := OCF.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	OCF.Start()

	return nil
}

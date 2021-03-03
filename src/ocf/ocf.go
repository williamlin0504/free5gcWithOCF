package main

import (
	"free5gcWithOCFWithOCF/src/app"
	"free5gcWithOCFWithOCF/src/ocf/logger"
	"free5gcWithOCFWithOCF/src/ocf/service"
	"free5gcWithOCFWithOCF/src/ocf/version"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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
	app.Usage = "-free5gcWithOCFWithOCFcfg common configuration file -ocfcfg ocf configuration file"
	app.Action = action
	app.Flags = OCF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorf("OCF Run Error: %v", err)
	}
}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gcWithOCFWithOCFcfg"))
	OCF.Initialize(c)
	OCF.Start()
}

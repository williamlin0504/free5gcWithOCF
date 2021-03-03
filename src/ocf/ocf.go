package main

import (
	"free5gcWithOCF/src/app"
	"free5gcWithOCF/src/ocf/logger"
	"free5gcWithOCF/src/ocf/service"
	"free5gcWithOCF/src/ocf/version"
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
	app.Usage = "-free5gcWithOCFcfg common configuration file -ocfcfg ocf configuration file"
	app.Action = action
	app.Flags = OCF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorf("OCF Run Error: %v", err)
	}
}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gcWithOCFcfg"))
	OCF.Initialize(c)
	OCF.Start()
}

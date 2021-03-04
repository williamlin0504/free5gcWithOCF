package main

import (
	"free5gcWithOCF/src/amf/logger"
	"free5gcWithOCF/src/amf/service"
	"free5gcWithOCF/src/amf/version"
	"free5gcWithOCF/src/app"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var AMF = &service.AMF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "amf"
	appLog.Infoln(app.Name)
	appLog.Infoln("AMF version: ", version.GetVersion())
	app.Usage = "-free5gchfg common configuration file -amfcfg amf configuration file"
	app.Action = action
	app.Flags = AMF.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorf("AMF Run error: %v", err)
	}
}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gchfg"))
	AMF.Initialize(c)
	AMF.Start()
}

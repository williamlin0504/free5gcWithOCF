package main

import (
	"free5gcWithOCF/src/app"
	"free5gcWithOCF/webconsole/backend/logger"
	"free5gcWithOCF/webconsole/backend/webui_service"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var WEBUI = &webui_service.WEBUI{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "webui"
	appLog.Infoln(app.Name)
	app.Usage = "-free5gchfg common configuration file -webuicfg webui configuration file"
	app.Action = action
	app.Flags = WEBUI.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Warnf("Error args: %v", err)
	}
}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gchfg"))
	WEBUI.Initialize(c)
	WEBUI.Start()
}

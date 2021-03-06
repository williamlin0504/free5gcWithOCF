package main

import (
	"fmt"
	"free5gcWithOCF/src/app"
	"free5gcWithOCF/src/udm/logger"
	"free5gcWithOCF/src/udm/service"
	"free5gcWithOCF/src/udm/version"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var UDM = &service.UDM{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "udm"
	fmt.Print(app.Name, "\n")
	appLog.Infoln("UDM version: ", version.GetVersion())
	app.Usage = "-free5gcWithOCFcfg common configuration file -udmcfg udm configuration file"
	app.Action = action
	app.Flags = UDM.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("UDM Run error: %v", err)
	}

	// appLog.Infoln(app.Name)

}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gcWithOCFcfg"))
	UDM.Initialize(c)
	UDM.Start()
}

/*
 *
 * AUSF Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"fmt"
	"free5gcWithOCF/src/app"
	"free5gcWithOCF/src/ausf/logger"
	"free5gcWithOCF/src/ausf/service"
	"free5gcWithOCF/src/ausf/version"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var AUSF = &service.AUSF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "ausf"
	fmt.Print(app.Name, "\n")
	appLog.Infoln("AUSF version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -ausfcfg ausf configuration file"
	app.Action = action
	app.Flags = AUSF.GetCliCmd()

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gccfg"))
	AUSF.Initialize(c)
	AUSF.Start()
}

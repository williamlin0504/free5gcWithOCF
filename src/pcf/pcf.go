/*
 * Nccf_BDTPolicyControl Service API
 *
 * The Nccf_BDTPolicyControl Service is used by an NF service consumer to
 * retrieve background data transfer policies from the ccf and to update the ccf with
 * the background data transfer policy selected by the NF service consumer.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"fmt"
	"free5gc/src/app"
	"free5gc/src/ccf/logger"
	"free5gc/src/ccf/service"
	"free5gc/src/ccf/version"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var ccf = &service.ccf{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "ccf"
	fmt.Print(app.Name, "\n")
	appLog.Infoln("ccf version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -ccfcfg ccf configuration file"
	app.Action = action
	app.Flags = ccf.GetCliCmd()

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("ccf Run err: %v", err)
	}

}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String("free5gccfg"))
	ccf.Initialize(c)
	ccf.Start()
}

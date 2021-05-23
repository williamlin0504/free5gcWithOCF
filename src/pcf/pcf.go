/*
 * Npcf_BDTPolicyControl Service API
 *
 * The Npcf_BDTPolicyControl Service is used by an NF service consumer to
 * retrieve background data transfer policies from the PCF and to update the PCF with
 * the background data transfer policy selected by the NF service consumer.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"fmt"
	" free5gcWithOCF/src/app"
	" free5gcWithOCF/src/pcf/logger"
	" free5gcWithOCF/src/pcf/service"
	" free5gcWithOCF/src/pcf/version"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var PCF = &service.PCF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "pcf"
	fmt.Print(app.Name, "\n")
	appLog.Infoln("PCF version: ", version.GetVersion())
	app.Usage = "- free5gccfg common configuration file -pcfcfg pcf configuration file"
	app.Action = action
	app.Flags = PCF.GetCliCmd()

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("PCF Run err: %v", err)
	}

}

func action(c *cli.Context) {
	app.AppInitializeWillInitialize(c.String(" free5gccfg"))
	PCF.Initialize(c)
	PCF.Start()
}

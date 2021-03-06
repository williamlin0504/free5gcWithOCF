//+build debug

package util

import (
	"free5gcWithOCF/lib/path_util"
)

var UdrLogPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/udrsslkey.log")
var UdrPemPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/_debug.pem")
var UdrKeyPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/_debug.key")
var DefaultUdrConfigPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/config/udrcfg.conf")

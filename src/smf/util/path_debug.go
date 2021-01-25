//+build debug

package util

import (
	"free5gcWithOCF/lib/path_util"
)

var SmfLogPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/smfsslkey.log")
var SmfPemPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/_debug.pem")
var SmfKeyPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/_debug.key")
var DefaultSmfConfigPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/config/smfcfg.conf")

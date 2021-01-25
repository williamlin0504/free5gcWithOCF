//+build debug

package util

import (
	"free5gcWithOCF/lib/path_util"
)

var AusfLogPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/ausfsslkey.log")
var AusfPemPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/ausf.pem")
var AusfKeyPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/ausf.key")
var DefaultAusfConfigPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/config/ausfcfg.conf")

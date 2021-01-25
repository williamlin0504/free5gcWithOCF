//+build !debug

package util

import (
	"free5gcWithOCF/lib/path_util"
)

var OcfLogPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/ocfsslkey.log")
var OcfPemPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/ocf.pem")
var OcfKeyPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/ocf.key")
var DefaultOcfConfigPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/config/ocfcfg.conf")

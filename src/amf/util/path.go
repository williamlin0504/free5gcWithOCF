//+build !debug

package util

import (
	"free5gcWithOCF/lib/path_util"
)

var AmfLogPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/amfsslkey.log")
var AmfPemPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/amf.pem")
var AmfKeyPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/amf.key")
var DefaultAmfConfigPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/config/amfcfg.conf")

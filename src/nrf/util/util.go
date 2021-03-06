//+build !debug

package util

import (
	"free5gcWithOCF/lib/path_util"
)

// Path of HTTP2 key and log file

var NrfLogPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/nrfsslkey.log")
var NrfPemPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/nrf.pem")
var NrfKeyPath = path_util.Gofree5gcWithOCFPath("free5gcWithOCF/support/TLS/nrf.key")

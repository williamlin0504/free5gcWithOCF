//+build debug

package util

import (
	"free5gcWithOCF/lib/path_util"
)

// Path of HTTP2 key and log file

var NrfLogPath = path_util.Gofree5gcPath("free5gcWithOCF/nrfsslkey.log")
var NrfPemPath = path_util.Gofree5gcPath("free5gcWithOCF/support/TLS/_debug.pem")
var NrfKeyPath = path_util.Gofree5gcPath("free5gcWithOCF/support/TLS/_debug.key")

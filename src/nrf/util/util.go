//+build !debug

package util

import (
	" free5gcWithOCF/lib/path_util"
)

// Path of HTTP2 key and log file

var NrfLogPath = path_util.Go free5gcPath(" free5gcWithOCF/nrfsslkey.log")
var NrfPemPath = path_util.Go free5gcPath(" free5gcWithOCF/support/TLS/nrf.pem")
var NrfKeyPath = path_util.Go free5gcPath(" free5gcWithOCF/support/TLS/nrf.key")

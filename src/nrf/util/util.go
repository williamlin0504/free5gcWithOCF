//+build !debug

package util

import (
	" free5gc/lib/path_util"
)

// Path of HTTP2 key and log file

var NrfLogPath = path_util.Go free5gcPath(" free5gckey.log")
var NrfPemPath = path_util.Go free5gcPath(" free5gct/TLS/nrf.pem")
var NrfKeyPath = path_util.Go free5gcPath(" free5gct/TLS/nrf.key")

//+build !debug

package util

import (
	" free5gc/lib/path_util"
)

var SmfLogPath = path_util.Go free5gcPath(" free5gckey.log")
var SmfPemPath = path_util.Go free5gcPath(" free5gct/TLS/smf.pem")
var SmfKeyPath = path_util.Go free5gcPath(" free5gct/TLS/smf.key")
var DefaultSmfConfigPath = path_util.Go free5gcPath(" free5gc/smfcfg.conf")

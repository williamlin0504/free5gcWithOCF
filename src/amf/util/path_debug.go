//+build debug

package util

import (
	" free5gc/lib/path_util"
)

var AmfLogPath = path_util.Go free5gcPath(" free5gckey.log")
var AmfPemPath = path_util.Go free5gcPath(" free5gct/TLS/_debug.pem")
var AmfKeyPath = path_util.Go free5gcPath(" free5gct/TLS/_debug.key")
var DefaultAmfConfigPath = path_util.Go free5gcPath(" free5gc/amfcfg.conf")

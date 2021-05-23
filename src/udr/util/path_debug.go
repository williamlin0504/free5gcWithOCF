//+build debug

package util

import (
	" free5gc/lib/path_util"
)

var UdrLogPath = path_util.Go free5gcPath(" free5gckey.log")
var UdrPemPath = path_util.Go free5gcPath(" free5gct/TLS/_debug.pem")
var UdrKeyPath = path_util.Go free5gcPath(" free5gct/TLS/_debug.key")
var DefaultUdrConfigPath = path_util.Go free5gcPath(" free5gc/udrcfg.conf")

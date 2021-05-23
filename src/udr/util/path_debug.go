//+build debug

package util

import (
	" free5gcWithOCF/lib/path_util"
)

var UdrLogPath = path_util.Go free5gcWithOCFPath(" free5gcWithOCF/udrsslkey.log")
var UdrPemPath = path_util.Go free5gcWithOCFPath(" free5gcWithOCF/support/TLS/_debug.pem")
var UdrKeyPath = path_util.Go free5gcWithOCFPath(" free5gcWithOCF/support/TLS/_debug.key")
var DefaultUdrConfigPath = path_util.Go free5gcWithOCFPath(" free5gcWithOCF/config/udrcfg.conf")

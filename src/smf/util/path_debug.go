//+build debug

package util

import (
	" free5gcWithOCF/lib/path_util"
)

var SmfLogPath = path_util.Go free5gcPath(" free5gcWithOCF/smfsslkey.log")
var SmfPemPath = path_util.Go free5gcPath(" free5gcWithOCF/support/TLS/_debug.pem")
var SmfKeyPath = path_util.Go free5gcPath(" free5gcWithOCF/support/TLS/_debug.key")
var DefaultSmfConfigPath = path_util.Go free5gcPath(" free5gcWithOCF/config/smfcfg.conf")

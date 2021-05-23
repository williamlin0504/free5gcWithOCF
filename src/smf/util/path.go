//+build !debug

package util

import (
	" free5gcWithOCF/lib/path_util"
)

var SmfLogPath = path_util.Go free5gcWithOCFPath(" free5gcWithOCF/smfsslkey.log")
var SmfPemPath = path_util.Go free5gcWithOCFPath(" free5gcWithOCF/support/TLS/smf.pem")
var SmfKeyPath = path_util.Go free5gcWithOCFPath(" free5gcWithOCF/support/TLS/smf.key")
var DefaultSmfConfigPath = path_util.Go free5gcWithOCFPath(" free5gcWithOCF/config/smfcfg.conf")

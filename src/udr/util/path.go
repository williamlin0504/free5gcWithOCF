//+build !debug

package util

import (
	" free5gcWithOCF/lib/path_util"
)

var UdrLogPath = path_util.Go free5gcPath(" free5gcWithOCF/udrsslkey.log")
var UdrPemPath = path_util.Go free5gcPath(" free5gcWithOCF/support/TLS/udr.pem")
var UdrKeyPath = path_util.Go free5gcPath(" free5gcWithOCF/support/TLS/udr.key")
var DefaultUdrConfigPath = path_util.Go free5gcPath(" free5gcWithOCF/config/udrcfg.conf")

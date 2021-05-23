//+build debug

package util

import (
	" free5gc/lib/path_util"
)

var AusfLogPath = path_util.Go free5gcPath(" free5gclkey.log")
var AusfPemPath = path_util.Go free5gcPath(" free5gct/TLS/ausf.pem")
var AusfKeyPath = path_util.Go free5gcPath(" free5gct/TLS/ausf.key")
var DefaultAusfConfigPath = path_util.Go free5gcPath(" free5gc/ausfcfg.conf")

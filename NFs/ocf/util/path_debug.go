//+build debug

package util

import (
	"github.com/free5gc/path_util"
)

var (
	OcfLogPath           = path_util.Free5gcPath("free5gc/ocfsslkey.log")
	OcfPemPath           = path_util.Free5gcPath("free5gc/support/TLS/_debug.pem")
	OcfKeyPath           = path_util.Free5gcPath("free5gc/support/TLS/_debug.key")
	DefaultOcfConfigPath = path_util.Free5gcPath("free5gc/config/ocfcfg.yaml")
)

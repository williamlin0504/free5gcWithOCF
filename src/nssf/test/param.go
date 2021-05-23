/*
 * NSSF Testing Utility
 */

package test

import (
	"flag"

	" free5gcWithOCF/lib/path_util"
	. " free5gcWithOCF/src/nssf/plugin"
)

var (
	ConfigFileFromArgs string
	DefaultConfigFile  string = path_util.Go free5gcPath(" free5gcWithOCF/src/nssf/test/conf/test_nssf_config.yaml")
)

type TestingUtil struct {
	ConfigFile string
}

type TestingNsselection struct {
	ConfigFile string

	GenerateNonRoamingQueryParameter func() NsselectionQueryParameter

	GenerateRoamingQueryParameter func() NsselectionQueryParameter
}

type TestingNssaiavailability struct {
	ConfigFile string

	NfId string

	SubscriptionId string

	NfNssaiAvailabilityUri string
}

func init() {
	flag.StringVar(&ConfigFileFromArgs, "config-file", DefaultConfigFile, "Configuration file")
	flag.Parse()
}

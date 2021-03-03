/*
 * OCF Configuration Factory
 */

package factory

import "free5gcWithOCF/src/ocf/context"

type Config struct {
	Info *Info `yaml:"info"`

	Configuration *Configuration `yaml:"configuration"`
}

type Info struct {
	Version string `yaml:"version,omitempty"`

	Description string `yaml:"description,omitempty"`
}

type Configuration struct {
	OCFInfo          context.OCFNFInfo          `yaml:"OCFInformation"`
	AMFSCTPAddresses []context.AMFSCTPAddresses `yaml:"AMFSCTPAddresses"`

	IKEBindAddr          string `yaml:"IKEBindAddress"`
	IPSecGatewayAddr     string `yaml:"IPSecInterfaceAddress"`
	GTPBindAddr          string `yaml:"GTPBindAddress"`
	TCPPort              uint16 `yaml:"NASTCPPort"`
	FQDN                 string `yaml:"FQDN"`                 // e.g. ocf.free5gcWithOCF.org
	PrivateKey           string `yaml:"PrivateKey"`           // file path
	CertificateAuthority string `yaml:"CertificateAuthority"` // file path
	Certificate          string `yaml:"Certificate"`          // file path
	UEIPAddressRange     string `yaml:"UEIPAddressRange"`     // e.g. 10.0.1.0/24
	InterfaceMark        uint32 `yaml:"IPSecInterfaceMark"`   // must != 0, if not specified, random one
}

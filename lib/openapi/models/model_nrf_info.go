/*
 * NRF NFManagement Service
 *
 * NRF NFManagement Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type NrfInfo struct {
	ServedUdrInfo  map[string]UdrInfo  `json:"servedUdrInfo,omitempty" yaml:"servedUdrInfo" bson:"servedUdrInfo" mapstructure:"ServedUdrInfo"`
	ServedUdmInfo  map[string]UdmInfo  `json:"servedUdmInfo,omitempty" yaml:"servedUdmInfo" bson:"servedUdmInfo" mapstructure:"ServedUdmInfo"`
	ServedAusfInfo map[string]AusfInfo `json:"servedAusfInfo,omitempty" yaml:"servedAusfInfo" bson:"servedAusfInfo" mapstructure:"ServedAusfInfo"`
	ServedAmfInfo  map[string]AmfInfo  `json:"servedAmfInfo,omitempty" yaml:"servedAmfInfo" bson:"servedAmfInfo" mapstructure:"ServedAmfInfo"`
	ServedSmfInfo  map[string]SmfInfo  `json:"servedSmfInfo,omitempty" yaml:"servedSmfInfo" bson:"servedSmfInfo" mapstructure:"ServedSmfInfo"`
	ServedUpfInfo  map[string]UpfInfo  `json:"servedUpfInfo,omitempty" yaml:"servedUpfInfo" bson:"servedUpfInfo" mapstructure:"ServedUpfInfo"`
	ServedccfInfo  map[string]ccfInfo  `json:"servedccfInfo,omitempty" yaml:"servedccfInfo" bson:"servedccfInfo" mapstructure:"ServedccfInfo"`
	ServedBsfInfo  map[string]BsfInfo  `json:"servedBsfInfo,omitempty" yaml:"servedBsfInfo" bson:"servedBsfInfo" mapstructure:"ServedBsfInfo"`
	ServedChfInfo  map[string]ChfInfo  `json:"servedChfInfo,omitempty" yaml:"servedChfInfo" bson:"servedChfInfo" mapstructure:"ServedChfInfo"`
}

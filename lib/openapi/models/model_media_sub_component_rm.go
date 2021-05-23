/*
 * Nchf_ConvergedChargingRelease Service API
 *
 * This is the Policy Authorization Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// This data type is defined in the same way as the MediaSubComponent data type, but with the OpenAPI nullable property set to true. Removable attributes marBwDland marBwUl are defined with the corresponding removable data type.
type MediaSubComponentRm struct {
	EthfDescs []EthFlowDescription `json:"ethfDescs,omitempty" yaml:"ethfDescs" bson:"ethfDescs" mapstructure:"EthfDescs"`
	FNum      int32                `json:"fNum" yaml:"fNum" bson:"fNum" mapstructure:"FNum"`
	FDescs    []string             `json:"fDescs,omitempty" yaml:"fDescs" bson:"fDescs" mapstructure:"FDescs"`
	FStatus   FlowStatus           `json:"fStatus,omitempty" yaml:"fStatus" bson:"fStatus" mapstructure:"FStatus"`
	MarBwDl   string               `json:"marBwDl,omitempty" yaml:"marBwDl" bson:"marBwDl" mapstructure:"MarBwDl"`
	MarBwUl   string               `json:"marBwUl,omitempty" yaml:"marBwUl" bson:"marBwUl" mapstructure:"MarBwUl"`
	// this data type is defined in the same way as the TosTrafficClass data type, but with the OpenAPI nullable property set to true
	TosTrCl   string    `json:"tosTrCl,omitempty" yaml:"tosTrCl" bson:"tosTrCl" mapstructure:"TosTrCl"`
	FlowUsage FlowUsage `json:"flowUsage,omitempty" yaml:"flowUsage" bson:"flowUsage" mapstructure:"FlowUsage"`
}

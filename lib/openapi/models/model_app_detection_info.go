/*
 * Nccf_SMPolicyControl
 *
 * Session Management Policy Control Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type AppDetectionInfo struct {
	// A reference to the application detection filter configured at the UPF
	AppId string `json:"appId" yaml:"appId" bson:"appId" mapstructure:"AppId"`
	// Identifier sent by the SMF in order to allow correlation of application Start and Stop events to the specific service data flow description, if service data flow descriptions are deducible.
	InstanceId string `json:"instanceId,omitempty" yaml:"instanceId" bson:"instanceId" mapstructure:"InstanceId"`
	// Contains the detected service data flow descriptions if they are deducible.
	SdfDescriptions []FlowInformation `json:"sdfDescriptions,omitempty" yaml:"sdfDescriptions" bson:"sdfDescriptions" mapstructure:"SdfDescriptions"`
}

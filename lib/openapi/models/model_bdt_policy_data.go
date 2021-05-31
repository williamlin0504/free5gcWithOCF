/*
 * Npcf_BDTPolicyControl Service API
 *
 * The Npcf_BDTPolicyControl Service is used by an NF service consumer to retrieve background data transfer policies from the pcf and to update the pcf with the background data transfer policy selected by the NF service consumer.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// Describes the authorization data of an Individual BDT policy resource.
type BdtPolicyData struct {
	// string identifying a BDT Reference ID as defined in subclause 5.3.3 of 3GPP TS 29.154.
	BdtRefId string `json:"bdtRefId" yaml:"bdtRefId" bson:"bdtRefId" mapstructure:"BdtRefId"`
	// Contains transfer policies.
	TransfPolicies []TransferPolicy `json:"transfPolicies" yaml:"transfPolicies" bson:"transfPolicies" mapstructure:"TransfPolicies"`
	// Contains an identity of the selected transfer policy.
	SelTransPolicyId int32  `json:"selTransPolicyId,omitempty" yaml:"selTransPolicyId" bson:"selTransPolicyId" mapstructure:"SelTransPolicyId"`
	SuppFeat         string `json:"suppFeat,omitempty" yaml:"suppFeat" bson:"suppFeat" mapstructure:"SuppFeat"`
}

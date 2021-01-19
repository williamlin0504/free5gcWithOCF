/*
 * Npcf_PolicyAuthorization Service API
 *
 * This is the Policy Authorization Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// Identifies the modifications to an Individual Application Session Context and may include the modifications to the sub-resource Events Subscription.
type AppSessionContextUpdateData struct {
	// Contains an AF application identifier.
	AfAppId   string                  `json:"afAppId,omitempty" yaml:"afAppId" bson:"afAppId" mapstructure:"AfAppId"`
	AfRoutReq *AfRoutingRequirementRm `json:"afRoutReq,omitempty" yaml:"afRoutReq" bson:"afRoutReq" mapstructure:"AfRoutReq"`
	// Contains an identity of an application service provider.
	AspId string `json:"aspId,omitempty" yaml:"aspId" bson:"aspId" mapstructure:"AspId"`
	// string identifying a BDT Reference ID as defined in subclause 5.3.3 of 3GPP TS 29.154.
	BdtRefId      string                      `json:"bdtRefId,omitempty" yaml:"bdtRefId" bson:"bdtRefId" mapstructure:"BdtRefId"`
	EvSubsc       *EventsSubscReqDataRm       `json:"evSubsc,omitempty" yaml:"evSubsc" bson:"evSubsc" mapstructure:"EvSubsc"`
	MedComponents map[string]MediaComponentRm `json:"medComponents,omitempty" yaml:"medComponents" bson:"medComponents" mapstructure:"MedComponents"`
	// indication of MPS service request
	MpsId string `json:"mpsId,omitempty" yaml:"mpsId" bson:"mpsId" mapstructure:"MpsId"`
	// Contains an identity of a sponsor.
	SponId     string           `json:"sponId,omitempty" yaml:"sponId" bson:"sponId" mapstructure:"SponId"`
	SponStatus SponsoringStatus `json:"sponStatus,omitempty" yaml:"sponStatus" bson:"sponStatus" mapstructure:"SponStatus"`
}

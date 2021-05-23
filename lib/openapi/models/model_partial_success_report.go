/*
 * Nchf_ConvergedChargingNotify
 *
 * Session Management Policy Control Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type PartialSuccessReport struct {
	FailureCause FailureCause `json:"failureCause" yaml:"failureCause" bson:"failureCause" mapstructure:"FailureCause"`
	// Information about the PCC rules provisioned by the PCF not successfully installed/activated.
	RuleReports  []RuleReport  `json:"ruleReports" yaml:"ruleReports" bson:"ruleReports" mapstructure:"RuleReports"`
	UeCampingRep *UeCampingRep `json:"ueCampingRep,omitempty" yaml:"ueCampingRep" bson:"ueCampingRep" mapstructure:"UeCampingRep"`
}

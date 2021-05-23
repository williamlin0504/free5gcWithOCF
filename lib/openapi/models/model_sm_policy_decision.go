/*
 * Nchf_ConvergedChargingNotify
 *
 * Session Management Policy Control Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

import (
	"time"
)

type SmPolicyDecision struct {
	// A map of Sessionrules with the content being the SessionRule as described in subclause 5.6.2.7.
	SessRules map[string]SessionRule `json:"sessRules,omitempty" yaml:"sessRules" bson:"sessRules" mapstructure:"SessRules"`
	// A map of PCC rules with the content being the PCCRule as described in subclause 5.6.2.6.
	PccRules map[string]PccRule `json:"pccRules,omitempty" yaml:"pccRules" bson:"pccRules" mapstructure:"PccRules"`
	// If it is included and set to true, it indicates the P-CSCF Restoration is requested.
	PcscfRestIndication bool `json:"pcscfRestIndication,omitempty" yaml:"pcscfRestIndication" bson:"pcscfRestIndication" mapstructure:"PcscfRestIndication"`
	// Map of QoS data policy decisions.
	QosDecs map[string]QosData `json:"qosDecs,omitempty" yaml:"qosDecs" bson:"qosDecs" mapstructure:"QosDecs"`
	// Map of Charging data policy decisions.
	ChgDecs      map[string]ChargingData `json:"chgDecs,omitempty" yaml:"chgDecs" bson:"chgDecs" mapstructure:"ChgDecs"`
	ChargingInfo *ChargingInformation    `json:"chargingInfo,omitempty" yaml:"chargingInfo" bson:"chargingInfo" mapstructure:"ChargingInfo"`
	// Map of Traffic Control data policy decisions.
	TraffContDecs map[string]TrafficControlData `json:"traffContDecs,omitempty" yaml:"traffContDecs" bson:"traffContDecs" mapstructure:"TraffContDecs"`
	// Map of Usage Monitoring data policy decisions.
	UmDecs map[string]UsageMonitoringData `json:"umDecs,omitempty" yaml:"umDecs" bson:"umDecs" mapstructure:"UmDecs"`
	// Map of QoS characteristics for non standard 5QIs. This map uses the 5QI values as keys.
	QosChars           map[string]QosCharacteristics `json:"qosChars,omitempty" yaml:"qosChars" bson:"qosChars" mapstructure:"QosChars"`
	ReflectiveQoSTimer int32                         `json:"reflectiveQoSTimer,omitempty" yaml:"reflectiveQoSTimer" bson:"reflectiveQoSTimer" mapstructure:"ReflectiveQoSTimer"`
	// A map of condition data with the content being as described in subclause 5.6.2.9.
	Conds            map[string]ConditionData `json:"conds,omitempty" yaml:"conds" bson:"conds" mapstructure:"Conds"`
	RevalidationTime *time.Time               `json:"revalidationTime,omitempty" yaml:"revalidationTime" bson:"revalidationTime" mapstructure:"RevalidationTime"`
	// Indicates the offline charging is applicable to the PDU session or PCC rule.
	Offline bool `json:"offline,omitempty" yaml:"offline" bson:"offline" mapstructure:"Offline"`
	// Indicates the online charging is applicable to the PDU session or PCC rule.
	Online bool `json:"online,omitempty" yaml:"online" bson:"online" mapstructure:"Online"`
	// Defines the policy control request triggers subscribed by the PCF.
	PolicyCtrlReqTriggers []PolicyControlRequestTrigger `json:"policyCtrlReqTriggers,omitempty" yaml:"policyCtrlReqTriggers" bson:"policyCtrlReqTriggers" mapstructure:"PolicyCtrlReqTriggers"`
	// Defines the last list of rule control data requested by the PCF.
	LastReqRuleData  []RequestedRuleData `json:"lastReqRuleData,omitempty" yaml:"lastReqRuleData" bson:"lastReqRuleData" mapstructure:"LastReqRuleData"`
	LastReqUsageData *RequestedUsageData `json:"lastReqUsageData,omitempty" yaml:"lastReqUsageData" bson:"lastReqUsageData" mapstructure:"LastReqUsageData"`
	// Map of PRA information.
	PraInfos     map[string]PresenceInfoRm `json:"praInfos,omitempty" yaml:"praInfos" bson:"praInfos" mapstructure:"PraInfos"`
	Ipv4Index    int32                     `json:"ipv4Index,omitempty" yaml:"ipv4Index" bson:"ipv4Index" mapstructure:"Ipv4Index"`
	Ipv6Index    int32                     `json:"ipv6Index,omitempty" yaml:"ipv6Index" bson:"ipv6Index" mapstructure:"Ipv6Index"`
	QosFlowUsage QosFlowUsage              `json:"qosFlowUsage,omitempty" yaml:"qosFlowUsage" bson:"qosFlowUsage" mapstructure:"QosFlowUsage"`
	SuppFeat     string                    `json:"suppFeat,omitempty" yaml:"suppFeat" bson:"suppFeat" mapstructure:"SuppFeat"`
}

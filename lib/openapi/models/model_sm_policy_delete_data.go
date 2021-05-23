/*
 * Npcf_SMPolicyControl
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

type SmPolicyDeleteData struct {
	UserLocationInfo     *UserLocation `json:"userLocationInfo,omitempty" yaml:"userLocationInfo" bson:"userLocationInfo" mapstructure:"UserLocationInfo"`
	UeTimeZone           string        `json:"ueTimeZone,omitempty" yaml:"ueTimeZone" bson:"ueTimeZone" mapstructure:"UeTimeZone"`
	ServingNetwork       *NetworkId    `json:"servingNetwork,omitempty" yaml:"servingNetwork" bson:"servingNetwork" mapstructure:"ServingNetwork"`
	UserLocationInfoTime *time.Time    `json:"userLocationInfoTime,omitempty" yaml:"userLocationInfoTime" bson:"userLocationInfoTime" mapstructure:"UserLocationInfoTime"`
	// Contains the RAN and/or NAS release cause.
	RanNasRelCauses []RanNasRelCause `json:"ranNasRelCauses,omitempty" yaml:"ranNasRelCauses" bson:"ranNasRelCauses" mapstructure:"RanNasRelCauses"`
	// Contains the usage report
	AccuUsageReports []AccuUsageReport `json:"accuUsageReports,omitempty" yaml:"accuUsageReports" bson:"accuUsageReports" mapstructure:"AccuUsageReports"`
}

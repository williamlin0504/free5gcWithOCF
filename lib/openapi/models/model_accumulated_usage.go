/*
 * Nchf_ConvergedChargingRelease Service API
 *
 * This is the Policy Authorization Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type AccumulatedUsage struct {
	// Unsigned integer identifying a period of time in units of seconds.
	Duration int32 `json:"duration,omitempty" yaml:"duration" bson:"duration" mapstructure:"Duration"`
	// Unsigned integer identifying a volume in units of bytes.
	TotalVolume int64 `json:"totalVolume,omitempty" yaml:"totalVolume" bson:"totalVolume" mapstructure:"TotalVolume"`
	// Unsigned integer identifying a volume in units of bytes.
	DownlinkVolume int64 `json:"downlinkVolume,omitempty" yaml:"downlinkVolume" bson:"downlinkVolume" mapstructure:"DownlinkVolume"`
	// Unsigned integer identifying a volume in units of bytes.
	UplinkVolume int64 `json:"uplinkVolume,omitempty" yaml:"uplinkVolume" bson:"uplinkVolume" mapstructure:"UplinkVolume"`
}

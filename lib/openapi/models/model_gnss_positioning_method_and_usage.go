/*
 * Nocf_Location
 *
 * OCF Location Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type GnssPositioningMethodAndUsage struct {
	Mode  PositioningMode `json:"mode" yaml:"mode" bson:"mode"`
	Gnss  GnssId          `json:"gnss" yaml:"gnss" bson:"gnss"`
	Usage Usage           `json:"usage" yaml:"usage" bson:"usage"`
}

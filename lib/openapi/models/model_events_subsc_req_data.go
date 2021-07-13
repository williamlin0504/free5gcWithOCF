/*
 * Nccf_PolicyAuthorization Service API
 *
 * This is the Policy Authorization Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// Identifies the events the application subscribes to.
type EventsSubscReqData struct {
	Events []AfEventSubscription `json:"events" yaml:"events" bson:"events" mapstructure:"Events"`
	// string providing an URI formatted according to IETF RFC 3986.
	NotifUri string          `json:"notifUri,omitempty" yaml:"notifUri" bson:"notifUri" mapstructure:"NotifUri"`
	UsgThres *UsageThreshold `json:"usgThres,omitempty" yaml:"usgThres" bson:"usgThres" mapstructure:"UsgThres"`
}

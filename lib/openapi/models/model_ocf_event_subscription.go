/*
 * Nocf_Communication
 *
 * OCF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type OcfEventSubscription struct {
	EventList                     *[]OcfEvent   `json:"eventList,omitempty"`
	EventNotifyUri                string        `json:"eventNotifyUri"`
	NotifyCorrelationId           string        `json:"notifyCorrelationId"`
	NfId                          string        `json:"nfId"`
	SubsChangeNotifyUri           string        `json:"subsChangeNotifyUri,omitempty"`
	SubsChangeNotifyCorrelationId string        `json:"subsChangeNotifyCorrelationId,omitempty"`
	Supi                          string        `json:"supi,omitempty"`
	GroupId                       string        `json:"groupId,omitempty"`
	Gpsi                          string        `json:"gpsi,omitempty"`
	Pei                           string        `json:"pei,omitempty"`
	AnyUE                         bool          `json:"anyUE,omitempty"`
	Options                       *OcfEventMode `json:"options,omitempty"`
}

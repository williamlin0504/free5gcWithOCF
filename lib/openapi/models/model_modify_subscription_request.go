/*
 * Namf_EventExposure
 *
 * AMF Event Exposure Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type ModifySubscriptionRequest struct {
	SubscriptionItemInner    *AmfUpdateEventSubscriptionItemInner
	OptionItem               *AmfUpdateEventOptionItem
	SubscriptionItemInnerOCF *OcfUpdateEventSubscriptionItemInner
	OptionItemOCF            *OcfUpdateEventOptionItem
}

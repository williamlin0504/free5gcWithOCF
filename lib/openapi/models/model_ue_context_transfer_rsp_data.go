/*
 * Nocf_Communication
 *
 * OCF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type UeContextTransferRspData struct {
	UeContext         *UeContext     `json:"ueContext"`
	UeRadioCapability *N2InfoContent `json:"ueRadioCapability,omitempty"`
	SupportedFeatures string         `json:"supportedFeatures,omitempty"`
}

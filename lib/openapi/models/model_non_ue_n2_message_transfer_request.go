/*
 * Nocf_Communication
 *
 * OCF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type NonUeN2MessageTransferRequest struct {
	JsonData                *N2InformationTransferReqData `json:"jsonData,omitempty" multipart:"contentType:application/json"`
	BinaryDataN2Information []byte                        `json:"binaryDataN2Information,omitempty" multipart:"contentType:application/vnd.3gpp.ngap,class:JsonData.N2Information.N2InformationClass,ref:(N2InfoContent).NgapData.ContentId"`
}

/*
 * Namf_Communication
 *
 * AMF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type SeafData struct {
	NgKsi                *NgKsi  `json:"ngKsi"`
	KeyAmf               *KeyAmf `json:"keyAmf"`
	KeyOcf               *KeyOcf `json:"keyOcf"`
	Nh                   string  `json:"nh,omitempty"`
	Ncc                  int32   `json:"ncc,omitempty"`
	KeyAmfChangeInd      bool    `json:"keyAmfChangeInd,omitempty"`
	KeyAmfHDerivationInd bool    `json:"keyAmfHDerivationInd,omitempty"`
	KeyOcfChangeInd      bool    `json:"keyOcfChangeInd,omitempty"`
	KeyOcfHDerivationInd bool    `json:"keyOcfHDerivationInd,omitempty"`
}

/*
 * Nccf_PolicyAuthorization Service API
 *
 * This is the Policy Authorization Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// Describes the authorization data of an Individual Application Session Context created by the ccf.
type AppSessionContextRespData struct {
	ServAuthInfo ServAuthInfo `json:"servAuthInfo,omitempty" yaml:"servAuthInfo" bson:"servAuthInfo" mapstructure:"ServAuthInfo"`
	SuppFeat     string       `json:"suppFeat,omitempty" yaml:"suppFeat" bson:"suppFeat" mapstructure:"SuppFeat"`
}

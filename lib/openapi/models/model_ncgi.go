/*
 * Nocf_Location
 *
 * OCF Location Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type Ncgi struct {
	PlmnId   *PlmnId `json:"plmnId" yaml:"plmnId" bson:"plmnId"`
	NrCellId string  `json:"nrCellId" yaml:"nrCellId" bson:"nrCellId"`
}

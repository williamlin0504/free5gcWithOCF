/*
 * Nccf_BDTPolicyControl Service API
 *
 * The Nccf_BDTPolicyControl Service is used by an NF service consumer to retrieve background data transfer policies from the ccf and to update the ccf with the background data transfer policy selected by the NF service consumer.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// Describes a network area information in which the NF service consumer requests the number of UEs.
type NetworkAreaInfo struct {
	// Contains a list of E-UTRA cell identities.
	Ecgis []Ecgi `json:"ecgis,omitempty" yaml:"ecgis" bson:"ecgis" mapstructure:"Ecgis"`
	// Contains a list of NR cell identities.
	Ncgis []Ncgi `json:"ncgis,omitempty" yaml:"ncgis" bson:"ncgis" mapstructure:"Ncgis"`
	// Contains a list of NG RAN nodes.
	GRanNodeIds []GlobalRanNodeId `json:"gRanNodeIds,omitempty" yaml:"gRanNodeIds" bson:"gRanNodeIds" mapstructure:"GRanNodeIds"`
	// Contains a list of tracking area identities.
	Tais []Tai `json:"tais,omitempty" yaml:"tais" bson:"tais" mapstructure:"Tais"`
}

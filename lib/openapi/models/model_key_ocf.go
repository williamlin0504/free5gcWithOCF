/*
 * Nocf_Communication
 *
 * OCF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type KeyOcf struct {
	KeyType KeyOcfType `json:"keyType"`
	KeyVal  string     `json:"keyVal"`
}
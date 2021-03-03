/*
 * Nocf_Communication
 *
 * OCF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type N1N2MessageTransferError struct {
	Error   *ProblemDetails       `json:"error"`
	ErrInfo *N1N2MsgTxfrErrDetail `json:"errInfo,omitempty"`
}

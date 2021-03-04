/*
 * Nchf_SMPolicyControl
 *
 * Session Management Policy Control Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package Nchf_SMPolicyControl

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/net/http2"
)

// APIClient manages communication with the Nchf_SMPolicyControl API v1.0.1
// In most cases there should be only one, shared, APIClient.
type APIClient struct {
	cfg    *Configuration
	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// API Services
	DefaultApi         *DefaultApiService
	DefaultCallbackApi *DefaultCallbackApi
}

type service struct {
	client *APIClient
}

// NewAPIClient creates a new API client. Requires a userAgent string describing your application.
// optionally a custom http.Client to allow for advanced features such as caching.
func NewAPIClient(cfg *Configuration) *APIClient {
	if cfg.httpClient == nil {
		cfg.httpClient = http.DefaultClient
		cfg.httpClient.Transport = &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	c := &APIClient{}
	c.cfg = cfg
	c.common.client = c

	// API Services
	c.DefaultApi = (*DefaultApiService)(&c.common)
	c.DefaultCallbackApi = (*DefaultCallbackApi)(&c.common)

	return c
}

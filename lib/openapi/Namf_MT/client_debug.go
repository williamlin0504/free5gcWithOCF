 //+build debug 

/*
 * Namf_MT
 *
 * AMF Mobile Termination Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package Namf_MT

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"net/http"
	" free5gc/lib/http2_util"
)

// APIClient manages communication with the Namf_MT API v1.0.0
// In most cases there should be only one, shared, APIClient.
type APIClient struct {
	cfg    *Configuration
	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// API Services
	UeContextDocumentApi  *UeContextDocumentApiService
	UeReachIndDocumentApi *UeReachIndDocumentApiService
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
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Rand:         http2_util.ZeroSource{},
			},
	}

	c := &APIClient{}
	c.cfg = cfg
	c.common.client = c

	// API Services
	c.UeContextDocumentApi = (*UeContextDocumentApiService)(&c.common)
	c.UeReachIndDocumentApi = (*UeReachIndDocumentApiService)(&c.common)

	return c
}

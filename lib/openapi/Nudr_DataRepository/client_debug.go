//+build debug

/*
 * Nudr_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package Nudr_DataRepository

import (
	"crypto/tls"
	"free5gc/lib/http2_util"
	"net/http"

	"golang.org/x/net/http2"
)

// APIClient manages communication with the Nudr_DataRepository API OpenAPI file API v1.0.0
// In most cases there should be only one, shared, APIClient.
type APIClient struct {
	cfg    *Configuration
	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// API Services
	OCF3GPPAccessRegistrationDocumentApi         *OCF3GPPAccessRegistrationDocumentApiService
	OCFNon3GPPAccessRegistrationDocumentApi      *OCFNon3GPPAccessRegistrationDocumentApiService
	AccessAndMobilityDataApi                     *AccessAndMobilityDataApiService
	AccessAndMobilitySubscriptionDataDocumentApi *AccessAndMobilitySubscriptionDataDocumentApiService
	OcfSubscriptionInfoDocumentApi               *OcfSubscriptionInfoDocumentApiService
	AuthEventDocumentApi                         *AuthEventDocumentApiService
	AuthenticationDataDocumentApi                *AuthenticationDataDocumentApiService
	AuthenticationSoRDocumentApi                 *AuthenticationSoRDocumentApiService
	AuthenticationStatusDocumentApi              *AuthenticationStatusDocumentApiService
	CreateOCFSubscriptionInfoDocumentApi         *CreateOCFSubscriptionInfoDocumentApiService
	DefaultApi                                   *DefaultApiService
	EventOCFSubscriptionInfoDocumentApi          *EventOCFSubscriptionInfoDocumentApiService
	EventExposureDataDocumentApi                 *EventExposureDataDocumentApiService
	EventExposureGroupSubscriptionsCollectionApi *EventExposureGroupSubscriptionsCollectionApiService
	EventExposureSubscriptionDocumentApi         *EventExposureSubscriptionDocumentApiService
	EventExposureSubscriptionsCollectionApi      *EventExposureSubscriptionsCollectionApiService
	OperatorSpecificDataContainerDocumentApi     *OperatorSpecificDataContainerDocumentApiService
	ParameterProvisionDocumentApi                *ParameterProvisionDocumentApiService
	PduSessionManagementDataApi                  *PduSessionManagementDataApiService
	ProvisionedDataDocumentApi                   *ProvisionedDataDocumentApiService
	ProvisionedParameterDataDocumentApi          *ProvisionedParameterDataDocumentApiService
	QueryOCFSubscriptionInfoDocumentApi          *QueryOCFSubscriptionInfoDocumentApiService
	QueryIdentityDataBySUPIOrGPSIDocumentApi     *QueryIdentityDataBySUPIOrGPSIDocumentApiService
	QueryODBDataBySUPIOrGPSIDocumentApi          *QueryODBDataBySUPIOrGPSIDocumentApiService
	RetrievalOfSharedDataApi                     *RetrievalOfSharedDataApiService
	SDMSubscriptionDocumentApi                   *SDMSubscriptionDocumentApiService
	SDMSubscriptionsCollectionApi                *SDMSubscriptionsCollectionApiService
	SMFRegistrationDocumentApi                   *SMFRegistrationDocumentApiService
	SMFRegistrationsCollectionApi                *SMFRegistrationsCollectionApiService
	SMFSelectionSubscriptionDataDocumentApi      *SMFSelectionSubscriptionDataDocumentApiService
	SMSF3GPPRegistrationDocumentApi              *SMSF3GPPRegistrationDocumentApiService
	SMSFNon3GPPRegistrationDocumentApi           *SMSFNon3GPPRegistrationDocumentApiService
	SMSManagementSubscriptionDataDocumentApi     *SMSManagementSubscriptionDataDocumentApiService
	SMSSubscriptionDataDocumentApi               *SMSSubscriptionDataDocumentApiService
	SessionManagementSubscriptionDataApi         *SessionManagementSubscriptionDataApiService
	SubsToNofifyCollectionApi                    *SubsToNofifyCollectionApiService
	SubsToNotifyDocumentApi                      *SubsToNotifyDocumentApiService
	TraceDataDocumentApi                         *TraceDataDocumentApiService
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
				Rand:               http2_util.ZeroSource{},
			},
		}
	}

	c := &APIClient{}
	c.cfg = cfg
	c.common.client = c

	// API Services
	c.OCF3GPPAccessRegistrationDocumentApi = (*OCF3GPPAccessRegistrationDocumentApiService)(&c.common)
	c.OCFNon3GPPAccessRegistrationDocumentApi = (*OCFNon3GPPAccessRegistrationDocumentApiService)(&c.common)
	c.AccessAndMobilityDataApi = (*AccessAndMobilityDataApiService)(&c.common)
	c.AccessAndMobilitySubscriptionDataDocumentApi = (*AccessAndMobilitySubscriptionDataDocumentApiService)(&c.common)
	c.OcfSubscriptionInfoDocumentApi = (*OcfSubscriptionInfoDocumentApiService)(&c.common)
	c.AuthEventDocumentApi = (*AuthEventDocumentApiService)(&c.common)
	c.AuthenticationDataDocumentApi = (*AuthenticationDataDocumentApiService)(&c.common)
	c.AuthenticationSoRDocumentApi = (*AuthenticationSoRDocumentApiService)(&c.common)
	c.AuthenticationStatusDocumentApi = (*AuthenticationStatusDocumentApiService)(&c.common)
	c.CreateOCFSubscriptionInfoDocumentApi = (*CreateOCFSubscriptionInfoDocumentApiService)(&c.common)
	c.DefaultApi = (*DefaultApiService)(&c.common)
	c.EventOCFSubscriptionInfoDocumentApi = (*EventOCFSubscriptionInfoDocumentApiService)(&c.common)
	c.EventExposureDataDocumentApi = (*EventExposureDataDocumentApiService)(&c.common)
	c.EventExposureGroupSubscriptionsCollectionApi = (*EventExposureGroupSubscriptionsCollectionApiService)(&c.common)
	c.EventExposureSubscriptionDocumentApi = (*EventExposureSubscriptionDocumentApiService)(&c.common)
	c.EventExposureSubscriptionsCollectionApi = (*EventExposureSubscriptionsCollectionApiService)(&c.common)
	c.OperatorSpecificDataContainerDocumentApi = (*OperatorSpecificDataContainerDocumentApiService)(&c.common)
	c.ParameterProvisionDocumentApi = (*ParameterProvisionDocumentApiService)(&c.common)
	c.PduSessionManagementDataApi = (*PduSessionManagementDataApiService)(&c.common)
	c.ProvisionedDataDocumentApi = (*ProvisionedDataDocumentApiService)(&c.common)
	c.ProvisionedParameterDataDocumentApi = (*ProvisionedParameterDataDocumentApiService)(&c.common)
	c.QueryOCFSubscriptionInfoDocumentApi = (*QueryOCFSubscriptionInfoDocumentApiService)(&c.common)
	c.QueryIdentityDataBySUPIOrGPSIDocumentApi = (*QueryIdentityDataBySUPIOrGPSIDocumentApiService)(&c.common)
	c.QueryODBDataBySUPIOrGPSIDocumentApi = (*QueryODBDataBySUPIOrGPSIDocumentApiService)(&c.common)
	c.RetrievalOfSharedDataApi = (*RetrievalOfSharedDataApiService)(&c.common)
	c.SDMSubscriptionDocumentApi = (*SDMSubscriptionDocumentApiService)(&c.common)
	c.SDMSubscriptionsCollectionApi = (*SDMSubscriptionsCollectionApiService)(&c.common)
	c.SMFRegistrationDocumentApi = (*SMFRegistrationDocumentApiService)(&c.common)
	c.SMFRegistrationsCollectionApi = (*SMFRegistrationsCollectionApiService)(&c.common)
	c.SMFSelectionSubscriptionDataDocumentApi = (*SMFSelectionSubscriptionDataDocumentApiService)(&c.common)
	c.SMSF3GPPRegistrationDocumentApi = (*SMSF3GPPRegistrationDocumentApiService)(&c.common)
	c.SMSFNon3GPPRegistrationDocumentApi = (*SMSFNon3GPPRegistrationDocumentApiService)(&c.common)
	c.SMSManagementSubscriptionDataDocumentApi = (*SMSManagementSubscriptionDataDocumentApiService)(&c.common)
	c.SMSSubscriptionDataDocumentApi = (*SMSSubscriptionDataDocumentApiService)(&c.common)
	c.SessionManagementSubscriptionDataApi = (*SessionManagementSubscriptionDataApiService)(&c.common)
	c.SubsToNofifyCollectionApi = (*SubsToNofifyCollectionApiService)(&c.common)
	c.SubsToNotifyDocumentApi = (*SubsToNotifyDocumentApiService)(&c.common)
	c.TraceDataDocumentApi = (*TraceDataDocumentApiService)(&c.common)

	return c
}

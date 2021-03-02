package consumer

import (
	"context"
	"fmt"
	"net/http"

	ocf_context "github.com/free5gc/ocf/context"
	"github.com/free5gc/ocf/logger"
	"github.com/free5gc/ocf/util"
	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/models"
)

func SendSearchNFInstances(nrfUri string, targetNfType, requestNfType models.NfType,
	param *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (models.SearchResult, error) {
	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	result, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requestNfType, param)
	if res != nil && res.StatusCode == http.StatusTemporaryRedirect {
		err = fmt.Errorf("Temporary Redirect For Non NRF Consumer")
	}
	defer func() {
		if bodyCloseErr := res.Body.Close(); bodyCloseErr != nil {
			err = fmt.Errorf("SearchNFInstances' response body cannot close: %+w", bodyCloseErr)
		}
	}()
	return result, err
}

func SearchUdmSdmInstance(ue *ocf_context.OcfUe, nrfUri string, targetNfType, requestNfType models.NfType,
	param *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) error {
	resp, localErr := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, param)
	if localErr != nil {
		return localErr
	}

	// select the first UDM_SDM, TODO: select base on other info
	var sdmUri string
	for _, nfProfile := range resp.NfInstances {
		ue.UdmId = nfProfile.NfInstanceId
		sdmUri = util.SearchNFServiceUri(nfProfile, models.ServiceName_NUDM_SDM, models.NfServiceStatus_REGISTERED)
		if sdmUri != "" {
			break
		}
	}
	ue.NudmSDMUri = sdmUri
	if ue.NudmSDMUri == "" {
		err := fmt.Errorf("OCF can not select an UDM by NRF")
		logger.ConsumerLog.Errorf(err.Error())
		return err
	}
	return nil
}

func SearchNssfNSSelectionInstance(ue *ocf_context.OcfUe, nrfUri string, targetNfType, requestNfType models.NfType,
	param *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) error {
	resp, localErr := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, param)
	if localErr != nil {
		return localErr
	}

	// select the first NSSF, TODO: select base on other info
	var nssfUri string
	for _, nfProfile := range resp.NfInstances {
		ue.NssfId = nfProfile.NfInstanceId
		nssfUri = util.SearchNFServiceUri(nfProfile, models.ServiceName_NNSSF_NSSELECTION, models.NfServiceStatus_REGISTERED)
		if nssfUri != "" {
			break
		}
	}
	ue.NssfUri = nssfUri
	if ue.NssfUri == "" {
		return fmt.Errorf("OCF can not select an NSSF by NRF")
	}
	return nil
}

func SearchOcfCommunicationInstance(ue *ocf_context.OcfUe, nrfUri string, targetNfType,
	requestNfType models.NfType, param *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (err error) {
	resp, localErr := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, param)
	if localErr != nil {
		err = localErr
		return
	}

	// select the first OCF, TODO: select base on other info
	var ocfUri string
	for _, nfProfile := range resp.NfInstances {
		ue.TargetOcfProfile = &nfProfile
		ocfUri = util.SearchNFServiceUri(nfProfile, models.ServiceName_NOCF_COMM, models.NfServiceStatus_REGISTERED)
		if ocfUri != "" {
			break
		}
	}
	ue.TargetOcfUri = ocfUri
	if ue.TargetOcfUri == "" {
		err = fmt.Errorf("OCF can not select an target OCF by NRF")
	}
	return
}
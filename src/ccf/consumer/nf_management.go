package consumer

import (
	"context"
	"fmt"
	" free5gcWithOCF/lib/openapi/Nnrf_NFManagement"
	" free5gcWithOCF/lib/openapi/models"
	ccf_context " free5gcWithOCF/src/ccf/context"
	"net/http"
	"strings"
	"time"
)

func BuildNFInstance(context *ccf_context.CCFContext) (profile models.NfProfile, err error) {
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NfType_CCF
	profile.NfStatus = models.NfStatus_REGISTERED
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, context.RegisterIPv4)
	service := []models.NfService{}
	for _, nfService := range context.NfService {
		service = append(service, nfService)
	}
	profile.NfServices = &service
	profile.CcfInfo = &models.CcfInfo{
		DnnList: []string{
			" free5gcWithOCF",
			"internet",
		},
		// SupiRanges: &[]models.SupiRange{
		// 	{
		// 		//from TS 29.510 6.1.6.2.9 example2
		//		//no need to set supirange in this moment 2019/10/4
		// 		Start:   "123456789040000",
		// 		End:     "123456789059999",
		// 		Pattern: "^imsi-12345678904[0-9]{4}$",
		// 	},
		// },
	}
	return
}

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (
	resouceNrfUri string, retrieveNfInstanceID string, err error) {

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response
	for {
		_, res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, profile)
		if err != nil || res == nil {
			//TODO : add log
			fmt.Println(fmt.Errorf("CCF register to NRF Error[%v]", err.Error()))
			time.Sleep(2 * time.Second)
			continue
		}
		status := res.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			break
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := res.Header.Get("Location")
			resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			retrieveNfInstanceID = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			break
		} else {
			fmt.Println("NRF return wrong status code", status)
		}
	}
	return resouceNrfUri, retrieveNfInstanceID, err
}

package consumer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	ocf_context "github.com/free5gc/ocf/context"
	"github.com/free5gc/ocf/logger"
	"github.com/free5gc/ocf/util"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/models"
)

func BuildNFInstance(context *ocf_context.OCFContext) (profile models.NfProfile, err error) {
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NfType_OCF
	profile.NfStatus = models.NfStatus_REGISTERED
	var plmns []models.PlmnId
	for _, plmnItem := range context.PlmnSupportList {
		plmns = append(plmns, plmnItem.PlmnId)
	}
	if len(plmns) > 0 {
		profile.PlmnList = &plmns
		// TODO: change to Per Plmn Support Snssai List
		profile.SNssais = &context.PlmnSupportList[0].SNssaiList
	}
	ocfInfo := models.OcfInfo{}
	if len(context.ServedGuamiList) == 0 {
		err = fmt.Errorf("Gumai List is Empty in OCF")
		return
	}
	regionId, setId, _, err1 := util.SeperateOcfId(context.ServedGuamiList[0].OcfId)
	if err1 != nil {
		err = err1
		return
	}
	ocfInfo.OcfRegionId = regionId
	ocfInfo.OcfSetId = setId
	ocfInfo.GuamiList = &context.ServedGuamiList
	if len(context.SupportTaiLists) == 0 {
		err = fmt.Errorf("SupportTaiList is Empty in OCF")
		return
	}
	ocfInfo.TaiList = &context.SupportTaiLists
	profile.OcfInfo = &ocfInfo
	if context.RegisterIPv4 == "" {
		err = fmt.Errorf("OCF Address is empty")
		return
	}
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, context.RegisterIPv4)
	service := []models.NfService{}
	for _, nfService := range context.NfService {
		service = append(service, nfService)
	}
	if len(service) > 0 {
		profile.NfServices = &service
	}

	defaultNotificationSubscription := models.DefaultNotificationSubscription{
		CallbackUri:      fmt.Sprintf("%s/nocf-callback/v1/n1-message-notify", context.GetIPv4Uri()),
		NotificationType: models.NotificationType_N1_MESSAGES,
		N1MessageClass:   models.N1MessageClass__5_GMM,
	}
	profile.DefaultNotificationSubscriptions =
		append(profile.DefaultNotificationSubscriptions, defaultNotificationSubscription)
	return profile, err
}

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (
	resouceNrfUri string, retrieveNfInstanceId string, err error) {
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response
	for {
		_, res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, profile)
		if err != nil || res == nil {
			// TODO : add log
			fmt.Println(fmt.Errorf("OCF register to NRF Error[%s]", err.Error()))
			time.Sleep(2 * time.Second)
			continue
		}
		defer func() {
			if bodyCloseErr := res.Body.Close(); bodyCloseErr != nil {
				err = fmt.Errorf("SearchNFInstances' response body cannot close: %+w", bodyCloseErr)
			}
		}()
		status := res.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			break
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := res.Header.Get("Location")
			resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			break
		} else {
			fmt.Println(fmt.Errorf("handler returned wrong status code %d", status))
			fmt.Println(fmt.Errorf("NRF return wrong status code %d", status))
		}
	}
	return resouceNrfUri, retrieveNfInstanceId, err
}

func SendDeregisterNFInstance() (problemDetails *models.ProblemDetails, err error) {
	logger.ConsumerLog.Infof("[OCF] Send Deregister NFInstance")

	ocfSelf := ocf_context.OCF_Self()
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(ocfSelf.NrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response

	res, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(context.Background(), ocfSelf.NfId)
	if err == nil {
		return
	} else if res != nil {
		defer func() {
			if bodyCloseErr := res.Body.Close(); bodyCloseErr != nil {
				err = fmt.Errorf("SearchNFInstances' response body cannot close: %+w", bodyCloseErr)
			}
		}()
		if res.Status != err.Error() {
			return
		}
		problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = openapi.ReportError("server no response")
	}
	return
}
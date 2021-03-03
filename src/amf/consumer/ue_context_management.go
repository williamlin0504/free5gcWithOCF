package consumer

import (
	"context"

	"free5gc/lib/openapi"
	"free5gc/lib/openapi/Nudm_UEContextManagement"
	"free5gc/lib/openapi/models"
	ocf_context "free5gc/src/ocf/context"
)

func UeCmRegistration(ue *ocf_context.OcfUe, accessType models.AccessType, initialRegistrationInd bool) (
	*models.ProblemDetails, error) {

	configuration := Nudm_UEContextManagement.NewConfiguration()
	configuration.SetBasePath(ue.NudmUECMUri)
	client := Nudm_UEContextManagement.NewAPIClient(configuration)

	ocfSelf := ocf_context.OCF_Self()

	switch accessType {
	case models.AccessType__3_GPP_ACCESS:
		registrationData := models.Ocf3GppAccessRegistration{
			OcfInstanceId:          ocfSelf.NfId,
			InitialRegistrationInd: initialRegistrationInd,
			Guami:                  &ocfSelf.ServedGuamiList[0],
			RatType:                ue.RatType,
			// TODO: not support Homogenous Support of IMS Voice over PS Sessions this stage
			ImsVoPs: models.ImsVoPs_HOMOGENEOUS_NON_SUPPORT,
		}

		_, httpResp, localErr := client.OCFRegistrationFor3GPPAccessApi.Registration(context.Background(),
			ue.Supi, registrationData)
		if localErr == nil {
			return nil, nil
		} else if httpResp != nil {
			if httpResp.Status != localErr.Error() {
				return nil, localErr
			}
			problem := localErr.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
			return &problem, nil
		} else {
			return nil, openapi.ReportError("server no response")
		}
	case models.AccessType_NON_3_GPP_ACCESS:
		registrationData := models.OcfNon3GppAccessRegistration{
			OcfInstanceId: ocfSelf.NfId,
			Guami:         &ocfSelf.ServedGuamiList[0],
			RatType:       ue.RatType,
		}

		_, httpResp, localErr :=
			client.OCFRegistrationForNon3GPPAccessApi.Register(context.Background(), ue.Supi, registrationData)
		if localErr == nil {
			return nil, nil
		} else if httpResp != nil {
			if httpResp.Status != localErr.Error() {
				return nil, localErr
			}
			problem := localErr.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
			return &problem, nil
		} else {
			return nil, openapi.ReportError("server no response")
		}
	}

	return nil, nil
}

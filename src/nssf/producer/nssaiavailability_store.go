/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package producer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	jsonpatch "github.com/evanphx/json-patch"

	"free5gc/lib/openapi/models"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/logger"
	"free5gc/src/nssf/plugin"
	"free5gc/src/nssf/util"
)

// NSSAIAvailability DELETE method
func NSSAIAvailabilityDeleteProcedure(nfId string) *models.ProblemDetails {
	var problemDetails *models.ProblemDetails
	for i, ocfConfig := range factory.NssfConfig.Configuration.OcfList {
		if ocfConfig.NfId == nfId {
			factory.NssfConfig.Configuration.OcfList = append(
				factory.NssfConfig.Configuration.OcfList[:i],
				factory.NssfConfig.Configuration.OcfList[i+1:]...)
			return nil
		}
	}

	*problemDetails = models.ProblemDetails{
		Title:  util.UNSUPPORTED_RESOURCE,
		Status: http.StatusNotFound,
		Detail: fmt.Sprintf("OCF ID '%s' does not exist", nfId),
	}
	return problemDetails
}

// NSSAIAvailability PATCH method
func NSSAIAvailabilityPatchProcedure(nssaiAvailabilityUpdateInfo plugin.PatchDocument, nfId string) (
	*models.AuthorizedNssaiAvailabilityInfo, *models.ProblemDetails) {
	var (
		response       *models.AuthorizedNssaiAvailabilityInfo = &models.AuthorizedNssaiAvailabilityInfo{}
		problemDetails *models.ProblemDetails
	)

	var ocfIdx int
	var original []byte
	hitOcf := false
	factory.ConfigLock.RLock()
	for ocfIdx, ocfConfig := range factory.NssfConfig.Configuration.OcfList {
		if ocfConfig.NfId == nfId {
			// Since json-patch package does not have idea of optional field of datatype,
			// provide with null or empty value instead of omitting the field
			temp := factory.NssfConfig.Configuration.OcfList[ocfIdx].SupportedNssaiAvailabilityData
			const dummyString string = "DUMMY"
			for i := range temp {
				for j := range temp[i].SupportedSnssaiList {
					if temp[i].SupportedSnssaiList[j].Sd == "" {
						temp[i].SupportedSnssaiList[j].Sd = dummyString
					}
				}
			}
			var err error
			original, err = json.Marshal(temp)
			if err != nil {
				logger.Nssaiavailability.Errorf("Marshal error in NSSAIAvailabilityPatchProcedure: %+v", err)
			}
			original = bytes.ReplaceAll(original, []byte(dummyString), []byte(""))

			// original, _ = json.Marshal(factory.NssfConfig.Configuration.OcfList[ocfIdx].SupportedNssaiAvailabilityData)

			hitOcf = true
			break
		}
	}
	factory.ConfigLock.RUnlock()
	if !hitOcf {
		*problemDetails = models.ProblemDetails{
			Title:  util.UNSUPPORTED_RESOURCE,
			Status: http.StatusNotFound,
			Detail: fmt.Sprintf("OCF ID '%s' does not exist", nfId),
		}
		return nil, problemDetails
	}

	// TODO: Check if returned HTTP status codes or problem details are proper when errors occur

	// Provide JSON string with null or empty value in `Value` of `PatchItem`
	for i, patchItem := range nssaiAvailabilityUpdateInfo {
		if reflect.ValueOf(patchItem.Value).Kind() == reflect.Map {
			_, exist := patchItem.Value.(map[string]interface{})["sst"]
			_, notExist := patchItem.Value.(map[string]interface{})["sd"]
			if exist && !notExist {
				nssaiAvailabilityUpdateInfo[i].Value.(map[string]interface{})["sd"] = ""
			}
		}
	}
	patchJSON, err := json.Marshal(nssaiAvailabilityUpdateInfo)
	if err != nil {
		logger.Nssaiavailability.Errorf("Marshal error in NSSAIAvailabilityPatchProcedure: %+v", err)
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		*problemDetails = models.ProblemDetails{
			Title:  util.MALFORMED_REQUEST,
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		return nil, problemDetails
	}

	modified, err := patch.Apply(original)
	if err != nil {
		*problemDetails = models.ProblemDetails{
			Title:  util.INVALID_REQUEST,
			Status: http.StatusConflict,
			Detail: err.Error(),
		}
		return nil, problemDetails
	}

	factory.ConfigLock.Lock()
	err = json.Unmarshal(modified, &factory.NssfConfig.Configuration.OcfList[ocfIdx].SupportedNssaiAvailabilityData)
	factory.ConfigLock.Unlock()
	if err != nil {
		*problemDetails = models.ProblemDetails{
			Title:  util.INVALID_REQUEST,
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		return nil, problemDetails
	}

	// Return all authorized NSSAI availability information
	response.AuthorizedNssaiAvailabilityData, err = util.AuthorizeOfOcfFromConfig(nfId)
	if err != nil {
		logger.Nssaiavailability.Errorf("util AuthorizeOfOcfFromConfig error in NSSAIAvailabilityPatchProcedure: %+v", err)
	}

	// TODO: Return authorized NSSAI availability information of updated TAI only

	return response, nil
}

// NSSAIAvailability PUT method
func NSSAIAvailabilityPutProcedure(nssaiAvailabilityInfo models.NssaiAvailabilityInfo, nfId string) (
	*models.AuthorizedNssaiAvailabilityInfo, *models.ProblemDetails) {
	var (
		response       *models.AuthorizedNssaiAvailabilityInfo = &models.AuthorizedNssaiAvailabilityInfo{}
		problemDetails *models.ProblemDetails
	)

	for _, s := range nssaiAvailabilityInfo.SupportedNssaiAvailabilityData {
		if !util.CheckSupportedNssaiInPlmn(s.SupportedSnssaiList, *s.Tai.PlmnId) {
			problemDetails = &models.ProblemDetails{
				Title:  util.UNSUPPORTED_RESOURCE,
				Status: http.StatusForbidden,
				Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
				Cause:  "SNSSAI_NOT_SUPPORTED",
			}
			return nil, problemDetails
		}
	}

	// TODO: Currently authorize all the provided S-NSSAIs
	//       Take some issue into consideration e.g. operator policies

	hitOcf := false
	// Find OCF configuration of given NfId
	// If found, then update the SupportedNssaiAvailabilityData
	factory.ConfigLock.Lock()
	for i, ocfConfig := range factory.NssfConfig.Configuration.OcfList {
		if ocfConfig.NfId == nfId {
			factory.NssfConfig.Configuration.OcfList[i].SupportedNssaiAvailabilityData =
				nssaiAvailabilityInfo.SupportedNssaiAvailabilityData

			hitOcf = true
			break
		}
	}
	factory.ConfigLock.Unlock()

	// If no OCF record is found, create a new one
	if !hitOcf {
		var ocfConfig factory.OcfConfig
		ocfConfig.NfId = nfId
		ocfConfig.SupportedNssaiAvailabilityData = nssaiAvailabilityInfo.SupportedNssaiAvailabilityData
		factory.ConfigLock.Lock()
		factory.NssfConfig.Configuration.OcfList = append(factory.NssfConfig.Configuration.OcfList, ocfConfig)
		factory.ConfigLock.Unlock()
	}

	// Return all authorized NSSAI availability information
	// a.AuthorizedNssaiAvailabilityData, _ = authorizeOfOcfFromConfig(nfId)

	// Return authorized NSSAI availability information of updated TAI only
	for _, s := range nssaiAvailabilityInfo.SupportedNssaiAvailabilityData {
		authorizedNssaiAvailabilityData, err := util.AuthorizeOfOcfTaFromConfig(nfId, *s.Tai)
		if err == nil {
			response.AuthorizedNssaiAvailabilityData =
				append(response.AuthorizedNssaiAvailabilityData, authorizedNssaiAvailabilityData)
		} else {
			logger.Nssaiavailability.Warnf(err.Error())
		}
	}

	return response, problemDetails
}

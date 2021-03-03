package util

import (
	"encoding/hex"
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/ocf/logger"
	"strconv"
)

func SnssaiHexToModels(hexString string) (*models.Snssai, error) {
	sst, err := strconv.ParseInt(hexString[:2], 16, 32)
	if err != nil {
		return nil, err
	}
	sNssai := models.Snssai{
		Sst: int32(sst),
		Sd:  hexString[2:],
	}
	return &sNssai, nil
}

func SnssaiModelsToHex(snssai models.Snssai) string {
	sst := fmt.Sprintf("%02x", snssai.Sst)
	return sst + snssai.Sd
}

func SeperateOcfId(ocfid string) (regionId, setId, ptrId string, err error) {
	if len(ocfid) != 6 {
		err = fmt.Errorf("len of ocfId[%s] != 6", ocfid)
		return
	}
	// regionId: 16bits, setId: 10bits, ptrId: 6bits
	regionId = ocfid[:2]
	byteArray, err1 := hex.DecodeString(ocfid[2:])
	if err1 != nil {
		err = err1
		return
	}
	byteSetId := []byte{byteArray[0] >> 6, byteArray[0]<<2 | byteArray[1]>>6}
	setId = hex.EncodeToString(byteSetId)[1:]
	bytePtrId := []byte{byteArray[1] & 0x3f}
	ptrId = hex.EncodeToString(bytePtrId)
	return
}

func PlmnIdStringToModels(plmnId string) (plmnID models.PlmnId) {
	plmnID.Mcc = plmnId[:3]
	plmnID.Mnc = plmnId[3:]
	return
}

func TACConfigToModels(intString string) (hexString string) {
	tmp, err := strconv.ParseUint(intString, 10, 32)
	if err != nil {
		logger.UtilLog.Errorf("ParseUint error: %+v", err)
	}
	hexString = fmt.Sprintf("%06x", tmp)
	return
}

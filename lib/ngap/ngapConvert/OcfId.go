package ngapConvert

import (
	"encoding/hex"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/logger"
)

func OcfIdToNgap(ocfId string) (regionId, setId, ptrId aper.BitString) {
	regionId = HexToBitString(ocfId[:2], 8)
	setId = HexToBitString(ocfId[2:5], 10)
	tmpByte, err := hex.DecodeString(ocfId[4:])
	if err != nil {
		logger.NgapLog.Warningln("OcfId From Models To NGAP Error: ", err.Error())
		return
	}
	shiftByte, err := aper.GetBitString(tmpByte, 2, 6)
	if err != nil {
		logger.NgapLog.Warningln("OcfId From Models To NGAP Error: ", err.Error())
		return
	}
	ptrId.BitLength = 6
	ptrId.Bytes = shiftByte
	return
}

func OcfIdToModels(regionId, setId, ptrId aper.BitString) (ocfId string) {
	regionHex := BitStringToHex(&regionId)
	tmpByte := []byte{setId.Bytes[0], (setId.Bytes[1] & 0xc0) | (ptrId.Bytes[0] >> 2)}
	restHex := hex.EncodeToString(tmpByte)
	ocfId = regionHex + restHex
	return
}

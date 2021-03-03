package util

import (
	"encoding/binary"
	"encoding/hex"
	"free5gcWithOCF/lib/aper"
	"free5gcWithOCF/lib/ngap/ngapType"
	"free5gcWithOCF/src/ocf/context"
	"free5gcWithOCF/src/ocf/logger"
	"strings"
)

func PlmnIdToNgap(plmnId context.PLMNID) (ngapPlmnId ngapType.PLMNIdentity) {
	var hexString string
	mcc := strings.Split(plmnId.Mcc, "")
	mnc := strings.Split(plmnId.Mnc, "")
	if len(plmnId.Mnc) == 2 {
		hexString = mcc[1] + mcc[0] + "f" + mcc[2] + mnc[1] + mnc[0]
	} else {
		hexString = mcc[1] + mcc[0] + mnc[0] + mcc[2] + mnc[2] + mnc[1]
	}
	var err error
	ngapPlmnId.Value, err = hex.DecodeString(hexString)
	if err != nil {
		logger.UtilLog.Errorf("DecodeString error: %+v", err)
	}
	return
}

func OcfIdToNgap(ocfId uint16) (ngapOcfId *aper.BitString) {
	ngapOcfId = new(aper.BitString)
	ngapOcfId.Bytes = make([]byte, 2)
	binary.BigEndian.PutUint16(ngapOcfId.Bytes, ocfId)
	ngapOcfId.BitLength = 16
	return
}

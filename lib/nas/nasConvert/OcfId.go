package nasConvert

import (
	"encoding/hex"
	"log"
)

func OcfIdToNas(ocfId string) (ocfRegionId uint8, ocfSetId uint16, ocfPointer uint8) {

	ocfIdBytes, err := hex.DecodeString(ocfId)
	if err != nil {
		log.Printf("ocfId decode failed: %+v", err)
	}

	ocfRegionId = uint8(ocfIdBytes[0])
	ocfSetId = uint16(ocfIdBytes[1])<<2 + (uint16(ocfIdBytes[2])&0x00c0)>>6
	ocfPointer = uint8(ocfIdBytes[2]) & 0x3f
	return
}

func OcfIdToModels(ocfRegionId uint8, ocfSetId uint16, ocfPointer uint8) (ocfId string) {

	tmpBytes := []uint8{ocfRegionId, uint8(ocfSetId>>2) & 0xff, uint8(ocfSetId&0x03) + ocfPointer&0x3f}
	ocfId = hex.EncodeToString(tmpBytes)
	return
}

package ngapType

import " free5gcWithOCF/lib/aper"

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceModifyItemModReq struct {
	PDUSessionID                            PDUSessionID
	NASPDU                                  *NASPDU `aper:"optional"`
	PDUSessionResourceModifyRequestTransfer aper.OctetString
	IEExtensions                            *ProtocolExtensionContainerPDUSessionResourceModifyItemModReqExtIEs `aper:"optional"`
}

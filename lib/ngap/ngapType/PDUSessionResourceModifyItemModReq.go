package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceModifyItemModReq struct {
	PDUSessionID                            PDUSessionID
	NASPDU                                  *NASPDU `aper:"optional"`
	PDUSessionResourceModifyRequestTransfer aper.OctetString
	IEExtensions                            *ProtocolExtensionContainerPDUSessionResourceModifyItemModReqExtIEs `aper:"optional"`
}

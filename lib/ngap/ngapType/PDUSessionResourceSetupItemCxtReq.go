package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceSetupItemCxtReq struct {
	PDUSessionID                           PDUSessionID
	NASPDU                                 *NASPDU `aper:"optional"`
	SNSSAI                                 SNSSAI  `aper:"valueExt"`
	PDUSessionResourceSetupRequestTransfer aper.OctetString
	IEExtensions                           *ProtocolExtensionContainerPDUSessionResourceSetupItemCxtReqExtIEs `aper:"optional"`
}

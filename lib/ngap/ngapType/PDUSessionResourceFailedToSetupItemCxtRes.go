package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceFailedToSetupItemCxtRes struct {
	PDUSessionID                                PDUSessionID
	PDUSessionResourceSetupUnsuccessfulTransfer aper.OctetString
	IEExtensions                                *ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemCxtResExtIEs `aper:"optional"`
}

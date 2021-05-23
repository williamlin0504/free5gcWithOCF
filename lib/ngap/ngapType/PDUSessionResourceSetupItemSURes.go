package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceSetupItemSURes struct {
	PDUSessionID                            PDUSessionID
	PDUSessionResourceSetupResponseTransfer aper.OctetString
	IEExtensions                            *ProtocolExtensionContainerPDUSessionResourceSetupItemSUResExtIEs `aper:"optional"`
}

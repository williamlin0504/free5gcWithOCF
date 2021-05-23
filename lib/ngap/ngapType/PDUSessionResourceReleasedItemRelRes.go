package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceReleasedItemRelRes struct {
	PDUSessionID                              PDUSessionID
	PDUSessionResourceReleaseResponseTransfer aper.OctetString
	IEExtensions                              *ProtocolExtensionContainerPDUSessionResourceReleasedItemRelResExtIEs `aper:"optional"`
}

package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceToReleaseItemRelCmd struct {
	PDUSessionID                             PDUSessionID
	PDUSessionResourceReleaseCommandTransfer aper.OctetString
	IEExtensions                             *ProtocolExtensionContainerPDUSessionResourceToReleaseItemRelCmdExtIEs `aper:"optional"`
}

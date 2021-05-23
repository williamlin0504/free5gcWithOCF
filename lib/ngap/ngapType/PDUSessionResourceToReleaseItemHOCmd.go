package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceToReleaseItemHOCmd struct {
	PDUSessionID                            PDUSessionID
	HandoverPreparationUnsuccessfulTransfer aper.OctetString
	IEExtensions                            *ProtocolExtensionContainerPDUSessionResourceToReleaseItemHOCmdExtIEs `aper:"optional"`
}

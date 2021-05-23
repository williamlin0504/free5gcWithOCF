package ngapType

import " free5gcWithOCF/lib/aper"

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceToReleaseItemHOCmd struct {
	PDUSessionID                            PDUSessionID
	HandoverPreparationUnsuccessfulTransfer aper.OctetString
	IEExtensions                            *ProtocolExtensionContainerPDUSessionResourceToReleaseItemHOCmdExtIEs `aper:"optional"`
}

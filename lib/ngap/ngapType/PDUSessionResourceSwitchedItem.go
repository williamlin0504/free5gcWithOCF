package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceSwitchedItem struct {
	PDUSessionID                         PDUSessionID
	PathSwitchRequestAcknowledgeTransfer aper.OctetString
	IEExtensions                         *ProtocolExtensionContainerPDUSessionResourceSwitchedItemExtIEs `aper:"optional"`
}

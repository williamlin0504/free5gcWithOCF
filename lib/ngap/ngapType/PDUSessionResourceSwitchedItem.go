package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceSwitchedItem struct {
	PDUSessionID                         PDUSessionID
	PathSwitchRequestAcknowledgeTransfer aper.OctetString
	IEExtensions                         *ProtocolExtensionContainerPDUSessionResourceSwitchedItemExtIEs `aper:"optional"`
}

package ngapType

import " free5gcWithOCF/lib/aper"

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceToBeSwitchedDLItem struct {
	PDUSessionID              PDUSessionID
	PathSwitchRequestTransfer aper.OctetString
	IEExtensions              *ProtocolExtensionContainerPDUSessionResourceToBeSwitchedDLItemExtIEs `aper:"optional"`
}

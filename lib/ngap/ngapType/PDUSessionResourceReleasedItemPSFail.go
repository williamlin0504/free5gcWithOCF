package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceReleasedItemPSFail struct {
	PDUSessionID                          PDUSessionID
	PathSwitchRequestUnsuccessfulTransfer aper.OctetString
	IEExtensions                          *ProtocolExtensionContainerPDUSessionResourceReleasedItemPSFailExtIEs `aper:"optional"`
}

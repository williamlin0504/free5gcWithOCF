package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceItemHORqd struct {
	PDUSessionID             PDUSessionID
	HandoverRequiredTransfer aper.OctetString
	IEExtensions             *ProtocolExtensionContainerPDUSessionResourceItemHORqdExtIEs `aper:"optional"`
}

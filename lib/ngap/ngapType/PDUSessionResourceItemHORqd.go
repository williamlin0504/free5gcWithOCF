package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceItemHORqd struct {
	PDUSessionID             PDUSessionID
	HandoverRequiredTransfer aper.OctetString
	IEExtensions             *ProtocolExtensionContainerPDUSessionResourceItemHORqdExtIEs `aper:"optional"`
}

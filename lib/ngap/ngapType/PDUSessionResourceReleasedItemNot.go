package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceReleasedItemNot struct {
	PDUSessionID                             PDUSessionID
	PDUSessionResourceNotifyReleasedTransfer aper.OctetString
	IEExtensions                             *ProtocolExtensionContainerPDUSessionResourceReleasedItemNotExtIEs `aper:"optional"`
}

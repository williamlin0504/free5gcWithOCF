package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type PDUSessionResourceFailedToModifyItemModRes struct {
	PDUSessionID                                 PDUSessionID
	PDUSessionResourceModifyUnsuccessfulTransfer aper.OctetString
	IEExtensions                                 *ProtocolExtensionContainerPDUSessionResourceFailedToModifyItemModResExtIEs `aper:"optional"`
}

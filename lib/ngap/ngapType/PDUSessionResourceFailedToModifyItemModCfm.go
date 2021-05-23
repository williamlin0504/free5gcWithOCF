package ngapType

import " free5gcWithOCF/lib/aper"

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceFailedToModifyItemModCfm struct {
	PDUSessionID                                           PDUSessionID
	PDUSessionResourceModifyIndicationUnsuccessfulTransfer aper.OctetString
	IEExtensions                                           *ProtocolExtensionContainerPDUSessionResourceFailedToModifyItemModCfmExtIEs `aper:"optional"`
}

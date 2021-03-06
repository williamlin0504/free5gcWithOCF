package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceFailedToSetupItemCxtRes struct {
	PDUSessionID                                PDUSessionID
	PDUSessionResourceSetupUnsuccessfulTransfer aper.OctetString
	IEExtensions                                *ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemCxtResExtIEs `aper:"optional"`
}

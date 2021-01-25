package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceItemCxtRelCpl struct {
	PDUSessionID PDUSessionID
	IEExtensions *ProtocolExtensionContainerPDUSessionResourceItemCxtRelCplExtIEs `aper:"optional"`
}

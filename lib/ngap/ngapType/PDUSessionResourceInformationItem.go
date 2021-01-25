package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionResourceInformationItem struct {
	PDUSessionID              PDUSessionID
	QosFlowInformationList    QosFlowInformationList
	DRBsToQosFlowsMappingList *DRBsToQosFlowsMappingList                                         `aper:"optional"`
	IEExtensions              *ProtocolExtensionContainerPDUSessionResourceInformationItemExtIEs `aper:"optional"`
}

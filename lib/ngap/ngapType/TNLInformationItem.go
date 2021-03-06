package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type TNLInformationItem struct {
	QosFlowPerTNLInformation QosFlowPerTNLInformation                            `aper:"valueExt"`
	IEExtensions             *ProtocolExtensionContainerTNLInformationItemExtIEs `aper:"optional"`
}

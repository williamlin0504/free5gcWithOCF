package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type TNLAssociationItem struct {
	TNLAssociationAddress CPTransportLayerInformation                         `aper:"valueLB:0,valueUB:1"`
	Cause                 Cause                                               `aper:"valueLB:0,valueUB:5"`
	IEExtensions          *ProtocolExtensionContainerTNLAssociationItemExtIEs `aper:"optional"`
}

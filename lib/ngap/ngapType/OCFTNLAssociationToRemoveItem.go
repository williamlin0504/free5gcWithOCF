package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type OCFTNLAssociationToRemoveItem struct {
	OCFTNLAssociationAddress CPTransportLayerInformation                                    `aper:"valueLB:0,valueUB:1"`
	IEExtensions             *ProtocolExtensionContainerOCFTNLAssociationToRemoveItemExtIEs `aper:"optional"`
}

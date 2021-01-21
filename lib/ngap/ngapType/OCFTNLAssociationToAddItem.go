package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type OCFTNLAssociationToAddItem struct {
	OCFTNLAssociationAddress CPTransportLayerInformation `aper:"valueLB:0,valueUB:1"`
	TNLAssociationUsage      *TNLAssociationUsage        `aper:"optional"`
	TNLAddressWeightFactor   TNLAddressWeightFactor
	IEExtensions             *ProtocolExtensionContainerOCFTNLAssociationToAddItemExtIEs `aper:"optional"`
}

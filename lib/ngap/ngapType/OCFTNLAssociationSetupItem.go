package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type OCFTNLAssociationSetupItem struct {
	OCFTNLAssociationAddress CPTransportLayerInformation                                 `aper:"valueLB:0,valueUB:1"`
	IEExtensions             *ProtocolExtensionContainerOCFTNLAssociationSetupItemExtIEs `aper:"optional"`
}
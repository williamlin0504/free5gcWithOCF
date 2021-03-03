package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type UEAssociatedLogicalNGConnectionItem struct {
	AmfUENGAPID  *AmfUENGAPID                                                         `aper:"optional"`
	RANUENGAPID  *RANUENGAPID                                                         `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerUEAssociatedLogicalNGConnectionItemExtIEs `aper:"optional"`
}

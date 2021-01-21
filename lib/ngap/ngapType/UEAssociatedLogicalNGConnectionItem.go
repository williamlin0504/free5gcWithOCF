package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type UEAssociatedLogicalNGConnectionItem struct {
	AMFUENGAPID  *AMFUENGAPID                                                         `aper:"optional"`
	OCFUENGAPID  *OCFUENGAPID                                                         `aper:"optional"`
	RANUENGAPID  *RANUENGAPID                                                         `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerUEAssociatedLogicalNGConnectionItemExtIEs `aper:"optional"`
}

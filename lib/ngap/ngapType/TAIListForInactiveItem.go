package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type TAIListForInactiveItem struct {
	TAI          TAI                                                     `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerTAIListForInactiveItemExtIEs `aper:"optional"`
}

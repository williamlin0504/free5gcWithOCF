package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type AreaOfInterestItem struct {
	AreaOfInterest               AreaOfInterest `aper:"valueExt"`
	LocationReportingReferenceID LocationReportingReferenceID
	IEExtensions                 *ProtocolExtensionContainerAreaOfInterestItemExtIEs `aper:"optional"`
}

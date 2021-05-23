package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type CancelledCellsInTAIEUTRAItem struct {
	EUTRACGI           EUTRACGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCancelledCellsInTAIEUTRAItemExtIEs `aper:"optional"`
}

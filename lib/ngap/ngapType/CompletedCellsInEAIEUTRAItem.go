package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type CompletedCellsInEAIEUTRAItem struct {
	EUTRACGI     EUTRACGI                                                      `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCompletedCellsInEAIEUTRAItemExtIEs `aper:"optional"`
}

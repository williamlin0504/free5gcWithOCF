package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type CompletedCellsInTAINRItem struct {
	NRCGI        NRCGI                                                      `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCompletedCellsInTAINRItemExtIEs `aper:"optional"`
}

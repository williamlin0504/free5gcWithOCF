package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type TAICancelledEUTRAItem struct {
	TAI                      TAI `aper:"valueExt"`
	CancelledCellsInTAIEUTRA CancelledCellsInTAIEUTRA
	IEExtensions             *ProtocolExtensionContainerTAICancelledEUTRAItemExtIEs `aper:"optional"`
}

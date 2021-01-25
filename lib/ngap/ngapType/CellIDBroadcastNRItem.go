package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type CellIDBroadcastNRItem struct {
	NRCGI        NRCGI                                                  `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCellIDBroadcastNRItemExtIEs `aper:"optional"`
}

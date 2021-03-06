package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type CellIDCancelledNRItem struct {
	NRCGI              NRCGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCellIDCancelledNRItemExtIEs `aper:"optional"`
}

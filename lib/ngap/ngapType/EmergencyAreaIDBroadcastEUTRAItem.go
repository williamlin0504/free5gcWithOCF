package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type EmergencyAreaIDBroadcastEUTRAItem struct {
	EmergencyAreaID          EmergencyAreaID
	CompletedCellsInEAIEUTRA CompletedCellsInEAIEUTRA
	IEExtensions             *ProtocolExtensionContainerEmergencyAreaIDBroadcastEUTRAItemExtIEs `aper:"optional"`
}

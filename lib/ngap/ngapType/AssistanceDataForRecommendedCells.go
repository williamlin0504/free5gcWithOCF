package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type AssistanceDataForRecommendedCells struct {
	RecommendedCellsForPaging RecommendedCellsForPaging                                          `aper:"valueExt"`
	IEExtensions              *ProtocolExtensionContainerAssistanceDataForRecommendedCellsExtIEs `aper:"optional"`
}

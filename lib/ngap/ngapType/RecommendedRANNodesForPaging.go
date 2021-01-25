package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type RecommendedRANNodesForPaging struct {
	RecommendedRANNodeList RecommendedRANNodeList
	IEExtensions           *ProtocolExtensionContainerRecommendedRANNodesForPagingExtIEs `aper:"optional"`
}

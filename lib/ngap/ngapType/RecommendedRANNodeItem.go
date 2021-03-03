package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type RecommendedRANNodeItem struct {
	OCFPagingTarget OCFPagingTarget                                         `aper:"valueLB:0,valueUB:2"`
	IEExtensions    *ProtocolExtensionContainerRecommendedRANNodeItemExtIEs `aper:"optional"`
}

package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type UEAggregateMaximumBitRate struct {
	UEAggregateMaximumBitRateDL BitRate
	UEAggregateMaximumBitRateUL BitRate
	IEExtensions                *ProtocolExtensionContainerUEAggregateMaximumBitRateExtIEs `aper:"optional"`
}

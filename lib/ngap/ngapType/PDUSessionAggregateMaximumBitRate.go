package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type PDUSessionAggregateMaximumBitRate struct {
	PDUSessionAggregateMaximumBitRateDL BitRate
	PDUSessionAggregateMaximumBitRateUL BitRate
	IEExtensions                        *ProtocolExtensionContainerPDUSessionAggregateMaximumBitRateExtIEs `aper:"optional"`
}

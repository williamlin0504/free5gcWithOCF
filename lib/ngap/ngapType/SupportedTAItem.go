package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type SupportedTAItem struct {
	TAC               TAC
	BroadcastPLMNList BroadcastPLMNList
	IEExtensions      *ProtocolExtensionContainerSupportedTAItemExtIEs `aper:"optional"`
}

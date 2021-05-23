package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type FiveGSTMSI struct {
	AMFSetID     AMFSetID
	AMFPointer   AMFPointer
	FiveGTMSI    FiveGTMSI
	IEExtensions *ProtocolExtensionContainerFiveGSTMSIExtIEs `aper:"optional"`
}

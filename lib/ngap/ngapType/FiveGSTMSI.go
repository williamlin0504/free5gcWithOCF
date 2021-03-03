package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type FiveGSTMSI struct {
	OCFSetID     OCFSetID
	OCFPointer   OCFPointer
	FiveGTMSI    FiveGTMSI
	IEExtensions *ProtocolExtensionContainerFiveGSTMSIExtIEs `aper:"optional"`
}

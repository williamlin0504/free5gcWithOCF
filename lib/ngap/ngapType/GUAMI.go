package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type GUAMI struct {
	PLMNIdentity PLMNIdentity
	OCFRegionID  OCFRegionID
	OCFSetID     OCFSetID
	OCFPointer   OCFPointer
	IEExtensions *ProtocolExtensionContainerGUAMIExtIEs `aper:"optional"`
}

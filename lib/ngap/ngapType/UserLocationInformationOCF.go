package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type UserLocationInformationOCF struct {
	IPAddress    TransportLayerAddress
	PortNumber   PortNumber
	IEExtensions *ProtocolExtensionContainerUserLocationInformationOCFExtIEs `aper:"optional"`
}

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type UENGAPIDPair struct {
	OCFUENGAPID  OCFUENGAPID
	RANUENGAPID  RANUENGAPID
	IEExtensions *ProtocolExtensionContainerUENGAPIDPairExtIEs `aper:"optional"`
}

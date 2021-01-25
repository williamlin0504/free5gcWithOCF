package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type UENGAPIDPair struct {
	AMFUENGAPID  AMFUENGAPID
	OCFUENGAPID  OCFUENGAPID
	RANUENGAPID  RANUENGAPID
	IEExtensions *ProtocolExtensionContainerUENGAPIDPairExtIEs `aper:"optional"`
}

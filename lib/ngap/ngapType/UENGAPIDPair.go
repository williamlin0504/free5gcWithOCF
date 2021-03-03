package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type UENGAPIDPair struct {
	AmfUENGAPID  AmfUENGAPID
	RANUENGAPID  RANUENGAPID
	IEExtensions *ProtocolExtensionContainerUENGAPIDPairExtIEs `aper:"optional"`
}

package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type GlobalOCFID struct {
	PLMNIdentity PLMNIdentity
	OCFID        OCFID                                        `aper:"valueLB:0,valueUB:1"`
	IEExtensions *ProtocolExtensionContainerGlobalOCFIDExtIEs `aper:"optional"`
}

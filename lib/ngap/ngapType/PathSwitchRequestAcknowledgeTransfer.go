package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type PathSwitchRequestAcknowledgeTransfer struct {
	ULNGUUPTNLInformation *UPTransportLayerInformation                                          `aper:"valueLB:0,valueUB:1,optional"`
	SecurityIndication    *SecurityIndication                                                   `aper:"valueExt,optional"`
	IEExtensions          *ProtocolExtensionContainerPathSwitchRequestAcknowledgeTransferExtIEs `aper:"optional"`
}

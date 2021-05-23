package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type TraceActivation struct {
	NGRANTraceID                   NGRANTraceID
	InterfacesToTrace              InterfacesToTrace
	TraceDepth                     TraceDepth
	TraceCollectionEntityIPAddress TransportLayerAddress
	IEExtensions                   *ProtocolExtensionContainerTraceActivationExtIEs `aper:"optional"`
}

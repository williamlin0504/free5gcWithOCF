package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type CriticalityDiagnosticsIEItem struct {
	IECriticality Criticality
	IEID          ProtocolIEID
	TypeOfError   TypeOfError
	IEExtensions  *ProtocolExtensionContainerCriticalityDiagnosticsIEItemExtIEs `aper:"optional"`
}

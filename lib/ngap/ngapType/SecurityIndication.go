package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type SecurityIndication struct {
	IntegrityProtectionIndication       IntegrityProtectionIndication
	ConfidentialityProtectionIndication ConfidentialityProtectionIndication
	MaximumIntegrityProtectedDataRate   *MaximumIntegrityProtectedDataRate                  `aper:"optional"`
	IEExtensions                        *ProtocolExtensionContainerSecurityIndicationExtIEs `aper:"optional"`
}

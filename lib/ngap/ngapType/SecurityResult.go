package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type SecurityResult struct {
	IntegrityProtectionResult       IntegrityProtectionResult
	ConfidentialityProtectionResult ConfidentialityProtectionResult
	IEExtensions                    *ProtocolExtensionContainerSecurityResultExtIEs `aper:"optional"`
}

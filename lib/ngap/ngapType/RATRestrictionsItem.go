package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type RATRestrictionsItem struct {
	PLMNIdentity              PLMNIdentity
	RATRestrictionInformation RATRestrictionInformation
	IEExtensions              *ProtocolExtensionContainerRATRestrictionsItemExtIEs `aper:"optional"`
}

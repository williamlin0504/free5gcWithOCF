package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

const (
	UEPagingIdentityPresentNothing int = iota /* No components present */
	UEPagingIdentityPresentFiveGSTMSI
	UEPagingIdentityPresentChoiceExtensions
)

type UEPagingIdentity struct {
	Present          int
	FiveGSTMSI       *FiveGSTMSI `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerUEPagingIdentityExtIEs
}

package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

const (
	UENGAPIDsPresentNothing int = iota /* No components present */
	UENGAPIDsPresentUENGAPIDPair
	UENGAPIDsPresentAMFUENGAPID
	UENGAPIDsPresentOCFUENGAPID
	UENGAPIDsPresentChoiceExtensions
)

type UENGAPIDs struct {
	Present          int
	UENGAPIDPair     *UENGAPIDPair `aper:"valueExt"`
	AMFUENGAPID      *AMFUENGAPID
	OCFUENGAPID      *OCFUENGAPID
	ChoiceExtensions *ProtocolIESingleContainerUENGAPIDsExtIEs
}

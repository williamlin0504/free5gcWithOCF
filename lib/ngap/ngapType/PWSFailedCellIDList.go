package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

const (
	PWSFailedCellIDListPresentNothing int = iota /* No components present */
	PWSFailedCellIDListPresentEUTRACGIPWSFailedList
	PWSFailedCellIDListPresentNRCGIPWSFailedList
	PWSFailedCellIDListPresentChoiceExtensions
)

type PWSFailedCellIDList struct {
	Present               int
	EUTRACGIPWSFailedList *EUTRACGIList
	NRCGIPWSFailedList    *NRCGIList
	ChoiceExtensions      *ProtocolIESingleContainerPWSFailedCellIDListExtIEs
}

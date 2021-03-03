package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

const (
	OCFIDPresentNothing int = iota /* No components present */
	OCFIDPresentOCFID
	OCFIDPresentChoiceExtensions
)

type OCFID struct {
	Present          int
	OCFID            *aper.BitString `aper:"sizeLB:16,sizeUB:16"`
	ChoiceExtensions *ProtocolIESingleContainerOCFIDExtIEs
}

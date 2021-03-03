package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

const (
	N3IWFIDPresentNothing int = iota /* No components present */
	N3IWFIDPresentN3IWFID
	N3IWFIDPresentChoiceExtensions
	OCFIDPresentNothing int = iota /* No components present */
	OCFIDPresentN3IWFID
	OCFIDPresentChoiceExtensions
)

type N3IWFID struct {
	Present          int
	N3IWFID          *aper.BitString `aper:"sizeLB:16,sizeUB:16"`
	ChoiceExtensions *ProtocolIESingleContainerN3IWFIDExtIEs
}

type OCFID struct {
	Present          int
	OCFID            *aper.BitString `aper:"sizeLB:16,sizeUB:16"`
	ChoiceExtensions *ProtocolIESingleContainerOCFIDExtIEs
}

package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

const (
	MICOModeIndicationPresentTrue aper.Enumerated = 0
)

type MICOModeIndication struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}

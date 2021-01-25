package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

const (
	TypeOfErrorPresentNotUnderstood aper.Enumerated = 0
	TypeOfErrorPresentMissing       aper.Enumerated = 1
)

type TypeOfError struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}

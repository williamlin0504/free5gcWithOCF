package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

const (
	TypeOfErrorPresentNotUnderstood aper.Enumerated = 0
	TypeOfErrorPresentMissing       aper.Enumerated = 1
)

type TypeOfError struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}

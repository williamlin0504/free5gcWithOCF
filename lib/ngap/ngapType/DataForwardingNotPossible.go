package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

const (
	DataForwardingNotPossiblePresentDataForwardingNotPossible aper.Enumerated = 0
)

type DataForwardingNotPossible struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}

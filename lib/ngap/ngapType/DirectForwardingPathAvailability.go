package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

const (
	DirectForwardingPathAvailabilityPresentDirectPathAvailable aper.Enumerated = 0
)

type DirectForwardingPathAvailability struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}

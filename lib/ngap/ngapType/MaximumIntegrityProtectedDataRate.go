package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

const (
	MaximumIntegrityProtectedDataRatePresentBitrate64kbs  aper.Enumerated = 0
	MaximumIntegrityProtectedDataRatePresentMaximumUERate aper.Enumerated = 1
)

type MaximumIntegrityProtectedDataRate struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}

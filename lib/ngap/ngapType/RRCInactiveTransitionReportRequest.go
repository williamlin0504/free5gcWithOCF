package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

const (
	RRCInactiveTransitionReportRequestPresentSubsequentStateTransitionReport aper.Enumerated = 0
	RRCInactiveTransitionReportRequestPresentSingleRrcConnectedStateReport   aper.Enumerated = 1
	RRCInactiveTransitionReportRequestPresentCancelReport                    aper.Enumerated = 2
)

type RRCInactiveTransitionReportRequest struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:2"`
}

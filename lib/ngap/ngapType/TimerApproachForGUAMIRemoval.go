package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

const (
	TimerApproachForGUAMIRemovalPresentApplyTimer aper.Enumerated = 0
)

type TimerApproachForGUAMIRemoval struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}

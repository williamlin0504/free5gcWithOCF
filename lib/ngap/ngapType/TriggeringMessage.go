package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

const (
	TriggeringMessagePresentInitiatingMessage    aper.Enumerated = 0
	TriggeringMessagePresentSuccessfulOutcome    aper.Enumerated = 1
	TriggeringMessagePresentUnsuccessfullOutcome aper.Enumerated = 2
)

type TriggeringMessage struct {
	Value aper.Enumerated `aper:"valueLB:0,valueUB:2"`
}

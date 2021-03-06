package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

const (
	PreEmptionCapabilityPresentShallNotTriggerPreEmption aper.Enumerated = 0
	PreEmptionCapabilityPresentMayTriggerPreEmption      aper.Enumerated = 1
)

type PreEmptionCapability struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}

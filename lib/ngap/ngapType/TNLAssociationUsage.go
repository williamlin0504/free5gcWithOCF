package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

const (
	TNLAssociationUsagePresentUe    aper.Enumerated = 0
	TNLAssociationUsagePresentNonUe aper.Enumerated = 1
	TNLAssociationUsagePresentBoth  aper.Enumerated = 2
)

type TNLAssociationUsage struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:2"`
}

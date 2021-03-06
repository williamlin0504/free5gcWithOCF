package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type WarningType struct {
	Value aper.OctetString `aper:"sizeLB:2,sizeUB:2"`
}

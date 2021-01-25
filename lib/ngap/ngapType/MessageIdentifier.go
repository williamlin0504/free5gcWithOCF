package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type MessageIdentifier struct {
	Value aper.BitString `aper:"sizeLB:16,sizeUB:16"`
}

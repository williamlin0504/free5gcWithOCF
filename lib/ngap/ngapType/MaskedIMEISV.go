package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type MaskedIMEISV struct {
	Value aper.BitString `aper:"sizeLB:64,sizeUB:64"`
}

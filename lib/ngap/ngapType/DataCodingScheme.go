package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type DataCodingScheme struct {
	Value aper.BitString `aper:"sizeLB:8,sizeUB:8"`
}

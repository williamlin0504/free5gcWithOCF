package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type SST struct {
	Value aper.OctetString `aper:"sizeLB:1,sizeUB:1"`
}

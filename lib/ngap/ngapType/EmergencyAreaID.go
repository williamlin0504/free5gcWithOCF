package ngapType

import "free5gcWithOCF/lib/aper"

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type EmergencyAreaID struct {
	Value aper.OctetString `aper:"sizeLB:3,sizeUB:3"`
}

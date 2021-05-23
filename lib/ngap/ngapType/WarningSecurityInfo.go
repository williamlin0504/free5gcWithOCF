package ngapType

import " free5gcWithOCF/lib/aper"

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type WarningSecurityInfo struct {
	Value aper.OctetString `aper:"sizeLB:50,sizeUB:50"`
}

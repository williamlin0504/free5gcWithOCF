package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type WarningSecurityInfo struct {
	Value aper.OctetString `aper:"sizeLB:50,sizeUB:50"`
}

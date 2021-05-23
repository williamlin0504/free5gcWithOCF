package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type FiveGTMSI struct {
	Value aper.OctetString `aper:"sizeLB:4,sizeUB:4"`
}

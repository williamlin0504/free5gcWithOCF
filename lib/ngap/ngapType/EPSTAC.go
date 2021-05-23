package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type EPSTAC struct {
	Value aper.OctetString `aper:"sizeLB:2,sizeUB:2"`
}

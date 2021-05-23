package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type SerialNumber struct {
	Value aper.BitString `aper:"sizeLB:16,sizeUB:16"`
}

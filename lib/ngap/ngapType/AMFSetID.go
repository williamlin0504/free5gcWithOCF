package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type AMFSetID struct {
	Value aper.BitString `aper:"sizeLB:10,sizeUB:10"`
}

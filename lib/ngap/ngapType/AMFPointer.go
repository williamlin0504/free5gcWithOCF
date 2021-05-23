package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type AMFPointer struct {
	Value aper.BitString `aper:"sizeLB:6,sizeUB:6"`
}

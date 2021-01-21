package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

type OCFPointer struct {
	Value aper.BitString `aper:"sizeLB:6,sizeUB:6"`
}

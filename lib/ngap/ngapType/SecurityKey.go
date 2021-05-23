package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type SecurityKey struct {
	Value aper.BitString `aper:"sizeLB:256,sizeUB:256"`
}

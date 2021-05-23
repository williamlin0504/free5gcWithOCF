package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type NRCellIdentity struct {
	Value aper.BitString `aper:"sizeLB:36,sizeUB:36"`
}

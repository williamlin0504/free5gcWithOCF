package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type EUTRACellIdentity struct {
	Value aper.BitString `aper:"sizeLB:28,sizeUB:28"`
}

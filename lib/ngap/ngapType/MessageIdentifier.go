package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type MessageIdentifier struct {
	Value aper.BitString `aper:"sizeLB:16,sizeUB:16"`
}

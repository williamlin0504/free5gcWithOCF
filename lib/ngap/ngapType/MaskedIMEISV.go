package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type MaskedIMEISV struct {
	Value aper.BitString `aper:"sizeLB:64,sizeUB:64"`
}

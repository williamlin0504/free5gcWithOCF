package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type RATRestrictionInformation struct {
	Value aper.BitString `aper:"sizeExt,sizeLB:8,sizeUB:8"`
}

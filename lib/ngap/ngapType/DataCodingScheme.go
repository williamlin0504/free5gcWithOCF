package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type DataCodingScheme struct {
	Value aper.BitString `aper:"sizeLB:8,sizeUB:8"`
}

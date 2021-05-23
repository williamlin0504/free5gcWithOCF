package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type TransportLayerAddress struct {
	Value aper.BitString `aper:"sizeExt,sizeLB:1,sizeUB:160"`
}

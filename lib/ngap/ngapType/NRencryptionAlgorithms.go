package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type NRencryptionAlgorithms struct {
	Value aper.BitString `aper:"sizeExt,sizeLB:16,sizeUB:16"`
}

package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type NGRANTraceID struct {
	Value aper.OctetString `aper:"sizeLB:8,sizeUB:8"`
}

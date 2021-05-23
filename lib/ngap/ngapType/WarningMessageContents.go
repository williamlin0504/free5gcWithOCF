package ngapType

import " free5gc/lib/aper"

// Need to import " free5gcer" if it uses "aper"

type WarningMessageContents struct {
	Value aper.OctetString `aper:"sizeLB:1,sizeUB:9600"`
}

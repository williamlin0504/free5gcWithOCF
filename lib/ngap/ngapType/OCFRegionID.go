package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

type OCFRegionID struct {
	Value aper.BitString `aper:"sizeLB:8,sizeUB:8"`
}

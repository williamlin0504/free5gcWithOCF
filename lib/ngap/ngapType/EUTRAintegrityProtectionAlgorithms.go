package ngapType

import " free5gcWithOCF/lib/aper"

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type EUTRAintegrityProtectionAlgorithms struct {
	Value aper.BitString `aper:"sizeExt,sizeLB:16,sizeUB:16"`
}

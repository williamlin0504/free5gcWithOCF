package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type NetworkInstance struct {
	Value int64 `aper:"valueExt,valueLB:1,valueUB:256"`
}

package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type OCFUENGAPID struct {
	Value int64 `aper:"valueLB:0,valueUB:1099511627775"`
}

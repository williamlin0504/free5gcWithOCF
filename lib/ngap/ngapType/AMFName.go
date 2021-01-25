package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type AMFName struct {
	Value string `aper:"sizeExt,sizeLB:1,sizeUB:150"`
}

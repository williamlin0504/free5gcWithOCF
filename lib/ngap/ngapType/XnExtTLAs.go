package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct XnExtTLAs */
/* XnExtTLAItem */
type XnExtTLAs struct {
	List []XnExtTLAItem `aper:"valueExt,sizeLB:1,sizeUB:2"`
}

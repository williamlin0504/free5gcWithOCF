package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct SliceOverloadList */
/* SliceOverloadItem */
type SliceOverloadList struct {
	List []SliceOverloadItem `aper:"valueExt,sizeLB:1,sizeUB:1024"`
}

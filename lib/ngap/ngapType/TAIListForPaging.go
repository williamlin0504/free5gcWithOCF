package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct TAIListForPaging */
/* TAIListForPagingItem */
type TAIListForPaging struct {
	List []TAIListForPagingItem `aper:"valueExt,sizeLB:1,sizeUB:16"`
}
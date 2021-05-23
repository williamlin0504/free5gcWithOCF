package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct NR_CGIListForWarning */
/* NRCGI */
type NRCGIListForWarning struct {
	List []NRCGI `aper:"valueExt,sizeLB:1,sizeUB:65535"`
}

package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct CompletedCellsInTAI_NR */
/* CompletedCellsInTAINRItem */
type CompletedCellsInTAINR struct {
	List []CompletedCellsInTAINRItem `aper:"valueExt,sizeLB:1,sizeUB:65535"`
}

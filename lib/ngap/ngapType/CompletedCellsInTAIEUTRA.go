package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct CompletedCellsInTAI_EUTRA */
/* CompletedCellsInTAIEUTRAItem */
type CompletedCellsInTAIEUTRA struct {
	List []CompletedCellsInTAIEUTRAItem `aper:"valueExt,sizeLB:1,sizeUB:65535"`
}
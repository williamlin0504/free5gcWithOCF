package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct ForbiddenTACs */
/* TAC */
type ForbiddenTACs struct {
	List []TAC `aper:"sizeLB:1,sizeUB:4096"`
}

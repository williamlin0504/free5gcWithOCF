package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceFailedToModifyListModRes */
/* PDUSessionResourceFailedToModifyItemModRes */
type PDUSessionResourceFailedToModifyListModRes struct {
	List []PDUSessionResourceFailedToModifyItemModRes `aper:"valueExt,sizeLB:1,sizeUB:256"`
}

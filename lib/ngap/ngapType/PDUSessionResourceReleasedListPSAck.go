package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceReleasedListPSAck */
/* PDUSessionResourceReleasedItemPSAck */
type PDUSessionResourceReleasedListPSAck struct {
	List []PDUSessionResourceReleasedItemPSAck `aper:"valueExt,sizeLB:1,sizeUB:256"`
}

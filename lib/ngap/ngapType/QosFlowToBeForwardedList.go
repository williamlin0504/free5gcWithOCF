package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct QosFlowToBeForwardedList */
/* QosFlowToBeForwardedItem */
type QosFlowToBeForwardedList struct {
	List []QosFlowToBeForwardedItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}

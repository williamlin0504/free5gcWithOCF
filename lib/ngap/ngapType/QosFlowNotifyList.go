package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct QosFlowNotifyList */
/* QosFlowNotifyItem */
type QosFlowNotifyList struct {
	List []QosFlowNotifyItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}

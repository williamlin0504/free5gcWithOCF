package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct QosFlowSetupResponseListHOReqAck */
/* QosFlowSetupResponseItemHOReqAck */
type QosFlowSetupResponseListHOReqAck struct {
	List []QosFlowSetupResponseItemHOReqAck `aper:"valueExt,sizeLB:1,sizeUB:64"`
}

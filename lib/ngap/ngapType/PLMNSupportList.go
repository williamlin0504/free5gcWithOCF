package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PLMNSupportList */
/* PLMNSupportItem */
type PLMNSupportList struct {
	List []PLMNSupportItem `aper:"valueExt,sizeLB:1,sizeUB:12"`
}

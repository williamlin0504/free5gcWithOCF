package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct OCF_TNLAssociationSetupList */
/* OCFTNLAssociationSetupItem */
type OCFTNLAssociationSetupList struct {
	List []OCFTNLAssociationSetupItem `aper:"valueExt,sizeLB:1,sizeUB:32"`
}

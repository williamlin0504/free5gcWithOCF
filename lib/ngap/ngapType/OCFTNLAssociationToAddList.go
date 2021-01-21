package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct OCF_TNLAssociationToAddList */
/* OCFTNLAssociationToAddItem */
type OCFTNLAssociationToAddList struct {
	List []OCFTNLAssociationToAddItem `aper:"valueExt,sizeLB:1,sizeUB:32"`
}

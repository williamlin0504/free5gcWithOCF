package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct OCF_TNLAssociationToRemoveList */
/* OCFTNLAssociationToRemoveItem */
type OCFTNLAssociationToRemoveList struct {
	List []OCFTNLAssociationToRemoveItem `aper:"valueExt,sizeLB:1,sizeUB:32"`
}

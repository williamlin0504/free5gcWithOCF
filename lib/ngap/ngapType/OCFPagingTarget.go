package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

const (
	OCFPagingTargetPresentNothing int = iota /* No components present */
	OCFPagingTargetPresentGlobalRANNodeID
	OCFPagingTargetPresentTAI
	OCFPagingTargetPresentChoiceExtensions
)

type OCFPagingTarget struct {
	Present          int
	GlobalRANNodeID  *GlobalRANNodeID `aper:"valueLB:0,valueUB:3"`
	TAI              *TAI             `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerOCFPagingTargetExtIEs
}

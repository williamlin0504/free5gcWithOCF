package ngapType

// Need to import " free5gcWithOCF/lib/aper" if it uses "aper"

type PriorityLevelARP struct {
	Value int64 `aper:"valueLB:1,valueUB:15"`
}

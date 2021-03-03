package context

type OCFNFInfo struct {
	GlobalOCFID     GlobalOCFID       `yaml:"GlobalOCFID"`
	RanNodeName     string            `yaml:"Name,omitempty"`
	SupportedTAList []SupportedTAItem `yaml:"SupportedTAList"`
}

type GlobalOCFID struct {
	PLMNID PLMNID `yaml:"PLMNID"`
	OCFID  uint16 `yaml:"OCFID"` // with length 2 bytes
}

type SupportedTAItem struct {
	TAC               string              `yaml:"TAC"`
	BroadcastPLMNList []BroadcastPLMNItem `yaml:"BroadcastPLMNList"`
}

type BroadcastPLMNItem struct {
	PLMNID              PLMNID             `yaml:"PLMNID"`
	TAISliceSupportList []SliceSupportItem `yaml:"TAISliceSupportList"`
}

type PLMNID struct {
	Mcc string `yaml:"MCC"`
	Mnc string `yaml:"MNC"`
}

type SliceSupportItem struct {
	SNSSAI SNSSAIItem `yaml:"SNSSAI"`
}

type SNSSAIItem struct {
	SST string `yaml:"SST"`
	SD  string `yaml:"SD,omitempty"`
}

type OCFSCTPAddresses struct {
	IPAddresses []string `yaml:"IP"`
	Port        int      `yaml:"Port,omitempty"`
}

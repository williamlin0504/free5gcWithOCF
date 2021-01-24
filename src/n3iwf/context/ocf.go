package context

import (
	"bytes"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/ngapConvert"
	"free5gc/lib/ngap/ngapType"

	"git.cs.nctu.edu.tw/calee/sctp"
)

type N3IWFOCF struct {
	SCTPAddr              string
	SCTPConn              *sctp.SCTPConn
	OCFName               *ngapType.OCFName
	ServedGUAMIList       *ngapType.ServedGUAMIList
	RelativeOCFCapacity   *ngapType.RelativeOCFCapacity
	PLMNSupportList       *ngapType.PLMNSupportList
	OCFTNLAssociationList map[string]*OCFTNLAssociationItem // v4+v6 as key
	// Overload related
	OCFOverloadContent *OCFOverloadContent
	// Relative Context
	N3iwfUeList map[int64]*N3IWFUe // ranUeNgapId as key
}

type OCFTNLAssociationItem struct {
	Ipv4                   string
	Ipv6                   string
	TNLAssociationUsage    *ngapType.TNLAssociationUsage
	TNLAddressWeightFactor *int64
}

type OCFOverloadContent struct {
	Action     *ngapType.OverloadAction
	TrafficInd *int64
	NSSAIList  []SliceOverloadItemOCF
}
type SliceOverloadItemOCF struct {
	SNssaiList []ngapType.SNSSAI
	Action     *ngapType.OverloadAction
	TrafficInd *int64
}

func (ocf *N3IWFOCF) init(sctpAddr string, conn *sctp.SCTPConn) {
	ocf.SCTPAddr = sctpAddr
	ocf.SCTPConn = conn
	ocf.OCFTNLAssociationList = make(map[string]*OCFTNLAssociationItem)
	ocf.N3iwfUeList = make(map[int64]*N3IWFUe)
}

func (ocf *N3IWFOCF) FindUeByOcfUeNgapID(id int64) *N3IWFUe {
	for _, n3iwfUe := range ocf.N3iwfUeList {
		if n3iwfUe.OcfUeNgapId == id {
			return n3iwfUe
		}
	}
	return nil
}

func (ocf *N3IWFOCF) RemoveAllRelatedUe() {
	for _, ue := range ocf.N3iwfUeList {
		ue.Remove()
	}
}

func (ocf *N3IWFOCF) AddOCFTNLAssociationItem(info ngapType.CPTransportLayerInformation) *OCFTNLAssociationItem {
	item := &OCFTNLAssociationItem{}
	item.Ipv4, item.Ipv6 = ngapConvert.IPAddressToString(*info.EndpointIPAddress)
	ocf.OCFTNLAssociationList[item.Ipv4+item.Ipv6] = item
	return item
}

func (ocf *N3IWFOCF) FindOCFTNLAssociationItem(info ngapType.CPTransportLayerInformation) *OCFTNLAssociationItem {
	v4, v6 := ngapConvert.IPAddressToString(*info.EndpointIPAddress)
	return ocf.OCFTNLAssociationList[v4+v6]
}

func (ocf *N3IWFOCF) DeleteOCFTNLAssociationItem(info ngapType.CPTransportLayerInformation) {
	v4, v6 := ngapConvert.IPAddressToString(*info.EndpointIPAddress)
	delete(ocf.OCFTNLAssociationList, v4+v6)
}

func (ocf *N3IWFOCF) StartOverload(
	resp *ngapType.OverloadResponse, trafloadInd *ngapType.TrafficLoadReductionIndication,
	nssai *ngapType.OverloadStartNSSAIList) *OCFOverloadContent {
	if resp == nil && trafloadInd == nil && nssai == nil {
		return nil
	}
	content := OCFOverloadContent{}
	if resp != nil {
		content.Action = resp.OverloadAction
	}
	if trafloadInd != nil {
		content.TrafficInd = &trafloadInd.Value
	}
	if nssai != nil {
		for _, item := range nssai.List {
			sliceItem := SliceOverloadItemOCF{}
			for _, item2 := range item.SliceOverloadList.List {
				sliceItem.SNssaiList = append(sliceItem.SNssaiList, item2.SNSSAI)
			}
			if item.SliceOverloadResponse != nil {
				sliceItem.Action = item.SliceOverloadResponse.OverloadAction
			}
			if item.SliceTrafficLoadReductionIndication != nil {
				sliceItem.TrafficInd = &item.SliceTrafficLoadReductionIndication.Value
			}
			content.NSSAIList = append(content.NSSAIList, sliceItem)
		}
	}
	ocf.OCFOverloadContent = &content
	return ocf.OCFOverloadContent
}
func (ocf *N3IWFOCF) StopOverload() {
	ocf.OCFOverloadContent = nil
}

// FindAvalibleOCFByCompareGUAMI compares the incoming GUAMI with OCF served GUAMI
// and return if this OCF is avalible for UE
func (ocf *N3IWFOCF) FindAvalibleOCFByCompareGUAMI(ueSpecifiedGUAMI *ngapType.GUAMI) bool {
	for _, ocfServedGUAMI := range ocf.ServedGUAMIList.List {
		codedOCFServedGUAMI, err := aper.MarshalWithParams(&ocfServedGUAMI.GUAMI, "valueExt")
		if err != nil {
			return false
		}
		codedUESpecifiedGUAMI, err := aper.MarshalWithParams(ueSpecifiedGUAMI, "valueExt")
		if err != nil {
			return false
		}
		if !bytes.Equal(codedOCFServedGUAMI, codedUESpecifiedGUAMI) {
			continue
		}
		return true
	}
	return false
}

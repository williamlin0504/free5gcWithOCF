package context

import (
	"encoding/binary"
	"errors"
	"fmt"
	"free5gcWithOCF/lib/ngap/ngapType"
	ike_message "free5gcWithOCF/src/ocf/ike/message"
	"net"

	gtpv1 "github.com/wmnsk/go-gtp/v1"
)

const (
	AmfUENGAPIDUnspecified int64 = 0xffffffffff
)

type OCFUe struct {
	/* UE identity*/
	RanUeNgapId           int64
	AmfUENGAPID           int64
	IPAddrv4              string
	IPAddrv6              string
	PortNumber            int32
	MaskedIMEISV          *ngapType.MaskedIMEISV // TS 38.413 9.3.1.54
	Guti                  string
	RRCEstablishmentCause int16
	IPSecInnerIP          string

	/* Relative Context */
	AMF *OCFAMF

	/* PDU Session */
	PduSessionList map[int64]*PDUSession // pduSessionId as key

	/* PDU Session Setup Temporary Data */
	TemporaryPDUSessionSetupData *PDUSessionSetupTemporaryData
	/* Temporary cached NAS message */
	// Used when NAS registration accept arrived before
	// UE setup NAS TCP connection with OCF
	TemporaryCachedNASMessage []byte

	/* Security */
	Kocf                 []uint8                          // 32 bytes (256 bits), value is from NGAP IE "Security Key"
	SecurityCapabilities *ngapType.UESecurityCapabilities // TS 38.413 9.3.1.86

	/* IKE Security Association */
	OCFIKESecurityAssociation   *IKESecurityAssociation
	OCFChildSecurityAssociation *ChildSecurityAssociation

	/* NAS IKE Connection */
	IKEConnection *UDPSocketInfo
	/* NAS TCP Connection */
	TCPConnection net.Conn

	/* Others */
	Guami                            *ngapType.GUAMI
	IndexToRfsp                      int64
	Ambr                             *ngapType.UEAggregateMaximumBitRate
	AllowedNssai                     *ngapType.AllowedNSSAI
	RadioCapability                  *ngapType.UERadioCapability                // TODO: This is for RRC, can be deleted
	CoreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation // TS 38.413 9.3.1.15
	IMSVoiceSupported                int32
}

type PDUSession struct {
	Id                               int64 // PDU Session ID
	Type                             *ngapType.PDUSessionType
	Ambr                             *ngapType.PDUSessionAggregateMaximumBitRate
	Snssai                           ngapType.SNSSAI
	NetworkInstance                  *ngapType.NetworkInstance
	SecurityCipher                   bool
	SecurityIntegrity                bool
	MaximumIntegrityDataRateUplink   *ngapType.MaximumIntegrityProtectedDataRate
	MaximumIntegrityDataRateDownlink *ngapType.MaximumIntegrityProtectedDataRate
	GTPConnection                    *GTPConnectionInfo
	QFIList                          []uint8
	QosFlows                         map[int64]*QosFlow // QosFlowIdentifier as key
}

type PDUSessionSetupTemporaryData struct {
	// Slice of unactivated PDU session
	UnactivatedPDUSession []int64 // PDUSessionID as content
	// NGAPProcedureCode is used to identify which type of
	// response shall be used
	NGAPProcedureCode ngapType.ProcedureCode
	// PDU session setup list response
	SetupListCxtRes  *ngapType.PDUSessionResourceSetupListCxtRes
	FailedListCxtRes *ngapType.PDUSessionResourceFailedToSetupListCxtRes
	SetupListSURes   *ngapType.PDUSessionResourceSetupListSURes
	FailedListSURes  *ngapType.PDUSessionResourceFailedToSetupListSURes
}

type QosFlow struct {
	Identifier int64
	Parameters ngapType.QosFlowLevelQosParameters
}

type GTPConnectionInfo struct {
	UPFIPAddr           string
	UPFUDPAddr          net.Addr
	IncomingTEID        uint32
	OutgoingTEID        uint32
	UserPlaneConnection *gtpv1.UPlaneConn
}

type IKESecurityAssociation struct {
	// SPI
	RemoteSPI uint64
	LocalSPI  uint64

	// Message ID
	MessageID uint32

	// Transforms for IKE SA
	EncryptionAlgorithm    *ike_message.Transform
	PseudorandomFunction   *ike_message.Transform
	IntegrityAlgorithm     *ike_message.Transform
	DiffieHellmanGroup     *ike_message.Transform
	ExpandedSequenceNumber *ike_message.Transform

	// Used for key generating
	ConcatenatedNonce      []byte
	DiffieHellmanSharedKey []byte

	// Keys
	SK_d  []byte // used for child SA key deriving
	SK_ai []byte // used by initiator for integrity checking
	SK_ar []byte // used by responder for integrity checking
	SK_ei []byte // used by initiator for encrypting
	SK_er []byte // used by responder for encrypting
	SK_pi []byte // used by initiator for IKE authentication
	SK_pr []byte // used by responder for IKE authentication

	// State for IKE_AUTH
	State uint8

	// Temporary data stored for the use in later exchange
	InitiatorID              *ike_message.IdentificationInitiator
	InitiatorCertificate     *ike_message.Certificate
	IKEAuthResponseSA        *ike_message.SecurityAssociation
	TrafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	TrafficSelectorResponder *ike_message.TrafficSelectorResponder
	LastEAPIdentifier        uint8

	// Authentication data
	LocalUnsignedAuthentication  []byte
	RemoteUnsignedAuthentication []byte

	// UE context
	ThisUE *OCFUe
}

type ChildSecurityAssociation struct {
	// SPI
	SPI uint32

	// IP address
	PeerPublicIPAddr  net.IP
	LocalPublicIPAddr net.IP

	// Traffic selector
	SelectedIPProtocol    uint8
	TrafficSelectorLocal  net.IPNet
	TrafficSelectorRemote net.IPNet

	// Security
	EncryptionAlgorithm               uint16
	InitiatorToResponderEncryptionKey []byte
	ResponderToInitiatorEncryptionKey []byte
	IntegrityAlgorithm                uint16
	InitiatorToResponderIntegrityKey  []byte
	ResponderToInitiatorIntegrityKey  []byte
	ESN                               bool

	// UE context
	ThisUE *OCFUe
}

type UDPSocketInfo struct {
	Conn    *net.UDPConn
	OCFAddr *net.UDPAddr
	UEAddr  *net.UDPAddr
}

func (ue *OCFUe) init(ranUeNgapId int64) {
	ue.RanUeNgapId = ranUeNgapId
	ue.AmfUENGAPID = AmfUENGAPIDUnspecified
	ue.PduSessionList = make(map[int64]*PDUSession)
}

func (ue *OCFUe) Remove() {
	// remove from AMF context
	ue.DetachAMF()
	// remove from OCF context
	ocfSelf := OCFSelf()
	ocfSelf.DeleteOcfUe(ue.RanUeNgapId)
	ocfSelf.DeleteIKESecurityAssociation(ue.OCFIKESecurityAssociation.LocalSPI)
	ocfSelf.DeleteInternalUEIPAddr(ue.IPSecInnerIP)
	for _, pduSession := range ue.PduSessionList {
		ocfSelf.DeleteTEID(pduSession.GTPConnection.IncomingTEID)
	}
}

func (ue *OCFUe) FindPDUSession(pduSessionID int64) *PDUSession {
	if pduSession, ok := ue.PduSessionList[pduSessionID]; ok {
		return pduSession
	} else {
		return nil
	}
}

func (ue *OCFUe) CreatePDUSession(pduSessionID int64, snssai ngapType.SNSSAI) (*PDUSession, error) {
	if _, exists := ue.PduSessionList[pduSessionID]; exists {
		return nil, fmt.Errorf("PDU Session[ID:%d] is already exists", pduSessionID)
	}
	pduSession := &PDUSession{
		Id:       pduSessionID,
		Snssai:   snssai,
		QosFlows: make(map[int64]*QosFlow),
	}
	ue.PduSessionList[pduSessionID] = pduSession
	return pduSession, nil
}

func (ue *OCFUe) CreateIKEChildSecurityAssociation(
	chosenSecurityAssociation *ike_message.SecurityAssociation) (*ChildSecurityAssociation, error) {
	childSecurityAssociation := new(ChildSecurityAssociation)

	if chosenSecurityAssociation == nil {
		return nil, errors.New("chosenSecurityAssociation is nil")
	}

	if len(chosenSecurityAssociation.Proposals) == 0 {
		return nil, errors.New("No proposal")
	}

	childSecurityAssociation.SPI = binary.BigEndian.Uint32(chosenSecurityAssociation.Proposals[0].SPI)

	if len(chosenSecurityAssociation.Proposals[0].EncryptionAlgorithm) != 0 {
		childSecurityAssociation.EncryptionAlgorithm =
			chosenSecurityAssociation.Proposals[0].EncryptionAlgorithm[0].TransformID
	}
	if len(chosenSecurityAssociation.Proposals[0].IntegrityAlgorithm) != 0 {
		childSecurityAssociation.IntegrityAlgorithm =
			chosenSecurityAssociation.Proposals[0].IntegrityAlgorithm[0].TransformID
	}
	if len(chosenSecurityAssociation.Proposals[0].ExtendedSequenceNumbers) != 0 {
		if chosenSecurityAssociation.Proposals[0].ExtendedSequenceNumbers[0].TransformID == 0 {
			childSecurityAssociation.ESN = false
		} else {
			childSecurityAssociation.ESN = true
		}
	}

	// Link UE context
	childSecurityAssociation.ThisUE = ue
	// Record to OCF context
	ocfContext.ChildSA.Store(childSecurityAssociation.SPI, childSecurityAssociation)

	ue.OCFChildSecurityAssociation = childSecurityAssociation

	return childSecurityAssociation, nil
}

func (ue *OCFUe) AttachAMF(sctpAddr string) bool {
	if amf, ok := ocfContext.AMFPoolLoad(sctpAddr); ok {
		amf.OcfUeList[ue.RanUeNgapId] = ue
		ue.AMF = amf
		return true
	} else {
		return false
	}
}
func (ue *OCFUe) DetachAMF() {
	if ue.AMF == nil {
		return
	}
	delete(ue.AMF.OcfUeList, ue.RanUeNgapId)
}

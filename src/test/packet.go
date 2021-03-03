package test

import (
	"free5gcWithOCF/lib/nas"
	"free5gcWithOCF/lib/nas/nasMessage"

	// Nausf_UEAU_Client "free5gcWithOCF/lib/openapi/Nausf_UEAuthentication"
	"free5gcWithOCF/lib/ngap"
	"free5gcWithOCF/src/test/ngapTestpacket"
	// "free5gcWithOCF/lib/openapi/models"
)

func GetNGSetupRequest(gnbId []byte, bitlength uint64, name string) ([]byte, error) {
	message := ngapTestpacket.BuildNGSetupRequest()
	// GlobalRANNodeID
	ie := message.InitiatingMessage.Value.NGSetupRequest.ProtocolIEs.List[0]
	gnbID := ie.Value.GlobalRANNodeID.GlobalGNBID.GNBID.GNBID
	gnbID.Bytes = gnbId
	gnbID.BitLength = bitlength
	// RANNodeName
	ie = message.InitiatingMessage.Value.NGSetupRequest.ProtocolIEs.List[1]
	ie.Value.RANNodeName.Value = name

	return ngap.Encoder(message)
}

func GetInitialUEMessage(ranUeNgapID int64, nasPdu []byte, fiveGSTmsi string) ([]byte, error) {
	message := ngapTestpacket.BuildInitialUEMessage(ranUeNgapID, nasPdu, fiveGSTmsi)
	return ngap.Encoder(message)
}

func GetUplinkNASTransport(AmfUENGAPID, ranUeNgapID int64, nasPdu []byte) ([]byte, error) {
	message := ngapTestpacket.BuildUplinkNasTransport(AmfUENGAPID, ranUeNgapID, nasPdu)
	return ngap.Encoder(message)
}

func GetInitialContextSetupResponse(AmfUENGAPID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildInitialContextSetupResponseForRegistraionTest(AmfUENGAPID, ranUeNgapID)

	return ngap.Encoder(message)
}

func GetInitialContextSetupResponseForServiceRequest(
	AmfUENGAPID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
	message := ngapTestpacket.BuildInitialContextSetupResponse(AmfUENGAPID, ranUeNgapID, ipv4, nil)
	return ngap.Encoder(message)
}

func GetPDUSessionResourceSetupResponse(AmfUENGAPID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
	message := ngapTestpacket.BuildPDUSessionResourceSetupResponseForRegistrationTest(AmfUENGAPID, ranUeNgapID, ipv4)
	return ngap.Encoder(message)
}
func EncodeNasPduWithSecurity(ue *RanUeContext, pdu []byte, securityHeaderType uint8,
	securityContextAvailable, newSecurityContext bool) ([]byte, error) {
	m := nas.NewMessage()
	err := m.PlainNasDecode(&pdu)
	if err != nil {
		return nil, err
	}
	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    securityHeaderType,
	}
	return NASEncode(ue, m, securityContextAvailable, newSecurityContext)
}

func GetUEContextReleaseComplete(AmfUENGAPID int64, ranUeNgapID int64, pduSessionIDList []int64) ([]byte, error) {
	message := ngapTestpacket.BuildUEContextReleaseComplete(AmfUENGAPID, ranUeNgapID, pduSessionIDList)
	return ngap.Encoder(message)
}

func GetUEContextReleaseRequest(AmfUENGAPID int64, ranUeNgapID int64, pduSessionIDList []int64) ([]byte, error) {
	message := ngapTestpacket.BuildUEContextReleaseRequest(AmfUENGAPID, ranUeNgapID, pduSessionIDList)
	return ngap.Encoder(message)
}

func GetPDUSessionResourceReleaseResponse(AmfUENGAPID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildPDUSessionResourceReleaseResponseForReleaseTest(AmfUENGAPID, ranUeNgapID)
	return ngap.Encoder(message)
}
func GetPathSwitchRequest(AmfUENGAPID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildPathSwitchRequest(AmfUENGAPID, ranUeNgapID)
	message.InitiatingMessage.Value.PathSwitchRequest.ProtocolIEs.List =
		message.InitiatingMessage.Value.PathSwitchRequest.ProtocolIEs.List[0:5]
	return ngap.Encoder(message)
}

func GetHandoverRequired(
	AmfUENGAPID int64, ranUeNgapID int64, targetGNBID []byte, targetCellID []byte) ([]byte, error) {
	message := ngapTestpacket.BuildHandoverRequired(AmfUENGAPID, ranUeNgapID, targetGNBID, targetCellID)
	return ngap.Encoder(message)
}

func GetHandoverRequestAcknowledge(AmfUENGAPID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildHandoverRequestAcknowledge(AmfUENGAPID, ranUeNgapID)
	return ngap.Encoder(message)
}

func GetHandoverNotify(AmfUENGAPID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildHandoverNotify(AmfUENGAPID, ranUeNgapID)
	return ngap.Encoder(message)
}

func GetPDUSessionResourceSetupResponseForPaging(AmfUENGAPID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
	message := ngapTestpacket.BuildPDUSessionResourceSetupResponseForPaging(AmfUENGAPID, ranUeNgapID, ipv4)
	return ngap.Encoder(message)
}

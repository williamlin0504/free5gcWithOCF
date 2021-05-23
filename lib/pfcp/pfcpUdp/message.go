package pfcpUdp

import (
	"net"

	" free5gcWithOCF/lib/pfcp"
)

type Message struct {
	RemoteAddr  *net.UDPAddr
	PfcpMessage *pfcp.Message
}

func NewMessage(remoteAddr *net.UDPAddr, pfcpMessage *pfcp.Message) (msg Message) {
	msg = Message{}
	msg.RemoteAddr = remoteAddr
	msg.PfcpMessage = pfcpMessage
	return
}

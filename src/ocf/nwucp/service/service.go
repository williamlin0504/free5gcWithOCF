package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"

	"free5gcWithOCF/src/ocf/context"
	"free5gcWithOCF/src/ocf/logger"
	"free5gcWithOCF/src/ocf/ngap/message"

	"github.com/sirupsen/logrus"
)

var nwucpLog *logrus.Entry

func init() {
	nwucpLog = logger.NWuCPLog
}

// Run setup OCF NAS for UE to forward NAS message
// to OCF
func Run() error {
	// OCF context
	ocfSelf := context.OCFSelf()
	tcpAddr := fmt.Sprintf("%s:%d", ocfSelf.IPSecGatewayAddress, ocfSelf.TCPPort)

	tcpListener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		nwucpLog.Errorf("Listen TCP address failed: %+v", err)
		return errors.New("Listen failed")
	}

	nwucpLog.Tracef("Successfully listen %+v", tcpAddr)

	go listenAndServe(tcpListener)

	return nil
}

// listenAndServe handle TCP listener and accept incoming
// requests. It also stores accepted connection into UE
// context, and finally, call serveConn() to serve the messages
// received from the connection.
func listenAndServe(tcpListener net.Listener) {
	defer func() {
		err := tcpListener.Close()
		if err != nil {
			nwucpLog.Errorf("Error closing tcpListener: %+v", err)
		}
	}()

	for {
		connection, err := tcpListener.Accept()
		if err != nil {
			nwucpLog.Error("TCP server accept failed. Close the listener...")
			return
		}

		nwucpLog.Tracef("Accepted one UE from %+v", connection.RemoteAddr())

		// Find UE context and store this connection in to it, then check if
		// there is any cached NAS message for this UE. If yes, send to it.
		ocfSelf := context.OCFSelf()

		ueIP := strings.Split(connection.RemoteAddr().String(), ":")[0]
		ue, ok := ocfSelf.AllocatedUEIPAddressLoad(ueIP)
		if !ok {
			nwucpLog.Errorf("UE context not found for peer %+v", ueIP)
			continue
		}

		// Store connection
		ue.TCPConnection = connection

		if ue.TemporaryCachedNASMessage != nil {
			// Send to UE
			if n, err := connection.Write(ue.TemporaryCachedNASMessage); err != nil {
				nwucpLog.Errorf("Writing via IPSec signalling SA failed: %+v", err)
			} else {
				nwucpLog.Trace("Forward NWu <- N2")
				nwucpLog.Tracef("Wrote %d bytes", n)
			}
			// Clean the cached message
			ue.TemporaryCachedNASMessage = nil
		}

		go serveConn(ue, connection)
	}
}

// serveConn handle accepted TCP connection. It reads NAS packets
// from the connection and call forward() to forward NAS messages
// to OCF
func serveConn(ue *context.OCFUe, connection net.Conn) {
	defer func() {
		err := connection.Close()
		if err != nil {
			nwucpLog.Errorf("Error closing connection: %+v", err)
		}
	}()

	data := make([]byte, 65535)
	for {
		n, err := connection.Read(data)
		if err != nil {
			if err.Error() == "EOF" {
				nwucpLog.Warn("Connection close by peer")
				ue.TCPConnection = nil
				return
			} else {
				nwucpLog.Errorf("Read TCP connection failed: %+v", err)
			}
		}
		nwucpLog.Tracef("Get NAS PDU from UE:\nNAS length: %d\nNAS content:\n%s", n, hex.Dump(data[:n]))

		forwardData := make([]byte, n)
		copy(forwardData, data[:n])

		go forward(ue, forwardData)
	}
}

// forward forwards NAS messages sent from UE to the
// associated OCF
func forward(ue *context.OCFUe, packet []byte) {
	nwucpLog.Trace("Forward NWu -> N2")
	message.SendUplinkNASTransport(ue.OCF, ue, packet)
}
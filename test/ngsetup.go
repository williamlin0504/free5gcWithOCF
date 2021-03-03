package test

import (
	"fmt"
	"net"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/calee0219/fatal"
)

const ngapPPID uint32 = 0x3c000000

func getNgapIp(ocfIP, ranIP string, ocfPort, ranPort int) (ocfAddr, ranAddr *sctp.SCTPAddr, err error) {
	ips := []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", ocfIP); err1 != nil {
		err = fmt.Errorf("Error resolving address '%s': %v", ocfIP, err1)
		return nil, nil, err
	} else {
		ips = append(ips, *ip)
	}
	ocfAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    ocfPort,
	}
	ips = []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", ranIP); err1 != nil {
		err = fmt.Errorf("Error resolving address '%s': %v", ranIP, err1)
		return nil, nil, err
	} else {
		ips = append(ips, *ip)
	}
	ranAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    ranPort,
	}
	return ocfAddr, ranAddr, nil
}

func ConnectToOcf(ocfIP, ranIP string, ocfPort, ranPort int) (*sctp.SCTPConn, error) {
	ocfAddr, ranAddr, err := getNgapIp(ocfIP, ranIP, ocfPort, ranPort)
	if err != nil {
		return nil, err
	}
	conn, err := sctp.DialSCTP("sctp", ranAddr, ocfAddr)
	if err != nil {
		return nil, err
	}
	info, err := conn.GetDefaultSentParam()
	if err != nil {
		fatal.Fatalf("conn GetDefaultSentParam error in ConnectToOcf: %+v", err)
	}
	info.PPID = ngapPPID
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

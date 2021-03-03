package context

import (
	"crypto/rand"
	"crypto/rsa"
	"math"
	"math/big"
	"net"
	"sync"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/sirupsen/logrus"
	gtpv1 "github.com/wmnsk/go-gtp/v1"
	"golang.org/x/net/ipv4"

	"free5gcWithOCF/lib/idgenerator"
	"free5gcWithOCF/lib/ngap/ngapType"
	"free5gcWithOCF/src/ocf/logger"
)

var contextLog *logrus.Entry

var ocfContext = OCFContext{}

type OCFContext struct {
	NFInfo           OCFNFInfo
	OCFSCTPAddresses []*sctp.SCTPAddr

	// ID generator
	RANUENGAPIDGenerator *idgenerator.IDGenerator
	TEIDGenerator        *idgenerator.IDGenerator

	// Pools
	UePool                 sync.Map // map[int64]*OCFUe, RanUeNgapID as key
	OCFPool                sync.Map // map[string]*OCF, SCTPAddr as key
	OCFReInitAvailableList sync.Map // map[string]bool, SCTPAddr as key
	IKESA                  sync.Map // map[uint64]*IKESecurityAssociation, SPI as key
	ChildSA                sync.Map // map[uint32]*ChildSecurityAssociation, SPI as key
	GTPConnectionWithUPF   sync.Map // map[string]*gtpv1.UPlaneConn, UPF address as key
	AllocatedUEIPAddress   sync.Map // map[string]*OCFUe, IPAddr as key
	AllocatedUETEID        sync.Map // map[uint32]*OCFUe, TEID as key

	// OCF FQDN
	FQDN string

	// Security data
	CertificateAuthority []byte
	OCFCertificate       []byte
	OCFPrivateKey        *rsa.PrivateKey

	// UEIPAddressRange
	Subnet *net.IPNet

	// Network interface mark for xfrm
	Mark uint32

	// OCF local address
	IKEBindAddress      string
	IPSecGatewayAddress string
	GTPBindAddress      string
	TCPPort             uint16

	// OCF NWu interface raw socket
	NWuRawSocket *ipv4.RawConn
}

func init() {
	// init log
	contextLog = logger.ContextLog

	// init ID generator
	ocfContext.RANUENGAPIDGenerator = idgenerator.NewGenerator(0, math.MaxInt64)
	ocfContext.TEIDGenerator = idgenerator.NewGenerator(1, math.MaxUint32)
}

// Create new OCF context
func OCFSelf() *OCFContext {
	return &ocfContext
}

func (context *OCFContext) NewOcfUe() *OCFUe {
	ranUeNgapId, err := context.RANUENGAPIDGenerator.Allocate()
	if err != nil {
		contextLog.Errorf("New OCF UE failed: %+v", err)
		return nil
	}
	ocfUe := new(OCFUe)
	ocfUe.init(ranUeNgapId)
	context.UePool.Store(ranUeNgapId, ocfUe)
	return ocfUe
}

func (context *OCFContext) DeleteOcfUe(ranUeNgapId int64) {
	context.UePool.Delete(ranUeNgapId)
}

func (context *OCFContext) UePoolLoad(ranUeNgapId int64) (*OCFUe, bool) {
	ue, ok := context.UePool.Load(ranUeNgapId)
	if ok {
		return ue.(*OCFUe), ok
	} else {
		return nil, ok
	}
}

func (context *OCFContext) NewOcfOcf(sctpAddr string, conn *sctp.SCTPConn) *OCF {
	amf := new(OCF)
	amf.init(sctpAddr, conn)
	if item, loaded := context.OCFPool.LoadOrStore(sctpAddr, amf); loaded {
		contextLog.Warn("[Context] NewOcfOcf(): OCF entry already exists.")
		return item.(*OCF)
	} else {
		return amf
	}
}

func (context *OCFContext) DeleteOcfOcf(sctpAddr string) {
	context.OCFPool.Delete(sctpAddr)
}

func (context *OCFContext) OCFPoolLoad(sctpAddr string) (*OCF, bool) {
	amf, ok := context.OCFPool.Load(sctpAddr)
	if ok {
		return amf.(*OCF), ok
	} else {
		return nil, ok
	}
}

func (context *OCFContext) DeleteOCFReInitAvailableFlag(sctpAddr string) {
	context.OCFReInitAvailableList.Delete(sctpAddr)
}

func (context *OCFContext) OCFReInitAvailableListLoad(sctpAddr string) (bool, bool) {
	flag, ok := context.OCFReInitAvailableList.Load(sctpAddr)
	if ok {
		return flag.(bool), ok
	} else {
		return true, ok
	}
}

func (context *OCFContext) OCFReInitAvailableListStore(sctpAddr string, flag bool) {
	context.OCFReInitAvailableList.Store(sctpAddr, flag)
}

func (context *OCFContext) NewIKESecurityAssociation() *IKESecurityAssociation {
	ikeSecurityAssociation := new(IKESecurityAssociation)

	var maxSPI *big.Int = new(big.Int).SetUint64(math.MaxUint64)
	var localSPIuint64 uint64

	for {
		localSPI, err := rand.Int(rand.Reader, maxSPI)
		if err != nil {
			contextLog.Error("[Context] Error occurs when generate new IKE SPI")
			return nil
		}
		localSPIuint64 = localSPI.Uint64()
		if _, duplicate := context.IKESA.LoadOrStore(localSPIuint64, ikeSecurityAssociation); !duplicate {
			break
		}
	}

	ikeSecurityAssociation.LocalSPI = localSPIuint64

	return ikeSecurityAssociation
}

func (context *OCFContext) DeleteIKESecurityAssociation(spi uint64) {
	context.IKESA.Delete(spi)
}

func (context *OCFContext) IKESALoad(spi uint64) (*IKESecurityAssociation, bool) {
	securityAssociation, ok := context.IKESA.Load(spi)
	if ok {
		return securityAssociation.(*IKESecurityAssociation), ok
	} else {
		return nil, ok
	}
}

func (context *OCFContext) DeleteGTPConnection(upfAddr string) {
	context.GTPConnectionWithUPF.Delete(upfAddr)
}

func (context *OCFContext) GTPConnectionWithUPFLoad(upfAddr string) (*gtpv1.UPlaneConn, bool) {
	conn, ok := context.GTPConnectionWithUPF.Load(upfAddr)
	if ok {
		return conn.(*gtpv1.UPlaneConn), ok
	} else {
		return nil, ok
	}
}

func (context *OCFContext) GTPConnectionWithUPFStore(upfAddr string, conn *gtpv1.UPlaneConn) {
	context.GTPConnectionWithUPF.Store(upfAddr, conn)
}

func (context *OCFContext) NewInternalUEIPAddr(ue *OCFUe) net.IP {
	var ueIPAddr net.IP

	// TODO: Check number of allocated IP to detect running out of IPs
	for {
		ueIPAddr = generateRandomIPinRange(context.Subnet)
		if ueIPAddr != nil {
			if ueIPAddr.String() == context.IPSecGatewayAddress {
				continue
			}
			if _, ok := context.AllocatedUEIPAddress.LoadOrStore(ueIPAddr.String(), ue); !ok {
				break
			}
		}
	}

	return ueIPAddr
}

func (context *OCFContext) DeleteInternalUEIPAddr(ipAddr string) {
	context.AllocatedUEIPAddress.Delete(ipAddr)
}

func (context *OCFContext) AllocatedUEIPAddressLoad(ipAddr string) (*OCFUe, bool) {
	ue, ok := context.AllocatedUEIPAddress.Load(ipAddr)
	if ok {
		return ue.(*OCFUe), ok
	} else {
		return nil, ok
	}
}

func (context *OCFContext) NewTEID(ue *OCFUe) uint32 {
	teid64, err := context.TEIDGenerator.Allocate()
	if err != nil {
		contextLog.Errorf("New TEID failed: %+v", err)
		return 0
	}
	teid32 := uint32(teid64)

	context.AllocatedUETEID.Store(teid32, ue)

	return teid32
}

func (context *OCFContext) DeleteTEID(teid uint32) {
	context.AllocatedUETEID.Delete(teid)
}

func (context *OCFContext) AllocatedUETEIDLoad(teid uint32) (*OCFUe, bool) {
	ue, ok := context.AllocatedUETEID.Load(teid)
	if ok {
		return ue.(*OCFUe), ok
	} else {
		return nil, ok
	}
}

func (context *OCFContext) OCFSelection(ueSpecifiedGUAMI *ngapType.GUAMI) *OCF {
	var availableOCF *OCF
	context.OCFPool.Range(func(key, value interface{}) bool {
		amf := value.(*OCF)
		if amf.FindAvalibleOCFByCompareGUAMI(ueSpecifiedGUAMI) {
			availableOCF = amf
			return false
		} else {
			return true
		}
	})
	return availableOCF
}

func generateRandomIPinRange(subnet *net.IPNet) net.IP {
	ipAddr := make([]byte, 4)
	randomNumber := make([]byte, 4)

	_, err := rand.Read(randomNumber)
	if err != nil {
		contextLog.Errorf("Generate random number for IP address failed: %+v", err)
		return nil
	}

	// TODO: elimenate network name, gateway, and broadcast
	for i := 0; i < 4; i++ {
		alter := randomNumber[i] & (subnet.Mask[i] ^ 255)
		ipAddr[i] = subnet.IP[i] + alter
	}

	return net.IPv4(ipAddr[0], ipAddr[1], ipAddr[2], ipAddr[3])
}

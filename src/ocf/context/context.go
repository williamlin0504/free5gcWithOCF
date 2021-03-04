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

var ocfContext = CCFContext{}

type CCFContext struct {
	NFInfo           CCFNFInfo
	AMFSCTPAddresses []*sctp.SCTPAddr

	// ID generator
	RANUENGAPIDGenerator *idgenerator.IDGenerator
	TEIDGenerator        *idgenerator.IDGenerator

	// Pools
	UePool                 sync.Map // map[int64]*CCFUe, RanUeNgapID as key
	AMFPool                sync.Map // map[string]*CCFAMF, SCTPAddr as key
	AMFReInitAvailableList sync.Map // map[string]bool, SCTPAddr as key
	IKESA                  sync.Map // map[uint64]*IKESecurityAssociation, SPI as key
	ChildSA                sync.Map // map[uint32]*ChildSecurityAssociation, SPI as key
	GTPConnectionWithUPF   sync.Map // map[string]*gtpv1.UPlaneConn, UPF address as key
	AllocatedUEIPAddress   sync.Map // map[string]*CCFUe, IPAddr as key
	AllocatedUETEID        sync.Map // map[uint32]*CCFUe, TEID as key

	// CCF FQDN
	FQDN string

	// Security data
	CertificateAuthority []byte
	CCFCertificate       []byte
	CCFPrivateKey        *rsa.PrivateKey

	// UEIPAddressRange
	Subnet *net.IPNet

	// Network interface mark for xfrm
	Mark uint32

	// CCF local address
	IKEBindAddress      string
	IPSecGatewayAddress string
	GTPBindAddress      string
	TCPPort             uint16

	// CCF NWu interface raw socket
	NWuRawSocket *ipv4.RawConn
}

func init() {
	// init log
	contextLog = logger.ContextLog

	// init ID generator
	ocfContext.RANUENGAPIDGenerator = idgenerator.NewGenerator(0, math.MaxInt64)
	ocfContext.TEIDGenerator = idgenerator.NewGenerator(1, math.MaxUint32)
}

// Create new CCF context
func CCFSelf() *CCFContext {
	return &ocfContext
}

func (context *CCFContext) NewCcfUe() *CCFUe {
	ranUeNgapId, err := context.RANUENGAPIDGenerator.Allocate()
	if err != nil {
		contextLog.Errorf("New CCF UE failed: %+v", err)
		return nil
	}
	ocfUe := new(CCFUe)
	ocfUe.init(ranUeNgapId)
	context.UePool.Store(ranUeNgapId, ocfUe)
	return ocfUe
}

func (context *CCFContext) DeleteCcfUe(ranUeNgapId int64) {
	context.UePool.Delete(ranUeNgapId)
}

func (context *CCFContext) UePoolLoad(ranUeNgapId int64) (*CCFUe, bool) {
	ue, ok := context.UePool.Load(ranUeNgapId)
	if ok {
		return ue.(*CCFUe), ok
	} else {
		return nil, ok
	}
}

func (context *CCFContext) NewCcfAmf(sctpAddr string, conn *sctp.SCTPConn) *CCFAMF {
	amf := new(CCFAMF)
	amf.init(sctpAddr, conn)
	if item, loaded := context.AMFPool.LoadOrStore(sctpAddr, amf); loaded {
		contextLog.Warn("[Context] NewCcfAmf(): AMF entry already exists.")
		return item.(*CCFAMF)
	} else {
		return amf
	}
}

func (context *CCFContext) DeleteCcfAmf(sctpAddr string) {
	context.AMFPool.Delete(sctpAddr)
}

func (context *CCFContext) AMFPoolLoad(sctpAddr string) (*CCFAMF, bool) {
	amf, ok := context.AMFPool.Load(sctpAddr)
	if ok {
		return amf.(*CCFAMF), ok
	} else {
		return nil, ok
	}
}

func (context *CCFContext) DeleteAMFReInitAvailableFlag(sctpAddr string) {
	context.AMFReInitAvailableList.Delete(sctpAddr)
}

func (context *CCFContext) AMFReInitAvailableListLoad(sctpAddr string) (bool, bool) {
	flag, ok := context.AMFReInitAvailableList.Load(sctpAddr)
	if ok {
		return flag.(bool), ok
	} else {
		return true, ok
	}
}

func (context *CCFContext) AMFReInitAvailableListStore(sctpAddr string, flag bool) {
	context.AMFReInitAvailableList.Store(sctpAddr, flag)
}

func (context *CCFContext) NewIKESecurityAssociation() *IKESecurityAssociation {
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

func (context *CCFContext) DeleteIKESecurityAssociation(spi uint64) {
	context.IKESA.Delete(spi)
}

func (context *CCFContext) IKESALoad(spi uint64) (*IKESecurityAssociation, bool) {
	securityAssociation, ok := context.IKESA.Load(spi)
	if ok {
		return securityAssociation.(*IKESecurityAssociation), ok
	} else {
		return nil, ok
	}
}

func (context *CCFContext) DeleteGTPConnection(upfAddr string) {
	context.GTPConnectionWithUPF.Delete(upfAddr)
}

func (context *CCFContext) GTPConnectionWithUPFLoad(upfAddr string) (*gtpv1.UPlaneConn, bool) {
	conn, ok := context.GTPConnectionWithUPF.Load(upfAddr)
	if ok {
		return conn.(*gtpv1.UPlaneConn), ok
	} else {
		return nil, ok
	}
}

func (context *CCFContext) GTPConnectionWithUPFStore(upfAddr string, conn *gtpv1.UPlaneConn) {
	context.GTPConnectionWithUPF.Store(upfAddr, conn)
}

func (context *CCFContext) NewInternalUEIPAddr(ue *CCFUe) net.IP {
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

func (context *CCFContext) DeleteInternalUEIPAddr(ipAddr string) {
	context.AllocatedUEIPAddress.Delete(ipAddr)
}

func (context *CCFContext) AllocatedUEIPAddressLoad(ipAddr string) (*CCFUe, bool) {
	ue, ok := context.AllocatedUEIPAddress.Load(ipAddr)
	if ok {
		return ue.(*CCFUe), ok
	} else {
		return nil, ok
	}
}

func (context *CCFContext) NewTEID(ue *CCFUe) uint32 {
	teid64, err := context.TEIDGenerator.Allocate()
	if err != nil {
		contextLog.Errorf("New TEID failed: %+v", err)
		return 0
	}
	teid32 := uint32(teid64)

	context.AllocatedUETEID.Store(teid32, ue)

	return teid32
}

func (context *CCFContext) DeleteTEID(teid uint32) {
	context.AllocatedUETEID.Delete(teid)
}

func (context *CCFContext) AllocatedUETEIDLoad(teid uint32) (*CCFUe, bool) {
	ue, ok := context.AllocatedUETEID.Load(teid)
	if ok {
		return ue.(*CCFUe), ok
	} else {
		return nil, ok
	}
}

func (context *CCFContext) AMFSelection(ueSpecifiedGUAMI *ngapType.GUAMI) *CCFAMF {
	var availableAMF *CCFAMF
	context.AMFPool.Range(func(key, value interface{}) bool {
		amf := value.(*CCFAMF)
		if amf.FindAvalibleAMFByCompareGUAMI(ueSpecifiedGUAMI) {
			availableAMF = amf
			return false
		} else {
			return true
		}
	})
	return availableAMF
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

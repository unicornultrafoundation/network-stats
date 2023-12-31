package crawler

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"time"

	"github.com/ethereum/node-crawler/pkg/common"
	ethCommon "github.com/unicornultrafoundation/go-u2u/libs/common"
	"github.com/unicornultrafoundation/go-u2u/libs/crypto"
	"github.com/unicornultrafoundation/go-u2u/libs/log"
	"github.com/unicornultrafoundation/go-u2u/libs/p2p"
	"github.com/unicornultrafoundation/go-u2u/libs/p2p/enode"
	"github.com/unicornultrafoundation/go-u2u/libs/p2p/rlpx"

	"github.com/pkg/errors"
)

var (
	_status *Status
)

func getClientInfo(genesisHash ethCommon.Hash, networkID uint64, nodeURL string, n *enode.Node) (*common.ClientInfo, error) {
	var info common.ClientInfo

	conn, sk, err := dial(n)
	if err != nil {
		log.Error("dial error", "error", err)
		return nil, err
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return nil, errors.Wrap(err, "cannot set conn deadline")
	}

	if err = writeHello(conn, sk); err != nil {
		log.Error("write handshake failed", "error", err)
		return nil, err
	}
	if err = readHello(conn, &info); err != nil {
		log.Error("read handshake failed", "error", err)
		return nil, err
	}

	// If node provides no eth version, we can skip it.
	if conn.negotiatedProtoVersion == 0 {
		return &info, nil
	}

	if err = conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return nil, errors.Wrap(err, "cannot set conn deadline")
	}

	s := getStatus(uint32(conn.negotiatedProtoVersion), genesisHash, networkID, nodeURL)
	if err = conn.Write(s); err != nil {
		log.Error("conn write failed", "error", err)
		return nil, err
	}

	// Regardless of whether we wrote a status message or not, the remote side
	// might still send us one.

	if err = readStatus(conn, &info); err != nil {
		log.Error("read status", "error", err)
		return nil, err
	}

	// Disconnect from client
	_ = conn.Write(Disconnect{Reason: p2p.DiscQuitting})

	return &info, nil
}

// dial attempts to dial the given node and perform a handshake,
func dial(n *enode.Node) (*Conn, *ecdsa.PrivateKey, error) {
	var conn Conn

	// dial
	dialer := net.Dialer{Timeout: 10 * time.Second}
	fd, err := dialer.Dial("tcp", fmt.Sprintf("%v:%d", n.IP(), n.TCP()))
	if err != nil {
		log.Error("dialer dial failed", "error", err)
		return nil, nil, err
	}

	conn.Conn = rlpx.NewConn(fd, n.Pubkey())

	if err = conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return nil, nil, errors.Wrap(err, "cannot set conn deadline")
	}

	// do encHandshake
	ourKey, _ := crypto.GenerateKey()

	_, err = conn.Handshake(ourKey)
	if err != nil {
		log.Error("handshake failed", "error", err)
		return nil, nil, err
	}

	return &conn, ourKey, nil
}

func writeHello(conn *Conn, priv *ecdsa.PrivateKey) error {
	pub0 := crypto.FromECDSAPub(&priv.PublicKey)[1:]

	h := &Hello{
		Version: 5,
		Caps: []p2p.Cap{
			/* For U2U compabilities, below protocols are applied. More detail at:
			https://github.com/unicornultrafoundation/go-u2u/blob/d03dea550c200226620424a8a27497eaf9d6021a/gossip/service.go#L398-L403
			*/
			{Name: "u2u", Version: 1},
			/* fsnap should be supported however package we sent from this protocol is not handled in go-u2u correctly */
			// {Name: "fsnap", Version: 1},
		},
		ID: pub0,
	}

	conn.ourHighestProtoVersion = 68
	conn.ourHighestSnapProtoVersion = 1

	return conn.Write(h)
}

func readHello(conn *Conn, info *common.ClientInfo) error {
	switch msg := conn.Read().(type) {
	case *Hello:
		// set snappy if version is at least 5
		if msg.Version >= 5 {
			conn.SetSnappy(true)
		}
		info.Capabilities = msg.Caps
		info.SoftwareVersion = msg.Version
		info.ClientType = msg.Name

		conn.negotiateEthProtocol(info.Capabilities)

		return nil
	case *Disconnect:
		return fmt.Errorf("bad hello handshake disconnect: %v", msg.Reason.Error())
	case *Error:
		return fmt.Errorf("bad hello handshake error: %v", msg.Error())
	default:
		return fmt.Errorf("bad hello handshake code: %v", msg.Code())
	}
}

func getStatus(version uint32, genesis ethCommon.Hash, network uint64, nodeURL string) *Status {
	if _status == nil {
		_status = &Status{
			ProtocolVersion: version,
			NetworkID:       network,
			Genesis:         genesis,
		}
	}

	return _status
}

func readStatus(conn *Conn, info *common.ClientInfo) error {
	switch msg := conn.Read().(type) {
	case *Status:
		info.NetworkID = msg.NetworkID
		// m.ProtocolVersion
	case *Disconnect:
		return fmt.Errorf("bad status handshake disconnect: %v", msg.Reason.Error())
	case *Error:
		return fmt.Errorf("bad status handshake error: %v", msg.Error())
	default:
		return fmt.Errorf("bad status handshake code: %v", msg.Code())
	}
	return nil
}

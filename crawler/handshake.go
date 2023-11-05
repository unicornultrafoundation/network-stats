package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/ethereum/go-ethereum/params"

	"github.com/pkg/errors"
)

var (
	_status          *Status
	lastStatusUpdate time.Time
)

type clientInfo struct {
	ClientType      string
	ClientVersion   string
	ClientDesc      string
	OsType          string
	GoVersion       string
	SoftwareVersion uint64
	Capabilities    []p2p.Cap
	NetworkID       uint64
	ForkID          forkid.ID
	Blockheight     string
	TotalDifficulty *big.Int
	HeadHash        common.Hash
}

func getClientInfo(genesis *core.Genesis, networkID uint64, nodeURL string, n *enode.Node) (*clientInfo, error) {
	var info clientInfo

	conn, sk, err := dial(n)
	if err != nil {
		return &info, errors.Wrap(err, "couldNotDial: ")
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return &info, errors.Wrap(err, "cannot set conn deadline for hello")
	}

	if err = writeHello(conn, sk); err != nil {
		return &info, errors.Wrap(err, "writeHelloFailure")
	}
	if err = readHello(conn, &info); err != nil {
		return &info, errors.Wrap(err, "readHelloFailure")
	}

	// If node provides no eth version, we can skip it.
	if conn.negotiatedProtoVersion == 0 {
		return &info, nil
	}

	if err = conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		log.Warn("SetDeadline-2: " + err.Error())
		return &info, errors.Wrap(err, "cannot set conn deadline for status")
	}

	s := getStatus(genesis.Config, uint32(conn.negotiatedProtoVersion), genesis.ToBlock().Hash(), networkID, nodeURL)
	if err = conn.Write(s); err != nil {
		return &info, errors.Wrap(err, "getStatusError")
	}

	// Regardless of whether we wrote a status message or not, the remote side
	// might still send us one.

	if err = readStatus(conn, &info); err != nil {
		return &info, errors.Wrap(err, "readStatusError")
	}

	// Disconnect from client
	_ = conn.Write(Disconnect{Reason: p2p.DiscQuitting})

	return &info, nil
}

// dial attempts to dial the given node and perform a handshake,
func dial(n *enode.Node) (*Conn, *ecdsa.PrivateKey, error) {
	var conn Conn

	// dial
	fd, err := net.Dial("tcp", fmt.Sprintf("%v:%d", n.IP(), n.TCP()))
	if err != nil {
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
		return nil, nil, err
	}

	return &conn, ourKey, nil
}

func writeHello(conn *Conn, priv *ecdsa.PrivateKey) error {
	pub0 := crypto.FromECDSAPub(&priv.PublicKey)[1:]

	h := &Hello{
		Version: 5,
		Caps: []p2p.Cap{
			{Name: "diff", Version: 1},
			{Name: "eth", Version: 64},
			{Name: "eth", Version: 65},
			{Name: "eth", Version: 66},
			{Name: "les", Version: 2},
			{Name: "les", Version: 3},
			{Name: "les", Version: 4},
			{Name: "snap", Version: 1},
		},
		ID: pub0,
	}

	conn.ourHighestProtoVersion = 66

	return conn.Write(h)
}

func readHello(conn *Conn, info *clientInfo) error {
	switch msg := conn.Read().(type) {
	case *Hello:
		// set snappy if version is at least 5
		if msg.Version >= 5 {
			conn.SetSnappy(true)
		}
		info.Capabilities = msg.Caps
		info.SoftwareVersion = msg.Version

		splitClient := strings.Split(msg.Name, "/")
		if len(splitClient) == 4 {
			info.ClientType = splitClient[0]
			info.ClientDesc = ""
			info.ClientVersion = splitClient[1]
			info.OsType = splitClient[2]
			info.GoVersion = splitClient[3]
		} else if len(splitClient) == 5 {
			info.ClientType = splitClient[0]
			info.ClientDesc = splitClient[1]
			info.ClientVersion = splitClient[2]
			info.OsType = splitClient[3]
			info.GoVersion = splitClient[4]
		} else if len(splitClient) == 6 {
			info.ClientType = splitClient[0]
			info.ClientDesc = fmt.Sprintf("%v/%v", splitClient[2], splitClient[2])
			info.ClientVersion = splitClient[1]
			info.OsType = splitClient[4]
			info.GoVersion = splitClient[5]
		} else {
			info.ClientType = msg.Name
			info.ClientDesc = ""
			info.ClientVersion = ""
			info.OsType = ""
			info.GoVersion = ""
		}
	case *Disconnect:
		return fmt.Errorf("bad hello handshake: %v", msg.Reason.Error())
	case *Error:
		return fmt.Errorf("bad hello handshake: %v", msg.Error())
	default:
		return fmt.Errorf("bad hello handshake: %v", msg.Code())
	}

	conn.negotiateEthProtocol(info.Capabilities)

	return nil
}

func getStatus(config *params.ChainConfig, version uint32, genesis common.Hash, network uint64, nodeURL string) *Status {
	if _status == nil {
		_status = &Status{
			ProtocolVersion: version,
			NetworkID:       network,
			TD:              big.NewInt(0),
			Head:            genesis,
			Genesis:         genesis,
			//ForkID:          forkid.NewID(config, genesis, 0),
		}
	}

	if nodeURL != "" && time.Since(lastStatusUpdate) > 15*time.Second {
		cl, err := ethclient.Dial(nodeURL)
		if err != nil {
			log.Error("ethclient.Dial", "err", err)
			return _status
		}

		header, err := cl.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Error("cannot get header by number", "err", err)
			return _status
		}

		_status.Head = header.Hash()
		//_status.ForkID = forkid.NewID(config, genesis, header.Number.Uint64())
	}

	return _status
}

func readStatus(conn *Conn, info *clientInfo) error {
	switch msg := conn.Read().(type) {
	case *Status:
		info.ForkID = msg.ForkID
		info.HeadHash = msg.Head
		info.NetworkID = msg.NetworkID
		// m.ProtocolVersion
		info.TotalDifficulty = msg.TD
		// Set correct TD if received TD is higher
		if msg.TD.Cmp(_status.TD) > 0 {
			_status.TD = msg.TD
		}
	case *Disconnect:
		return fmt.Errorf("bad status handshake: %v", msg.Reason.Error())
	case *Error:
		return fmt.Errorf("bad status handshake: %v", msg.Error())
	default:
		return fmt.Errorf("bad status handshake: %v", msg.Code())
	}
	return nil
}

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enr"

	beacon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"

	"github.com/oschwald/geoip2-golang"
)

func updateNodes(db *sql.DB, geoipDB *geoip2.Reader, nodes []nodeJSON) error {
	log.Info("Writing nodes to db", "nodes", len(nodes))
	now := time.Now()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT OR REPLACE into nodes(ID, 
			Now,
			ClientType,
            ClientDesc,
            ClientVersion,
            OsType,
            GoVersion,
			PK,
			SoftwareVersion,
			Capabilities,
			NetworkID,
			ForkID,
			Blockheight,
			TotalDifficulty,
			HeadHash,
			IP,
			Country,
			City,
			Coordinates,
			FirstSeen,
			LastSeen,
			Seq,
			Score,
			ConnType,
            ErrorReason,
            ErrorString) 
			values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range nodes {
		fmt.Printf("node:%+v \n", n)

		info := &clientInfo{}
		if n.Info != nil {
			info = n.Info
		}

		if info.ClientType == "" {
			if n.ErrorReason == -1 {
				info.ClientType = "NA"
			}
		}
		connType := ""
		var portUDP enr.UDP
		if n.N.Load(&portUDP) == nil {
			connType = "UDP"
		}
		var portTCP enr.TCP
		if n.N.Load(&portTCP) == nil {
			connType = "TCP"
		}
		fid := fmt.Sprintf("Hash: %v, Next %v", info.ForkID.Hash, info.ForkID.Next)

		var eth2 ETH2
		if n.N.Load(&eth2) == nil {
			info.ClientType = "eth2"
			var dat beacon.Eth2Data
			if err := dat.Deserialize(codec.NewDecodingReader(bytes.NewReader(eth2), uint64(len(eth2)))); err == nil {
				fid = fmt.Sprintf("Hash: %v, Next %v", dat.ForkDigest, dat.NextForkEpoch)
			}
		}
		var caps string
		for _, c := range info.Capabilities {
			caps = fmt.Sprintf("%v, %v", caps, c.String())
		}
		var pk string
		if n.N.Pubkey() != nil {
			pk = fmt.Sprintf("X: %v, Y: %v", n.N.Pubkey().X.String(), n.N.Pubkey().Y.String())
		}

		var country, city, loc string
		if geoipDB != nil {
			// parse GeoIp info
			ipRecord, err := geoipDB.City(n.N.IP())
			if err != nil {
				return err
			}
			country, city, loc =
				ipRecord.Country.Names["en"],
				ipRecord.City.Names["en"],
				fmt.Sprintf("%v,%v", ipRecord.Location.Latitude, ipRecord.Location.Longitude)
		}

		_, err = stmt.Exec(
			n.N.ID().String(),
			now.String(),
			info.ClientType,
			info.ClientDesc,
			info.ClientVersion,
			info.OsType,
			info.GoVersion,
			pk,
			info.SoftwareVersion,
			caps,
			info.NetworkID,
			fid,
			info.Blockheight,
			info.TotalDifficulty.String(),
			info.HeadHash.String(),
			n.N.IP().String(),
			country,
			city,
			loc,
			n.FirstResponse.String(),
			n.LastResponse.String(),
			n.Seq,
			n.Score,
			connType,
			n.ErrorReason,
			n.ErrorString,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func createDB(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE nodes (
		ID text not null, 
		Now text not null,
		ClientType text,
		ClientDesc text,
		ClientVersion text,
		OsType text,
        GoVersion text,
		PK text,
		SoftwareVersion text,
		Capabilities text,
		NetworkID number,
		ForkID text,
		Blockheight text,
		TotalDifficulty text,
		HeadHash text,
		IP text,
		Country text,
		City text,
		Coordinates text,
		FirstSeen text,
		LastSeen text,
		Seq number,
		Score number,
		ConnType text,
		ErrorReason number,
		ErrorString text,
		PRIMARY KEY (ID)
	);
	delete from nodes;
	`
	_, err := db.Exec(sqlStmt)
	return err
}

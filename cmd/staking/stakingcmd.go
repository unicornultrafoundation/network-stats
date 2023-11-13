package main

import (
	"github.com/unicornultrafoundation/go-u2u/libs/cmd/utils"
	"github.com/unicornultrafoundation/network-stats/pkg/crawl"
	"gopkg.in/urfave/cli.v1"
)

var (
	crawlCommand = &cli.Command{
		Name:   "crawl",
		Usage:  "Crawl network block",
		Action: crawlNetwork,
		Flags: []cli.Flag{
			bootnodesFlag,
			nodeURLFlag,
			workersFlag,
			utils.NetworkIdFlag,
			genesisHashFlag,
		},
	}
)

func crawlNetwork(ctx *cli.Context) error {

	//nodesFile := ctx.String(nodeFileFlag.Name)
	//
	//if nodesFile != "" && gethCommon.FileExist(nodesFile) {
	//	inputSet = common.LoadNodesJSON(nodesFile)
	//}
	//
	//var db *sql.DB
	//if ctx.IsSet(crawlerDBFlag.Name) {
	//	name := ctx.String(crawlerDBFlag.Name)
	//	shouldInit := false
	//	if _, err := os.Stat(name); os.IsNotExist(err) {
	//		shouldInit = true
	//	}
	//
	//	var err error
	//	db, err = openSQLiteDB(
	//		name,
	//		ctx.String(autovacuumFlag.Name),
	//		ctx.Uint64(busyTimeoutFlag.Name),
	//	)
	//	if err != nil {
	//		panic(err)
	//	}
	//	log.Info("Connected to db")
	//	if shouldInit {
	//		log.Info("DB did not exist, init")
	//		if err := crawlerdb.CreateDB(db); err != nil {
	//			panic(err)
	//		}
	//	}
	//}

	crawler := crawl.Crawler{
		GenesisHash: ctx.String(genesisHashFlag.Name),
		NetworkID:   ctx.Uint64(utils.NetworkIdFlag.Name),
		NodeURL:     ctx.String(nodeURLFlag.Name),
		Bootnodes:   ctx.StringSlice(bootnodesFlag.Name),
		Workers:     ctx.Uint64(workersFlag.Name),
	}

	for {
		updatedSet := crawler.CrawlNetwork()
		if nodesFile != "" {
			updatedSet.WriteNodesJSON(nodesFile)
		}
	}
}

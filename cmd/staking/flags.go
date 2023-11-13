package main

import (
	"gopkg.in/urfave/cli.v1"
)

var (
	targetNetworkFlag = &cli.StringFlag{
		Name:  "api-db",
		Usage: "API SQLite file name",
		Value: "mainnet",
	}
	apiDBFlag = &cli.StringFlag{
		Name:  "api-db",
		Usage: "API SQLite file name",
	}
	apiListenAddrFlag = &cli.StringFlag{
		Name:  "addr",
		Usage: "Listening address",
		Value: "0.0.0.0:10000",
	}
	bootnodesFlag = &cli.StringSliceFlag{
		Name: "bootnodes",
		Usage: "Comma separated nodes used for bootstrapping. " +
			"Defaults to hard-coded values for the selected network",
	}
	nodeURLFlag = &cli.StringFlag{
		Name:  "nodeURL",
		Usage: "URL of the node you want to connect to",
	}
	workersFlag = &cli.Uint64Flag{
		Name:  "workers",
		Usage: "Number of workers to start for updating nodes",
		Value: 16,
	}
	genesisHashFlag = &cli.StringFlag{
		Name:  "genesisHash",
		Usage: "Genesis hash in hex. (Mainnet: 0x54e033c612a9b1a8ac8c6cb131f513202648f19b3a2756f8e2e40877d280606c)",
		Value: "0x54e033c612a9b1a8ac8c6cb131f513202648f19b3a2756f8e2e40877d280606c",
	}
)

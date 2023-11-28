package grabber

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/unicornultrafoundation/network-stats/pkg/jsonrpc"

	"github.com/unicornultrafoundation/go-u2u/libs/core/types"
)

var (
	defaultBlockInterval time.Duration = 2 * time.Second
)

type Config struct {
	rpcURL string
}

type Grabber struct {
	web3Client *jsonrpc.Web3
	blockCh    chan *types.Block
	closed     chan struct{}
}

func NewGrabber() (*Grabber, error) {
	g := &Grabber{}
	return g, nil
}

func (g *Grabber) Run(ctx context.Context, interval time.Duration) {
	if interval.String() == "0s" {
		interval = defaultBlockInterval
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			latest, err := g.web3Client.U2U.GetBlockNumber()
			if err != nil {
				//srv.Logger.Error("failed to get latest block number", zap.Error(err))
				//	continue
				//}
				//// delay listener for 1 block for correct responses of kardiaCall
				//if latest != 0 {
				//	latest--
				//}
				//lgr := srv.Logger.With(zap.Uint64("block", latest))
				//if latest <= prevHeader {
				//	continue
				//}
				//if prevHeader < latest {
				//	startTime = time.Now()
				//	block, err := srv.BlockByHeight(ctx, latest)
				//	if err != nil {
				//		lgr.Error("Failed to get block from RPC", zap.Error(err))
				//		continue
				//	}
				//	lgr.Info("Block info from network", zap.Any("Block", block))
				//	endTime = time.Since(startTime)
				//	srv.Metrics().RecordScrapingTime(endTime)
				//	lgr.Info("scraping block time", zap.Duration("TimeConsumed", endTime), zap.String("Avg", srv.Metrics().GetScrapingTime()))
				//	if block == nil {
				//		lgr.Error("Block not found")
				//		continue
				//	}
				//	// insert current block height to cache for re-verifying later
				//	// temp remove insert new unverified blocks
				//	err = srv.InsertUnverifiedBlocks(ctx, latest)
				//	if err != nil {
				//		lgr.Error("Failed to insert unverified block", zap.Error(err))
				//	}
				//	// import this latest block to cache and database
				//	totalImportTime := time.Now()
				//	if err := srv.ImportBlock(ctx, block, true); err != nil {
				//		lgr.Debug("Failed to import block", zap.Error(err))
				//		continue
				//	}
				//
				//	if err := srv.ProcessTxs(ctx, block, true); err != nil {
				//		lgr.Debug("Failed to process txs", zap.Error(err))
				//	}
				//
				//	go func() {
				//		if err := srv.ProcessLogsOfTxs(ctx, block.Txs, block.Time); err != nil {
				//			lgr.Debug("cannot process logs", zap.Error(err))
				//		}
				//
				//		if err := srv.FilterProposalEvent(ctx, block.Txs); err != nil {
				//			lgr.Debug("filter proposal event failed", zap.Error(err))
				//		}
				//		if err := srv.ProcessActiveAddress(ctx, block.Txs); err != nil {
				//			lgr.Debug("failed to process active address", zap.Error(err))
				//		}
				//	}()

				//lgr.Debug("Total import block time", zap.Duration("TotalTime", time.Since(totalImportTime)))
				//if latest-1 > prevHeader {
				//	lgr.Warn("we are behind network, inserting error blocks", zap.Uint64("from", prevHeader), zap.Uint64("to", latest))
				//	err := srv.InsertErrorBlocks(ctx, prevHeader, latest)
				//	if err != nil {
				//		lgr.Error("failed to insert error block height", zap.Error(err))
				//		continue
				//	}
				//}
				//prevHeader = latest
				//if latest%cfg.UpdateStatsInterval == 0 {
				//	_ = srv.UpdateCurrentStats(ctx)
				//}
			}
			fmt.Println("latest block number", latest)
			block, err := g.web3Client.U2U.GetBlockByNumber(new(big.Int).SetUint64(latest), true)
			if err != nil {
				continue
			}
			fmt.Printf("Latest block info \n")
			fmt.Printf("Block header: %+v \n", block.Header())
			fmt.Printf("Block body: %+v \n", block.Body())
			// Do insert

		}
	}
}

func (g *Grabber) fetchBlock(done chan<- types.Block, it types.Block) {
	defer func() { done <- it }()
	//for it.Next() {
	//	select {
	//	case c.ch <- it.Node():
	//	case <-c.closed:
	//		return
	//	}
	//}
}

func (g *Grabber) updateBlock() {

}

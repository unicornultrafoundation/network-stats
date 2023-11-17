package block

import (
	"context"
	"github.com/unicornultrafoundation/network-stats/pkg/jsonrpc"
	"time"
)

type RealtimeFetcher struct {
	Web3Client *jsonrpc.Web3
}

func (f *RealtimeFetcher) importLatestBlock(ctx context.Context) {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			latest, err := f.Web3Client.U2U.GetBlockNumber()
			if err != nil {
				continue
			}

			if latest <= prevHeader {
				continue
			}
			if prevHeader < latest {
				startTime = time.Now()
				block, err := srv.BlockByHeight(ctx, latest)
				if err != nil {
					lgr.Error("Failed to get block from RPC", zap.Error(err))
					continue
				}
				lgr.Info("Block info from network", zap.Any("Block", block))
				endTime = time.Since(startTime)
				srv.Metrics().RecordScrapingTime(endTime)
				lgr.Info("scraping block time", zap.Duration("TimeConsumed", endTime), zap.String("Avg", srv.Metrics().GetScrapingTime()))
				if block == nil {
					lgr.Error("Block not found")
					continue
				}
				// insert current block height to cache for re-verifying later
				// temp remove insert new unverified blocks
				err = srv.InsertUnverifiedBlocks(ctx, latest)
				if err != nil {
					lgr.Error("Failed to insert unverified block", zap.Error(err))
				}
				// import this latest block to cache and database
				totalImportTime := time.Now()
				if err := srv.ImportBlock(ctx, block, true); err != nil {
					lgr.Debug("Failed to import block", zap.Error(err))
					continue
				}

				if err := srv.ProcessTxs(ctx, block, true); err != nil {
					lgr.Debug("Failed to process txs", zap.Error(err))
				}

				go func() {
					if err := srv.ProcessLogsOfTxs(ctx, block.Txs, block.Time); err != nil {
						lgr.Debug("cannot process logs", zap.Error(err))
					}

					if err := srv.FilterProposalEvent(ctx, block.Txs); err != nil {
						lgr.Debug("filter proposal event failed", zap.Error(err))
					}
					if err := srv.ProcessActiveAddress(ctx, block.Txs); err != nil {
						lgr.Debug("failed to process active address", zap.Error(err))
					}
				}()

				lgr.Debug("Total import block time", zap.Duration("TotalTime", time.Since(totalImportTime)))
				if latest-1 > prevHeader {
					lgr.Warn("we are behind network, inserting error blocks", zap.Uint64("from", prevHeader), zap.Uint64("to", latest))
					err := srv.InsertErrorBlocks(ctx, prevHeader, latest)
					if err != nil {
						lgr.Error("failed to insert error block height", zap.Error(err))
						continue
					}
				}
				prevHeader = latest
				if latest%cfg.UpdateStatsInterval == 0 {
					_ = srv.UpdateCurrentStats(ctx)
				}
			}
		}
	}

}

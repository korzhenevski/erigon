package observer

import (
	"context"
	"errors"
	"fmt"
	"github.com/ledgerwatch/erigon/core/forkid"
	"github.com/ledgerwatch/erigon/p2p/enode"
	"github.com/ledgerwatch/erigon/params"
	"github.com/ledgerwatch/log/v3"
)

type Crawler struct {
	transport  DiscV4Transport
	bootnodes  []*enode.Node
	forkFilter forkid.Filter
	log        log.Logger
}

func NewCrawler(
	transport DiscV4Transport,
	bootnodes []*enode.Node,
	chain string,
	logger log.Logger,
) (*Crawler, error) {
	chainConfig := params.ChainConfigByChainName(chain)
	genesisHash := params.GenesisHashByChainName(chain)
	if (chainConfig == nil) || (genesisHash == nil) {
		return nil, fmt.Errorf("unknown chain %s", chain)
	}

	forkFilter := forkid.NewStaticFilter(chainConfig, *genesisHash)

	instance := Crawler{
		transport,
		bootnodes,
		forkFilter,
		logger,
	}
	return &instance, nil
}

func (crawler *Crawler) startSelectCandidates(ctx context.Context) <-chan *enode.Node {
	nodes := make(chan *enode.Node)
	go func() {
		err := crawler.selectCandidates(ctx, nodes)
		if (err != nil) && !errors.Is(err, context.Canceled) {
			crawler.log.Error("Failed to select candidates", "err", err)
		}
	}()
	return nodes
}

func (crawler *Crawler) selectCandidates(ctx context.Context, nodes chan<- *enode.Node) error {
	for _, node := range crawler.bootnodes {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case nodes <- node:
		}
	}
	return nil
}

func (crawler *Crawler) Run(ctx context.Context) error {
	nodes := crawler.startSelectCandidates(ctx)

	for node := range nodes {
		nodeDesc := node.String()
		logger := crawler.log.New("node", nodeDesc)

		interrogator, err := NewInterrogator(node, crawler.transport, crawler.forkFilter, logger)
		if err != nil {
			return fmt.Errorf("failed to create Interrogator for node %s: %w", nodeDesc, err)
		}

		go func() {
			peers, err := interrogator.Run(ctx)
			if err != nil {
				logger.Warn("Failed to interrogate node", "err", err)
				return
			}

			// TODO save to DB
			logger.Debug(fmt.Sprintf("Got %d peers", len(peers)))
		}()
	}
	return nil
}

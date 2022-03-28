package observer

import (
	"context"
	"fmt"
	"github.com/ledgerwatch/erigon/p2p/enode"
	"github.com/ledgerwatch/log/v3"
	"time"
)

type Interrogator struct {
	node      *enode.Node
	transport DiscV4Transport
	log       log.Logger
}

func NewInterrogator(node *enode.Node, transport DiscV4Transport, logger log.Logger) (*Interrogator, error) {
	instance := Interrogator{
		node,
		transport,
		logger,
	}
	return &instance, nil
}

func (interrogator *Interrogator) Run(ctx context.Context) ([]*enode.Node, error) {
	interrogator.log.Debug("Interrogating a node")

	err := interrogator.transport.Ping(interrogator.node)
	if err != nil {
		return nil, fmt.Errorf("ping-pong failed: %w", err)
	}

	enr, err := interrogator.transport.RequestENR(interrogator.node)
	if err != nil {
		return nil, fmt.Errorf("ENR request failed: %w", err)
	}

	// TODO filter enr
	interrogator.log.Debug("Got ENR", "enr", enr)

	keys := keygen(ctx, interrogator.node.Pubkey(), 10*time.Second, interrogator.log)
	interrogator.log.Debug(fmt.Sprintf("Generated %d keys", len(keys)))

	peersByID := make(map[enode.ID]*enode.Node)
	for _, key := range keys {
		neighbors := interrogator.transport.LookupPubkey(key)

		for _, node := range neighbors {
			peersByID[node.ID()] = node
		}
	}

	peers := valuesOfIDToNodeMap(peersByID)

	return peers, nil
}

func valuesOfIDToNodeMap(m map[enode.ID]*enode.Node) []*enode.Node {
	values := make([]*enode.Node, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

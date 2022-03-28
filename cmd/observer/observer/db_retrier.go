package observer

import (
	"context"
	"github.com/ledgerwatch/erigon/p2p/enode"
	"github.com/ledgerwatch/log/v3"
	"math/rand"
	"time"
)

type DBRetrier struct {
	db  DB
	log log.Logger
}

func NewDBRetrier(db DB, logger log.Logger) DBRetrier {
	return DBRetrier{db, logger}
}

const retryCount = 32

func retryBackoffTime(attempt int) time.Duration {
	if attempt <= 0 { return 0 }
	jitter := rand.Int63n(20 * time.Millisecond.Nanoseconds() * int64(attempt))
	var ns int64
	if attempt <= 5 {
		ns = ((50 * time.Millisecond.Nanoseconds()) << (attempt - 1)) + jitter
	} else {
		ns = time.Second.Nanoseconds() + jitter
	}
	return time.Duration(ns)
}

func (db DBRetrier) UpsertNode(ctx context.Context, node *enode.Node) error {
	var err error
	for i := 0; i <= retryCount; i += 1 {
		if i > 0 {
			db.log.Trace("retrying UpsertNode", "attempt", i, "err", err)
		}
		sleep(ctx, retryBackoffTime(i))
		err = db.db.UpsertNode(ctx, node)
		if (err == nil) || !db.db.IsConflictError(err) {
			break
		}
	}
	return err
}

func (db DBRetrier) UpdateForkCompatibility(ctx context.Context, node *enode.Node, isCompatFork bool) error {
	var err error
	for i := 0; i <= retryCount; i += 1 {
		if i > 0 {
			db.log.Trace("retrying UpdateForkCompatibility", "attempt", i, "err", err)
		}
		sleep(ctx, retryBackoffTime(i))
		err = db.db.UpdateForkCompatibility(ctx, node, isCompatFork)
		if (err == nil) || !db.db.IsConflictError(err) {
			break
		}
	}
	return err
}

func (db DBRetrier) UpdateClientID(ctx context.Context, node *enode.Node, clientID string) error {
	var err error
	for i := 0; i <= retryCount; i += 1 {
		if i > 0 {
			db.log.Trace("retrying UpdateClientID", "attempt", i, "err", err)
		}
		sleep(ctx, retryBackoffTime(i))
		err = db.db.UpdateClientID(ctx, node, clientID)
		if (err == nil) || !db.db.IsConflictError(err) {
			break
		}
	}
	return err
}

func (db DBRetrier) UpdateHandshakeError(ctx context.Context, node *enode.Node, handshakeErr string) error {
	var err error
	for i := 0; i <= retryCount; i += 1 {
		if i > 0 {
			db.log.Trace("retrying UpdateHandshakeError", "attempt", i, "err", err)
		}
		sleep(ctx, retryBackoffTime(i))
		err = db.db.UpdateHandshakeError(ctx, node, handshakeErr)
		if (err == nil) || !db.db.IsConflictError(err) {
			break
		}
	}
	return err
}

func (db DBRetrier) TakeCandidates(ctx context.Context, minUnusedDuration time.Duration, limit uint) ([]*enode.Node, error) {
	var result []*enode.Node
	var err error
	for i := 0; i <= retryCount; i += 1 {
		if i > 0 {
			db.log.Trace("retrying TakeCandidates", "attempt", i, "err", err)
		}
		sleep(ctx, retryBackoffTime(i))
		result, err = db.db.TakeCandidates(ctx, minUnusedDuration, limit)
		if (err == nil) || !db.db.IsConflictError(err) {
			break
		}
	}
	return result, err
}

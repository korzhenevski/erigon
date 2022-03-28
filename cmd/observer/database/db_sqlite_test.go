package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"path/filepath"
	"testing"
	"time"
)

func TestDBSQLiteInsertAndFind(t *testing.T) {
	ctx := context.Background()
	db, err := NewDBSQLite(filepath.Join(t.TempDir(), "observer.sqlite"))
	require.Nil(t, err)

	var id NodeID = "ba85011c70bcc5c04d8607d3a0ed29aa6179c092cbdda10d5d32684fb33ed01bd94f588ca8f91ac48318087dcb02eaf36773a7a453f0eedd6742af668097b29c"
	var addr NodeAddr
	addr.IP = net.ParseIP("10.0.1.16")
	addr.PortRLPx = 30303
	addr.PortDisc = 30304

	err = db.UpsertNodeAddr(ctx, id, addr)
	require.Nil(t, err)

	candidates, err := db.FindCandidates(ctx, time.Second, 1)
	require.Nil(t, err)
	require.Equal(t, 1, len(candidates))

	var candidateID NodeID
	var candidate NodeAddr
	for candidateID, candidate = range candidates {
	}

	assert.Equal(t, id, candidateID)
	assert.Equal(t, addr.IP, candidate.IP)
	assert.Equal(t, addr.PortDisc, candidate.PortDisc)
	assert.Equal(t, addr.PortRLPx, candidate.PortRLPx)
}

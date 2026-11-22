package engine

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jkusuda/orbital/backend/internal/orbit"
)

func TestGRPCMatchesLocal(t *testing.T) {
	addr := os.Getenv("ENGINE_ADDR")
	if addr == "" {
		t.Skip("ENGINE_ADDR not set")
	}
	g, err := NewGRPC(addr)
	require.NoError(t, err)
	defer g.Close()

	els := []orbit.Elements{{
		Epoch:         2461200.5,
		Eccentricity:  0.1911492279663492,
		SemiMajorAxis: 0.9223592206975018,
		Inclination:   3.340996879880978,
		AscendingNode: 203.8936514240762,
		PerihelionArg: 126.6795706895841,
		MeanAnomaly:   175.3304026592739,
	}}
	start := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)
	times := []time.Time{start, start.Add(24 * time.Hour), start.Add(48 * time.Hour)}

	want, _ := Local{}.Propagate(context.Background(), els, times)
	got, err := g.Propagate(context.Background(), els, times)
	require.NoError(t, err)
	for i := range times {
		assert.InDelta(t, want[0][i].X, got[0][i].X, 1e-12)
		assert.InDelta(t, want[0][i].Y, got[0][i].Y, 1e-12)
		assert.InDelta(t, want[0][i].Z, got[0][i].Z, 1e-12)
	}
}

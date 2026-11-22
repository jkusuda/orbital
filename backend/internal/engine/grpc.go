package engine

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orbitalv1 "github.com/jkusuda/orbital/backend/internal/gen/orbitalv1"
	"github.com/jkusuda/orbital/backend/internal/orbit"
)

// GRPC calls the C++ engine.
type GRPC struct {
	conn   *grpc.ClientConn
	client orbitalv1.PropagatorServiceClient
}

// maxResponseBytes raises gRPC's default 4 MB limit on received messages. Each
// point is ~29 bytes on the wire, so 2000 samples × ~70 asteroids would hit 4 MB.
const maxResponseBytes = 64 << 20

func NewGRPC(addr string) (*GRPC, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxResponseBytes)),
	)
	if err != nil {
		return nil, fmt.Errorf("engine: connect %s: %w", addr, err)
	}
	return &GRPC{conn: conn, client: orbitalv1.NewPropagatorServiceClient(conn)}, nil
}

func (g *GRPC) Close() error { return g.conn.Close() }

func (g *GRPC) Propagate(ctx context.Context, els []orbit.Elements, times []time.Time) ([][]orbit.Vec3, error) {
	req := &orbitalv1.PropagateRequest{
		Asteroids: make([]*orbitalv1.Elements, len(els)),
		TimesJd:   make([]float64, len(times)),
	}
	for j, el := range els {
		req.Asteroids[j] = &orbitalv1.Elements{
			EpochJd:          el.Epoch,
			Eccentricity:     el.Eccentricity,
			SemiMajorAxisAu:  el.SemiMajorAxis,
			InclinationDeg:   el.Inclination,
			AscendingNodeDeg: el.AscendingNode,
			PerihelionArgDeg: el.PerihelionArg,
			MeanAnomalyDeg:   el.MeanAnomaly,
		}
	}
	for i, t := range times {
		req.TimesJd[i] = orbit.JulianDate(t)
	}

	resp, err := g.client.Propagate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("engine: propagate: %w", err)
	}

	trajs := resp.GetTrajectories()
	if len(trajs) != len(els) {
		return nil, fmt.Errorf("engine: got %d trajectories, want %d", len(trajs), len(els))
	}
	out := make([][]orbit.Vec3, len(els))
	for j, tr := range trajs {
		pts := tr.GetPoints()
		if len(pts) != len(times) {
			return nil, fmt.Errorf("engine: trajectory %d has %d points, want %d", j, len(pts), len(times))
		}
		out[j] = make([]orbit.Vec3, len(times))
		for i, p := range pts {
			out[j][i] = orbit.Vec3{X: p.GetX(), Y: p.GetY(), Z: p.GetZ()}
		}
	}
	return out, nil
}

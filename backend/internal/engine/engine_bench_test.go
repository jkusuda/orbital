package engine

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jkusuda/orbital/backend/internal/orbit"
)

// Batch sizes as asteroids × times. 50×169 is roughly a real /api/trajectories
// call (a week of hourly samples); 50×2000 is the API's max sample count.
var benchSizes = []struct{ asteroids, times int }{
	{1, 1},
	{10, 100},
	{50, 169},
	{50, 2000},
	{1000, 1000},
}

// Run with the engine up
func BenchmarkPropagate(b *testing.B) {
	impls := []struct {
		name string
		prop Propagator
	}{{"go", Local{}}}

	if addr := os.Getenv("ENGINE_ADDR"); addr != "" {
		g, err := NewGRPC(addr)
		if err != nil {
			b.Fatal(err)
		}
		defer g.Close()
		impls = append(impls, struct {
			name string
			prop Propagator
		}{"cpp-grpc", g})
	}

	ctx := context.Background()
	for _, size := range benchSizes {
		els, times := benchInputs(size.asteroids, size.times)
		points := size.asteroids * size.times

		for _, impl := range impls {
			name := fmt.Sprintf("%s/%dx%d", impl.name, size.asteroids, size.times)
			b.Run(name, func(b *testing.B) {
				if _, err := impl.prop.Propagate(ctx, els, times); err != nil {
					b.Fatal(err)
				}
				for b.Loop() {
					if _, err := impl.prop.Propagate(ctx, els, times); err != nil {
						b.Fatal(err)
					}
				}
				b.ReportMetric(float64(points)*float64(b.N)/b.Elapsed().Seconds(), "points/s")
			})
		}
	}
}

func benchInputs(nAsteroids, nTimes int) ([]orbit.Elements, []time.Time) {
	els := make([]orbit.Elements, nAsteroids)
	for j := range els {
		els[j] = orbit.Elements{
			Epoch:         2461200.5,
			Eccentricity:  0.05 + float64(j%90)/100, // 0.05–0.94
			SemiMajorAxis: 0.8 + float64(j%25)/10,   // 0.8–3.2 AU
			Inclination:   float64(j % 40),
			AscendingNode: float64((j * 37) % 360),
			PerihelionArg: float64((j * 53) % 360),
			MeanAnomaly:   float64(j) * 7.3,
		}
	}
	start := time.Date(2026, 11, 22, 0, 0, 0, 0, time.UTC)
	times := make([]time.Time, nTimes)
	for i := range times {
		times[i] = start.Add(time.Duration(i) * time.Hour)
	}
	return els, times
}

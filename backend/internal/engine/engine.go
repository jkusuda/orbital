package engine

import (
	"context"
	"time"

	"github.com/jkusuda/orbital/backend/internal/orbit"
)

type Propagator interface {
	Propagate(ctx context.Context, els []orbit.Elements, times []time.Time) ([][]orbit.Vec3, error)
}

type Local struct{}

func (Local) Propagate(_ context.Context, els []orbit.Elements, times []time.Time) ([][]orbit.Vec3, error) {
	out := make([][]orbit.Vec3, len(els))
	for j, el := range els {
		out[j] = make([]orbit.Vec3, len(times))
		for i, t := range times {
			out[j][i] = orbit.Position(el, t)
		}
	}
	return out, nil
}

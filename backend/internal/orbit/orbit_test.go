package orbit_test

import (
	"testing"
	"time"

	"github.com/jkusuda/orbital/backend/internal/orbit"
	"github.com/stretchr/testify/assert"
)

func TestPosition(t *testing.T) {
	type testCase struct {
		el       orbit.Elements
		at       time.Time
		expected orbit.Vec3
	}

	t.Run("Position matches NASA Horizons data", func(t *testing.T) {
		tests := []testCase{
			{
				// Apophis (2099942) elements from NeoWs
				el: orbit.Elements{
					Epoch:         2461200.5,
					Eccentricity:  0.1911492279663492,
					SemiMajorAxis: 0.9223592206975018,
					Inclination:   3.340996879880978,
					AscendingNode: 203.8936514240762,
					PerihelionArg: 126.6795706895841,
					MeanAnomaly:   175.3304026592739,
				},
				at:       time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC),
				expected: orbit.Vec3{X: 0.5062187696, Y: -0.5546101423, Z: 0.0415725799},
			},
		}

		for _, test := range tests {
			actual := orbit.Position(test.el, test.at)
			assert.InDelta(t, test.expected.X, actual.X, 5e-4)
			assert.InDelta(t, test.expected.Y, actual.Y, 5e-4)
			assert.InDelta(t, test.expected.Z, actual.Z, 5e-4)
		}
	})
}

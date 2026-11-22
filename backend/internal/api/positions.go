package api

import (
	"net/http"
	"time"

	"github.com/jkusuda/orbital/backend/internal/orbit"
)

type positionsResponse struct {
	Time      time.Time          `json:"time"`
	Sun       [3]float64         `json:"sun"`
	Positions []asteroidPosition `json:"positions"`
}

type asteroidPosition struct {
	asteroidInfo
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (s *Server) handlePositions(w http.ResponseWriter, r *http.Request) {
	t, err := timeParam(r, "time", time.Now().UTC())
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	asteroids, ok := s.asteroidsApproaching(w, r, t.Add(-approachingWindow), t.Add(approachingWindow))
	if !ok {
		return
	}

	pos, ok := s.propagate(w, r, asteroids, []time.Time{t})
	if !ok {
		return
	}

	earth := orbit.EarthPosition(t)
	resp := positionsResponse{
		Time:      t,
		Sun:       vec(orbit.Vec3{}.Sub(earth)),
		Positions: make([]asteroidPosition, 0, len(asteroids)),
	}
	for j, a := range asteroids {
		p := pos[j][0].Sub(earth)
		resp.Positions = append(resp.Positions, asteroidPosition{
			asteroidInfo: infoOf(a),
			X:            p.X,
			Y:            p.Y,
			Z:            p.Z,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

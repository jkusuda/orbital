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

	earth := orbit.EarthPosition(t)
	resp := positionsResponse{
		Time:      t,
		Sun:       vec(orbit.Vec3{}.Sub(earth)),
		Positions: make([]asteroidPosition, 0, len(s.asteroids)),
	}
	for _, a := range s.asteroids {
		p := orbit.Position(toElements(a.OrbitalData), t).Sub(earth)
		resp.Positions = append(resp.Positions, asteroidPosition{
			asteroidInfo: infoOf(a),
			X:            p.X,
			Y:            p.Y,
			Z:            p.Z,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

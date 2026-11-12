package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jkusuda/orbital/backend/internal/orbit"
)

const maxSamples = 2000

type trajectoriesResponse struct {
	Times     []time.Time          `json:"times"`
	Sun       [][3]float64         `json:"sun"`
	Asteroids []asteroidTrajectory `json:"asteroids"`
}

type asteroidTrajectory struct {
	asteroidInfo
	Points [][3]float64 `json:"points"`
}

func (s *Server) handleTrajectories(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	start, err := timeParam(r, "start", now.Add(-7*24*time.Hour))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	end, err := timeParam(r, "end", now.Add(7*24*time.Hour))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	step, err := durationParam(r, "step", 2*time.Hour)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !end.After(start) {
		writeError(w, http.StatusBadRequest, "end must be after start")
		return
	}
	n := int(end.Sub(start)/step) + 1
	if n > maxSamples {
		writeError(w, http.StatusBadRequest,
			fmt.Sprintf("range/step gives %d samples, max is %d: use a larger step or a shorter range", n, maxSamples))
		return
	}

	resp := trajectoriesResponse{
		Times:     make([]time.Time, n),
		Sun:       make([][3]float64, n),
		Asteroids: make([]asteroidTrajectory, len(s.asteroids)),
	}
	elements := make([]orbit.Elements, len(s.asteroids))
	for j, a := range s.asteroids {
		elements[j] = toElements(a.OrbitalData)
		resp.Asteroids[j] = asteroidTrajectory{
			asteroidInfo: infoOf(a),
			Points:       make([][3]float64, n),
		}
	}

	for i := range n {
		t := start.Add(time.Duration(i) * step)
		earth := orbit.EarthPosition(t)
		resp.Times[i] = t
		resp.Sun[i] = vec(orbit.Vec3{}.Sub(earth))
		for j, el := range elements {
			resp.Asteroids[j].Points[i] = vec(orbit.Position(el, t).Sub(earth))
		}
	}
	writeJSON(w, http.StatusOK, resp)
}

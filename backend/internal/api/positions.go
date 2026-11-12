package api

import (
	"net/http"
	"time"

	"github.com/jkusuda/orbital/backend/internal/nasa"
	"github.com/jkusuda/orbital/backend/internal/orbit"
)

type positionsResponse struct {
	Time      time.Time          `json:"time"`
	Positions []asteroidPosition `json:"positions"`
}

type asteroidPosition struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Hazardous  bool    `json:"hazardous"`
	DiameterKm float64 `json:"diameter_km"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Z          float64 `json:"z"`
}

func (s *Server) handlePositions(w http.ResponseWriter, r *http.Request) {
	t := time.Now().UTC()
	if q := r.URL.Query().Get("time"); q != "" {
		parsed, err := time.Parse(time.RFC3339, q)
		if err != nil {
			writeError(w, http.StatusBadRequest, "time must be RFC3339, e.g. 2026-11-11T00:00:00Z")
			return
		}
		t = parsed
	}

	earth := orbit.EarthPosition(t)
	resp := positionsResponse{
		Time:      t,
		Positions: make([]asteroidPosition, 0, len(s.asteroids)),
	}
	for _, a := range s.asteroids {
		p := orbit.Position(toElements(a.OrbitalData), t).Sub(earth)
		resp.Positions = append(resp.Positions, asteroidPosition{
			ID:         a.ID,
			Name:       a.Name,
			Hazardous:  a.Hazardous,
			DiameterKm: (a.Diameter.Kilometers.Min + a.Diameter.Kilometers.Max) / 2,
			X:          p.X,
			Y:          p.Y,
			Z:          p.Z,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

func toElements(od nasa.OrbitalData) orbit.Elements {
	return orbit.Elements{
		Epoch:         od.Epoch,
		Eccentricity:  od.Eccentricity,
		SemiMajorAxis: od.SemiMajorAxis,
		Inclination:   od.Inclination,
		AscendingNode: od.AscendingNode,
		PerihelionArg: od.PerihelionArg,
		MeanAnomaly:   od.MeanAnomaly,
	}
}

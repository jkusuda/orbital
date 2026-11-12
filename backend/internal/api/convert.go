package api

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/jkusuda/orbital/backend/internal/nasa"
	"github.com/jkusuda/orbital/backend/internal/orbit"
)

type asteroidInfo struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Hazardous  bool    `json:"hazardous"`
	DiameterKm float64 `json:"diameter_km"`
	JPLURL     string  `json:"jpl_url"` // JPL Small-Body Database page
}

func infoOf(a nasa.Asteroid) asteroidInfo {
	return asteroidInfo{
		ID:         a.ID,
		Name:       a.Name,
		Hazardous:  a.Hazardous,
		DiameterKm: (a.Diameter.Kilometers.Min + a.Diameter.Kilometers.Max) / 2,
		JPLURL:     a.JPLURL,
	}
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

func vec(v orbit.Vec3) [3]float64 {
	round := func(f float64) float64 { return math.Round(f*1e8) / 1e8 }
	return [3]float64{round(v.X), round(v.Y), round(v.Z)}
}

func timeParam(r *http.Request, name string, def time.Time) (time.Time, error) {
	q := r.URL.Query().Get(name)
	if q == "" {
		return def, nil
	}
	t, err := time.Parse(time.RFC3339, q)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s must be RFC3339, e.g. 2026-11-11T00:00:00Z", name)
	}
	return t, nil
}

func durationParam(r *http.Request, name string, def time.Duration) (time.Duration, error) {
	q := r.URL.Query().Get(name)
	if q == "" {
		return def, nil
	}
	d, err := time.ParseDuration(q)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration, e.g. 2h", name)
	}
	return d, nil
}

package ingest

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jkusuda/orbital/backend/internal/nasa"
	"github.com/jkusuda/orbital/backend/internal/orbit"
	"github.com/jkusuda/orbital/backend/internal/store"
)

const feedDays = 7

type Ingester struct {
	nasa  *nasa.Client
	store *store.Store
}

func New(client *nasa.Client, st *store.Store) *Ingester {
	return &Ingester{nasa: client, store: st}
}

func (in *Ingester) Refresh(ctx context.Context) error {
	started := time.Now()
	from := started.UTC().Format("2006-01-02")
	to := started.UTC().AddDate(0, 0, feedDays-1).Format("2006-01-02")

	feed, err := in.nasa.GetFeed(ctx, from, to)
	if err != nil {
		return err
	}

	saved, failed := 0, 0
	for _, a := range feed {
		if err := in.saveOne(ctx, a); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("ingest: skipping %s: %v", a.ID, err)
			failed++
			continue
		}
		saved++
	}

	log.Printf("ingest: saved %d, failed %d, took %s", saved, failed, time.Since(started).Round(time.Millisecond))
	if failed > 0 {
		return fmt.Errorf("ingest: %d of %d asteroids failed", failed, len(feed))
	}
	return nil
}

func (in *Ingester) saveOne(ctx context.Context, fromFeed nasa.Asteroid) error {
	neo, err := in.nasa.GetNEO(ctx, fromFeed.ID)
	if err != nil {
		return err
	}
	return in.store.SaveAsteroid(ctx, toStore(*neo), earthApproaches(fromFeed.CloseApproaches))
}

func toStore(a nasa.Asteroid) store.Asteroid {
	od := a.OrbitalData
	return store.Asteroid{
		ID:         a.ID,
		Name:       a.Name,
		Hazardous:  a.Hazardous,
		DiameterKm: (a.Diameter.Kilometers.Min + a.Diameter.Kilometers.Max) / 2,
		JPLURL:     a.JPLURL,
		Elements: orbit.Elements{
			Epoch:         od.Epoch,
			Eccentricity:  od.Eccentricity,
			SemiMajorAxis: od.SemiMajorAxis,
			Inclination:   od.Inclination,
			AscendingNode: od.AscendingNode,
			PerihelionArg: od.PerihelionArg,
			MeanAnomaly:   od.MeanAnomaly,
		},
	}
}

func earthApproaches(cas []nasa.CloseApproach) []store.CloseApproach {
	var out []store.CloseApproach
	for _, ca := range cas {
		if ca.OrbitingBody != "Earth" {
			continue
		}
		out = append(out, store.CloseApproach{
			Time:           time.UnixMilli(ca.EpochMs).UTC(),
			MissDistanceKm: ca.MissDistance.Km,
			SpeedKmS:       ca.Velocity.KmPerSec,
		})
	}
	return out
}

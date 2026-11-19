package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jkusuda/orbital/backend/internal/orbit"
)

type Store struct {
	pool *pgxpool.Pool
}

type Asteroid struct {
	ID         string
	Name       string
	Hazardous  bool
	DiameterKm float64
	JPLURL     string
	Elements   orbit.Elements
}

type CloseApproach struct {
	Time           time.Time
	MissDistanceKm float64
	SpeedKmS       float64
}

func New(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("store: connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("store: ping: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close()                         { s.pool.Close() }
func (s *Store) Ping(ctx context.Context) error { return s.pool.Ping(ctx) }

func (s *Store) SaveAsteroid(ctx context.Context, a Asteroid, approaches []CloseApproach) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("store: begin: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO asteroids (id, name, hazardous, diameter_km, jpl_url)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			name         = EXCLUDED.name,
			hazardous    = EXCLUDED.hazardous,
			diameter_km  = EXCLUDED.diameter_km,
			jpl_url      = EXCLUDED.jpl_url,
			last_updated = now()`,
		a.ID, a.Name, a.Hazardous, a.DiameterKm, a.JPLURL)
	if err != nil {
		return fmt.Errorf("store: upsert asteroid %s: %w", a.ID, err)
	}

	el := a.Elements
	_, err = tx.Exec(ctx, `
		INSERT INTO orbital_elements (asteroid_id, epoch_jd, eccentricity, semi_major_axis_au,
			inclination_deg, ascending_node_deg, perihelion_arg_deg, mean_anomaly_deg)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (asteroid_id) DO UPDATE SET
			epoch_jd           = EXCLUDED.epoch_jd,
			eccentricity       = EXCLUDED.eccentricity,
			semi_major_axis_au = EXCLUDED.semi_major_axis_au,
			inclination_deg    = EXCLUDED.inclination_deg,
			ascending_node_deg = EXCLUDED.ascending_node_deg,
			perihelion_arg_deg = EXCLUDED.perihelion_arg_deg,
			mean_anomaly_deg   = EXCLUDED.mean_anomaly_deg`,
		a.ID, el.Epoch, el.Eccentricity, el.SemiMajorAxis,
		el.Inclination, el.AscendingNode, el.PerihelionArg, el.MeanAnomaly)
	if err != nil {
		return fmt.Errorf("store: upsert orbit %s: %w", a.ID, err)
	}

	for _, ca := range approaches {
		_, err = tx.Exec(ctx, `
			INSERT INTO close_approaches (asteroid_id, approach_time, miss_distance_km, speed_km_s)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (asteroid_id, approach_time) DO UPDATE SET
				miss_distance_km = EXCLUDED.miss_distance_km,
				speed_km_s       = EXCLUDED.speed_km_s`,
			a.ID, ca.Time, ca.MissDistanceKm, ca.SpeedKmS)
		if err != nil {
			return fmt.Errorf("store: upsert approach %s: %w", a.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("store: commit %s: %w", a.ID, err)
	}
	return nil
}

func (s *Store) AsteroidsApproaching(ctx context.Context, from, to time.Time) ([]Asteroid, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT a.id, a.name, a.hazardous, a.diameter_km, a.jpl_url,
		       o.epoch_jd, o.eccentricity, o.semi_major_axis_au, o.inclination_deg,
		       o.ascending_node_deg, o.perihelion_arg_deg, o.mean_anomaly_deg
		FROM asteroids a
		JOIN orbital_elements o ON o.asteroid_id = a.id
		WHERE EXISTS (
			SELECT 1 FROM close_approaches c
			WHERE c.asteroid_id = a.id AND c.approach_time BETWEEN $1 AND $2
		)
		ORDER BY a.id`, from, to)
	if err != nil {
		return nil, fmt.Errorf("store: query asteroids: %w", err)
	}
	defer rows.Close()

	var out []Asteroid
	for rows.Next() {
		var a Asteroid
		el := &a.Elements
		if err := rows.Scan(&a.ID, &a.Name, &a.Hazardous, &a.DiameterKm, &a.JPLURL,
			&el.Epoch, &el.Eccentricity, &el.SemiMajorAxis, &el.Inclination,
			&el.AscendingNode, &el.PerihelionArg, &el.MeanAnomaly); err != nil {
			return nil, fmt.Errorf("store: scan asteroid: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

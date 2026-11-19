-- +goose Up
CREATE TABLE asteroids (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    hazardous    BOOLEAN NOT NULL,
    diameter_km  DOUBLE PRECISION NOT NULL,
    jpl_url      TEXT NOT NULL,
    first_seen   TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_updated TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE orbital_elements (
    asteroid_id        TEXT PRIMARY KEY REFERENCES asteroids(id) ON DELETE CASCADE,
    epoch_jd           DOUBLE PRECISION NOT NULL,
    eccentricity       DOUBLE PRECISION NOT NULL,
    semi_major_axis_au DOUBLE PRECISION NOT NULL,
    inclination_deg    DOUBLE PRECISION NOT NULL,
    ascending_node_deg DOUBLE PRECISION NOT NULL,
    perihelion_arg_deg DOUBLE PRECISION NOT NULL,
    mean_anomaly_deg   DOUBLE PRECISION NOT NULL
);

CREATE TABLE close_approaches (
    asteroid_id      TEXT NOT NULL REFERENCES asteroids(id) ON DELETE CASCADE,
    approach_time    TIMESTAMPTZ NOT NULL,
    miss_distance_km DOUBLE PRECISION NOT NULL,
    speed_km_s       DOUBLE PRECISION NOT NULL,
    PRIMARY KEY (asteroid_id, approach_time)
);

CREATE INDEX close_approaches_time_idx ON close_approaches (approach_time);

-- +goose Down
DROP TABLE close_approaches;
DROP TABLE orbital_elements;
DROP TABLE asteroids;

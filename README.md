# Orbital

A 3D viewer for the asteroids passing close to Earth this week. Orbital pulls data from NASA, computes where each asteroid is relative to Earth, and draws them in the browser with a timeline you can scrub ±7 days.

## How it works

```
NASA ──► Go API ──► Postgres
            │
            ├──► Browser (Three.js)
            └──► C++ engine over gRPC (optional)
```

- **Go API**: refreshes from NASA every 6 hours, stores asteroids in Postgres, and serves positions over REST.
- **Orbit math**: Kepler's equation turns NASA's orbital elements into positions. Checked against JPL Horizons.
- **C++ engine**: the same math as a gRPC service. Set `ENGINE_ADDR` and the API uses it instead of Go.
- **Frontend**: Three.js and Vite.

## Run it

```bash
docker compose up -d                          # Postgres
cd backend && make migrate-up && make run     # API on :8080 (needs backend/.env)
cd frontend && npm install && npm run dev     # app on :5173
cd engine && make run                         # optional C++ engine on :50051
```

Copy `backend/.env.example` to `backend/.env` and add a free NASA key from [api.nasa.gov](https://api.nasa.gov).

## Benchmark

Go in-process vs C++ over gRPC (`make bench` in `backend/`):

| Points | Go | C++ (gRPC) |
|---|---|---|
| 1 | 0.32 µs | 96 µs |
| 8.4k (typical request) | 1.9 ms | 2.3 ms |
| 1M | 277 ms | 247 ms |

Each gRPC call costs about 95 µs, so C++ only wins on very large batches. Go stays the default.

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jkusuda/orbital/backend/internal/store"
)

// Server holds the dependencies
type Server struct {
	store *store.Store
}

func NewServer(st *store.Store) *Server {
	return &Server{store: st}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", s.handleHealth)
	mux.HandleFunc("GET /api/positions", s.handlePositions)
	mux.HandleFunc("GET /api/trajectories", s.handleTrajectories)
	return mux
}

// approachingWindow is how far either side of the requested time an asteroid's
// close approach can be for it to be shown.
const approachingWindow = 7 * 24 * time.Hour

// asteroidsApproaching loads asteroids with a close approach between from and to.
// On failure it writes a 500 (logging the real error, not sending it) and returns false.
func (s *Server) asteroidsApproaching(w http.ResponseWriter, r *http.Request, from, to time.Time) ([]store.Asteroid, bool) {
	asteroids, err := s.store.AsteroidsApproaching(r.Context(), from, to)
	if err != nil {
		log.Printf("%s: %v", r.URL.Path, err)
		writeError(w, http.StatusInternalServerError, "could not load asteroids")
		return nil, false
	}
	return asteroids, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jkusuda/orbital/backend/internal/nasa"
)

// Server holds the dependencies
type Server struct {
	asteroids []nasa.Asteroid
}

func NewServer(asteroids []nasa.Asteroid) *Server {
	return &Server{asteroids: asteroids}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("GET /positions", s.handlePositions)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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

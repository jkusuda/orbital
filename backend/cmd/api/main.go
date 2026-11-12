package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jkusuda/orbital/backend/internal/api"
	"github.com/jkusuda/orbital/backend/internal/nasa"
)

func main() {
	// 1. config
	key := os.Getenv("NASA_API_KEY")
	if key == "" {
		log.Fatal("NASA_API_KEY is not set")
	}

	// 2. dependencies
	client := nasa.NewClient(key)
	asteroids, err := loadAsteroids(client)
	if err != nil {
		log.Fatal(err)
	}
	srv := api.NewServer(asteroids)

	// 3. start
	log.Printf("loaded %d asteroids, listening on :8080", len(asteroids))
	log.Fatal(http.ListenAndServe(":8080", srv.Routes()))
}

// loadAsteroids fetches this week's feed
func loadAsteroids(client *nasa.Client) ([]nasa.Asteroid, error) {
	start := time.Now().Format("2006-01-02")
	feed, err := client.GetFeed(start, start)
	if err != nil {
		return nil, err
	}

	var out []nasa.Asteroid
	for _, a := range feed {
		neo, err := client.GetNEO(a.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, *neo)
	}
	return out, nil
}

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jkusuda/orbital/backend/internal/nasa"
	"github.com/jkusuda/orbital/backend/internal/orbit"
)

func main() {
	nasaKey := os.Getenv("NASA_API_KEY")
	if nasaKey == "" {
		log.Fatal("NASA_API_KEY is not set")
	}

	client := nasa.NewClient(nasaKey)
	fmt.Println("client ready")
	asteroid, err := client.GetNEO("2099942") // Apophis
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", asteroid)
	t := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)
	p := orbit.Position(toElements(asteroid.OrbitalData), t)
	fmt.Printf("%.6f %.6f %.6f\n", p.X, p.Y, p.Z)

	asteroids, err := client.GetFeed("2026-11-11", "2026-11-11")
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range asteroids {
		neo, err := client.GetNEO(a.ID)
		if err != nil {
			log.Fatal(err)
		}
		t := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)
		p := orbit.Position(toElements(neo.OrbitalData), t)
		fmt.Printf("%.6f %.6f %.6f\n", p.X, p.Y, p.Z)
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

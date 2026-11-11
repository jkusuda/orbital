package main

import (
	"log"
	"os"
)

func main() {
	var nasa_key = os.Getenv("NASA_API_KEY")
	if nasa_key == "" {
		log.Fatal("NASA_API_KEY is not set")
	}
}

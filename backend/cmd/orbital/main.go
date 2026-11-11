package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jkusuda/orbital/backend/internal/nasa"
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
}

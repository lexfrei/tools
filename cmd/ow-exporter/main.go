package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("ow-exporter - Overwatch 2 Statistics Exporter")
	fmt.Println("Status: In Development")
	fmt.Println("")
	fmt.Println("See README.md for planned features and API documentation.")
	fmt.Println("Track progress: https://github.com/lexfrei/tools/issues/439")

	log.Println("ow-exporter starting...")

	// TODO: Implement full application
	// - REST API server
	// - SQLite user storage
	// - Background profile scraper
	// - Prometheus metrics exporter

	os.Exit(0)
}
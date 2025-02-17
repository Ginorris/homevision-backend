// TODO use containers to run the application
// TODO automate test run with github actions
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Ginorris/homevision-backend/internal/api"
	"github.com/Ginorris/homevision-backend/internal/config"
	"github.com/Ginorris/homevision-backend/internal/downloader"
	"github.com/Ginorris/homevision-backend/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	// TODO: add actual logging of hardware performance
	start := time.Now()
	// Load environment variables from .env file.
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: no .env file found")
	}

	// Load configuration from file and environment variables.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	fmt.Print(cfg)
	
	// TODO fecth pages concurrently to handle higher workloads
	// Fetch pages sequentially.
	fetchStart := time.Now()
	var houses []models.House
	for i := 1; i <= cfg.PageCount; i++ {
		page, err := api.FetchPage(i, cfg.MaxRetries, cfg.RetryDelay, cfg.APIURL)
		if err != nil {
			// If there's an error (e.g., JSON parsing), log it and continue.
			log.Printf("Error fetching page %d: %v", i, err)
			continue
		}
		// Log number of houses for this page.
		fmt.Printf("Fetched %d houses from page %d\n", len(page.Houses), i)
		houses = append(houses, page.Houses...)
	}
	fmt.Printf("Fetching took %s\n", time.Since(fetchStart))

	// Download photos concurrently.
	downloadStart := time.Now()
	if err := downloader.DownloadPhotos(houses, cfg.Concurrency, cfg.MediaDir); err != nil {
		log.Fatalf("Error downloading photos: %v", err)
	}
	fmt.Printf("Photo downloads took %s\n", time.Since(downloadStart))

	fmt.Printf("Total processing time: %s\n", time.Since(start))
}

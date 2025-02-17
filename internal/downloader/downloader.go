package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/Ginorris/homevision-backend/internal/models"
	"github.com/Ginorris/homevision-backend/internal/utils"
)

// TODO implement retry logic for photo dowload in case of network error
// DownloadPhoto downloads a single photo and saves it using the specified naming format.
func DownloadPhoto(house models.House, mediaDir string) error {
	// Ensure the media directory exists
	if err := os.MkdirAll(mediaDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create media directory: %w", err)
	}

	// Determine file name
	ext := filepath.Ext(house.PhotoURL)
	filename := fmt.Sprintf("%s/%d-%s%s", mediaDir, house.ID, utils.SanitizeFilename(house.Address), ext)

	// Skip downloading if file already exists
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("Skipping download, file already exists: %s\n", filename)
		return nil
	}

	// Perform the request
	resp, err := http.Get(house.PhotoURL)
	if err != nil {
		return fmt.Errorf("error downloading photo for house %d: %w", house.ID, err)
	}
	defer resp.Body.Close()

	// Create the output file
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filename, err)
	}
	defer out.Close()

	// Write the photo to disk
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving photo for house %d: %w", house.ID, err)
	}

	fmt.Printf("Downloaded: %s\n", filename)
	return nil
}

// DownloadPhotos concurrently downloads photos from a list of houses using a worker pool.
func DownloadPhotos(houses []models.House, concurrency int, mediaDir string) error {
	// Create channels for tasks and errors.
	tasks := make(chan models.House, len(houses)) // Buffered channel prevents blocking
	errCh := make(chan error, len(houses))

	// Use a WaitGroup to ensure all workers finish
	var wg sync.WaitGroup

	// Launch worker goroutines
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// we dont use _, house because its a channel
			for house := range tasks {
				err := DownloadPhoto(house, mediaDir)
				if err != nil {
					fmt.Printf("Worker %d - Error downloading: %v\n", workerID, err)
					errCh <- err
				}
			}
		}(i)
	}

	// Enqueue tasks
	for _, house := range houses {
		tasks <- house
	}
	close(tasks)

	// Wait for all workers to finish
	wg.Wait()
	close(errCh)

	// Aggregate Errors & Return First Error
	var firstErr error
	for err := range errCh {
		if firstErr == nil {
			firstErr = err
		}
		fmt.Printf("Download Error: %v\n", err)
	}
	return firstErr
}

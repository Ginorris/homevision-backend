package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Ginorris/homevision-backend/internal/models"
)

// TODO: Use exponential backoff for retryDelay.
// TODO: should we add caching for the pages?
// FetchPage retrieves house data for the specified page.
// If the page fails to load after maxRetries, it logs the error and returns an empty PageResponse.
func FetchPage(page int, maxRetries int, retryDelay time.Duration, baseURL string) (*models.PageResponse, error) {
	fmt.Printf("Fetching page %d\n", page)
	url := fmt.Sprintf("%s?page=%d&per_page=10", baseURL, page)

	var resp *http.Response
	var err error

	// Retry loop
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("Fetching page %d, attempt %d\n", page, i+1)
		resp, err = http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			break // Successful response
		}

		if err == nil {
			err = fmt.Errorf("received non-200 status: %d", resp.StatusCode)
			resp.Body.Close()
		}
		time.Sleep(retryDelay)
	}

	// If there is an error after retries, log it and return an empty PageResponse.
	if err != nil {
		fmt.Printf("Failed to fetch page %d after %d attempts: %v\n", page, maxRetries, err)
		return &models.PageResponse{
			Houses: []models.House{},
			OK:     false,
		}, nil
	}

	// Ensure response body is closed after reading.
	defer resp.Body.Close()

	// Read and parse the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pageResp models.PageResponse
	if err = json.Unmarshal(body, &pageResp); err != nil {
		return nil, err
	}
	return &pageResp, nil
}

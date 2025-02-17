package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Ginorris/homevision-backend/internal/downloader"
	"github.com/Ginorris/homevision-backend/internal/models"
	"github.com/Ginorris/homevision-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

// mockImageServer returns a server that always responds with "mock_image_data".
func mockImageServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some network delay
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "mock_image_data")
	}))
}

func TestDownloadPhoto_Success(t *testing.T) {
	server := mockImageServer()
	defer server.Close()

	// Ensure directory exists
	mediaDir := "media"
	os.MkdirAll(mediaDir, os.ModePerm)

	house := models.House{
		ID:       1,
		Address:  "123 Main St",
		PhotoURL: server.URL,
	}

	err := downloader.DownloadPhoto(house, mediaDir)
	assert.NoError(t, err, "Expected no error for a valid image download")
	ext := filepath.Ext(server.URL)
	expectedFilename := fmt.Sprintf("%s/%d-%s%s", mediaDir, house.ID, utils.SanitizeFilename(house.Address), ext)

	// Check if file exists in the media directory
	_, err = os.Stat(expectedFilename)
	assert.NoError(t, err, "File should be created")

	// Clean up test files
	os.Remove(expectedFilename)
	os.RemoveAll(mediaDir)
}

func TestDownloadPhotos_Concurrency(t *testing.T) {
	server := mockImageServer()
	defer server.Close()

	// Ensure directory exists
	mediaDir := "media"
	os.MkdirAll(mediaDir, os.ModePerm)

	houses := []models.House{
		{ID: 1, Address: "House One", PhotoURL: server.URL},
		{ID: 2, Address: "House Two", PhotoURL: server.URL},
	}

	concurrency := 2
	err := downloader.DownloadPhotos(houses, concurrency, mediaDir)
	assert.NoError(t, err, "Expected no errors when downloading multiple photos concurrently")

	// Check if files exist in media directory
	for _, house := range houses {
		ext := filepath.Ext(server.URL)
		expectedFilename := fmt.Sprintf("%s/%d-%s%s", mediaDir, house.ID, utils.SanitizeFilename(house.Address), ext)
		_, err := os.Stat(expectedFilename)
		assert.NoError(t, err, fmt.Sprintf("File should be created for house %s", house.Address))
		os.Remove(expectedFilename) // Clean up individual files
	}

	// Clean up media directory
	os.RemoveAll(mediaDir)
}

func TestDownloadPhoto_Failure(t *testing.T) {
	// Define media directory
	mediaDir := "media"
	os.MkdirAll(mediaDir, os.ModePerm) // Ensure directory exists

	// Use a domain reserved for invalid use so it won't resolve
	house := models.House{
		ID:       1,
		Address:  "Bad House",
		PhotoURL: "http://nonexistent.invalid/image.jpg",
	}

	err := downloader.DownloadPhoto(house, mediaDir)
	assert.Error(t, err, "Expected error for an invalid URL")

	// Clean up media directory
	os.RemoveAll(mediaDir)
}

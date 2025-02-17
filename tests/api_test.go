package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ginorris/homevision-backend/internal/api"
	"github.com/stretchr/testify/assert"
)

// Mock server that returns a sample JSON response for testing.
func mockServer(status int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintln(w, response)
	}))
}

func TestFetchPage_Success(t *testing.T) {
	// Mock a valid API response.
	server := mockServer(http.StatusOK, `{"houses": [{"id": 1, "address": "123 Main St", "photoURL": "http://example.com/image.jpg"}], "ok": true}`)
	defer server.Close()

	maxRetries := 6
	retryDelay := 500 * time.Millisecond
	pageResponse, err := api.FetchPage(1, maxRetries, retryDelay, server.URL)

	assert.NoError(t, err, "Expected no error for a successful API call")
	assert.NotNil(t, pageResponse, "Page response should not be nil")
	assert.Len(t, pageResponse.Houses, 1, "Expected exactly 1 house")
}

func TestFetchPage_Non200Response(t *testing.T) {
	// Mock an API that returns 500 Internal Server Error.
	server := mockServer(http.StatusInternalServerError, `Internal Server Error`)
	defer server.Close()

	maxRetries := 6
	retryDelay := 500 * time.Millisecond
	pageResponse, err := api.FetchPage(1, maxRetries, retryDelay, server.URL)

	assert.NoError(t, err, "Error should be logged but function should return an empty response")
	assert.Empty(t, pageResponse.Houses, "Should return an empty house list on failure")
}

func TestFetchPage_NetworkFailure(t *testing.T) {
	maxRetries := 3
	retryDelay := 100 * time.Millisecond

	// Close server immediately to simulate a network error.
	server := mockServer(http.StatusOK, `{}`)
	server.Close()

	pageResponse, err := api.FetchPage(1, maxRetries, retryDelay, server.URL)

	assert.NoError(t, err, "Error should be logged but function should return an empty response")
	assert.Empty(t, pageResponse.Houses, "Should return an empty house list on network failure")
}

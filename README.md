# HomeVision Backend Takehome

This repository contains a Go-based solution for the HomeVision photo downloader challenge. The project fetches house data from an unstable, paginated API, downloads associated photos concurrently, and saves them locally with a specific naming format.

## Features

- **Paginated API Fetching:** Retrieves the first 10 pages of house listings.
- **Robust Error Handling:** Implements retries to handle unstable API responses.
- **Concurrent Photo Downloads:** Uses Go's goroutines and channels for efficient concurrent processing.
- **Configuration Management:** Loads settings from a JSON configuration file and environment variables (using [godotenv](https://github.com/joho/godotenv) and [viper](https://github.com/spf13/viper)).
- **Scalable Structure:** Organized project layout for maintainability and future enhancements.
- **Unit Testing:** Includes basic unit tests for critical modules.

## Getting Started

### Prerequisites

- **Go:** Version 1.16 or later is required.
- **Git:** For cloning the repository.
- **WSL/Linux/MacOS:** Recommended for a consistent development environment.

### Setup

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/yourusername/homevision-backend-takehome.git
   cd homevision-backend-takehome
   ```

2. **Set Up Environment Variables:**
   - Copy the example file to create your own `.env` and set API_URL:
    ```bash
    cp .env.example .env
    ```

3. **Review and Edit Configuration:**
   - Open `config/config.json` to verify default settings (e.g., `maxRetries`, `retryDelaySeconds`, `pageCount`, `concurrency`, `mediaDir`).
   - Adjust these values as needed for your environment.

## Running the Application

To run the application:
```bash
go run ./cmd/main.go
```
This command will:
- Load configuration from `config/config.json` and environment variables.
- Fetch the specified number of pages of house data.
- Concurrently download house photos into the `media/` folder.
- Log timing information for both the API fetching and photo downloading phases.

## Running Tests

To run all tests:
```bash
go test ./...
```


## Project Structure
```bash
homevision-backend-takehome/
├── cmd/
│   └── main.go             # Main workflow
├── config/
│   └── config.json
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── tests/
|   ├── downloader_test.go
|   └── api_test.go
└── internal/
    ├── api/
    │   └── api.go          # API fetching logic with retry mechanism
    ├── config/
    │   └── app_config.go   # Configuration loader using Viper
    ├── downloader/
    │   └── downloader.go   # Concurrent photo download logic
    ├── models/
    │   └── types.go        # Data models (House, PageResponse)
    └── utils/
        └── utils.go
```
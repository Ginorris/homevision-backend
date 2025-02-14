# HomeVision Backend Takehome

This repository contains a Go-based solution for the HomeVision photo downloader challenge. The project fetches house data from an unstable, paginated API, downloads associated photos concurrently, and saves them locally with a specific naming format.

## Features

- **Paginated API Fetching:** Retrieves the first 10 pages of house listings.
- **Robust Error Handling:** Implements retries to handle unstable API responses.
- **Concurrent Photo Downloads:** Uses Go's goroutines and channels for efficient concurrent processing.
- **Scalable Structure:** Organized project layout for maintainability and future enhancements.
- **Unit Testing:** Includes basic unit tests for critical modules.

## Getting Started

### Prerequisites

- **Go:** Ensure you have Go installed in your WSL environment. Verify your installation with:
    ```bash
    go version
    ```

### Instalation

1. Clone the Repository:
    ```bash
    git clone git@github.com:Ginorris/homevision-backend.git
    cd homevision-backend
    ```
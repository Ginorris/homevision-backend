package models

// House represents a single house entry from the API response.
type House struct {
	ID        int    `json:"id"`
	Address   string `json:"address"`
	Homeowner string `json:"homeowner"`
	Price     int    `json:"price"`
	PhotoURL  string `json:"photoURL"`
}

// PageResponse represents the structure of the API response for a single page.
type PageResponse struct {
	Houses []House `json:"houses"`
	OK     bool    `json:"ok"`
}

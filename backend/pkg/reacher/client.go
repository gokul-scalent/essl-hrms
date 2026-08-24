package reacher

import (
	"net/http"
	"time"
)

// Client is used to communicate with the Reacher email verification service.
type Client struct {
	httpClient *http.Client // HTTP client used to send API requests
	baseURL    string       // Base URL of the Reacher API (http://localhost:8080)
}

// NewClient creates and returns a new Reacher client.
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			// Automatically cancels the request if it takes longer than the timeout.
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

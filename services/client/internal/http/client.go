package http

import (
	"net/http"
	"time"
)

// NewClient creates and returns a new HTTP client.
func NewClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

package utils

import (
	"net/http"
	"time"
)

// HTTPClient provides a shared HTTP client with connection pooling
// for better performance across all services
var HTTPClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	},
}

// ShortTimeoutClient for quick validation requests
var ShortTimeoutClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     60 * time.Second,
	},
}

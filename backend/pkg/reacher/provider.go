package reacher

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// NewReacherClient creates a Reacher client using the application configuration.
func NewReacherClient(cfg ReacherConfig) (*Client, error) {
	// Convert the timeout string (e.g. "10s", "30s") into a time.Duration.
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return nil, err
	}

	// if err := validateReacherBaseURL(cfg.BaseURL, timeout); err != nil {
	// 	return nil, err
	// }

	// Create and return a new Reacher client with the configured base URL and timeout.
	return NewClient(cfg.BaseURL, timeout), nil
}

func validateReacherBaseURL(baseURL string, timeout time.Duration) error {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid Reacher baseURL %q: %w", baseURL, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("invalid Reacher baseURL %q", baseURL)
	}

	host := parsed.Host
	if !strings.Contains(host, ":") {
		switch parsed.Scheme {
		case "http":
			host += ":80"
		case "https":
			host += ":443"
		default:
			return fmt.Errorf("unsupported scheme for Reacher baseURL %q", baseURL)
		}
	}

	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return fmt.Errorf("could not connect to Reacher baseURL %q: %w", baseURL, err)
	}
	_ = conn.Close()
	return nil
}

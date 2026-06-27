package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultGoAPIPort        = "8080"
	defaultNodeAPIURL       = "http://localhost:3000"
	defaultHTTPTimeoutMS    = 3000
	defaultAlgorithm        = "modified-gram-schmidt"
	defaultNumericTolerance = 1e-9
)

type Config struct {
	Port              string
	NodeAPIURL        string
	HTTPClientTimeout time.Duration
	Algorithm         string
	Tolerance         float64
}

func Load() Config {
	timeoutMS := defaultHTTPTimeoutMS
	if rawTimeout := os.Getenv("HTTP_CLIENT_TIMEOUT_MS"); rawTimeout != "" {
		if parsed, err := strconv.Atoi(rawTimeout); err == nil && parsed > 0 {
			timeoutMS = parsed
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("GO_API_PORT")
	}
	if port == "" {
		port = defaultGoAPIPort
	}

	nodeAPIURL := os.Getenv("NODE_API_URL")
	if nodeAPIURL == "" {
		nodeAPIURL = defaultNodeAPIURL
	}

	return Config{
		Port:              port,
		NodeAPIURL:        nodeAPIURL,
		HTTPClientTimeout: time.Duration(timeoutMS) * time.Millisecond,
		Algorithm:         defaultAlgorithm,
		Tolerance:         defaultNumericTolerance,
	}
}

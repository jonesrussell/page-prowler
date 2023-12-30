package crawler

import (
	"net/url"

	"github.com/jonesrussell/page-prowler/internal/logger"
)

// GetHostFromURL extracts the host from the given URL.
func GetHostFromURL(inputURL string, appLogger logger.Logger) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		appLogger.Error("Failed to parse URL", "url", inputURL, "error", err)
		return "", err
	}

	host := parsedURL.Hostname()

	return host, nil
}

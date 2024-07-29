package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type URLParseError struct {
	URL string
	Err error
}

func (e *URLParseError) Error() string {
	return fmt.Sprintf("Failed to parse URL: %s, Error: %v", e.URL, e.Err)
}

func ExtractLastSegmentFromURL(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		return ""
	}

	if u.Path == "" || u.Path == "/" {
		return ""
	}

	segments := strings.Split(u.Path, "/")
	title := segments[len(segments)-1]

	return title
}

func GetHostFromURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", &URLParseError{URL: inputURL, Err: err}
	}

	host := parsedURL.Hostname()
	if host == "" {
		return "", errors.New("failed to extract host from URL")
	}

	return host, nil
}

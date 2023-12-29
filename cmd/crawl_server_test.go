package cmd

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPostArticlesStart(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"URL": "https://example.com", "SearchTerms": "test", "CrawlSiteID": "1", "MaxDepth": 1, "Debug": false}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/articles/start")

	// Create an instance of CrawlServer
	server := &CrawlServer{
		CrawlManager: &crawler.CrawlManager{
			Logger: logger.New(true), // create a new ZapLoggerWrapper
			Client: &mockRedisClient{},
		},
	}

	// Assertions
	if assert.NoError(t, server.PostArticlesStart(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"message":"Crawling started successfully"}`, strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestGetPing(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/ping")

	// Create an instance of CrawlServer
	server := &CrawlServer{}

	// Assertions
	if assert.NoError(t, server.GetPing(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Pong", rec.Body.String())
	}
}

func TestGetHostFromURL(t *testing.T) {
	// Create a mock Logger
	zapLogger, _ := zap.NewDevelopment()
	logger := &logger.ZapLoggerWrapper{Logger: zapLogger.Sugar()}

	// Define test cases
	testCases := []struct {
		url      string
		expected string
	}{
		{"http://example.com/path", "example.com"},
		{"https://www.example.com", "www.example.com"},
		// add more test cases here
	}

	// Run test cases
	for _, tc := range testCases {
		host, err := getHostFromURL(tc.url, logger)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if host != tc.expected {
			t.Errorf("Expected host %v, but got %v", tc.expected, host)
		}
	}
}

func TestConfigureCollector(t *testing.T) {
	allowedDomains := []string{"test.com"}
	maxDepth := 2

	collector := configureCollector(allowedDomains, maxDepth)

	if collector == nil {
		t.Errorf("Expected collector to be not nil")
		return
	}

	if !reflect.DeepEqual(collector.AllowedDomains, allowedDomains) {
		t.Errorf("Expected allowed domains to be %v, but got %v", allowedDomains, collector.AllowedDomains)
	}

	if collector.MaxDepth != maxDepth {
		t.Errorf("Expected max depth to be %d, but got %d", maxDepth, collector.MaxDepth)
	}
}

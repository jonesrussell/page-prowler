package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jonesrussell/page-prowler/internal/crawler"
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
			Logger: &ZapLoggerWrapper{
				logger: zap.NewExample().Sugar(),
			},
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

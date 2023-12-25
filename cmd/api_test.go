package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock/v8"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupMockRedis() (context.Context, *redis.Client, *redismock.ClientMock) {
	mini, _ := miniredis.Run()
	client := redis.NewClient(&redis.Options{Addr: mini.Addr()})
	mock := redismock.NewNiceMock(client)
	ctx := context.Background()
	return ctx, client, mock
}

func TestGetPing(t *testing.T) {
	e := echo.New()
	ctx, _, mock := setupMockRedis()
	defer mock.FlushAll(ctx)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Initialize your manager here...
	manager, err := initializeManager(ctx, true)
	if err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	handler := &ServerInterfaceWrapper{
		Handler: &CrawlServer{
			CrawlManager: manager,
		},
	}

	if assert.NoError(t, handler.GetPing(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestPostArticlesStart(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://",
		httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("GET", "https://example.ca",
		httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("GET", "https://example.ca/robots.txt",
		httpmock.NewStringResponder(200, ""))

	e := echo.New()
	ctx, _, _ := setupMockRedis()

	// Test with a valid request body.
	req := httptest.NewRequest(http.MethodPost, "/articles/start", strings.NewReader(`{"CrawlSiteID":"site1","Debug":true,"MaxDepth":2,"SearchTerms":"term1","URL":"https://example.ca"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	manager, _ := initializeManager(ctx, true)
	handler := &ServerInterfaceWrapper{
		Handler: &CrawlServer{
			CrawlManager: manager,
		},
	}

	if assert.NoError(t, handler.PostArticlesStart(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"message":"Crawling started successfully"}`, strings.TrimSuffix(rec.Body.String(), "\n"))
	}

	// Test with an invalid request body.
	req = httptest.NewRequest(http.MethodPost, "/articles/start", strings.NewReader(`{"CrawlSiteID":"site1","Debug":true,"MaxDepth":2,"SearchTerms":"term1","URL":""}`)) // Empty URL
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err := handler.PostArticlesStart(c)
	if assert.Error(t, err) {
		assert.Equal(t, "URL cannot be empty", err.Error())
	}
}

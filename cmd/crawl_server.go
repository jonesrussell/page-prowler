package cmd

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/labstack/echo/v4"
)

// CrawlServer represents the server that handles the crawling process.
type CrawlServer struct {
	CrawlManager *crawler.CrawlManager
}

// PostArticlesStart starts the article posting process.
func (s *CrawlServer) PostArticlesStart(ctx echo.Context) error {
	var req PostArticlesStartJSONBody
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	// Ensure the URL is not empty
	if *req.URL == "" {
		return errors.New("URL cannot be empty")
	}

	// Initialize your crawler service here
	crawlerService, err := initializeManager(ctx.Request().Context(), *req.Debug)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Ensure the URL is correctly formatted
	url := strings.TrimSpace(*req.URL)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	err = StartCrawling(ctx.Request().Context(), url, *req.SearchTerms, *req.CrawlSiteID, *req.MaxDepth, *req.Debug, crawlerService, s)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Crawling started successfully"})
}

// GetPing handles the ping request.
func (s *CrawlServer) GetPing(ctx echo.Context) error {
	// Implement your logic here
	return ctx.String(http.StatusOK, "Pong")
}

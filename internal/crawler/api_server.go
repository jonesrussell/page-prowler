package crawler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// PostMatchlinksStart starts the article posting process.
func PostMatchlinksStart(ctx echo.Context) error {
	// Get the CrawlManager from the context
	manager := ctx.Get(string(echoManagerKey)).(*CrawlManager)
	if manager == nil {
		log.Fatalf("CrawlManager is not initialized")
	}

	var req PostMatchlinksJSONBody
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	// Ensure the URL is not empty
	if *req.URL == "" {
		return errors.New("URL cannot be empty")
	}

	// Ensure the URL is correctly formatted
	url := strings.TrimSpace(*req.URL)
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	err := manager.StartCrawling(
		ctx.Request().Context(),
		url,
		*req.SearchTerms,
		*req.CrawlSiteID,
		*req.MaxDepth,
		*req.Debug,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Crawling started successfully"})
}

// GetPing handles the ping request.
func GetPing(ctx echo.Context) error {
	// Implement your logic here
	return ctx.String(http.StatusOK, "Pong")
}

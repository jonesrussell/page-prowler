package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Echo server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()

		AttachRoutes(e)

		e.Start(":3000")
	},
}

func AttachRoutes(e *echo.Echo) {
	e.GET("/ping", Ping)

	crawlerService, err := initializeManager(context.Background(), false)
	if err != nil {
		log.Fatalf("Failed to initialize Crawl Manager: %v", err)
	}

	e.POST("/start-crawling", StartCrawling(crawlerService))
}

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "pong",
	})
}

func StartCrawling(crawlerService *crawler.CrawlManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data struct {
			URL         string
			SearchTerms string
			CrawlSiteID string
			MaxDepth    int
			Debug       bool
		}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Cannot parse JSON",
			})
		}

		ctx := context.Background()
		crawlerService.Logger.Info("Crawler started...")

		err := crawlerService.RedisClient.Ping(ctx) // Use RedisClient from CrawlManager
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to ping Redis",
			})
		}

		// Continue with your logic here

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Crawling started",
		})
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

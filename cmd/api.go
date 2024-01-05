package cmd

import (
	"log"

	"github.com/jonesrussell/page-prowler/internal/api"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

type echoContextKey string

const (
	echoManagerKey echoContextKey = "manager"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Long: `The 'api' command starts the API server. This server handles various requests related to matchlinks and crawling jobs.
		  For example, you can start a new article matching job with the '/matchlinks' endpoint, or retrieve the status of a job with the '/matchlinks/info/{id}' endpoint.
		  Similarly, you can start a new crawling job with the '/crawling/start' endpoint, or retrieve the status of a crawling job with the '/crawling/info/{id}' endpoint.
		  The server also includes a '/ping' endpoint for health checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()

		// Get the manager from the context
		manager, ok := cmd.Context().Value(managerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			log.Fatalf("CrawlManager is not initialized")
		}

		// Add the middleware to the Echo instance
		e.Use(CrawlManagerMiddleware(manager))

		apiServerInterface := &api.ApiServerInterface{}

		// Create a group with /v1 prefix
		v1 := e.Group("/v1")

		// Register handlers under /v1
		api.RegisterHandlers(v1, apiServerInterface)

		if err := e.Start(":3000"); err != nil {
			log.Fatalf("Error starting echo server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func CrawlManagerMiddleware(manager *crawler.CrawlManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set the CrawlManager in the context
			c.Set(string(echoManagerKey), manager)
			return next(c)
		}
	}
}

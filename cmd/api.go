package cmd

import (
	"log"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Long: `The 'api' command starts the API server. This server handles various requests related to articles and crawling jobs.
    For example, you can start a new article matching job with the '/articles/start' endpoint, or retrieve the status of a job with the '/articles/info/{id}' endpoint.
    Similarly, you can start a new crawling job with the '/crawling/start' endpoint, or retrieve the status of a crawling job with the '/crawling/info/{id}' endpoint.
    The server also includes a '/ping' endpoint for health checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()

		manager := cmd.Context().Value("manager").(*crawler.CrawlManager)
		if manager == nil {
			log.Fatalf("CrawlManager is not initialized")
		}

		// Add the middleware to the Echo instance
		e.Use(CrawlManagerMiddleware(manager))

		server := &ServerInterfaceWrapper{
			Handler: &CrawlServer{
				CrawlManager: manager,
			},
		}

		RegisterHandlers(e, server)

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
			c.Set("manager", manager)
			return next(c)
		}
	}
}

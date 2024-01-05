package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/tasks"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		// Retrieve the Redis connection details
		redisHost := viper.GetString("REDIS_HOST")
		redisPort := viper.GetString("REDIS_PORT")
		redisAuth := viper.GetString("REDIS_AUTH")

		// Add the middleware to the Echo instance
		e.Use(CrawlManagerMiddleware(manager))

		// Register handlers
		e.GET("/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "Pong")
		})

		e.POST("/matchlinks", func(c echo.Context) error {
			var req crawler.PostMatchlinksJSONBody
			if err := c.Bind(&req); err != nil {
				return err
			}

			// Validate the input parameters
			if req.URL == nil || *req.URL == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL cannot be empty"})
			}
			if req.SearchTerms == nil || *req.SearchTerms == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "SearchTerms cannot be empty"})
			}
			if req.CrawlSiteID == nil || *req.CrawlSiteID == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "CrawlSiteID cannot be empty"})
			}
			if req.MaxDepth == nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "MaxDepth cannot be null"})
			}

			// Default Debug to false if it is nil
			if req.Debug == nil {
				req.Debug = new(bool)
				*req.Debug = false
			}

			// Ensure the URL is correctly formatted
			url := strings.TrimSpace(*req.URL)
			if !strings.HasPrefix(url, "https://") {
				url = "https://" + url
			}

			// Create a new asynq.Client using the same Redis connection details
			client := asynq.NewClient(asynq.RedisClientOpt{
				Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
				Password: redisAuth,
			})

			payload := &tasks.CrawlTaskPayload{
				URL:         url,
				SearchTerms: *req.SearchTerms,
				CrawlSiteID: *req.CrawlSiteID,
				MaxDepth:    *req.MaxDepth,
				Debug:       *req.Debug,
			}

			tid, err := tasks.EnqueueCrawlTask(client, payload)
			if err != nil {
				log.Println("Error enqueuing crawl task: ", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			return c.JSON(http.StatusOK, map[string]string{"message": "Crawling started successfully", "task_id": tid})
		})

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

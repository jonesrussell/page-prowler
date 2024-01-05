package cmd

import (
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/api"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

		// Enable CORS
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		}))

		// Get the manager from the context
		manager, ok := cmd.Context().Value(managerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			log.Fatalf("CrawlManager is not initialized")
		}

		// Add the middleware to the Echo instance
		e.Use(CrawlManagerMiddleware(manager))

		redisHost := viper.GetString("REDIS_HOST")
		redisPort := viper.GetString("REDIS_PORT")
		redisAuth := viper.GetString("REDIS_AUTH")

		// Initialize the Inspector
		inspector := asynq.NewInspector(asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: redisAuth,
		})

		apiServerInterface := &api.ApiServerInterface{
			Inspector: inspector,
		}

		// Create a group with /v1 prefix
		v1 := e.Group("/v1")

		// Register handlers under /v1
		api.RegisterHandlers(v1, apiServerInterface)

		if err := e.StartTLS(":3000", "/ssl/cert.pem", "/ssl/key_unencrypted.pem"); err != nil {
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

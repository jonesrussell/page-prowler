package cmd

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/api"
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	CertPathEnvKey string = "SSL_CERT_PATH"
	KeyPathEnvKey  string = "SSL_KEY_PATH"
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
		manager, ok := cmd.Context().Value(common.CrawlManagerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			log.Fatalf("CrawlManager is not initialized")
		}

		// Add the middleware to the Echo instance
		e.Use(CrawlManagerMiddleware(manager))

		redisDetails := manager.Client.Options()
		redisAddr := redisDetails.Addr
		redisAuth := redisDetails.Password

		// Initialize the Inspector
		inspector := asynq.NewInspector(asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisAuth,
		})

		apiServerInterface := &api.ServerApiInterface{
			Inspector: inspector,
		}

		// Create a group with /v1 prefix
		v1 := e.Group("/v1")

		// Register handlers under /v1
		api.RegisterHandlers(v1, apiServerInterface)

		if err := e.StartTLS(":3000", viper.GetString(CertPathEnvKey), viper.GetString(KeyPathEnvKey)); err != nil {
			log.Fatalf("Error starting echo server: %v", err)
		}
	},
}

func CrawlManagerMiddleware(manager *crawler.CrawlManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set the CrawlManager in the context using the string constant as the key
			c.Set(string(common.CrawlManagerKey), manager)
			return next(c)
		}
	}
}

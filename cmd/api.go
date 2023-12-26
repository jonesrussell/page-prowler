package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		ctx := context.Background()
		debug := viper.GetBool("debug")
		manager, err := initializeManager(ctx, debug, nil)
		if err != nil {
			fmt.Println("Failed to initialize manager", err)
			return
		}

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

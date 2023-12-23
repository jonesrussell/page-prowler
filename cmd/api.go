package cmd

import (
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

		// Initialize your server
		server := ServerInterfaceWrapper{} // assuming you have a ServerInterfaceWrapper struct in api.gen.go

		// Register routes
		RegisterHandlers(e, &server) // assuming you have a RegisterHandlers function in api.gen.go

		// Start the server
		e.Start(":3000")
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

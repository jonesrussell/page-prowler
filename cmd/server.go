package cmd

import (
	"context"
	"net/http"

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

		e.GET("/ping", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"message": "pong",
			})
		})

		e.POST("/start-crawling", func(c echo.Context) error {
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
			crawlerService, err := initializeManager(ctx, data.Debug)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to initialize Crawl Manager",
				})
			}

			err = startCrawling(ctx, data.URL, data.SearchTerms, data.CrawlSiteID, data.MaxDepth, data.Debug, crawlerService)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to start crawling",
				})
			}

			return c.JSON(http.StatusOK, map[string]string{
				"message": "Crawling started",
			})
		})

		e.Start(":3000")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

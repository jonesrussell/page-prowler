package cmd

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Fiber server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := fiber.New()

		app.Get("/ping", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "pong",
			})
		})

		app.Post("/start-crawling", func(c *fiber.Ctx) error {
			var data struct {
				URL         string
				SearchTerms string
				CrawlSiteID string
				MaxDepth    int
				Debug       bool
			}
			if err := c.BodyParser(&data); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Cannot parse JSON",
				})
			}

			ctx := context.Background()
			crawlerService, err := initializeManager(ctx, data.Debug)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to initialize Crawl Manager",
				})
			}

			err = startCrawling(ctx, data.URL, data.SearchTerms, data.CrawlSiteID, data.MaxDepth, data.Debug, crawlerService)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to start crawling",
				})
			}

			return c.JSON(fiber.Map{
				"message": "Crawling started",
			})
		})

		app.Listen(":3000")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

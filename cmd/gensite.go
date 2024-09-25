package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

// NewsPage represents the structure of the news page
type NewsPage struct {
	Title    string
	Articles []Article
}

// Article represents a news article
type Article struct {
	Title string
	URL   string
}

// NewGenSiteCmd creates a new gensite command
func NewGenSiteCmd() *cobra.Command {
	var siteName string

	cmd := &cobra.Command{
		Use:   "gensite",
		Short: "Generate a static news site",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateSite(siteName)
		},
	}

	cmd.Flags().StringVarP(&siteName, "site", "s", "", "Name of the site to generate")
	cmd.MarkFlagRequired("site")

	return cmd
}

// generateSite generates a static HTML file for the specified site
func generateSite(siteName string) error {
	// Define articles based on the site
	var articles []Article
	switch siteName {
	case "cp24":
		articles = []Article{
			{"Site 1 Article 1", "https://example.com/site1/article1"},
			{"Site 1 Article 2", "https://example.com/site1/article2"},
		}
	case "site2":
		articles = []Article{
			{"Site 2 Article 1", "https://example.com/site2/article1"},
			{"Site 2 Article 2", "https://example.com/site2/article2"},
		}
	default:
		return fmt.Errorf("unknown site: %s", siteName)
	}

	// Create the output directory if it doesn't exist
	outputDir := fmt.Sprintf("static/%s", siteName)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create the HTML file
	file, err := os.Create(fmt.Sprintf("%s/index.html", outputDir))
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Execute the template and write to the file
	tmpl, err := template.ParseFiles("static/templates/cp24.gohtml")
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	page := NewsPage{
		Title:    fmt.Sprintf("%s News", siteName),
		Articles: articles,
	}

	if err := tmpl.Execute(file, page); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	fmt.Printf("Static site generated: %s/index.html\n", outputDir)
	return nil
}

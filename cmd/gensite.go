package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/jonesrussell/page-prowler/news"
	"github.com/spf13/cobra"
)

// NewsPage represents the data for the news page template
type NewsPage struct {
	Title         string
	TopStory      news.Article
	BreakingNews  []news.Article
	LatestUpdates []news.Article
}

func NewGenSiteCmd(newsService news.Service) *cobra.Command {
	var siteName string

	cmd := &cobra.Command{
		Use:   "gensite",
		Short: "Generate a static news site",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if siteName == "" {
				return fmt.Errorf("required flag \"site\" not set")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return generateSite(siteName, newsService)
		},
	}

	cmd.Flags().StringVarP(&siteName, "site", "s", "", "Name of the news site (required)")

	return cmd
}

func generateSite(siteName string, newsService news.Service) error {
	topStory, err := newsService.GetTopStory(siteName)
	if err != nil {
		return fmt.Errorf("failed to get top story: %v", err)
	}

	breakingNews, err := newsService.GetBreakingNews(siteName)
	if err != nil {
		return fmt.Errorf("failed to get breaking news: %v", err)
	}

	latestUpdates, err := newsService.GetLatestUpdates(siteName)
	if err != nil {
		return fmt.Errorf("failed to get latest updates: %v", err)
	}

	// Create the output directory with the site name
	outputDir := fmt.Sprintf("static/generated/%s", siteName)
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
	tmpl, err := template.ParseFiles(
		"static/templates/cp24.html",
		"static/templates/cp24/header.html",
		"static/templates/cp24/footer.html",
		"static/templates/cp24/top_story.html",
		"static/templates/cp24/top_story_bn.html",
		"static/templates/cp24/latest_updates.html",
	)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	page := NewsPage{
		Title:         fmt.Sprintf("%s News", siteName),
		TopStory:      topStory,
		BreakingNews:  breakingNews,
		LatestUpdates: latestUpdates,
	}

	if err := tmpl.Execute(file, page); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	fmt.Printf("Static site generated: %s/index.html\n", outputDir)
	return nil
}

package cmd

import (
	"fmt"
	"net/http"

	"github.com/jonesrussell/page-prowler/news"
	"github.com/spf13/cobra"
)

func NewServeCmd(newsService news.Service) *cobra.Command {
	var port int
	var siteName string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve a static news site",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if siteName == "" {
				return fmt.Errorf("required flag \"site\" not set")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return serveSite(siteName, port, newsService)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to serve the site on")
	cmd.Flags().StringVarP(&siteName, "site", "s", "", "Name of the site to serve (required)")

	return cmd
}

func serveSite(siteName string, port int, newsService news.Service) error {
	// Generate the site first
	if err := generateSite(siteName, newsService); err != nil {
		return fmt.Errorf("failed to generate site: %v", err)
	}

	// Serve the generated files
	outputDir := fmt.Sprintf("static/generated/%s", siteName)
	fs := http.FileServer(http.Dir(outputDir))
	http.Handle("/", fs)

	fmt.Printf("Serving %s site on http://localhost:%d\n", siteName, port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

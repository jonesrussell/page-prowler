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
	TopStory Article // Added TopStory field
	Articles []Article
}

// Article represents a news article
type Article struct {
	Title       string
	Link        string
	Image       string
	Description string // Ensure this field exists if you're using it
}

// NewGenSiteCmd creates a new gensite command
func NewGenSiteCmd() *cobra.Command {
	var siteName string

	cmd := &cobra.Command{
		Use:   "gensite",
		Short: "Generate a static news site",
		RunE: func(_ *cobra.Command, _ []string) error {
			return generateSite(siteName)
		},
	}

	cmd.Flags().StringVarP(&siteName, "site", "s", "", "Name of the site to generate")
	if err := cmd.MarkFlagRequired("site"); err != nil {
		fmt.Println("Error marking flag as required:", err)
		return nil
	}

	return cmd
}

// generateSite generates a static HTML file for the specified site
func generateSite(siteName string) error {
	// Define articles based on the site
	var topStory Article // Declare topStory variable
	var articles []Article
	switch siteName {
	case "cp24":
		topStory = Article{
			Title:       "Ontario Premier Doug Ford says he wants to build a tunnel under Hwy. 401",
			Link:        "https://www.cp24.com/news/ontario-premier-doug-ford-says-he-wants-to-build-a-tunnel-under-hwy-401-1.7051216",
			Image:       "https://www.cp24.com/polopoly_fs/1.7051238.1727270667!/httpImage/image.jpg_gen/derivatives/landscape_620/image.jpg",
			Description: "Premier Doug Ford says he wants to build a tunnel under Highway 401 that would stretch from Brampton to Scarborough.",
		} // Set top story
		articles = []Article{
			{
				Title:       "Tearful complainant alleges Jacob Hoggard raped, choked her after Hedley concert",
				Link:        "https://www.cp24.com/news/tearful-complainant-alleges-jacob-hoggard-raped-choked-her-after-hedley-concert-1.7051025",
				Image:       "https://www.cp24.com/polopoly_fs/1.7050029.1727194241!/image/image.jpeg_gen/derivatives/landscape_300/image.jpeg",
				Description: "Opening arguments are expected to get underway today in the sexual assault trial of Canadian musician Jacob Hoggard.",
			},
			{
				Title:       "Toronto teachers’ union accuses Ford of diverting attention away from Grassy Narrows",
				Link:        "https://www.cp24.com/news/toronto-teachers-union-accuses-ford-of-diverting-attention-away-from-grassy-narrows-as-province-begins-investigating-controversial-field-trip-1.7051645",
				Image:       "https://www.cp24.com/polopoly_fs/1.7051296.1727272549!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "Fallout over controversial field trip in Toronto.",
			},
			{
				Title:       "Thieves stole more than $2.2 million of merchandise from moving tractor trailers",
				Link:        "https://www.cp24.com/news/thieves-stole-more-than-2-2-million-of-merchandise-from-moving-tractor-trailers-police-1.7051521",
				Image:       "https://www.cp24.com/polopoly_fs/1.5752494.1648298841!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "A Peel Regional Police cruiser is seen in this undated image.",
			},
			{
				Title:       "'This a bright red warning light': Toronto's housing crisis to get worse",
				Link:        "https://www.cp24.com/news/this-a-bright-red-warning-light-toronto-s-housing-crisis-to-get-worse-as-development-applications-drop-off-bild-says-1.7051629",
				Image:       "https://www.cp24.com/polopoly_fs/1.7033313.1726047954!/httpImage/image.jpg_gen/derivatives/landscape_300/image.jpg",
				Description: "A new condo construction site is reflected in the window on an ongoing condo construction site in downtown Toronto.",
			},
		}
	case "site2":
		topStory = Article{
			Title:       "Site 2 Top Story",
			Link:        "https://example.com/site2/topstory",
			Image:       "",
			Description: "This is a dummy description for Site 2 Top Story.",
		} // Set top story for site2
		articles = []Article{
			{
				Title:       "Site 2 Article 1",
				Link:        "https://example.com/site2/article1",
				Image:       "",
				Description: "This is a dummy description for Site 2 Article 1.",
			},
			{
				Title:       "Site 2 Article 2",
				Link:        "https://example.com/site2/article2",
				Image:       "",
				Description: "This is a dummy description for Site 2 Article 2.",
			},
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
	tmpl, err := template.ParseFiles("static/templates/cp24.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	page := NewsPage{
		Title:    fmt.Sprintf("%s News", siteName),
		TopStory: topStory, // Include top story in the page
		Articles: articles, // Include articles in the page
	}

	if err := tmpl.Execute(file, page); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	fmt.Printf("Static site generated: %s/index.html\n", outputDir)
	return nil
}

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jonesrussell/page-prowler/redis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getLinksCmd = &cobra.Command{
	Use:   "getlinks",
	Short: "Get the list of links for a given crawlsiteid",
	Run: func(cmd *cobra.Command, args []string) {
		crawlsiteid := viper.GetString("crawlsiteid")
		if crawlsiteid == "" {
			fmt.Println("crawlsiteid is required")
			os.Exit(1)
		}

		ctx := context.Background()
		redisClient, err := redis.NewClient(viper.GetString("REDIS_HOST"), viper.GetString("REDIS_AUTH"), viper.GetString("REDIS_PORT"))
		if err != nil {
			fmt.Println("Failed to initialize Redis client", "error", err)
			os.Exit(1)
		}

		links, err := redisClient.SMembers(ctx, crawlsiteid)
		if err != nil {
			fmt.Println("Failed to get links from Redis", "error", err)
			os.Exit(1)
		}

		for _, link := range links {
			fmt.Println(link)
		}
	},
}

func init() {
	getLinksCmd.Flags().String("crawlsiteid", "", "CrawlSite ID")
	if err := viper.BindPFlag("crawlsiteid", getLinksCmd.Flags().Lookup("crawlsiteid")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
	rootCmd.AddCommand(getLinksCmd)
}

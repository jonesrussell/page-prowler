package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/rediswrapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume URLs from Redis",
	Long:  `Consume is a CLI tool designed to fetch URLs from a Redis set.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		crawlSiteID, _ := cmd.Flags().GetString("crawlsiteid")
		debug, _ := cmd.Flags().GetBool("debug")

		if crawlSiteID == "" {
			fmt.Println("CrawlSiteId is required")
			os.Exit(1)
		}

		fmt.Println("CrawlSiteId:", crawlSiteID)
		fmt.Println("Debug:", debug)
		startConsuming(ctx, crawlSiteID, debug)
	},
}

func init() {
	rootCmd.AddCommand(consumeCmd)

	consumeCmd.PersistentFlags().String("crawlsiteid", "", "CrawlSite ID")
	consumeCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	viper.BindPFlag("crawlsiteid", consumeCmd.PersistentFlags().Lookup("crawlsiteid"))
	viper.BindPFlag("debug", consumeCmd.PersistentFlags().Lookup("debug"))

	fmt.Println("All configuration keys and values:")
	for _, key := range viper.AllKeys() {
		fmt.Printf("%s: %v\n", key, viper.Get(key))
	}
}

func initializeConsumeManager(ctx context.Context, debug bool) (*crawler.CrawlManager, error) {
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	log := logger.New(debug)

	if debug {
		fmt.Println("Redis Host:", redisHost)
		fmt.Println("Redis Port:", redisPort)
		fmt.Println("Redis Auth:", redisAuth)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisAuth,
		DB:       0,
	})

	redisWrapper, err := rediswrapper.NewRedisWrapper(ctx, rdb)
	if err != nil {
		log.Error("Failed to initialize Redis", "error", err)
		return nil, err
	}

	return &crawler.CrawlManager{
		Logger:       log,
		RedisWrapper: redisWrapper,
	}, nil
}

func startConsuming(ctx context.Context, crawlSiteID string, debug bool) {
	crawlerService, err := initializeConsumeManager(ctx, debug)
	if err != nil {
		crawlerService.Logger.Error("Failed to initialize Consume Manager", "error", err)
		os.Exit(1)
	}

	urls, err := crawlerService.RedisWrapper.SMembers(ctx, crawlSiteID)
	if err != nil {
		crawlerService.Logger.Error("Error fetching URLs from Redis", "error", err)
		return
	}

	for _, url := range urls {
		crawlerService.Logger.Info("Fetched URL", "url", url)
	}
}

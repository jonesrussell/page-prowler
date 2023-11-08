package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume URLs from Redis",
	Long:  `Consume is a CLI tool designed to fetch URLs from a Redis set.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		startConsuming(ctx, viper.GetString("crawlsiteid"))
	},
}

func main() {
	// Cobra related code
	cobra.OnInitialize(initConfig)
	consumeCmd.PersistentFlags().String("crawlsiteid", "", "CrawlSite ID")

	consumeCmd.MarkPersistentFlagRequired("crawlsiteid")

	viper.BindPFlag("crawlsiteid", consumeCmd.PersistentFlags().Lookup("crawlsiteid"))

	// Execute the consume command
	if err := consumeCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // Automatically override values from the .env file with those from the environment.

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while reading config file", err)
	}
}

func startConsuming(ctx context.Context, crawlSiteID string) {
	// Fetch Redis configuration from Viper
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	// Create a new Redis client instance.
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisAuth, // no password set
		DB:       0,         // use default DB
	})

	// Initialize RedisWrapper
	redisWrapper, err := rediswrapper.NewRedisWrapper(ctx, rdb)
	if err != nil {
		fmt.Println("Failed to initialize Redis", "error", err)
		os.Exit(1)
	}

	// Fetch URLs from Redis
	urls, err := redisWrapper.SMembers(ctx, crawlSiteID)
	if err != nil {
		fmt.Println("Error fetching URLs from Redis", "error", err)
		return
	}

	// Print the fetched URLs
	for _, url := range urls {
		fmt.Println(url)
	}
}

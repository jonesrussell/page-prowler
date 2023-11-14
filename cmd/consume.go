/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume URLs from Redis",
	Long:  `Consume is a CLI tool designed to fetch URLs from a Redis set.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CrawlSiteId (from cmd):", cmd.Flag("crawlsiteid").Value.String())
		fmt.Println("Debug (from cmd):", cmd.Flag("debug").Value.String())
		fmt.Println("Args:", args)
		ctx := context.Background()
		crawlSiteID := viper.GetString("crawlsiteid")
		debug := viper.GetBool("debug")
		fmt.Println("CrawlSiteId:", crawlSiteID)
		fmt.Println("Debug:", debug) // Print the value of debug
		startConsuming(ctx, crawlSiteID, debug)
	},
}

func init() {
	rootCmd.AddCommand(consumeCmd)

	consumeCmd.PersistentFlags().String("crawlsiteid", "", "CrawlSite ID")
	consumeCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	consumeCmd.MarkPersistentFlagRequired("crawlsiteid")

	viper.BindPFlag("crawlsiteid", consumeCmd.PersistentFlags().Lookup("crawlsiteid"))
	viper.BindPFlag("debug", consumeCmd.PersistentFlags().Lookup("debug"))

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	// viper.AutomaticEnv() // Automatically override values from the .env file with those from the environment.

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while reading config file", err)
	}

	// Debugging: Print all configuration keys and values
	fmt.Println("All configuration keys and values:")
	for _, key := range viper.AllKeys() {
		fmt.Printf("%s: %v\n", key, viper.Get(key))
	}
}

func initializeConsumeManager(ctx context.Context, debug bool) (*crawler.CrawlManager, error) {
	// Fetch Redis configuration from Viper
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	// Initialize Logger with the new logger package
	log := logger.New(debug) // Use the logger package's New function to get a logger instance

	// Create a new Redis client instance.
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisAuth, // no password set
		DB:       0,         // use default DB
	})

	// Initialize RedisWrapper
	redisWrapper, err := rediswrapper.NewRedisWrapper(ctx, rdb)
	if err != nil {
		log.Error("Failed to initialize Redis", "error", err)
		return nil, err
	}

	// Return the CrawlManager instance
	return &crawler.CrawlManager{
		Logger:       log, // Use the new logger instance here
		RedisWrapper: redisWrapper,
	}, nil
}

func startConsuming(ctx context.Context, crawlSiteID string, debug bool) {
	crawlerService, err := initializeConsumeManager(ctx, debug)
	if err != nil {
		crawlerService.Logger.Error("Failed to initialize Consume Manager", "error", err)
		os.Exit(1)
	}

	// Fetch URLs from Redis
	urls, err := crawlerService.RedisWrapper.SMembers(ctx, crawlSiteID)
	if err != nil {
		crawlerService.Logger.Error("Error fetching URLs from Redis", "error", err)
		return
	}

	// Print the fetched URLs
	for _, url := range urls {
		crawlerService.Logger.Info("Fetched URL", "url", url)
	}
}

package cmd

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/rediswrapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "page-prowler",
	Short: "A web crawler for data extraction and URL consumption",
	Long: `Page Prowler is a CLI tool designed for web scraping and data extraction from websites, 
           as well as consuming URLs from a Redis set. It provides two main functionalities:

           1. The 'crawl' command: This command is used to crawl websites and extract information based on specified search terms.
           2. The 'consume' command: This command fetches URLs from a Redis set.

           Page Prowler is designed to be flexible and easy to use, making it a powerful tool for any data extraction needs.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("crawlsiteid", "", "CrawlSite ID")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	viper.BindPFlag("crawlsiteid", rootCmd.PersistentFlags().Lookup("crawlsiteid"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // Automatically override values from the .env file with those from the environment.

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while reading config file", err)
	}
}

func initializeManager(ctx context.Context, debug bool) (*crawler.CrawlManager, error) {
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	log := logger.New(debug)

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

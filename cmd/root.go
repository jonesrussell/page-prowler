package cmd

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/redis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "page-prowler",
	Short: "A tool for finding articles from websites",
	Long: `Page Prowler is a tool that finds articles from websites where the URL matches provided terms. It provides functionalities for:

1. Crawling specific websites and extracting articles that match the provided terms ('articles' command)
	2. Consuming URLs from a Redis set ('consume' command)

	In addition to the command line interface, Page Prowler also provides an HTTP API for interacting with the tool.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("crawlsiteid", "", "CrawlSite ID")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	if err := viper.BindPFlag("crawlsiteid", rootCmd.PersistentFlags().Lookup("crawlsiteid")); err != nil {
		log.Fatalf("Error binding crawlsiteid flag: %v", err)
	}

	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.Fatalf("Error binding debug flag: %v", err)
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

func initializeLogger(debug bool) logger.Logger {
	return logger.New(debug)
}

func initializeManager(ctx context.Context, debug bool) (*crawler.CrawlManager, error) {
	log := initializeLogger(debug)

	var redisClient *redis.Client
	var err error

	if !testing.Testing() {
		redisHost := viper.GetString("REDIS_HOST")
		redisAuth := viper.GetString("REDIS_AUTH")
		redisPort := viper.GetString("REDIS_PORT")
		redisClient, err = redis.NewClient(redisHost, redisAuth, redisPort)
		if err != nil {
			log.Error("Failed to initialize Redis", "error", err)
			return nil, err
		}
	}

	mongoDBWrapper, err := mongodbwrapper.NewMongoDBWrapper(ctx, "mongodb://localhost:27017")

	if err != nil {
		log.Error("Failed to initialize MongoDB", "error", err)
		return nil, err
	}

	return &crawler.CrawlManager{
		Logger:         log,
		Client:         redisClient, // Use the Client field from the Client struct
		MongoDBWrapper: mongoDBWrapper,
	}, nil
}

package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonesrussell/loggo"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/cmd"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/viper"
)

func InitializeManager(
	redisClient prowlredis.ClientInterface,
	appLogger *loggo.Logger,
) (*crawler.CrawlManager, error) {
	if redisClient == nil {
		return nil, errors.New("redisClient cannot be nil")
	}
	if appLogger == nil {
		return nil, errors.New("appLogger cannot be nil")
	}

	// Create an instance of CrawlOptions
	options := &crawler.CrawlOptions{}

	// Pass the options instance to NewCrawlManager
	collector := crawler.NewCollectorWrapper(colly.NewCollector())
	return crawler.NewCrawlManager(appLogger, redisClient, collector, options), nil
}

func main() {
	// Create a new logger instance
	loggerInterface, err := loggo.NewLogger("./loggo.log")
	if err != nil {
		fmt.Println("Error creating logger:", err)
		return
	}

	// Perform a type assertion to convert loggerInterface to *loggo.Logger
	logger, ok := loggerInterface.(*loggo.Logger)
	if !ok {
		fmt.Println("Error: logger is not of type *loggo.Logger")
		return
	}

	// Initialize Viper
	viper.AutomaticEnv() // Read environment variables
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	if err != nil {
		logger.Error("Could not read config file", err)
		return
	}

	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	if redisHost == "" {
		logger.Error("REDIS_HOST is not set but is required", nil)
		return
	}

	if redisPort == "" {
		logger.Error("REDIS_PORT is not set but is required", nil)
		return
	}

	cfg := &prowlredis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisAuth,
		DB:       0, // TODO: redisDB
	}

	ctx := context.Background()
	redisClient, err := prowlredis.NewClient(ctx, cfg)
	if err != nil {
		logger.Error("Failed to initialize Redis client", err)
		return
	}

	// Initialize the manager
	manager, err := InitializeManager(redisClient, logger)
	if err != nil {
		logger.Error("Error initializing manager", err)
		return
	}

	// Create a new root command with the manager
	rootCmd := cmd.NewRootCmd(manager)

	// Execute the root command with the logger
	err = rootCmd.Execute()
	if err != nil {
		logger.Error("root command execute failed", err)
	}
}

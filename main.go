package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jonesrussell/loggo"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/cmd"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/viper"
)

func InitializeManager(
	redisClient prowlredis.ClientInterface,
	appLogger loggo.LoggerInterface,
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

	// No need to assert loggerInterface here; pass it directly if InitializeManager expects loggo.LoggerInterface
	return crawler.NewCrawlManager(appLogger, redisClient, collector, options), nil
}

func main() {
	// Create a new logger instance with debug level
	loggerInterface, err := loggo.NewLogger("./loggo.log", slog.LevelDebug)
	if err != nil {
		fmt.Println("Error creating logger:", err)
		return
	}

	// Initialize Viper
	viper.AutomaticEnv() // Read environment variables
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Could not read config file", err) // Use fmt.Println for simplicity here since logger isn't ready yet
		return
	}

	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	if redisHost == "" || redisPort == "" {
		fmt.Println("REDIS_HOST or REDIS_PORT is not set but is required") // Simplified error handling for brevity
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
		fmt.Println("Failed to initialize Redis client:", err) // Simplified error handling
		return
	}

	// Initialize the manager with loggerInterface directly, no need for type assertion
	manager, err := InitializeManager(redisClient, loggerInterface)
	if err != nil {
		fmt.Println("Error initializing manager:", err) // Simplified error handling
		return
	}

	// Create a new root command with the manager
	rootCmd := cmd.NewRootCmd(manager)

	// Execute the root command
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println("root command execute failed:", err) // Simplified error handling
		return
	}
}

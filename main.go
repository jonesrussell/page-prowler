package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"regexp"

	"github.com/jonesrussell/loggo"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/redisstorage"
	"github.com/jonesrussell/page-prowler/cmd"
	"github.com/jonesrussell/page-prowler/crawler"
	"github.com/jonesrussell/page-prowler/dbmanager"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/spf13/viper"
)

func InitializeManager(
	dbManager dbmanager.DatabaseManagerInterface, // Add dbManager as a parameter
	appLogger loggo.LoggerInterface,
	cfg *prowlredis.Options,
) (*crawler.CrawlManager, error) {
	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	debugger := &debug.LogDebugger{
		Output: file,
	}

	// Define your allowed URLs
	URLFilters := []*regexp.Regexp{
		// www.cp24.com
		// regexp.MustCompile(`news`),
	}

	// Create a new Colly collector
	collector := colly.NewCollector(
		colly.Debugger(debugger),
		colly.MaxDepth(1),
		colly.URLFilters(URLFilters...),
	)

	collectorWrapper := crawler.NewCollectorWrapper(collector, appLogger, URLFilters)

	// Create the Redis storage
	storage := &redisstorage.Storage{
		Address:  cfg.Addr,
		Password: cfg.Password,
		DB:       1, // TODO: redis storage db
		Prefix:   "prowl",
	}

	// Set the storage for the collector
	err = collector.SetStorage(storage)
	if err != nil {
		return nil, fmt.Errorf("failed to set storage: %v", err)
	}

	// delete previous data from storage
	if err := storage.Clear(); err != nil {
		log.Fatal(err)
	}

	contentProcessor := termmatcher.NewDefaultContentProcessor()
	termMatcher := termmatcher.NewTermMatcher(appLogger, 0.8, contentProcessor)

	return crawler.NewCrawlManager(
		appLogger,
		dbManager,
		collectorWrapper,
		&crawler.CrawlOptions{},
		storage,
		termMatcher,
	), nil
}

func main() {
	// Create a new logger instance with debug level
	logger, err := loggo.NewLogger("./loggo.log", slog.LevelDebug)
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
		fmt.Println("REDIS_HOST or REDIS_PORT is not set but is required")
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
		fmt.Println("Failed to initialize Redis client:", err)
		return
	}

	dbManager := dbmanager.NewRedisManager(redisClient, logger) // Create a new DatabaseManager instance

	// Initialize the manager with loggerInterface directly, no need for type assertion
	manager, err := InitializeManager(dbManager, logger, cfg) // Pass dbManager to InitializeManager
	if err != nil {
		fmt.Println("Error initializing manager:", err)
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

package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Crawlsiteid string
	Debug       bool
)

type key int

const (
	managerKey key = iota
)

var rootCmd = &cobra.Command{
	Use:   "page-prowler",
	Short: "A tool for finding matchlinks from websites",
	Long: `Page Prowler is a tool that finds matchlinks from websites where the URL matches provided terms. It provides functionalities for:

1. Crawling specific websites and extracting matchlinks that match the provided terms ('matchlinks' command)
	2. Consuming URLs from a Redis set ('consume' command)

	In addition to the command line interface, Page Prowler also provides an HTTP API for interacting with the tool.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize your dependencies here
		ctx := context.Background()

		appLogger := initializeLogger(viper.GetBool("debug"), logger.DefaultLogLevel)

		redisHost := viper.GetString("REDIS_HOST")
		redisPort := viper.GetString("REDIS_PORT")
		mongodbUri := viper.GetString("MONGODB_URI")

		if redisHost == "" {
			log.Println("REDIS_HOST is not set but is required")
			return fmt.Errorf("REDIS_HOST is not set but is required")
		}

		if redisPort == "" {
			log.Println("REDIS_PORT is not set but is required")
			return fmt.Errorf("REDIS_PORT is not set but is required")
		}

		if mongodbUri == "" {
			log.Println("MONGODB_URI is not set but is required")
			return fmt.Errorf("MONGODB_URI is not set but is required")
		}

		redisClient, err := prowlredis.NewClient(
			ctx,
			redisHost,
			"", // No auth needed
			redisPort,
		)
		if err != nil {
			return fmt.Errorf("failed to initialize Redis client: %v", err)
		}

		mongoDBWrapper, err := mongodbwrapper.NewMongoDB(ctx, mongodbUri)
		if err != nil {
			return fmt.Errorf("failed to initialize MongoDB wrapper: %v", err)
		}

		// Now you can pass them to the initializeManager function
		manager, err := initializeManager(redisClient, appLogger, mongoDBWrapper)
		if err != nil {
			log.Println("Error initializing manager:", err)
			return err
		}

		// Set the manager to the context
		ctx = context.WithValue(ctx, managerKey, manager)

		// Set the context of the command
		cmd.SetContext(ctx)

		return nil
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Enable debug mode")

	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.Fatalf("Error binding debug flag: %v", err)
	}
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // Automatically override values from the .env file with those from the environment.

	_ = viper.ReadInConfig()
}

func initializeLogger(debug bool, level logger.LogLevel) logger.Logger {
	return logger.New(debug, level)
}

func initializeManager(
	redisClient prowlredis.ClientInterface,
	appLogger logger.Logger,
	mongoDBWrapper mongodbwrapper.MongoDBInterface,
) (*crawler.CrawlManager, error) {
	return crawler.NewCrawlManager(appLogger, redisClient, mongoDBWrapper), nil
}

package cmd

import (
	"context"
	"fmt"
	"go.uber.org/zap/zapcore"
	"log"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/redis"
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
	Short: "A tool for finding articles from websites",
	Long: `Page Prowler is a tool that finds articles from websites where the URL matches provided terms. It provides functionalities for:

1. Crawling specific websites and extracting articles that match the provided terms ('articles' command)
	2. Consuming URLs from a Redis set ('consume' command)

	In addition to the command line interface, Page Prowler also provides an HTTP API for interacting with the tool.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize your dependencies here
		ctx := context.Background()

		redisWrapper, err := redis.NewClient(viper.GetString("REDIS_HOST"), viper.GetString("REDIS_AUTH"), viper.GetString("REDIS_PORT"))
		if err != nil {
			return fmt.Errorf("failed to initialize Redis client: %v", err)
		}

		appLogger := initializeLogger(viper.GetBool("debug"))

		mongoDBWrapper, err := mongodbwrapper.NewMongoDB(ctx, viper.GetString("MONGODB_URI"))
		if err != nil {
			return fmt.Errorf("failed to initialize MongoDB wrapper: %v", err)
		}

		// Now you can pass them to the initializeManager function
		manager, err := initializeManager(redisWrapper, appLogger, mongoDBWrapper)
		if err != nil {
			return fmt.Errorf("failed to initialize manager: %v", err)
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

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while reading config file", err)
	}
}

func initializeLogger(_ bool) logger.Logger {
	return logger.New(false, zapcore.InfoLevel)
}

func initializeManager(
	redisWrapper *redis.ClientWrapper,
	appLogger logger.Logger,
	mongoDBWrapper mongodbwrapper.MongoDBInterface,
) (*crawler.CrawlManager, error) {
	return &crawler.CrawlManager{
		Logger:         appLogger,
		Client:         redisWrapper,
		MongoDBWrapper: mongoDBWrapper,
	}, nil
}

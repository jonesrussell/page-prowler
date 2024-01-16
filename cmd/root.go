package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Debug bool
)

var ErrCrawlManagerNotInitialized = errors.New("CrawlManager is not initialized")
var ErrCrawlsiteidRequired = errors.New("crawlsiteid is required")

var rootCmd = &cobra.Command{
	Use:   "page-prowler",
	Short: "A tool for finding matchlinks from websites",
	Long: `Page Prowler is a tool that finds matchlinks from websites where the URL matches provided terms. It provides functionalities for:

1. Crawling specific websites and extracting matchlinks that match the provided terms ('matchlinks' command)

	In addition to the command line interface, Page Prowler also provides an HTTP API for interacting with the tool.`,
	SilenceErrors: false,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.Println("PersistentPreRunE function started")

		// Initialize your dependencies here
		ctx := context.Background()

		appLogger, err := initializeLogger(logger.DefaultLogLevel)
		if err != nil {
			log.Println("Error initializing logger:", err)
			return err
		}
		log.Println("Logger initialized successfully")

		redisHost := viper.GetString("REDIS_HOST")
		redisPort := viper.GetString("REDIS_PORT")
		redisAuth := viper.GetString("REDIS_AUTH")
		mongodbUri := viper.GetString("MONGODB_URI")

		log.Println("Checking REDIS_HOST")
		if redisHost == "" {
			log.Println("REDIS_HOST is not set but is required")
			return fmt.Errorf("REDIS_HOST is not set but is required")
		}
		log.Println("REDIS_HOST is set")

		log.Println("Checking REDIS_PORT")
		if redisPort == "" {
			log.Println("REDIS_PORT is not set but is required")
			return fmt.Errorf("REDIS_PORT is not set but is required")
		}
		log.Println("REDIS_PORT is set")

		if mongodbUri == "" {
			log.Println("MONGODB_URI is not set but is required")
			return fmt.Errorf("MONGODB_URI is not set but is required")
		}
		log.Println("MONGODB_URI is set")

		cfg := &prowlredis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: redisAuth,
			DB:       0, // TODO: redisDB
		}

		redisClient, err := prowlredis.NewClient(ctx, cfg)
		if err != nil {
			log.Printf("Failed to initialize Redis client: %v", err)
			return fmt.Errorf("failed to initialize Redis client: %v", err)
		}
		log.Println("Redis client initialized successfully")

		mongoDBWrapper, err := mongodbwrapper.NewMongoDB(ctx, mongodbUri)
		if err != nil {
			log.Printf("Failed to initialize MongoDB wrapper: %v", err)
			return fmt.Errorf("failed to initialize MongoDB wrapper: %v", err)
		}
		log.Println("MongoDB wrapper initialized successfully")

		manager, err := InitializeManager(redisClient, appLogger, mongoDBWrapper)
		if err != nil {
			log.Printf("Error initializing manager: %v", err)
			return err
		}
		log.Println("Manager initialized successfully")

		// Set the manager to the context
		ctx = context.WithValue(ctx, common.CrawlManagerKey, manager)

		// Set the context of the command
		cmd.SetContext(ctx)

		log.Println("PersistentPreRunE function finished")
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Enable debug mode")

	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.Fatalf("Error binding debug flag: %v", err)
	}

	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(matchlinksCmd)
	rootCmd.AddCommand(workerCmd)
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // Automatically override values from the .env file with those from the environment.

	_ = viper.ReadInConfig()

	// Bind the 'crawlsiteid' flag to Viper
	if err := viper.BindPFlag("crawlsiteid", GetLinksCmd.Flags().Lookup("crawlsiteid")); err != nil {
		log.Fatalf("Error binding crawlsiteid flag: %v", err)
	}
}

func initializeLogger(level logger.LogLevel) (logger.Logger, error) {
	initlog, err := logger.New(level)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}
	return initlog, nil
}

func InitializeManager(
	redisClient prowlredis.ClientInterface,
	appLogger logger.Logger,
	mongoDBWrapper mongodbwrapper.MongoDBInterface,
) (*crawler.CrawlManager, error) {
	if redisClient == nil {
		return nil, errors.New("redisClient cannot be nil")
	}
	if appLogger == nil {
		return nil, errors.New("appLogger cannot be nil")
	}
	if mongoDBWrapper == nil {
		return nil, errors.New("mongoDBWrapper cannot be nil")
	}
	return crawler.NewCrawlManager(appLogger, redisClient, mongoDBWrapper), nil
}

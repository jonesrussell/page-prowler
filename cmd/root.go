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
	"go.uber.org/zap/zapcore"
)

var (
	Debug  bool
	Siteid string
)

var ErrCrawlManagerNotInitialized = errors.New("CrawlManager is not initialized")
var ErrSiteidRequired = errors.New("siteid is required")

var RootCmd = &cobra.Command{
	Use:   "page-prowler",
	Short: "A tool for finding matchlinks from websites",
	Long: `Page Prowler is a tool that finds matchlinks from websites where the URL matches provided terms. It provides functionalities for:

1. Crawling specific websites and extracting matchlinks that match the provided terms ('matchlinks' command)

	In addition to the command line interface, Page Prowler also provides an HTTP API for interacting with the tool.`,
	SilenceErrors: false,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// Initialize your dependencies here
		ctx := context.Background()

		appLogger, err := initializeLogger()
		if err != nil {
			log.Println("Error initializing logger:", err)
			return err
		}

		redisHost := viper.GetString("REDIS_HOST")
		redisPort := viper.GetString("REDIS_PORT")
		redisAuth := viper.GetString("REDIS_AUTH")
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

		mongoDBWrapper, err := mongodbwrapper.NewMongoDB(ctx, mongodbUri)
		if err != nil {
			log.Printf("Failed to initialize MongoDB wrapper: %v", err)
			return fmt.Errorf("failed to initialize MongoDB wrapper: %v", err)
		}

		manager, err := InitializeManager(redisClient, appLogger, mongoDBWrapper)
		if err != nil {
			log.Printf("Error initializing manager: %v", err)
			return err
		}

		// Set the manager to the context
		ctx = context.WithValue(ctx, common.CrawlManagerKey, manager)

		// Set the context of the command
		cmd.SetContext(ctx)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Initialize Viper
	viper.AutomaticEnv() // Read environment variables
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Could not read config file")
	}

	err = viper.BindEnv("debug")
	if err != nil {
		log.Fatalf("Error binding debug flag: %v", err)
	} // Bind the DEBUG environment variable to a config key

	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", viper.GetBool("debug"), "Enable debug output")

	// Bind the environment variable to the flag
	err = viper.BindEnv("siteid")
	if err != nil {
		log.Fatalf("Failed to bind env var: %v", err)
	}

	// Define the siteid flag and set its default value from the environment variable
	RootCmd.PersistentFlags().StringVarP(&Siteid, "siteid", "s", viper.GetString("siteid"), "Set siteid for redis set key")
}

func initializeLogger() (logger.Logger, error) {
	var level zapcore.Level
	level = zapcore.InfoLevel
	if Debug {
		level = zapcore.DebugLevel
	}
	return logger.New(level) // Use the new logger constructor
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

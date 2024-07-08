package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
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
	RunE: func(_ *cobra.Command, _ []string) error {
		app := fx.New(
			fx.Provide(
				NewManager,
			),
		)

		if err := app.Start(context.Background()); err != nil {
			log.Fatal(err)
		}

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
	checkError(viper.ReadInConfig(), "Could not read config file")

	checkError(viper.BindEnv("debug"), "Error binding debug flag") // Bind the DEBUG environment variable to a config key

	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", viper.GetBool("debug"), "Enable debug output")

	// Bind the environment variable to the flag
	checkError(viper.BindEnv("siteid"), "Failed to bind env var")

	// Define the siteid flag and set its default value from the environment variable
	RootCmd.PersistentFlags().StringVarP(&Siteid, "siteid", "s", viper.GetString("siteid"), "Set siteid for redis set key")
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func NewRedisClient(lc fx.Lifecycle) (prowlredis.ClientInterface, error) {
	cfg, err := getRedisConfig()
	if err != nil {
		return nil, err
	}

	var client prowlredis.ClientInterface

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			client, err = prowlredis.NewClient(context.Background(), cfg)
			return err
		},
		OnStop: func(_ context.Context) error {
			if client != nil {
				return client.Close()
			}
			return nil
		},
	})

	return client, nil
}

func getRedisConfig() (*prowlredis.Options, error) {
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	if redisHost == "" {
		return nil, fmt.Errorf("REDIS_HOST is not set but is required")
	}

	if redisPort == "" {
		return nil, fmt.Errorf("REDIS_PORT is not set but is required")
	}

	return &prowlredis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisAuth,
		DB:       0, // TODO: redisDB
	}, nil
}

func NewLogger() (logger.Logger, error) {
	loggerInstance, err := logger.New(getLoggerLevel())
	if err != nil {
		return nil, err
	}
	return loggerInstance, nil
}

func getLoggerLevel() zapcore.Level {
	if Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func NewManager(lc fx.Lifecycle) (*crawler.CrawlManager, error) {
	appLogger, err := NewLogger()
	if err != nil {
		return nil, err
	}

	redisClient, err := NewRedisClient(lc)
	if err != nil {
		return nil, err
	}

	if redisClient == nil {
		return nil, errors.New("redisClient cannot be nil")
	}
	if appLogger == nil {
		return nil, errors.New("appLogger cannot be nil")
	}

	// Create an instance of CrawlOptions
	options := &crawler.CrawlOptions{}

	// Pass the options instance to NewCrawlManager
	return crawler.NewCrawlManager(appLogger, redisClient, options), nil
}

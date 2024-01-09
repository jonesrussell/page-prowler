package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocolly/colly"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Consumer represents a consumer that processes URLs from a Redis stream.
type Consumer struct {
	Client    *redis.Client
	Collector *colly.Collector
	Logger    *zap.SugaredLogger
	Stream    string
	Group     string
}

// NewConsumer creates a new Consumer instance.
func NewConsumer(redisClient *redis.Client, collector *colly.Collector, logger *zap.SugaredLogger, stream, group string) *Consumer {
	return &Consumer{
		Client:    redisClient,
		Collector: collector,
		Logger:    logger,
		Stream:    stream,
		Group:     group,
	}
}

// Consume starts listening to the Redis stream and processes URLs.
func (c *Consumer) Consume(ctx context.Context) error {
	// Create a consumer group if it doesn't exist
	if err := c.Client.XGroupCreateMkStream(ctx, c.Stream, c.Group, "$").Err(); err != nil && err != redis.Nil {
		return err
	}

	// Create a channel for receiving signals to gracefully stop the consumer
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			return nil // The context is canceled, so we should stop the consumer gracefully.
		default:
			// Read messages from the Redis stream
			messages, err := c.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    c.Group,
				Consumer: "crawler", // Your consumer name
				Streams:  []string{c.Stream, ">"},
				Block:    0,
			}).Result()

			if err != nil {
				c.Logger.Errorw("Error reading from Redis stream", "error", err)
				continue
			}

			// Process the messages
			for _, message := range messages {
				for _, xMessage := range message.Messages {
					href := xMessage.Values["href"].(string)
					c.Logger.Infow("Processing URL", "href", href)

					// Add your custom URL processing logic here using the Collector.
					// For example, you can use c.Collector.Visit(href) to crawl the URL.

					// Acknowledge the message to remove it from the stream
					if err := c.Client.XAck(ctx, c.Stream, c.Group, xMessage.ID).Err(); err != nil {
						c.Logger.Errorw("Error acknowledging message", "error", err)
					}
				}
			}
		case <-stopChan:
			// Received a signal to stop, so we exit the consumer loop
			return nil
		}
	}
}

func main() {
	// Initialize the logger, Redis client, and Colly collector here
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("Error syncing logger: %v", err)
		}
	}()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	collector := colly.NewCollector()

	// Create a new Consumer
	consumer := NewConsumer(redisClient, collector, logger.Sugar(), "my_stream", "my_group")

	// Create a context with cancellation support
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the consumer
	if err := consumer.Consume(ctx); err != nil {
		fmt.Printf("Consumer error: %v\n", err)
		os.Exit(1)
	}
}

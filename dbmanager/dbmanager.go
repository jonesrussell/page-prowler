package dbmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/models"
)

type DatabaseManagerInterface interface {
	SaveResults(ctx context.Context, results []models.PageData, key string) error
	ClearRedisSet(ctx context.Context, key string) error
	GetLinksFromRedis(ctx context.Context, key string) ([]string, error)
	RedisOptions() prowlredis.Options
}

type RedisManager struct {
	client prowlredis.ClientInterface
	logger loggo.LoggerInterface
}

func NewRedisManager(client prowlredis.ClientInterface, logger loggo.LoggerInterface) *RedisManager {
	return &RedisManager{
		client: client,
		logger: logger,
	}
}

func (rm *RedisManager) SaveResults(ctx context.Context, results []models.PageData, key string) error {
	// Log the key and the results at the top
	rm.logger.Debug("Redis", "key", key)
	rm.logger.Debug("Redis", "results", results)

	for _, result := range results {
		data, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("error marshaling PageData: %w", err)
		}
		str := string(data)

		err = rm.client.SAdd(ctx, key, str)
		if err != nil {
			return fmt.Errorf("error adding data to Redis: %w", err)
		}
	}

	return nil
}

func (rm *RedisManager) ClearRedisSet(ctx context.Context, key string) error {
	return rm.client.Del(ctx, key)
}

func (rm *RedisManager) GetLinksFromRedis(ctx context.Context, key string) ([]string, error) {
	return rm.client.SMembers(ctx, key)
}

func (rm *RedisManager) RedisOptions() prowlredis.Options {
	return *rm.client.Options()
}

type DBManager struct {
	redisClient prowlredis.ClientInterface
}

func NewDBManager(redisClient prowlredis.ClientInterface) *DBManager {
	return &DBManager{
		redisClient: redisClient,
	}
}

func (dm *DBManager) SaveResults(ctx context.Context, results []string) error {
	// Implement the logic to save results
	// For example:
	return dm.redisClient.Set(ctx, "results", strings.Join(results, ","), 0)
}

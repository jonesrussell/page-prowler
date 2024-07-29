// dbmanager/dbmanager.go
package dbmanager

import (
	"context"
	"encoding/json"

	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/models"
)

type DatabaseManagerInterface interface {
	SaveResultsToRedis(ctx context.Context, results []models.PageData, key string) error
	ClearRedisSet(ctx context.Context, key string) error
	GetLinksFromRedis(ctx context.Context, key string) ([]string, error)
	RedisOptions() prowlredis.Options
}

type RedisManager struct {
	client prowlredis.ClientInterface
}

func NewRedisManager(client prowlredis.ClientInterface) *RedisManager {
	return &RedisManager{
		client: client,
	}
}

func (rm *RedisManager) SaveResultsToRedis(ctx context.Context, results []models.PageData, key string) error {
	for _, result := range results {
		data, err := json.Marshal(result)
		if err != nil {
			return err
		}
		str := string(data)

		err = rm.client.SAdd(ctx, key, str)
		if err != nil {
			return err
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

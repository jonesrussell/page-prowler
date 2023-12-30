package mongodbwrapper

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBInterface interface {
	Connect(ctx context.Context) error
}

type MongoDB struct {
	Client *mongo.Client
}

func NewMongoDB(ctx context.Context, uri string) (*MongoDB, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &MongoDB{Client: client}, nil
}

func (m *MongoDB) Connect(ctx context.Context) error {
	err := m.Client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	return nil
}

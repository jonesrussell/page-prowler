package mongodbwrapper

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBWrapper struct {
	Client *mongo.Client
}

func NewMongoDBWrapper(ctx context.Context, uri string) (*MongoDBWrapper, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &MongoDBWrapper{Client: client}, nil
}

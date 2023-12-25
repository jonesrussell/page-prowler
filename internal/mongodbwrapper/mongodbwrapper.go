package mongodbwrapper

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBWrapper wraps the MongoDB client.
type MongoDBWrapper struct {
	Client *mongo.Client
}

// NewMongoDBWrapper creates a new MongoDBWrapper.
func NewMongoDBWrapper(ctx context.Context, uri string) (*MongoDBWrapper, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &MongoDBWrapper{Client: client}, nil
}

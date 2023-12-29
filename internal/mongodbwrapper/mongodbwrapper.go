package mongodbwrapper

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBWrapperInterface interface {
	// Add all methods that you use from MongoDBWrapper in your code
	// For example, if you use a Connect method, you would add:
	Connect(ctx context.Context) error
}

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

func (m *MongoDBWrapper) Connect(ctx context.Context) error {
	// Your existing implementation
	// Make sure to return an error at the end
	return nil
}

package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabase represent mongodb database
type MongoDatabase struct {
	client *mongo.Client
}

// NewMongoDatabase a MongoDatabase's constructor
func NewMongoDatabase(client *mongo.Client) *MongoDatabase {
	return &MongoDatabase{
		client: client,
	}
}

// ShowList function
func (m *MongoDatabase) ShowList(ctx context.Context) error {
	fmt.Println("databases list name:")
	filter := make(bson.M)

	nameOnly := true
	opts := &options.ListDatabasesOptions{NameOnly: &nameOnly}
	databases, err := m.client.ListDatabaseNames(ctx, filter, opts)
	if err != nil {
		return fmt.Errorf("error showing database list : %s\n", err.Error())
	}

	for _, databaseName := range databases {
		fmt.Printf("- %s\n", databaseName)
	}

	return nil
}

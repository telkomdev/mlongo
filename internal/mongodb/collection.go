package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoCollection represent mongodb collection
type MongoCollection struct {
	db     *mongo.Database
	dbName string
}

// NewMongoCollection a MongoCollection's constructor
func NewMongoCollection(db *mongo.Database, dbName string) *MongoCollection {
	return &MongoCollection{
		db:     db,
		dbName: dbName,
	}
}

// ShowList function
func (m *MongoCollection) ShowList(ctx context.Context) error {
	fmt.Printf("collections name from database: %s\n", m.dbName)
	filter := make(bson.M)
	collections, err := m.db.ListCollectionNames(ctx, filter)
	if err != nil {
		return fmt.Errorf("error showing list index : %s\n", err.Error())
	}

	for _, collectionName := range collections {
		fmt.Printf("- %s\n", collectionName)
	}

	return nil
}

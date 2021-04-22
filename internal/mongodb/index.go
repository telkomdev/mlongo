package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// ASC index order type
	ASC = 1
	// DESC index order type
	DESC = -1
)

// MongoIndex represent mongodb index
type MongoIndex struct {
	db *mongo.Database
}

// NewMongoIndex a MongoIndex's constructor
func NewMongoIndex(db *mongo.Database) *MongoIndex {
	return &MongoIndex{
		db: db,
	}
}

// Create function
func (m *MongoIndex) Create(ctx context.Context, collectionName, fieldName, order string, unique bool) error {
	if order != "asc" && order != "desc" {
		return fmt.Errorf("invalid order type %s\n", order)
	}

	orderType := ASC

	if order == "desc" {
		orderType = DESC
	}

	indexOptions := &options.IndexOptions{}
	if unique {
		indexOptions.SetUnique(true)
	}

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			fieldName: orderType,
		},
		Options: indexOptions,
	}

	collection := m.db.Collection(collectionName)
	indexView := collection.Indexes()

	idx, err := indexView.CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("error execute create index : %s\n", err.Error())
	}

	fmt.Printf("create index success * index name: %s\n", idx)

	return nil
}

// Drop function
func (m *MongoIndex) Drop(ctx context.Context, collectionName, indexName string) error {
	opts := options.DropIndexes().SetMaxTime(2 * time.Second)

	collection := m.db.Collection(collectionName)
	indexView := collection.Indexes()
	raw, err := indexView.DropOne(ctx, indexName, opts)
	if err != nil {
		return fmt.Errorf("error execute drop index : %s\n", err.Error())
	}

	fmt.Printf("drop index success * %s\n", raw.String())

	return nil
}

// ShowList function
func (m *MongoIndex) ShowList(ctx context.Context) error {
	filter := make(bson.M)
	collections, err := m.db.ListCollectionNames(ctx, filter)
	if err != nil {
		return fmt.Errorf("error showing list index : %s\n", err.Error())
	}

	for _, collectionName := range collections {

		collection := m.db.Collection(collectionName)

		indexes := collection.Indexes()

		opts := options.ListIndexes().SetMaxTime(2 * time.Second)
		listIndex, err := indexes.List(ctx, opts)
		if err != nil {
			return fmt.Errorf("error execute list index : %s\n", err.Error())
		}

		fmt.Println()
		fmt.Printf("List Indexes from Collection : %s\n", collection.Name())

		for listIndex.Next(ctx) {
			var res bson.M
			err = listIndex.Decode(&res)
			if err != nil {
				return fmt.Errorf("error show list index : %s\n", err.Error())
			}

			var unique bool
			if res["unique"] != nil {
				unique = res["unique"].(bool)
			}

			fmt.Printf("- %s | unique = %t\n", res["name"], unique)

		}

	}

	return nil
}

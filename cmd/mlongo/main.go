package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ASC  = 1
	DESC = -1
)

func main() {
	var (
		host           string
		dbName         string
		username       string
		password       string
		collectionName string
		fieldName      string
		order          string
		unique         bool
		port           int
	)

	flag.StringVar(&host, "host", "localhost", "host")
	flag.IntVar(&port, "port", 27017, "port")
	flag.StringVar(&dbName, "database", "", "database name")
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&collectionName, "collection", "", "collection name")
	flag.StringVar(&fieldName, "field", "", "field name")
	flag.StringVar(&order, "order", "", "order type")

	flag.BoolVar(&unique, "unique", false, "unique index")

	flag.Usage = func() {
		fmt.Println("Usage:		mlongo [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("-host 			mongodb host eg: localhost, 127.0.0.1")
		fmt.Println("-port 			mongodb port eg: 27017")
		fmt.Println("-database 		database name")
		fmt.Println("-username 		mongodb server username")
		fmt.Println("-password 		mongodb server password")
		fmt.Println("-collection		collection name")
		fmt.Println("-field			field name")
		fmt.Println("-order			index order : asc or desc")

	}

	flag.Parse()

	if order != "asc" && order != "desc" {
		fmt.Printf("invalid order type %s\n", order)
		os.Exit(1)
	}

	address := fmt.Sprintf("mongodb://%s:%d", host, port)

	if username != "" {
		address = fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(address))
	if err != nil {
		fmt.Printf("error create client %s\n", err.Error())
		os.Exit(1)
	}

	ctx := context.Background()

	if err := client.Connect(ctx); err != nil {
		fmt.Printf("error create connection %s\n", err.Error())
		os.Exit(1)
	}

	defer func() { client.Disconnect(ctx) }()

	collection := client.Database(dbName).Collection(collectionName)

	orderType := ASC

	if order == "desc" {
		orderType = DESC
	}

	indexOptions := &options.IndexOptions{}
	if unique {
		indexOptions.SetUnique(true)
	}

	indexes := collection.Indexes()

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			fieldName: orderType,
		},
		Options: indexOptions,
	}

	idx, err := indexes.CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Printf("error create index %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("create index success * %s\n", idx)
}

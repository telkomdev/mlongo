package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	// Command type
	Command int
)

const (
	// Version of mlongo
	Version = "v1.0.0"

	// ASC index order type
	ASC = 1
	// DESC index order type
	DESC = -1

	// Create command
	Create Command = iota

	// Drop command
	Drop

	//SubscribeCommand command
	List
)

// String function
func (c Command) String() string {
	switch c {
	case Create:
		return "create"
	case List:
		return "list"
	default:
		panic("command not found")
	}
}

// CommandFromString function
func CommandFromString(c string) Command {
	switch c {
	case "create":
		return Create
	case "list":
		return List
	default:
		panic("command not found")
	}
}

// CreateIndex function
func CreateIndex(ctx context.Context, indexView mongo.IndexView, fieldName, order string, unique bool) error {
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

	idx, err := indexView.CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("error execute create index : %s\n", err.Error())
	}

	fmt.Printf("create index success * index name: %s\n", idx)

	return nil
}

// DropIndex function
func DropIndex(ctx context.Context, indexView mongo.IndexView, indexName string) error {
	opts := options.DropIndexes().SetMaxTime(2 * time.Second)
	raw, err := indexView.DropOne(ctx, indexName, opts)
	if err != nil {
		return fmt.Errorf("error execute drop index : %s\n", err.Error())
	}

	fmt.Printf("drop index success * %s\n", raw.String())

	return nil
}

// ShowListIndex function
func ShowListIndex(ctx context.Context, database *mongo.Database) error {
	filter := make(bson.M)
	collections, err := database.ListCollectionNames(ctx, filter)
	if err != nil {
		return fmt.Errorf("error showing list index : %s\n", err.Error())
	}

	for _, collectionName := range collections {

		collection := database.Collection(collectionName)

		indexes := collection.Indexes()

		opts := options.ListIndexes().SetMaxTime(2 * time.Second)
		listIndex, err := indexes.List(ctx, opts)
		if err != nil {
			return fmt.Errorf("error execute list index : %s\n", err.Error())
		}

		fmt.Printf("List Indexes from Collection : %s\n", collection.Name())

		for listIndex.Next(ctx) {
			var res bson.M
			err = listIndex.Decode(&res)
			if err != nil {
				return fmt.Errorf("error show list index : %s\n", err.Error())
			}

			var unique bool
			if res["unique"] != nil {
				unique = true
			}

			fmt.Printf("- %s | unique = %t\n", res["name"], unique)

		}
	}

	return nil
}

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
		indexName      string
		showVersion    bool
	)

	// sub command
	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	dropCommand := flag.NewFlagSet("drop", flag.ExitOnError)
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)

	// createCommand options
	createCommand.StringVar(&host, "host", "localhost", "host")
	createCommand.IntVar(&port, "port", 27017, "port")
	createCommand.StringVar(&dbName, "database", "", "database name")
	createCommand.StringVar(&username, "username", "", "username")
	createCommand.StringVar(&password, "password", "", "password")
	createCommand.StringVar(&collectionName, "collection", "", "collection name")
	createCommand.StringVar(&fieldName, "field", "", "field name")
	createCommand.StringVar(&order, "order", "", "order type")
	createCommand.BoolVar(&unique, "unique", false, "unique index")

	// dropCommand options
	dropCommand.StringVar(&host, "host", "localhost", "host")
	dropCommand.IntVar(&port, "port", 27017, "port")
	dropCommand.StringVar(&dbName, "database", "", "database name")
	dropCommand.StringVar(&username, "username", "", "username")
	dropCommand.StringVar(&password, "password", "", "password")
	dropCommand.StringVar(&collectionName, "collection", "", "collection name")
	dropCommand.StringVar(&indexName, "name", "", "index name that will be drop")

	// listCommand options
	listCommand.StringVar(&host, "host", "localhost", "host")
	listCommand.IntVar(&port, "port", 27017, "port")
	listCommand.StringVar(&dbName, "database", "", "database name")
	listCommand.StringVar(&username, "username", "", "username")
	listCommand.StringVar(&password, "password", "", "password")
	listCommand.StringVar(&collectionName, "collection", "", "collection name")

	flag.BoolVar(&showVersion, "version", false, "show version")

	flag.Usage = func() {
		fmt.Println("Usage:		mlongo [options]")
		fmt.Println()
		fmt.Println("Show List Index: ")
		fmt.Println("mlongo list -host localhost -port 27017 -username admin -password admin -database mydb -collection users")
		fmt.Println()
		fmt.Println("Create Index: ")
		fmt.Println("mlongo create -host localhost -port 27017 -username admin -password admin -database mydb -collection users -field name -order asc -unique")
		fmt.Println()
		fmt.Println("Drop Index: ")
		fmt.Println("mlongo drop -host localhost -port 27017 -username admin -password admin -database mydb -collection users -name index_name_1")
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
		fmt.Println("-version		show version")

	}

	flag.Parse()

	if showVersion {
		fmt.Printf("version %s\n", Version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("required sub command")
		os.Exit(1)
	}

	if !strings.Contains(os.Args[1], "version") {
		switch os.Args[1] {
		case "create":
			createCommand.Parse(os.Args[2:])
		case "drop":
			dropCommand.Parse(os.Args[2:])
		case "list":
			listCommand.Parse(os.Args[2:])
		default:
			fmt.Printf("invalid sub command %s\n", os.Args[1])
			os.Exit(1)
		}
	}

	address := fmt.Sprintf("mongodb://%s:%d/%s", host, port, dbName)

	if username != "" {
		address = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", username, password, host, port, dbName)
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

	disconnect := func() { client.Disconnect(ctx) }

	database := client.Database(dbName)

	collection := database.Collection(collectionName)

	indexes := collection.Indexes()

	// createCommand parsed
	if createCommand.Parsed() {
		err := CreateIndex(ctx, indexes, fieldName, order, unique)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		disconnect()
		os.Exit(0)

	}

	// dropCommand parsed
	if dropCommand.Parsed() {
		err := DropIndex(ctx, indexes, indexName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		disconnect()
		os.Exit(0)
	}

	// listCommand parsed
	if listCommand.Parsed() {
		err := ShowListIndex(ctx, database)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		disconnect()
		os.Exit(0)
	}

}

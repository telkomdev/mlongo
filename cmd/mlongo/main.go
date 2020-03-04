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
	ASC  = 1
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
func CreateIndex(indexView mongo.IndexView) {

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
		fmt.Println("Create Index: ")
		fmt.Println("mlongo create -host localhost -port 27017 -username admin -password admin -database mydb -collection users -field name -order asc -unique")
		fmt.Println()
		fmt.Println("Show List Index: ")
		fmt.Println("mlongo list -host localhost -port 27017 -username admin -password admin -database mydb -collection users")
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
		fmt.Printf("version %s\n", "v1.0.0")
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

	disconnect := func() { client.Disconnect(ctx) }

	collection := client.Database(dbName).Collection(collectionName)

	indexes := collection.Indexes()

	// createCommand parsed
	if createCommand.Parsed() {

		if order != "asc" && order != "desc" {
			fmt.Printf("invalid order type %s\n", order)
			os.Exit(1)
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

		idx, err := indexes.CreateOne(ctx, indexModel)
		if err != nil {
			fmt.Printf("error execute create index : %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("create index success * index name: %s\n", idx)

		disconnect()
		os.Exit(0)

	}

	// dropCommand parsed
	if dropCommand.Parsed() {
		opts := options.DropIndexes().SetMaxTime(2 * time.Second)
		raw, err := indexes.DropOne(ctx, indexName, opts)
		if err != nil {
			fmt.Printf("error execute drop index : %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("drop index success * %s\n", raw.String())

		disconnect()
		os.Exit(0)
	}

	// listCommand parsed
	if listCommand.Parsed() {
		opts := options.ListIndexes().SetMaxTime(2 * time.Second)
		listIndex, err := indexes.List(ctx, opts)
		if err != nil {
			fmt.Printf("error execute list index : %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("List Indexes: ")

		for listIndex.Next(ctx) {
			var res bson.M
			err = listIndex.Decode(&res)
			if err != nil {
				fmt.Printf("error show list index : %s\n", err.Error())
				os.Exit(1)
			}

			fmt.Printf("- %s\n", res["name"])

		}

		disconnect()
		os.Exit(0)
	}

}

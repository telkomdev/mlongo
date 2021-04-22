package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/telkomdev/mlongo/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Version of mlongo
	Version = "v1.0.0"
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
		indexName      string
		showVersion    bool
	)

	// collection sub command
	listCollectionCommand := flag.NewFlagSet("list", flag.ExitOnError)

	// index sub command
	createIndexCommand := flag.NewFlagSet("create", flag.ExitOnError)
	dropIndexCommand := flag.NewFlagSet("drop", flag.ExitOnError)
	listIndexCommand := flag.NewFlagSet("list", flag.ExitOnError)

	// createIndexCommand options
	createIndexCommand.StringVar(&host, "host", "localhost", "host")
	createIndexCommand.IntVar(&port, "port", 27017, "port")
	createIndexCommand.StringVar(&dbName, "database", "", "database name")
	createIndexCommand.StringVar(&username, "username", "", "username")
	createIndexCommand.StringVar(&password, "password", "", "password")
	createIndexCommand.StringVar(&collectionName, "collection", "", "collection name")
	createIndexCommand.StringVar(&fieldName, "field", "", "field name")
	createIndexCommand.StringVar(&order, "order", "", "order type")
	createIndexCommand.BoolVar(&unique, "unique", false, "unique index")

	// dropIndexCommand options
	dropIndexCommand.StringVar(&host, "host", "localhost", "host")
	dropIndexCommand.IntVar(&port, "port", 27017, "port")
	dropIndexCommand.StringVar(&dbName, "database", "", "database name")
	dropIndexCommand.StringVar(&username, "username", "", "username")
	dropIndexCommand.StringVar(&password, "password", "", "password")
	dropIndexCommand.StringVar(&collectionName, "collection", "", "collection name")
	dropIndexCommand.StringVar(&indexName, "name", "", "index name that will be drop")

	// listIndexCommand options
	listIndexCommand.StringVar(&host, "host", "localhost", "host")
	listIndexCommand.IntVar(&port, "port", 27017, "port")
	listIndexCommand.StringVar(&dbName, "database", "", "database name")
	listIndexCommand.StringVar(&username, "username", "", "username")
	listIndexCommand.StringVar(&password, "password", "", "password")

	// listCollectionCommand options
	listCollectionCommand.StringVar(&host, "host", "localhost", "host")
	listCollectionCommand.IntVar(&port, "port", 27017, "port")
	listCollectionCommand.StringVar(&dbName, "database", "", "database name")
	listCollectionCommand.StringVar(&username, "username", "", "username")
	listCollectionCommand.StringVar(&password, "password", "", "password")

	flag.BoolVar(&showVersion, "version", false, "show version")

	flag.Usage = func() {
		fmt.Println("Usage:		mlongo [options]")
		fmt.Println()

		fmt.Println("Show List Collection: ")
		fmt.Println("mlongo collection list -host localhost -port 27017 -username admin -password admin -database mydb")
		fmt.Println()

		fmt.Println("Show List Index: ")
		fmt.Println("mlongo index list -host localhost -port 27017 -username admin -password admin -database mydb")
		fmt.Println()

		fmt.Println("Create Index: ")
		fmt.Println("mlongo index create -host localhost -port 27017 -username admin -password admin -database mydb -collection users -field name -order asc -unique")
		fmt.Println()

		fmt.Println("Drop Index: ")
		fmt.Println("mlongo index drop -host localhost -port 27017 -username admin -password admin -database mydb -collection users -name index_name_1")
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
		case "collection":
			switch os.Args[2] {
			case "list":
				listCollectionCommand.Parse(os.Args[3:])
				break
			default:
				fmt.Printf("invalid sub command %s\n", os.Args[2])
				os.Exit(1)
			}
			break
		case "index":
			switch os.Args[2] {
			case "create":
				createIndexCommand.Parse(os.Args[3:])
				break
			case "drop":
				dropIndexCommand.Parse(os.Args[3:])
				break
			case "list":
				listIndexCommand.Parse(os.Args[3:])
				break
			default:
				fmt.Printf("invalid sub command %s\n", os.Args[2])
				os.Exit(1)
			}
			break

		default:
			fmt.Printf("invalid sub command %s\n", os.Args[1])
			os.Exit(1)

		}
	}

	address := fmt.Sprintf("mongodb://%s:%d/%s", host, port, dbName)

	if username != "" {
		address = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", username, password, host, port, dbName)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(address), options.Client().SetConnectTimeout(time.Second*4))
	if err != nil {
		fmt.Printf("error create client %s\n", err.Error())
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		fmt.Printf("error create connection %s\n", err.Error())
		os.Exit(1)
	}

	disconnect := func() { client.Disconnect(ctx) }

	database := client.Database(dbName)

	mongoIndex := mongodb.NewMongoIndex(database)
	mongoCollection := mongodb.NewMongoCollection(database, dbName)

	// collection command
	if listCollectionCommand.Parsed() {
		err := mongoCollection.ShowList(ctx)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		disconnect()
		os.Exit(0)
	}

	// index command
	// createIndexCommand parsed
	if createIndexCommand.Parsed() {
		err := mongoIndex.Create(ctx, collectionName, fieldName, order, unique)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		disconnect()
		os.Exit(0)

	}

	// dropIndexCommand parsed
	if dropIndexCommand.Parsed() {
		err := mongoIndex.Drop(ctx, collectionName, indexName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		disconnect()
		os.Exit(0)
	}

	// listIndexCommand parsed
	if listIndexCommand.Parsed() {
		err := mongoIndex.ShowList(ctx)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		disconnect()
		os.Exit(0)
	}

}

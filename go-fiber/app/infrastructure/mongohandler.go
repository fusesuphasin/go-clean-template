package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/casbin/casbin/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName string
var uri string

// NewSQLHandler returns connection and methos which is related to database handling.
func NewMongoHandler(ctx context.Context) (*casbin.Enforcer, error){
	dbName = os.Getenv("DB_NAME")
	uri = os.Getenv("URI")

	// Create a new client and connect to the server
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	enforcer := CasbinLoad(clientOptions, dbName)
	return enforcer, nil
}

func Open() (*mongo.Client , error) {
	dbName = os.Getenv("DB_NAME")
	uri = os.Getenv("URI")
	// Create a new client and connect to the server
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")


	return client, err
}

func GetMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := Open()

	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}
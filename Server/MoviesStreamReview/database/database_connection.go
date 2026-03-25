package database

import (
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Connect to database client
func Connect() *mongo.Client {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set!")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Connection failed", err)
	}

	return client
}

// Open a database collection
func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	databaseName := os.Getenv("DATABASE_NAME")
	return client.Database(databaseName).Collection(collectionName)
}

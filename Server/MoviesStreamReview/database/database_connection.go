package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client = Connect()

// Connect to database client
func Connect() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find the .env file")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set!")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Connection failed", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Ping failed", err)
	}

	return client
}

// Open a database collection
func OpenCollection(collectionName string) *mongo.Collection {
	databaseName := os.Getenv("DATABASE_NAME")
	return Client.Database(databaseName).Collection(collectionName)
}


package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Connect to database client
func Connect() *mongo.Client {
	loadEnv(".env")

	MongoDb := os.Getenv("MONGODB_URI")
	if MongoDb == "" {
		log.Fatal("MONGODB_URI not set!")
	}
	
	fmt.Println("MongoDb URI:", MongoDb)

	clientOptions := options.Client().ApplyURI(MongoDb)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}
	return client
}

var Client *mongo.Client = Connect()

// Open a database collection
func OpenCollection(collectionName string) *mongo.Collection {
	loadEnv(".env")

	databaseName := os.Getenv("DATABASE_NAME")

	fmt.Println("DATABASE_NAME:", databaseName)

	collection := Client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return nil
	}
	return collection
}

func loadEnv(fileName string) {
	if err := godotenv.Load(fileName); err != nil {
		log.Println("Warning: unable to find the .env file")
	}
}
package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Variable to store the database
var DB *mongo.Database

// Function to connect to the database
func ConnectDB() {
	// Get environment variables for mongodbUri and dbName
	mongodbUri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")

	if mongodbUri == "" || dbName == "" {
		log.Fatal("DB credentials are missing")
	}

	// Connect to MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(mongodbUri))
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	// Create context for Pinging to mongodb to test whether working fine
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB Ping failed: ", err)
	}

	fmt.Println("Successfully connected to MongoDB")

	// Get the database
	DB = client.Database(dbName)
}

// Function that returns pointer to collection
func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

// ConnectDB establishes connection to MongoDB
func ConnectDB(uri, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	Client = client
	Database = client.Database(dbName)

	log.Println("Successfully connected to MongoDB!")
	return nil
}

// DisconnectDB closes the MongoDB connection
func DisconnectDB() error {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := Client.Disconnect(ctx)
		if err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
		}
		log.Println("Disconnected from MongoDB!")
	}
	return nil
}

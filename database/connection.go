package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var Client *mongo.Client
var Database *mongo.Database

// ConnectDB establishes connection to MongoDB with optimized settings
func ConnectDB(uri, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure connection pool and performance options
	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(100).                        // Maximum number of connections in the pool
		SetMinPoolSize(10).                         // Minimum number of connections in the pool
		SetMaxConnIdleTime(30 * time.Second).       // Maximum time a connection can be idle
		SetConnectTimeout(10 * time.Second).        // Connection timeout
		SetSocketTimeout(30 * time.Second).         // Socket timeout
		SetServerSelectionTimeout(5 * time.Second). // Server selection timeout
		SetRetryWrites(true).                       // Enable retryable writes
		SetRetryReads(true).                        // Enable retryable reads
		SetReadConcern(readconcern.Majority()).     // Use majority read concern
		SetWriteConcern(writeconcern.Majority())    // Use majority write concern

	client, err := mongo.Connect(ctx, clientOptions)
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

	// Create indexes for better performance
	if err := createIndexes(); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	log.Println("Successfully connected to MongoDB with optimized settings!")
	return nil
}

// createIndexes creates database indexes for better query performance
func createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := Database.Collection("users")

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"email": 1,
			},
			Options: options.Index().SetUnique(true).SetName("email_unique"),
		},
		{
			Keys: map[string]interface{}{
				"name": 1,
			},
			Options: options.Index().SetName("name_index"),
		},
		{
			Keys: map[string]interface{}{
				"created_at": -1,
			},
			Options: options.Index().SetName("created_at_desc"),
		},
		{
			Keys: map[string]interface{}{
				"age": 1,
			},
			Options: options.Index().SetName("age_index"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
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

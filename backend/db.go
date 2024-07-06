package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://root:example@localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the MongoDB server to check if the connection is successful
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")

	// Insert dummy data
	insertDummyData()

}

func insertDummyData() {
	collection := client.Database("procurementdb").Collection("items")

	dummyItems := []interface{}{
		bson.D{{"name", "Item 1"}, {"description", "This is item 1"}},
		bson.D{{"name", "Item 2"}, {"description", "This is item 2"}},
		bson.D{{"name", "Item 3"}, {"description", "This is item 3"}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertMany(ctx, dummyItems)
	if err != nil {
		log.Fatalf("Failed to insert dummy data: %v", err)
	}

	log.Println("Successfully inserted dummy data")
}

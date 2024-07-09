package main

import (
	"context"
	"log"
	"time"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struc {
	MongoDB struct {
		URI string `json:"uri"`
	} 	`json:"mongodb"`
}

var client *mongo.Client

func init() {
	//read config file
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	clientOptions := options.Client().ApplyURI(config.MongoDB.URI)
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

}

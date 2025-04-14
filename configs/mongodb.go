package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDatabase establishes a connection to the MongoDB database using the URI from the environment variable.
// It returns a pointer to the mongo.Client and an error if any occurs during the connection process.
func ConnectDatabase() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
		return nil, err
	}

	fmt.Println("Connected to Database:", uri)
	return client, nil
}

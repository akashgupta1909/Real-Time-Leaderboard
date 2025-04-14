package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akashgupta1909/Real-Time-Leaderboard/configs"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		log.Fatal("MONGODB_DB environment variable is not set")
	}

	dbClient, err := configs.ConnectDatabase()
	if err != nil {
		panic(err)
	}
	redisClient, err := configs.ConnectRedis()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Database and Caching Database")
	fmt.Println("Starting server on port " + port + "...")

	router := initRoutes(
		dbClient.Database(dbName),
		redisClient,
	)
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	error := server.ListenAndServe()
	if error != nil {
		log.Fatalln(error)
	}

	defer func() {
		fmt.Println("Disconnecting Database...")
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		fmt.Println("Disconnecting from Caching Database...")
		if err := redisClient.Close(); err != nil {
			panic(err)
		}
	}()
}

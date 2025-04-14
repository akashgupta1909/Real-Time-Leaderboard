package configs

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis establishes a connection to the Redis database using the address, password, and database number from the environment variables.
// It returns a pointer to the redis.Client and an error if any occurs during the connection process.
func ConnectRedis() (*redis.Client, error) {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		return nil, fmt.Errorf("REDIS_ADDRESS environment variable is not set")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		return nil, fmt.Errorf("REDIS_DB environment variable is not set")
	}

	db, err := strconv.Atoi(redisDB)
	if err != nil {
		return nil, fmt.Errorf("error converting REDIS_DB to int: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       db,
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to caching database: %v", err)
	}
	fmt.Println("Connected to Caching Database:", redisClient.Options().Addr)
	return redisClient, nil
}

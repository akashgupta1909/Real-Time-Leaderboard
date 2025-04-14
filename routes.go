package main

import (
	"net/http"

	"github.com/akashgupta1909/Real-Time-Leaderboard/internal/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// pingHandler is a simple HTTP handler that responds with "pong" to any request.
// This is useful for testing if the server is running and reachable.
func pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}

// initRoutes initializes the HTTP router and sets up the routes for the application.
func initRoutes(database *mongo.Database, redisClient *redis.Client) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/ping", pingHandler)

	// Mount the user routes under the "/auth" path
	auth.MountUserRoutes(router, database, redisClient)

	return router
}

package event

import (
	"encoding/json"
	"net/http"

	"github.com/akashgupta1909/Real-Time-Leaderboard/internal/middleware"
	"github.com/akashgupta1909/Real-Time-Leaderboard/models"
	"github.com/akashgupta1909/Real-Time-Leaderboard/repository"
	"github.com/akashgupta1909/Real-Time-Leaderboard/utils"
	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func MountEventRoutes(router *chi.Mux, database *mongo.Database, redisClient *redis.Client) {
	eventRepository := &repository.EventRepository{
		RedisClient:     redisClient,
		MongoCollection: database.Collection("events"),
	}
	userRepsository := &repository.UserRepository{
		MongoCollection: database.Collection("users"),
		RedisClient:     redisClient,
	}

	eventRouter := chi.NewRouter()
	eventRouter.Group(func(route chi.Router) {
		route.Use(middleware.AuthMiddleware(userRepsository))
		route.Post("/create", eventRepository.EventDBRequestHandler(createEventHandler))
		route.Get("/get", eventRepository.EventDBRequestHandler(createEventHandler))
	})

	router.Mount("/event", eventRouter)
}

func createEventHandler(w http.ResponseWriter, req *http.Request, repo *repository.EventRepository) {
	user := req.Context().Value(middleware.UserContextKey).(models.User)
	eventRequest := models.CreateEventRequest{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&eventRequest)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if eventRequest.Name == "" || eventRequest.Duration == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "name and duration are required")
		return
	}

	newEvent, err := createEvent(user, repo, eventRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create event")
		return
	}
	utils.RespondWithJson(w, http.StatusCreated, newEvent)
}

package auth

import (
	"encoding/json"
	"net/http"

	"github.com/akashgupta1909/Real-Time-Leaderboard/models"
	"github.com/akashgupta1909/Real-Time-Leaderboard/repository"
	"github.com/akashgupta1909/Real-Time-Leaderboard/utils"
	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func MountUserRoutes(router *chi.Mux, database *mongo.Database, redisClient *redis.Client) {
	userRepository := &repository.UserRepository{
		RedisClient:     redisClient,
		MongoCollection: database.Collection("users"),
	}

	authRouter := chi.NewRouter()
	authRouter.Post("/sign-up", userRepository.UserDBRequestHandler(signupHandler))
	authRouter.Post("/login", userRepository.UserDBRequestHandler(loginHandler))

	router.Mount("/auth", authRouter)
}

func signupHandler(w http.ResponseWriter, req *http.Request, repo *repository.UserRepository) {
	decoder := json.NewDecoder(req.Body)
	userRequest := models.UserRequest{}

	err := decoder.Decode(&userRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	createdUser, err := createNewUser(userRequest, repo)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	utils.RespondWithJson(w, http.StatusCreated, models.NewUserResponse(createdUser))

}

func loginHandler(w http.ResponseWriter, req *http.Request, repo *repository.UserRepository) {
	decoder := json.NewDecoder(req.Body)
	userRequest := models.FetchUserRequest{}

	err := decoder.Decode(&userRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	user, err := fetchUser(userRequest.Email, userRequest.Username, repo)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	utils.RespondWithJson(w, http.StatusOK, models.NewUserResponse(user))
}

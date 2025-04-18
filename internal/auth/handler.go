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
	userRequest := models.CreateUserRequest{}

	err := decoder.Decode(&userRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if userRequest.Email == "" || userRequest.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}
	if userRequest.FirstName == "" && userRequest.LastName == "" {
		userRequest.FirstName = "Default"
		userRequest.LastName = "User"
	}

	createdUser, token, err := createNewUser(userRequest, repo)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, models.NewUserResponse(createdUser.Username, "user created successfully", token))
}

func loginHandler(w http.ResponseWriter, req *http.Request, repo *repository.UserRepository) {
	decoder := json.NewDecoder(req.Body)
	userRequest := models.LoginUserRequest{}

	err := decoder.Decode(&userRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if userRequest.Username == "" && userRequest.Email == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "username or email is required")
		return
	}
	if userRequest.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "password is required")
		return
	}

	user, token, err := fetchUser(userRequest, repo)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, models.NewUserResponse(user.Username, "user login successful", token))
}

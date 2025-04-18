package auth

import (
	"crypto"
	"encoding/base64"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/akashgupta1909/Real-Time-Leaderboard/models"
	"github.com/akashgupta1909/Real-Time-Leaderboard/repository"
)

// createNewUser creates a new user in the database.
// It takes a UserRequest object and a UserRepository object as parameters.
func createNewUser(userRequest models.CreateUserRequest, repo *repository.UserRepository) (models.User, string, error) {
	userModel := models.User{
		Email:     strings.ToLower(userRequest.Email),
		FirstName: strings.ToLower(userRequest.FirstName),
		LastName:  strings.ToLower(userRequest.LastName),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userModel.Username = strings.ToLower(userModel.FirstName + userModel.LastName + "_" + strconv.Itoa(rand.Intn(1000)))
	passwordHash := string(crypto.SHA256.New().Sum([]byte(userRequest.Password)))
	userModel.Password = base64.StdEncoding.EncodeToString([]byte(passwordHash))

	createdUser, err := repo.CreateUser(userModel)
	if err != nil {
		return models.User{}, "", err
	}

	// Generate JWT token
	token, err := GenerateJWT(createdUser.ID.Hex())
	if err != nil {
		return models.User{}, "", err
	}

	return createdUser, token, nil
}

// fetchUser will fetch user based on email or username. If both are provided, it will prioritize email.
// It will return the user object and an error if any.
func fetchUser(userRequest models.LoginUserRequest, repo *repository.UserRepository) (models.User, string, error) {
	var err error
	var user models.User

	if userRequest.Email != "" {
		user, err = repo.FindUserByEmail(userRequest.Email)
	} else if userRequest.Username != "" {
		user, err = repo.FindUserByUsername(userRequest.Username)
	} else {
		return models.User{}, "", errors.New("either email or username must be provided")
	}

	if err != nil {
		return models.User{}, "", err
	}

	passwordHash := string(crypto.SHA256.New().Sum([]byte(userRequest.Password)))
	if user.Password != base64.StdEncoding.EncodeToString([]byte(passwordHash)) {
		return models.User{}, "", errors.New("invalid password")
	}

	token, err := GenerateJWT(user.ID.Hex())
	if err != nil {
		return models.User{}, "", err
	}

	return user, token, nil
}

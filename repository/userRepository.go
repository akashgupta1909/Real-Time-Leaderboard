package repository

import (
	"context"
	"net/http"

	"github.com/akashgupta1909/Real-Time-Leaderboard/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	RedisClient     *redis.Client
	MongoCollection *mongo.Collection
}

type userHandler func(w http.ResponseWriter, req *http.Request, repo *UserRepository)

func (repo *UserRepository) UserDBRequestHandler(next userHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		next(w, req, repo)
	}
}

// CreateUser is a method that creates a new user in the database.
// It takes a UserRequest object as input and returns the created User object or an error.
func (repo *UserRepository) CreateUser(user models.User) (models.User, error) {
	createdUser, err := repo.MongoCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return models.User{}, err
	}
	user.ID = createdUser.InsertedID.(primitive.ObjectID)
	return user, nil
}

// FindUserByEmail is a method that finds a user by their email address.
// It takes an email string as input and returns the User object or an error.
func (repo *UserRepository) FindUserByEmail(email string) (models.User, error) {
	var user models.User
	err := repo.MongoCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// FindUserByUsername is a method that finds a user by their username.
// It takes a username string as input and returns the User object or an error.
func (repo *UserRepository) FindUserByUsername(username string) (models.User, error) {
	var user models.User
	err := repo.MongoCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// FindUserById is a method that finds a user by their ID.
// It takes an ID string as input and returns the User object or an error.
func (repo *UserRepository) FindUserByID(id string) (models.User, error) {
	var user models.User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}
	err = repo.MongoCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

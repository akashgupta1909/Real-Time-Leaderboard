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

type EventRepository struct {
	RedisClient     *redis.Client
	MongoCollection *mongo.Collection
}

type eventHandler func(w http.ResponseWriter, req *http.Request, repo *EventRepository)

// EventDBRequestHandler is a middleware function that wraps the event handler.
// It takes an eventHandler function as input and returns an http.HandlerFunc.
func (repo *EventRepository) EventDBRequestHandler(next eventHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		next(w, req, repo)
	}
}

// CreateEvent is a method that creates a new event in the database.
// It takes an Event object as input and returns the created Event object or an error.
func (repo *EventRepository) CreateEvent(event models.Event) (models.Event, error) {
	createdEvent, err := repo.MongoCollection.InsertOne(context.TODO(), event)
	if err != nil {
		return models.Event{}, err
	}
	event.ID = createdEvent.InsertedID.(primitive.ObjectID)
	return event, nil
}

// FindEventById is a method that finds an event by its ID.
// It takes an event ID as input and returns the Event object or an error.
func (repo *EventRepository) FindEventById(eventId primitive.ObjectID) (models.Event, error) {
	var event models.Event
	err := repo.MongoCollection.FindOne(context.TODO(), bson.M{"_id": eventId}).Decode(&event)
	if err != nil {
		return models.Event{}, err
	}
	return event, nil
}

package event

import (
	"time"

	"github.com/akashgupta1909/Real-Time-Leaderboard/models"
	"github.com/akashgupta1909/Real-Time-Leaderboard/repository"
)

func createEvent(user models.User, repo *repository.EventRepository, req models.CreateEventRequest) (models.Event, error) {
	event := models.Event{
		Name:         req.Name,
		Creator:      user.Username,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Duration:     req.Duration,
		IsActive:     true,
		Participants: req.Participants,
	}

	newEvent, err := repo.CreateEvent(event)
	if err != nil {
		return models.Event{}, err
	}
	return newEvent, nil
}

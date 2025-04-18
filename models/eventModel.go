package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Creator      string             `json:"creator" bson:"creator"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	Participants []string           `json:"participants" bson:"participants"`
	Duration     int                `json:"duration" bson:"duration"`
	IsActive     bool               `json:"is_active" bson:"is_active"`
	ScoreBoard   []ScoreBoardEntry  `json:"scoreboard" bson:"scoreboard"`
}

type CreateEventRequest struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
	Duration     int      `json:"duration"`
}

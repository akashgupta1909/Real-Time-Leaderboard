package models

type ScoreBoardEntry struct {
	Username  string `json:"username" bson:"username"`
	Email     string `json:"email" bson:"email"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Score     int    `json:"score" bson:"score"`
	Rank      int    `json:"rank" bson:"rank"`
}

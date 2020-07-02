package leaderboard

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"

	s "github.com/edlorenzo/leaderboards-service/config"
)

func collection() *mgo.Collection {
	return s.DB.C("leaderboard")
}

// LeaderboardForm struct
type LeaderboardForm struct {
	Score   int    `json:"score_to_add" validate:"required"`
	BoardID string `json:"board_id"`
	UserID  string `json:"user_id"`
}

// LeaderboardScore struct
type LeaderboardScore struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	BoardID     string        `json:"board_id" validate:"required"`
	UserID      string        `json:"user_id" validate:"required"`
	Score       int           `json:"score"`
	ScoredAt    time.Time     `json:"scored_at" bson:"scored_at"`
	DateAdded   time.Time     `json:"date_added" bson:"date_added"`
	DateUpdated time.Time     `json:"date_updated" bson:"date_updated"`
}

// Add New User
func Add(params *LeaderboardForm) (time.Time, error) {
	d := &LeaderboardScore{}

	d.Score = params.Score
	d.BoardID = params.BoardID
	d.UserID = params.UserID
	d.ScoredAt, d.DateAdded, d.DateUpdated = time.Now(), time.Now(), time.Now()

	return d.DateAdded, collection().Insert(d)
}

// GetOne accepts interface
func GetOne(condition interface{}) (*LeaderboardScore, error) {
	res := LeaderboardScore{}

	if err := collection().Find(condition).One(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// IsUsed check the existence of a admin id attached in a leaderboard
func IsUsed(id string) (bool, error) {
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"$or": []bson.M{
			bson.M{"_id": bson.ObjectIdHex(id)}}}},
		bson.M{"$group": bson.M{
			"_id":   nil,
			"count": bson.M{"$sum": 1}}},
	}

	pipe := s.DB.C("leaderboard").Pipe(pipeline)

	resp := []bson.M{}
	err := pipe.All(&resp)
	return len(resp) > 0, err
}

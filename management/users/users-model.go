package users

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"

	s "github.com/edlorenzo/leaderboards-service/config"
)

func init() {
	c := s.DB.C("users")

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func collection() *mgo.Collection {
	return s.DB.C("users")
}

// UserForm struct
type UserForm struct {
	Name string `json:"name" validate:"required"`
}

// Users struct
type Users struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `json:"name"`
	DateAdded   time.Time     `json:"date_added" bson:"date_added"`
	DateUpdated time.Time     `json:"date_updated" bson:"date_updated"`
}

// Add New User
func Add(params *UserForm) (time.Time, error) {
	d := &Users{}

	d.Name = params.Name
	d.DateAdded, d.DateUpdated = time.Now(), time.Now()

	return d.DateAdded, collection().Insert(d)
}

// GetOne accepts interface
func GetOne(condition interface{}) (*Users, error) {
	res := Users{}

	if err := collection().Find(condition).One(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

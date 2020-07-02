package admin

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	s "github.com/edlorenzo/leaderboards-service/config"
)

// AdminLeaderboards struct
type AdminLeaderboards struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `json:"name"`
	DateAdded   time.Time     `json:"date_added" bson:"date_added"`
	DateUpdated time.Time     `json:"date_updated" bson:"date_updated"`
}

// Properties struct
type Properties []struct {
	PropertyName  string `json:"property_name"`
	PropertyValue string `json:"property_value"`
}

// AdminForm struct
type AdminForm struct {
	Name string `json:"name" validate:"required"`
}

// KeyIDandValue struct
type KeyIDandValue []struct {
	PropertyName  string `json:"property_name"`
	PropertyValue string `json:"property_value"`
}

// Ids struct
type Ids struct {
	Data []string `json:"data"`
}

// LeaderboardList struct
type LeaderboardList struct {
	Score        int       `json:"score"`
	ScoredAt     time.Time `json:"scored_at" bson:"scored_at"`
	UserID       string    `json:"user_id" validate:"required"`
	Name         string    `json:"name"`
	Rank         int       `json:"rank"`
	Page         int       `json:"page"`
	ItemsPerPage int       `json:"perPage"`
	Search       string    `json:"search"`
	Order        string    `json:"order"`
	OrderBy      string    `json:"orderby"`
}

func init() {
	c := s.DB.C("admin_leaderboard")

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
	return s.DB.C("admin_leaderboard")
}

/*
// GetAll func - gets all admin
func GetAll(params LeaderboardList) ([]AdminLeaderboards, int, error) {
	var condition = bson.M{}
	filter := bson.M{}
	pack := []bson.M{}
	res := []AdminLeaderboards{}

	if len(params.HookName) > 0 {
		for _, v := range params.HookName {
			pack = append(pack, bson.M{"$regex": v, "$options": "i"})
		}
		filter = bson.M{"$or": pack}
	}

	order := params.Order
	if order != "" {
		if params.OrderBy == "desc" {
			order = "-" + order
		} else {
			order = params.Order
		}
	} else {
		order = "-date_added"
	}

	condition = bson.M{
		"$and": []bson.M{
			bson.M{"$or": []bson.M{
				bson.M{"hookname": bson.M{"$regex": params.Search, "$options": "i"}},
				bson.M{"hookurl": bson.M{"$regex": params.Search, "$options": "i"}},
				bson.M{"httpmethod": bson.M{"$regex": params.Search, "$options": "i"}},
			}},
			filter,
		},
	}

	page := params.Page
	itemsPerPage := params.ItemsPerPage
	skip := 0
	if page == 0 {
		page = 1
	}

	if page > 1 {
		skip = (page - 1) * itemsPerPage
	}

	if err := collection().Find(condition).Sort(order).All(&res); err != nil {
		return nil, 0, err
	}
	totalCount := len(res)

	if itemsPerPage != 0 {

		if err := collection().Find(condition).Sort(order).Skip(skip).Limit(itemsPerPage).All(&res); err != nil {
			return nil, 0, err
		}
	}

	return res, totalCount, nil
}
*/

// Add New Admin Leaderboard
func Add(params *AdminForm) (time.Time, error) {
	d := &AdminLeaderboards{}

	d.Name = params.Name
	d.DateAdded, d.DateUpdated = time.Now(), time.Now()

	return d.DateAdded, collection().Insert(d)
}

// GetOne accepts interface
func GetOne(condition interface{}) (*AdminLeaderboards, error) {
	res := AdminLeaderboards{}

	if err := collection().Find(condition).One(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Update admin
func Update(query map[string]interface{}, s *AdminLeaderboards) error {
	return collection().Update(query, s)
}

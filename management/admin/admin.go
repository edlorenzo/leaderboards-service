package admin

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	c "github.com/edlorenzo/leaderboards-service/config"
	h "github.com/edlorenzo/leaderboards-service/management/helpers"
)

var validate *validator.Validate

// AddAdminLeaderboard func.
func AddAdminLeaderboard(ctx context.Context) {
	validate = validator.New()
	params := &AdminForm{}

	if err := ctx.ReadJSON(params); err != nil {
		ctx.JSON(context.Map{"code": iris.StatusBadRequest, "message": err.Error()})
		return
	}

	err := validate.Struct(params)
	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		ctx.JSON(context.Map{"code": iris.StatusBadRequest, "message": h.FieldRequired, "error": errors})
		return
	}

	timedate, err := Add(params)

	if err != nil {
		response := err.Error()
		if mgo.IsDup(err) {
			response = h.DuplicateName
		}
		ctx.JSON(context.Map{"code": iris.StatusBadRequest, "message": response})
		return
	}
	rp, _ := GetOne(bson.M{"date_added": timedate})

	ctx.JSON(context.Map{"board": context.Map{"_id": rp.ID, "name": rp.Name}})
}

// GetAdminLeaderboard func - Get single admin
func GetAdminLeaderboard(ctx context.Context) {
	id := ctx.Params().Get("ID")
	//per_page := ctx.Params().Get("per_page")
	//page := ctx.Params().Get("page")

	response := ""
	code := iris.StatusBadRequest
	if bson.IsObjectIdHex(id) == false {
		response = h.ValidField
	} else {
		rp, err := GetOne(bson.M{"_id": bson.ObjectIdHex(id)})
		if err != nil {
			response = err.Error()
		} else {

			collection := c.DB.C("leaderboard")
			pipeline := []bson.M{
				bson.M{"$match": bson.M{"boardid": id}},
				bson.M{"$lookup": bson.M{"from": "admin_leaderboard", "localField": "boardid", "foreignField": "_id", "as": "score1"}},
				bson.M{"$lookup": bson.M{"from": "users", "localField": "userid", "foreignField": "_id", "as": "score2"}},
				bson.M{"$sort": bson.M{"score": -1}},
			}

			pipe := collection.Pipe(pipeline)
			resp := []bson.M{}
			err = pipe.All(&resp)

			if err != nil {
				fmt.Println("Errored: %#v \n", err)
			}
			fmt.Println(resp)



			code = iris.StatusOK
			ctx.JSON(context.Map{"board": context.Map{"_id": rp.ID, "name": rp.Name, "test": rp.ID}})
			return
		}
	}

	ctx.JSON(context.Map{"code": code, "board": response})

}

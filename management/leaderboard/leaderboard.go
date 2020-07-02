package leaderboard

import (
	"fmt"
	c "github.com/edlorenzo/leaderboards-service/config"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	h "github.com/edlorenzo/leaderboards-service/management/helpers"
)

var validate *validator.Validate

// AddLeaderboardScore func.
func AddLeaderboardScore(ctx context.Context) {
	lid := ctx.Params().Get("LID")
	uid := ctx.Params().Get("UID")
	response := ""
	code := iris.StatusBadRequest

	if bson.IsObjectIdHex(lid) == false && bson.IsObjectIdHex(uid) == false {
		response = h.ValidField
	} else {
		validate = validator.New()
		params := &LeaderboardForm{}
		params.BoardID = lid
		params.UserID = uid

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
		rp, _ := GetOne(bson.M{"scored_at": timedate})

		ctx.JSON(context.Map{"entry": context.Map{"_id": rp.ID, "board_id": rp.BoardID, "score": rp.Score, "scored_at": rp.ScoredAt, "user_id": rp.UserID}})
		return
	}

	ctx.JSON(context.Map{"code": code, "entry": response})
}

// GetLeaderboardScore func.
func GetLeaderboardScore(ctx context.Context) {
	id := ctx.Params().Get("ID")
	response := ""
	code := iris.StatusBadRequest
	if bson.IsObjectIdHex(id) == false {
		response = h.ValidField
	} else {

		collection := c.DB.C("leaderboard")
		pipelines := []bson.M{
			bson.M{"$match": bson.M{"boardid": bson.ObjectIdHex("5efb46478cf68220d80bbf31")}},
			bson.M{"$lookup": bson.M{ "from": "localCollection", "localField": "localField", "foreignField": "foreignField", "as": "resultField"}},
		}
		pipe := collection.Pipe(pipelines)
		resp := []bson.M{}
		var err error
		err = pipe.All(&resp)

		if err != nil {
			fmt.Println("Errored: %#v \n", err)
		}
		fmt.Println(resp)


		rp, err := GetOne(bson.M{"_id": bson.ObjectIdHex(id)})
		if err != nil {
			response = err.Error()
		} else {
			code = iris.StatusOK
			ctx.JSON(context.Map{"entry": context.Map{"_id": rp.ID, "board_id": rp.BoardID, "score": rp.Score, "scored_at": rp.ScoredAt, "user_id": rp.UserID}})
			return
		}
	}

	ctx.JSON(context.Map{"code": code, "user": response})
}

package users

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	h "github.com/edlorenzo/leaderboards-service/management/helpers"
)

var validate *validator.Validate

// AddUser func.
func AddUser(ctx context.Context) {
	validate = validator.New()
	params := &UserForm{}

	//ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

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

	ctx.JSON(context.Map{"user": context.Map{"_id": rp.ID, "name": rp.Name}})
}

// GetUser func - Get single user
func GetUser(ctx context.Context) {
	id := ctx.Params().Get("ID")
	response := ""
	code := iris.StatusBadRequest
	if bson.IsObjectIdHex(id) == false {
		response = h.ValidField
	} else {
		rp, err := GetOne(bson.M{"_id": bson.ObjectIdHex(id)})
		if err != nil {
			response = err.Error()
		} else {
			code = iris.StatusOK
			ctx.JSON(context.Map{"user": context.Map{"_id": rp.ID, "name": rp.Name}})
			return
		}
	}

	ctx.JSON(context.Map{"code": code, "user": response})
}

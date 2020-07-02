package main

import (
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/sessions"
	"log"

	a "github.com/edlorenzo/leaderboards-service/management/admin"
	l "github.com/edlorenzo/leaderboards-service/management/leaderboard"
	u "github.com/edlorenzo/leaderboards-service/management/users"
)

var sessionsManager *sessions.Sessions

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cookieName := "leaderbordid"
	hashKey := []byte("BD9C88804534D2F04E1FD96E687E6603")
	blockKey := []byte("0276E377D343D5DECDAAF103F8AF98D4")
	secureCookie := securecookie.New(hashKey, blockKey)

	sessionsManager = sessions.New(sessions.Config{
		Cookie: cookieName,
		Encode: secureCookie.Encode,
		Decode: secureCookie.Decode,
	})
}

func main() {

	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	crs := func(ctx iris.Context) {

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type")
		ctx.Header("X-Frame-Options", "sameorigin")
		ctx.Header("Content-Security-Policy", "self")
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Referrer-Policy", "origin")
		ctx.Header("X-XSS-Protection", "1; mode=block")

		if ctx.Method() != "OPTIONS" {
			ctx.Next()
		}

	}

	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions)
	{
		//User
		v1.Handle("POST", "/user", u.AddUser)
		v1.Handle("GET", "/user/{ID: string}", u.GetUser)

		//Admin Leaderboards
		v1.Handle("POST", "/admin/leaderboard", a.AddAdminLeaderboard)
		v1.Handle("GET", "/admin/leaderboard/{ID: string}", a.GetAdminLeaderboard)

		//Leaderboards Scores
		v1.Handle("POST", "/leaderboard", l.AddLeaderboardScore)
		v1.Handle("PUT", "/leaderboard/{LID: string}/user/{UID: string}/add_score", l.AddLeaderboardScore)
		v1.Handle("GET", "/leaderboard/{ID: string}", a.GetAdminLeaderboard)

	}

	app.Run(iris.Addr(":8081"), iris.WithoutServerError(iris.ErrServerClosed))
}

package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
)

// DB var
var DB *mgo.Database

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	conf := New()
	fmt.Println("DB HOST: ", conf.DB.DbHost)
	session, err := mgo.Dial(conf.DB.DbHost)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	DB = session.DB(conf.DB.DbName)
}

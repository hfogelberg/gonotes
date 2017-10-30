package main

import (
	"os"

	"github.com/gorilla/sessions"
	mgo "gopkg.in/mgo.v2"
)

var (
	MongoDbHost = os.Getenv("MONGO_DB_HOST")
	MongoDb     = "gonotes"
	HmacSecret  = "secret"
	CookieName  = "notes"
	Port        = os.Getenv("PORT")
	store       = sessions.NewCookieStore([]byte(HmacSecret))
)

type Connection struct {
	Db *mgo.Database
}

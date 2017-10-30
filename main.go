package main

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/hfogelberg/goat"
	"github.com/hfogelberg/toogo"
	"github.com/urfave/negroni"
)

func main() {
	session, err := mgo.Dial(MongoDbHost)
	if err != nil {
		fmt.Printf("Error connecting to MongoDB %s\n", err.Error())
		return
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	connection := Connection{session.DB(MongoDb)}

	// Initialize goat
	goat.New(sessions.NewCookieStore([]byte(HmacSecret)), "/notes/index", CookieName)

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", connection.indexHandler)
	r.HandleFunc("/about", aboutHandler)
	r.HandleFunc("/contact", contactHandler)
	r.HandleFunc("/googlelogin", goat.GoogleLoginHandler)
	r.HandleFunc("/gcallback", goat.GoogleCallbackHandler)
	r.HandleFunc("/favicon.ico", toogo.FaviconHandler)

	adm := r.PathPrefix("/notes/").Subrouter().StrictSlash(false)
	adm.HandleFunc("/index", adminHandler)
	adm.HandleFunc("/create", createHandler).Methods("GET")
	adm.HandleFunc("/create", connection.createPostHandler).Methods("POST")

	mux := http.NewServeMux()
	mux.Handle("/", r)
	mux.Handle("/notes/", negroni.New(
		negroni.HandlerFunc(connection.adminAuthMiddleware),
		negroni.Wrap(r),
	))

	static := http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
	r.PathPrefix("/public").Handler(static)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(Port, n)

}

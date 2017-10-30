package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/index.html", "templates/layout.html")
	err = tpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Fatalln("Error serving index template ", err.Error())
		return
	}
}

func (conn *Connection) adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("admin handler")
	session, err := store.Get(r, CookieName)
	if err != nil {
		log.Printf("Error getting cookie %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
	email := session.Values["email"].(string)

	notes, err := conn.getNotes(email)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

	tpl, err := template.New("").ParseFiles("templates/admin.html", "templates/layout.html")
	err = tpl.ExecuteTemplate(w, "layout", notes)
	if err != nil {
		log.Fatalln("Error serving admin template ", err.Error())
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/contact.html", "templates/layout.html")
	err = tpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Fatalln("Error serving contact template ", err.Error())
		return
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/create.html", "templates/layout.html")
	err = tpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Fatalln("Error serving create template ", err.Error())
		return
	}
}

func (conn *Connection) createPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST NOTE")
	text := r.FormValue("note")
	session, err := store.Get(r, CookieName)
	if err != nil {
		log.Printf("Error getting cookie %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
	email := session.Values["email"].(string)
	created := time.Now()

	note := Note{
		Text:    text,
		Email:   email,
		Created: created,
	}

	if err := conn.createNote(note); err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

	http.Redirect(w, r, "/notes/index", http.StatusPermanentRedirect)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/about.html", "templates/layout.html")
	err = tpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Fatalln("Error serving about template ", err.Error())
		return
	}
}

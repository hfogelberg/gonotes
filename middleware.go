package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (conn *Connection) adminAuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Admin Middleware")
	// 1. Get cookie. If there is no cookie, redirect to login
	session, err := store.Get(r, CookieName)
	if err != nil {
		fmt.Printf("Error getting cookie %s\n", err.Error())
		login(w, r)
	}

	if session.Values["accessToken"] == nil {
		fmt.Println("No token. Redirect to login")
		login(w, r)
	}

	log.Println("We have an access token")
	// 2. Get access token and email from cookie
	token := session.Values["accessToken"].(string)
	email := session.Values["email"].(string)

	// 3. Get user from Db by email.
	user, err := conn.getUserFromDb(email)
	if err != nil {
		// User not found or some other error. Redirect to Index
		fmt.Println("No user with that address, so create user")

		err := conn.createUser(email, token)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			next(w, r)
			return
		}
	}
	log.Println("Email is in Db")

	// // 3. Email is OK. Check if token is in Db, otherwise add it
	isOk := tokenIsValid(user, token)
	if isOk == false {
		err := conn.createToken(email, token)
		if err != nil {
			next(w, r)
		}
	}
	log.Println("Token is in Db")
	next(w, r)
}

func tokenIsValid(user User, token string) bool {
	for _, t := range user.Tokens {
		fmt.Println(t)

		if strings.Compare(token, t) == 0 {
			return true
		}
	}
	return false
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Redirect to login")
	http.Redirect(w, r, "/googlelogin", http.StatusPermanentRedirect)
}

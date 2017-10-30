package main

import "time"

type User struct {
	Email  string   `json:"email" bson:"email"`
	Tokens []string `json:"tokens" bson:"tokens"`
}

type Note struct {
	Text    string    `json:"text" bson:"text`
	Email   string    `json:"email" bson:"email"`
	Created time.Time `json:"created" bson:"created"`
}

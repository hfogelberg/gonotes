package main

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

func (conn *Connection) createNote(note Note) error {
	log.Printf("Save note")
	if err := conn.Db.C("notes").Insert(&note); err != nil {
		log.Printf("Error saving note to Db %s\n", err.Error())
		return err
	}

	return nil
}

func (conn *Connection) getNotes(email string) ([]*Note, error) {
	var notes []*Note

	err := conn.Db.C("notes").Find(bson.M{"email": email}).All(&notes)
	if err != nil {
		log.Printf("Error getting notes %s\n", err.Error())
		return notes, err
	}
	return notes, nil
}

func (conn *Connection) getUserFromDb(email string) (User, error) {
	user := User{}
	log.Printf("Checking email %s\n", email)
	if err := conn.Db.C("users").Find(bson.M{"email": email}).One(&user); err != nil {
		log.Printf("Error fetching user from Db %s\n", err.Error())
		return user, err
	}

	return user, nil
}

func (conn *Connection) createUser(email string, token string) error {
	log.Printf("Create user. Email: %s, Token: %s \n", email, token)

	user := User{
		Email:  email,
		Tokens: []string{token},
	}

	if err := conn.Db.C("users").Insert(&user); err != nil {
		log.Printf("Error inserting user %s\n", err.Error())
		return err
	}
	return nil
}

func (conn *Connection) createToken(email string, token string) error {
	log.Printf("Create token %s\n %s\n\n", email, token)
	query := bson.M{"email": email}
	update := bson.M{"$push": bson.M{"tokens": token}}
	if err := conn.Db.C("users").Update(query, update); err != nil {
		log.Printf("Error appending token %s\n", err.Error())
		return err
	}
	return nil
}

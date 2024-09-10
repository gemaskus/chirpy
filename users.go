package main

import (
	"fmt"
	"log"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	log.Printf("Creating User: %s", email)

	currentDBStructure, err := db.loadDB()

	if err != nil {
		return User{}, fmt.Errorf("Create User Err: %v", err)
	}

	id := len(currentDBStructure.Users) + 1

	newUser := User{
		ID:    id,
		Email: email,
	}

	currentDBStructure.Users[id] = newUser

	db.writeDB(currentDBStructure)

	return newUser, nil
}

func (db *DB) GetUsers() ([]User, error) {
	log.Printf("Retrieveing Users from DB")
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStruct.Users))

	for _, user := range dbStruct.Users {

		users = append(users, user)
	}

	return users, nil
}

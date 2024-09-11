package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	log.Printf("Creating User: %s", email)

	currentDBStructure, err := db.loadDB()

	if err != nil {
		return User{}, fmt.Errorf("Create User Err: %v", err)
	}

	id := len(currentDBStructure.Users) + 1

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return User{}, fmt.Errorf("Could not hash the password: %v", err)
	}

	newUser := User{
		ID:       id,
		Email:    email,
		Password: string(hashedPassword),
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
		newUser := User{
			Email:    user.Email,
			ID:       user.ID,
			Password: user.Password,
		}
		users = append(users, newUser)
	}

	return users, nil
}

func checkPassword(password, hashedPW string) (bool, error) {
	log.Printf("Password to be checked: %s", password)
	hashedNewPW, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	log.Printf("HashedNewPW: %s", string(hashedNewPW))
	log.Printf("HashedOldPW: %s", hashedPW)

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(password)); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

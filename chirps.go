package main

import (
	"fmt"
	"log"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	log.Printf("Creating Chirp: %s", body)
	currentDBStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, fmt.Errorf("Create Chirp Error: %v", err)
	}

	id := len(currentDBStructure.Chirps) + 1

	newChirp := Chirp{
		ID:   id,
		Body: body,
	}

	currentDBStructure.Chirps[id] = newChirp

	db.writeDB(currentDBStructure)

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	log.Printf("Retrieving Chirps from DB")
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStruct.Chirps))

	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

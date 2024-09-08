package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(filePath string) (*DB, error) {
	log.Printf("Creating new DB link")
	db := DB{
		path: filePath,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()

	return &db, err
}

func (db *DB) ensureDB() error {
	log.Printf("Ensuring the database exists")
	_, err := os.ReadFile(db.path)

	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("DB File does not exist, creating DB file")
			dbStructure := DBStructure{
				Chirps: map[int]Chirp{},
			}

			dat, err := json.Marshal(dbStructure)

			if err != nil {
				return err
			}

			if err := os.WriteFile(db.path, dat, 0666); err != nil {
				log.Fatal(err)
				return err
			}
		} else {
			//Some other error happened, log it and write out Internal server error
			return err
		}
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	log.Printf("Loading Chirp DB")
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}
	fileContents, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, fmt.Errorf("loadDB: Cannot read from given file: %s with error %v", db.path, err)
	}

	if err := json.Unmarshal(fileContents, &dbStructure); err != nil {
		return dbStructure, fmt.Errorf("loadDB: Cannot unmarshal the data read from the file: %s with error: %v", db.path, err)
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	log.Printf("Writing Chirp DB to file")
	db.mux.RLock()
	defer db.mux.RUnlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	if err := os.WriteFile(db.path, dat, 0666); err != nil {
		return err
	}
	return nil
}

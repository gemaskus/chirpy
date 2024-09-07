package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type postBody struct {
		MessageBody string `json:"body"`
	}

	badwords := [3]string{"kerfuffle", "sharbert", "fornax"}

	decoder := json.NewDecoder(r.Body)
	body := postBody{}
	if err := decoder.Decode(&body); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if len(body.MessageBody) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	splitRequest := strings.Split(body.MessageBody, " ")

	for index, element := range splitRequest {
		for _, badword := range badwords {
			lowerElement := strings.ToLower(element)
			if lowerElement == badword {
				splitRequest[index] = "****"
				break
			}
		}
	}

	//Cleaned Chirp Body String
	cleanedChirpBodyString := strings.Join(splitRequest, " ")

	log.Printf("Saving to file")

	newChirp, err := db.CreateChirp(cleanedChirpBodyString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, newChirp)
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

func (db *DB) handlerReturnChirps(w http.ResponseWriter, r *http.Request) {
	log.Printf("Retrieving Chirps")
	dbChirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	log.Printf("Creating Chirp")
	currentDBStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, fmt.Errorf("Create Chirp: %v", err)
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

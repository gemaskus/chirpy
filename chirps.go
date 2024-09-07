package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

	type success struct {
		MessageSuccess bool   `json:"valid"`
		CleanedBody    string `json:"cleaned_body"`
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

	newChirp := db.CreateChirp(cleanedChirpBodyString)

	respondWithJSON(w, http.StatusOK, newChirp)
}

func NewDB(filePath string) (*DB, error) {
	db := DB{
		path: filePath,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()

	return &db, err
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	if err != nil {
		if os.IsNotExist(err) {
			//file doesn't exist, so we much create the file.
			split := strings.Split(db.path, "/")
			fileName := split[len(split)-1]
			if err := os.WriteFile(fileName, []byte(""), 0666); err != nil {
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

}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	currentDBStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, fmt.Errorf("Create Chirp: %v", err)
	}

	return Chirp{}, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	return []Chirp{}, nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	split := strings.Split(db.path, "/")
	fileName := split[len(split)-1]

	dbStructure := DBStructure{}
	fileContents, err := os.ReadFile(fileName)

	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, fmt.Errorf("loadDB: Cannot read from given file: %s with error %v", fileName, err)
	}

	if err := json.Unmarshal(fileContents, &dbStructure); err != nil {
		return dbStructure, fmt.Errorf("loadDB: Cannot unmarshal the data read from the file: %s with error: %v".fileName, err)
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
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

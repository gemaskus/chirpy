package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type postBody struct {
	MessageBody string `json:"body"`
	EmailBody   string `json:"email"`
	PwBody      string `json:"password"`
}

func (db *DB) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
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

func (db *DB) handlerPostUsers(w http.ResponseWriter, r *http.Request) {
	body := postBody{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		respondWithError(w, http.StatusInternalServerError, "handlePostUser: Could not decode the message body")
		return
	}

	log.Printf("Saving new user to file")

	newUser, err := db.CreateUser(body.EmailBody, body.PwBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, newUser)
}

func (db *DB) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	body := postBody{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		respondWithError(w, http.StatusInternalServerError, "handlePostLogin: Could not decode message body")
	}

	users, err := db.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "handlePostLogin: Could not get users from database")
	}

	foundID := -1
	for index, user := range users {
		if body.EmailBody == user.Email {
			foundID = index
			break
		}
	}

	if foundID == -1 {
		respondWithError(w, http.StatusNotFound, "handlePostLogin: Could not find user")
		return
	}

	checkPW, err := checkPassword(body.PwBody, users[foundID].Password)
	if !checkPW {
		respondWithError(w, http.StatusUnauthorized, "handlePostLogin: Invalid Password")
		return
	} else {
		respondWithJSON(w, http.StatusOK, users[foundID])
	}
}

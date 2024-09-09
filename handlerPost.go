package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type postBody struct {
	MessageBody string `json:"body"`
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
	body := postBody{
		MessageBody: ""
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder(&body); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode the message body")
	}
	
}











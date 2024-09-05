package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
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

	newString := strings.Join(splitRequest, " ")

	log.Printf(body.MessageBody)

	respondWithJSON(w, http.StatusOK, success{
		MessageSuccess: true,
		CleanedBody:    newString,
	})
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

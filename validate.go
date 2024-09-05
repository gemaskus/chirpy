package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type postBody struct {
		MessageBody string `json:"body"`
	}

	type generalError struct {
		MessageError string `json:"error"`
	}

	type success struct {
		MessageSuccess bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	body := postBody{}
	if err := decoder.Decode(&body); err != nil {
		log.Printf("error decoding the body %v", err)
		errorBody := generalError{
			MessageError: "Something went wrong",
		}
		dat, _ := json.Marshal(errorBody)

		w.Header().Set("Content-Stype", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	if len(body.MessageBody) > 140 {
		errorBody := generalError{
			MessageError: "Chirp is too long",
		}

		dat, _ := json.Marshal(errorBody)
		w.Header().Set("Content-Stype", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	log.Printf(body.MessageBody)

	successbody := success{
		MessageSuccess: true,
	}

	dat, err := json.Marshal(successbody)

	if err != nil {
		log.Printf("error marshalling the valid body")
		errorBody := generalError{
			MessageError: "Something went wrong",
		}
		dat, _ := json.Marshal(errorBody)

		w.Header().Set("Content-Stype", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

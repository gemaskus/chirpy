package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
)

func (db *DB) handlerReturnChirps(w http.ResponseWriter, r *http.Request) {
	log.Printf("Retrieving Chirps")
	dbChirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}
	requestID, err := strconv.ParseInt(r.PathValue("chirpID"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse chirp ID")
		return
	}

	if requestID > int64(len(chirps)) {
		respondWithError(w, http.StatusNotFound, "Chirp ID could not be found")
		return
	}

	log.Printf("Requested Chirp ID: %d", requestID)

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	requestChirp := chirps[requestID-1]
	respondWithJSON(w, http.StatusOK, requestChirp)
}

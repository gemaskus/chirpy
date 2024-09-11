package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."
	const ChirpDBFileName = "ChirpsDB.json"
	const UserDBFileName = "Users.json"

	chirpDBFilePath := filepathRoot + "/" + ChirpDBFileName
	usersDBFilePath := filepathRoot + "/" + UserDBFileName

	config := &apiConfig{
		FileserverHits: 0,
	}

	chirpDB, err := NewDB(chirpDBFilePath)

	if err != nil {
		return //crash out for now
	}

	userDB, err := NewDB(usersDBFilePath)

	mux := http.NewServeMux()
	mux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer((http.Dir(filepathRoot))))))
	mux.HandleFunc("GET /api/healthz", config.handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", config.handlerHits)
	mux.HandleFunc("GET /api/reset", config.handlerReset)
	mux.HandleFunc("POST /api/chirps", chirpDB.handlerValidateChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", chirpDB.handlerReturnChirps)
	mux.HandleFunc("POST /api/users", userDB.handlerPostUsers)
	mux.HandleFunc("POST /api/login", userDB.handlePostLogin)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	FileserverHits int
	currentChirpID int
}

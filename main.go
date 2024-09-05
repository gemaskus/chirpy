package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	config := &apiConfig{
		FileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer((http.Dir(filepathRoot))))))
	mux.HandleFunc("GET /api/healthz", config.handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", config.handlerHits)
	mux.HandleFunc("GET /api/reset", config.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	FileserverHits int
}

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	config := &apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer((http.Dir(filepathRoot))))))
	mux.HandleFunc("GET /api/healthz", config.handlerReadiness)
	mux.HandleFunc("GET /api/metrics", config.handlerHits)
	mux.HandleFunc("GET /api/reset", config.handlerReset)
	//mux.HandleFunc("POST /healthz", handlerHealthzPost)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hits: %d\n", cfg.fileserverHits)
}

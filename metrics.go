package main

import (
	"html/template"
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits++

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, cfg); err != nil {
		log.Printf("Error executing template: %v", err)

		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

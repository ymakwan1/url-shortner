package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/ymakwan1/url-shortener/backend/database"
	"github.com/ymakwan1/url-shortener/backend/handlers"
	"github.com/ymakwan1/url-shortener/backend/middleware"
)

func main() {
	tokenBucket := middleware.NewTokenBucket(10, time.Minute)

	http.Handle("/", tokenBucket.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateShortURL(w, r, database.DB)
		case http.MethodGet:
			handlers.GetOriginalURL(w, r, database.DB)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	} else {
		log.Println("Server started")
	}
}

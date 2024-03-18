package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ymakwan1/url-shortener/backend/database"
	"github.com/ymakwan1/url-shortener/backend/handlers"
	"github.com/ymakwan1/url-shortener/backend/middleware"
)

func main() {

	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Llongfile)
	loggerError := log.New(os.Stdout, "ERROR: ", log.LstdFlags|log.Llongfile)
	logger.Print("Started")
	tokenBucket := middleware.NewTokenBucket(10, time.Minute)

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	http.Handle("/", corsMiddleware(tokenBucket.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			logger.Print("POST")
			handlers.CreateShortURL(w, r, database.DB)
		case http.MethodGet:
			logger.Print("GET")
			handlers.GetOriginalURL(w, r, database.DB)
		default:
			loggerError.Printf("Method not allowed (%d)", http.StatusMethodNotAllowed)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Start the HTTP server
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		loggerError.Printf("Error starting server: %v", err)
	} else {
		logger.Print("Server started")
	}
}

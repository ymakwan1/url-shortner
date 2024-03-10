package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ymakwan1/url-shortener/backend/jsonhandling"
	"github.com/ymakwan1/url-shortener/backend/validator"
)

type ShortURL struct {
	Key      string `json:"key"`
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
}

var ShortURLs = make(map[string]string)

func CreateShortURL(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		jsonhandling.Error(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonhandling.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if !validator.IsValidURL(req.URL) {
		jsonhandling.Error(w, http.StatusBadRequest, "Invalid URL")
	}

	for key, longURL := range ShortURLs {
		if longURL == req.URL {
			resp := ShortURL{
				Key:      key,
				LongURL:  longURL,
				ShortURL: "http://localhost/" + key,
			}
			jsonhandling.Response(w, http.StatusOK, resp)
			return
		}
	}

	key := generateKey(len(ShortURLs) + 1)
	ShortURLs[key] = req.URL
	resp := ShortURL{
		Key:      key,
		LongURL:  req.URL,
		ShortURL: "http://localhost/" + key,
	}
	jsonhandling.Response(w, http.StatusCreated, resp)
}

func GetOriginalURL(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func generateKey(i int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	key := make([]byte, i)
	for v := range key {
		key[v] = letters[v%len(letters)]
	}
	return string(key)
}

package handlers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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

	var existingKey string
	err := db.QueryRow("SELECT key FROM shortend_urls where long_url = $1", req.URL).Scan(&existingKey)

	if err == nil {
		resp := ShortURL{
			Key:      existingKey,
			LongURL:  req.URL,
			ShortURL: "http://localhost/" + existingKey,
		}
		jsonhandling.Response(w, http.StatusOK, resp)
		return
	} else if err != sql.ErrNoRows {
		jsonhandling.Error(w, http.StatusInternalServerError, "Failed to check existing URL")
		return
	}

	key := generateKey(req.URL)

	_, err = db.Exec("INSERT INTO shortened_urls (key, long_url) VALUES ($1, $2)", key, req.URL)
	if err != nil {
		jsonhandling.Error(w, http.StatusInternalServerError, "Failed to create shortened URL")
		return
	}

	resp := ShortURL{
		Key:      existingKey,
		LongURL:  req.URL,
		ShortURL: "http://localhost/" + existingKey,
	}

	jsonhandling.Response(w, http.StatusCreated, resp)
}

func GetOriginalURL(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[1:]

	var originalURL string
	err := db.QueryRow("SELECT long_url FROM shortend_urls where key = $1", key).Scan(&originalURL)

	if err != nil {
		jsonhandling.Error(w, http.StatusNotFound, "URL not found")
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func generateKey(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:6]
}

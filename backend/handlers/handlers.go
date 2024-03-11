package handlers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ymakwan1/url-shortener/backend/jsonhandling"
	lrucache "github.com/ymakwan1/url-shortener/backend/lru_cache"
	"github.com/ymakwan1/url-shortener/backend/validator"
)

type ShortURL struct {
	Key      string `json:"key"`
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
}

var cacheMutex sync.Mutex

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
		return
	}

	// Generate new key
	key := generateKey(req.URL)

	// Insert the URL into the database
	_, err := db.Exec("INSERT INTO shortened_urls (key, long_url) VALUES ($1, $2)", key, req.URL)
	if err != nil {
		jsonhandling.Error(w, http.StatusInternalServerError, "Failed to create shortened URL")
		return
	}

	cacheMutex.Lock()
	lrucache.LruCache.Set(key, req.URL)
	cacheMutex.Unlock()

	resp := ShortURL{
		Key:      key,
		LongURL:  req.URL,
		ShortURL: "http://localhost/" + key,
	}

	jsonhandling.Response(w, http.StatusCreated, resp)
}

func GetOriginalURL(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[1:]

	cacheMutex.Lock()
	longURL, ok := lrucache.LruCache.Get(key)
	cacheMutex.Unlock()

	if ok {
		http.Redirect(w, r, longURL, http.StatusFound)
		return
	}

	var originalURL string
	err := db.QueryRow("SELECT long_url FROM shortened_urls where key = $1", key).Scan(&originalURL)

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

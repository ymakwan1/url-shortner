package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ymakwan1/url-shortener/backend/handlers"
)

func TestCreateShortURL(t *testing.T) {
	router := http.NewServeMux()
	router.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateShortURL(w, r)
	})

	t.Run("CreateShortURL", func(t *testing.T) {
		reqBody := `{"url":"http://example.com"}`
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		var respBody handlers.ShortURL
		if err := json.NewDecoder(w.Body).Decode(&respBody); err != nil {
			t.Fatal(err)
		}

		if respBody.ShortURL == "" {
			t.Error("Expected non-empty ShortURL field in response")
		}
	})

	t.Run("InvalidURL", func(t *testing.T) {
		reqBody := `{"url":"example"}`
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("BindJSONFailure", func(t *testing.T) {
		reqBody := `{"invalid_json"`
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("ExistingLongURL", func(t *testing.T) {
		// Add a long URL to the ShortURLs map
		handlers.ShortURLs["existingKey"] = "http://examples.com"

		reqBody := `{"url":"http://examples.com"}`
		req, err := http.NewRequest("POST", "/shorten", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Decode the response body
		var respBody handlers.ShortURL
		if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
			t.Fatal(err)
		}

		// Check if the returned short URL matches the existing one
		if respBody.ShortURL != "http://localhost/existingKey" {
			t.Errorf("Expected short URL %s, got %s", "http://localhost/existingKey", respBody.ShortURL)
		}
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/shorten", nil) // Using POST method
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestGetOriginalURL(t *testing.T) {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOriginalURL(w, r)
	})

	t.Run("GetOriginalURL", func(t *testing.T) {
		handlers.ShortURLs["exampleKey"] = "http://example.com"

		req, err := http.NewRequest("GET", "/exampleKey", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusFound {
			t.Errorf("Expected status code %d, got %d", http.StatusFound, w.Code)
		}
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/exampleKey", nil) // Using POST method
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

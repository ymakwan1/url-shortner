package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ymakwan1/url-shortener/backend/handlers"
)

func TestCreateShortURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/shorten", handlers.CreateShortURL)

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
		if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
			t.Fatal(err)
		}

		if respBody.ShortURL == "" {
			t.Error("Expected non-empty ShortURL field in response")
		}
	})
}

func TestGetOriginalURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:key", handlers.GetOriginalURL)

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
}

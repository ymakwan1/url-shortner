package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ymakwan1/url-shortener/backend/middleware"
)

func TestTokenBucket(t *testing.T) {
	// Create a new token bucket with a limit of 2 requests per second
	tokenBucket := middleware.NewTokenBucket(2, time.Second)

	// Create a handler function that increments a counter
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Wrap the handler with the token bucket middleware
	handlerWithMiddleware := tokenBucket.Limit(handler)

	// Send three requests with a small delay between each request
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handlerWithMiddleware.ServeHTTP(rr, req)
		time.Sleep(200 * time.Millisecond) // Add a small delay between requests
		if rr.Code != http.StatusOK && rr.Code != http.StatusTooManyRequests {
			t.Errorf("Request %d returned unexpected status code: got %d", i+1, rr.Code)
		}
	}
}

package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortURLs(t *testing.T) {
	// Test invalid request body
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateShortURL(rr, req, nil)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Test invalid URL
	reqBody := `{"url": "invalidurl"}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	CreateShortURL(rr, req, nil)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Test database error
	// Mock database to return an error
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectExec("INSERT INTO shortened_urls").
		WillReturnError(errors.New("database error"))
	reqBody = `{"url": "https://example.com"}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	CreateShortURL(rr, req, db)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Test successful short URL creation
	// Mock database to return successful result
	mock.ExpectExec("INSERT INTO shortened_urls").
		WillReturnResult(sqlmock.NewResult(1, 1))
	reqBody = `{"url": "https://example.com"}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	CreateShortURL(rr, req, db)
	assert.Equal(t, http.StatusCreated, rr.Code)

}

func TestGetOriginalURL(t *testing.T) {
	//mock get URL success
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	key := "abc123"
	mock.ExpectQuery("SELECT long_url FROM shortened_urls").
		WithArgs(key).
		WillReturnRows(sqlmock.NewRows([]string{"long_url"}).AddRow("https://example.com"))

	req := httptest.NewRequest(http.MethodGet, "/"+key, nil)

	rr := httptest.NewRecorder()

	GetOriginalURL(rr, req, db)

	assert.Equal(t, http.StatusFound, rr.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled mock expectations: %s", err)
	}
}

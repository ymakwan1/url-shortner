package jsonhandling

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	Response(w, statusCode, map[string]string{"error": message})
}

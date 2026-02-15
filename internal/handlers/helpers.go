package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/phil-bot/rsyslog-rest-api/internal/models"
)

// respondJSON sends a JSON response with proper headers
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// respondError sends a JSON error response
func respondError(w http.ResponseWriter, status int, err *models.APIError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	
	if encodeErr := json.NewEncoder(w).Encode(err); encodeErr != nil {
		log.Printf("Error encoding error response: %v", encodeErr)
	}
}

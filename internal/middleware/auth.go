package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/phil-bot/rsyslog-rest-api/internal/models"
)

// Auth returns a middleware that validates API key authentication
func Auth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth if no API key is configured
			if apiKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check API key in header
			providedKey := r.Header.Get("X-API-Key")
			if providedKey == "" {
				respondError(w, http.StatusUnauthorized, models.NewAPIError(
					models.ErrCodeUnauthorized,
					"Missing API key").
					WithDetails("Provide X-API-Key header"))
				return
			}

			if providedKey != apiKey {
				respondError(w, http.StatusUnauthorized, models.NewAPIError(
					models.ErrCodeUnauthorized,
					"Invalid API key"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// respondError sends a JSON error response
func respondError(w http.ResponseWriter, status int, err *models.APIError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

package handlers

import (
	"net/http"
	"time"

	"github.com/phil-bot/rsyslog-rest-api/internal/database"
	"github.com/phil-bot/rsyslog-rest-api/internal/models"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db      *database.DB
	version string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB, version string) *HealthHandler {
	return &HealthHandler{
		db:      db,
		version: version,
	}
}

// ServeHTTP handles the health check endpoint
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, 
			models.NewAPIError("METHOD_NOT_ALLOWED", "Only GET method is allowed"))
		return
	}

	// Check database health
	dbStatus := "connected"
	if err := h.db.Health(); err != nil {
		dbStatus = "disconnected"
		respondJSON(w, http.StatusServiceUnavailable, models.HealthResponse{
			Status:    "unhealthy",
			Database:  dbStatus,
			Version:   h.version,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	respondJSON(w, http.StatusOK, models.HealthResponse{
		Status:    "healthy",
		Database:  dbStatus,
		Version:   h.version,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/phil-bot/rsyslox/internal/config"
	"github.com/phil-bot/rsyslox/internal/database"
)

// ServerDefaults contains all server-configured default values sent to the
// frontend on every /health response. The frontend applies these to new
// browser sessions that have not yet stored their own preferences.
type ServerDefaults struct {
	TimeRange           string `json:"time_range"`
	AutoRefreshInterval int    `json:"auto_refresh_interval"`
	Language            string `json:"language"`
	FontSize            string `json:"font_size"`
	TimeFormat          string `json:"time_format"`
}

// HealthResponse is the JSON body returned by GET /health.
type HealthResponse struct {
	Status    string          `json:"status"`
	Database  string          `json:"database,omitempty"`
	Version   string          `json:"version"`
	Timestamp string          `json:"timestamp"`
	SetupMode bool            `json:"setup_mode,omitempty"`
	Defaults  *ServerDefaults `json:"defaults,omitempty"`
}

// HealthHandler handles GET /health.
type HealthHandler struct {
	db      *database.DB
	version string
	cfg     *config.Config
}

// NewHealthHandler creates a HealthHandler.
func NewHealthHandler(db *database.DB, version string, cfg *config.Config) *HealthHandler {
	return &HealthHandler{db: db, version: version, cfg: cfg}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	cfgPath := config.ActiveConfigPath()
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusOK)
		if encErr := json.NewEncoder(w).Encode(HealthResponse{
			Status:    "setup",
			Version:   h.version,
			Timestamp: time.Now().Format(time.RFC3339),
			SetupMode: true,
		}); encErr != nil {
			log.Printf("health: encode error: %v", encErr)
		}
		return
	}

	if h.db == nil {
		w.WriteHeader(http.StatusOK)
		if encErr := json.NewEncoder(w).Encode(HealthResponse{
			Status:    "pending_restart",
			Version:   h.version,
			Timestamp: time.Now().Format(time.RFC3339),
		}); encErr != nil {
			log.Printf("health: encode error: %v", encErr)
		}
		return
	}

	if err := h.db.Health(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		if encErr := json.NewEncoder(w).Encode(HealthResponse{
			Status:    "unhealthy",
			Database:  "disconnected",
			Version:   h.version,
			Timestamp: time.Now().Format(time.RFC3339),
		}); encErr != nil {
			log.Printf("health: encode error: %v", encErr)
		}
		return
	}

	var defaults *ServerDefaults
	if h.cfg != nil {
		defaults = &ServerDefaults{
			TimeRange:           h.cfg.Server.DefaultTimeRange,
			AutoRefreshInterval: h.cfg.Server.AutoRefreshInterval,
			Language:            h.cfg.Server.DefaultLanguage,
			FontSize:            h.cfg.Server.DefaultFontSize,
			TimeFormat:          h.cfg.Server.DefaultTimeFormat,
		}
	}

	w.WriteHeader(http.StatusOK)
	if encErr := json.NewEncoder(w).Encode(HealthResponse{
		Status:    "healthy",
		Database:  "connected",
		Version:   h.version,
		Timestamp: time.Now().Format(time.RFC3339),
		Defaults:  defaults,
	}); encErr != nil {
		log.Printf("health: encode error: %v", encErr)
	}
}

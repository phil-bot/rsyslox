package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/phil-bot/rsyslog-rest-api/internal/models"
)

// RootHandler handles root endpoint requests
type RootHandler struct {
	installPath string
	version     string
}

// NewRootHandler creates a new root handler
func NewRootHandler(installPath, version string) *RootHandler {
	return &RootHandler{
		installPath: installPath,
		version:     version,
	}
}

// ServeHTTP handles the root endpoint
func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only handle root path
	if r.URL.Path != "/" {
		respondError(w, http.StatusNotFound, 
			models.NewAPIError(models.ErrCodeNotFound, "Endpoint not found"))
		return
	}

	// Only allow GET requests
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, 
			models.NewAPIError("METHOD_NOT_ALLOWED", "Only GET method is allowed"))
		return
	}

	// Try to serve index.html if it exists
	indexPath := filepath.Join(h.installPath, "index.html")
	if _, err := os.Stat(indexPath); err == nil {
		http.ServeFile(w, r, indexPath)
		return
	}

	// Otherwise show API info
	respondJSON(w, http.StatusOK, models.RootResponse{
		Name:    "rsyslog REST API",
		Version: h.version,
		Endpoints: map[string]string{
			"logs":   "/logs?limit=10&Priority=3",
			"meta":   "/meta or /meta/{column}",
			"health": "/health",
		},
	})
}

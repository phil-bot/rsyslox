package admin

import (
	"net/http"

	"github.com/phil-bot/rsyslox/internal/cleanup"
	"github.com/phil-bot/rsyslox/internal/models"
)

// CleanupStatusHandler handles GET /api/admin/cleanup/status.
// Returns the current partition mode, health, error (if any) and partition list.
type CleanupStatusHandler struct {
	cleaner *cleanup.Cleaner
}

// NewCleanupStatusHandler creates a new CleanupStatusHandler.
func NewCleanupStatusHandler(cleaner *cleanup.Cleaner) *CleanupStatusHandler {
	return &CleanupStatusHandler{cleaner: cleaner}
}

func (h *CleanupStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed,
			models.NewAPIError("METHOD_NOT_ALLOWED", "Only GET is allowed"))
		return
	}

	if h.cleaner == nil {
		respondJSON(w, http.StatusOK, cleanup.Status{
			Mode:       cleanup.ModeUnknown,
			Healthy:    false,
			Error:      "Cleanup service not initialised",
			Partitions: []cleanup.Partition{},
		})
		return
	}

	respondJSON(w, http.StatusOK, h.cleaner.Status())
}

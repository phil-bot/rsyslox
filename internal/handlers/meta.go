package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/phil-bot/rsyslog-rest-api/internal/database"
	"github.com/phil-bot/rsyslog-rest-api/internal/filters"
	"github.com/phil-bot/rsyslog-rest-api/internal/models"
)

// MetaHandler handles metadata requests
type MetaHandler struct {
	db *database.DB
}

// NewMetaHandler creates a new meta handler
func NewMetaHandler(db *database.DB) *MetaHandler {
	return &MetaHandler{
		db: db,
	}
}

// ServeHTTP handles the /meta endpoint
func (h *MetaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, 
			models.NewAPIError("METHOD_NOT_ALLOWED", "Only GET method is allowed"))
		return
	}

	// Check if this is the list endpoint
	if r.URL.Path == "/meta" {
		h.handleList(w, r)
		return
	}

	// Extract column from path
	column := strings.TrimPrefix(r.URL.Path, "/meta/")
	column = strings.TrimSpace(column)

	if column == "" {
		respondError(w, http.StatusBadRequest, 
			models.NewAPIError(models.ErrCodeInvalidColumn, "Column name required"))
		return
	}

	h.handleColumnValues(w, r, column)
}

// handleList returns the list of available columns
func (h *MetaHandler) handleList(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, models.MetaResponse{
		AvailableColumns: h.db.AvailableColumns,
		Usage:            "GET /meta/{column} to get distinct values for a column",
	})
}

// handleColumnValues returns distinct values for a specific column
func (h *MetaHandler) handleColumnValues(w http.ResponseWriter, r *http.Request, column string) {
	// Validate column exists
	if !h.db.IsValidColumn(column) {
		respondError(w, http.StatusBadRequest, 
			models.NewAPIError(models.ErrCodeInvalidColumn, 
				"Invalid column: "+column).
				WithDetails("Available columns: "+strings.Join(h.db.AvailableColumns, ", ")))
		return
	}

	query := r.URL.Query()

	// Validate date range
	startDate, endDate, err := filters.ValidateDateRange(
		query.Get("start_date"),
		query.Get("end_date"),
	)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			respondError(w, http.StatusBadRequest, apiErr)
		} else {
			respondError(w, http.StatusBadRequest, 
				models.NewAPIError(models.ErrCodeInvalidParameter, err.Error()))
		}
		return
	}

	// Validate priorities
	priorities, err := filters.ValidatePriorities(query["Priority"])
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			respondError(w, http.StatusBadRequest, apiErr)
		} else {
			respondError(w, http.StatusBadRequest, 
				models.NewAPIError(models.ErrCodeInvalidPriority, err.Error()))
		}
		return
	}

	// Validate facilities
	facilities, err := filters.ValidateFacilities(query["Facility"])
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			respondError(w, http.StatusBadRequest, apiErr)
		} else {
			respondError(w, http.StatusBadRequest, 
				models.NewAPIError(models.ErrCodeInvalidFacility, err.Error()))
		}
		return
	}

	// Validate messages
	messages, err := filters.ValidateMessages(query["Message"])
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			respondError(w, http.StatusBadRequest, apiErr)
		} else {
			respondError(w, http.StatusBadRequest, 
				models.NewAPIError(models.ErrCodeInvalidParameter, err.Error()))
		}
		return
	}

	// Build filter query
	builder := filters.New()
	builder.AddDateRange(startDate, endDate)
	builder.AddStringMultiValue("FromHost", query["FromHost"])
	builder.AddIntMultiValue("Priority", priorities)
	builder.AddIntMultiValue("Facility", facilities)
	builder.AddMessageSearch(messages)
	builder.AddStringMultiValue("SysLogTag", query["SysLogTag"])

	whereClause, args := builder.Build()

	// Query distinct values
	values, err := h.db.QueryDistinctValues(column, whereClause, args)
	if err != nil {
		log.Printf("Meta query error: %v", err)
		respondError(w, http.StatusInternalServerError, 
			models.NewAPIError(models.ErrCodeDatabaseError, "Failed to query metadata"))
		return
	}

	respondJSON(w, http.StatusOK, values)
}

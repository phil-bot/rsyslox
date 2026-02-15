package handlers

import (
	"log"
	"net/http"

	"github.com/phil-bot/rsyslog-rest-api/internal/database"
	"github.com/phil-bot/rsyslog-rest-api/internal/filters"
	"github.com/phil-bot/rsyslog-rest-api/internal/models"
)

// LogsHandler handles log retrieval requests
type LogsHandler struct {
	db *database.DB
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler(db *database.DB) *LogsHandler {
	return &LogsHandler{
		db: db,
	}
}

// ServeHTTP handles the /logs endpoint
func (h *LogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, 
			models.NewAPIError("METHOD_NOT_ALLOWED", "Only GET method is allowed"))
		return
	}

	query := r.URL.Query()

	// Validate pagination
	limit, offset, err := filters.ValidatePagination(
		query.Get("limit"),
		query.Get("offset"),
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

	// Count total matching entries
	total, err := h.db.CountLogs(whereClause, args)
	if err != nil {
		log.Printf("Count query error: %v", err)
		respondError(w, http.StatusInternalServerError, 
			models.NewAPIError(models.ErrCodeDatabaseError, "Failed to count logs"))
		return
	}

	// Query entries with pagination
	entries, err := h.db.QueryLogs(whereClause, args, limit, offset)
	if err != nil {
		log.Printf("Query error: %v", err)
		respondError(w, http.StatusInternalServerError, 
			models.NewAPIError(models.ErrCodeDatabaseError, "Failed to query logs"))
		return
	}

	respondJSON(w, http.StatusOK, models.LogsResponse{
		Total:  total,
		Offset: offset,
		Limit:  limit,
		Rows:   entries,
	})
}

package models

// LogsResponse represents the response structure for /logs endpoint
type LogsResponse struct {
	Total  int        `json:"total"`
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	Rows   []LogEntry `json:"rows"`
}

// MetaValue represents a meta value with optional label (for Priority/Facility)
type MetaValue struct {
	Val   int    `json:"val"`
	Label string `json:"label"`
}

// MetaResponse represents the response for /meta endpoint
type MetaResponse struct {
	AvailableColumns []string `json:"available_columns"`
	Usage            string   `json:"usage"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// RootResponse represents the root endpoint response
type RootResponse struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Endpoints map[string]string `json:"endpoints"`
}

// APIError represents a structured API error response (NEW in v0.2.3+)
type APIError struct {
	Code    string `json:"code"`              // Error code (e.g., "INVALID_PARAMETER")
	Message string `json:"message"`           // Human-readable message
	Details string `json:"details,omitempty"` // Optional additional details
	Field   string `json:"field,omitempty"`   // Optional field name for validation errors
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

// Common API error codes
const (
	ErrCodeInvalidParameter = "INVALID_PARAMETER"
	ErrCodeMissingParameter = "MISSING_PARAMETER"
	ErrCodeDatabaseError    = "DATABASE_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeInvalidColumn    = "INVALID_COLUMN"
	ErrCodeInvalidDateRange = "INVALID_DATE_RANGE"
	ErrCodeInvalidPriority  = "INVALID_PRIORITY"
	ErrCodeInvalidFacility  = "INVALID_FACILITY"
)

// NewAPIError creates a new APIError
func NewAPIError(code, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewValidationError creates a new validation error for a specific field
func NewValidationError(field, message string) *APIError {
	return &APIError{
		Code:    ErrCodeInvalidParameter,
		Field:   field,
		Message: message,
	}
}

// WithDetails adds details to an API error
func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

// WithField adds a field name to an API error
func (e *APIError) WithField(field string) *APIError {
	e.Field = field
	return e
}

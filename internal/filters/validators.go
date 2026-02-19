package filters

import (
	"fmt"
	"strconv"
	"time"

	"github.com/phil-bot/rsyslox/internal/models"
)

// ValidateSeverities validates severity values (0-7) and returns them as integers.
func ValidateSeverities(severities []string) ([]int, error) {
	if len(severities) == 0 {
		return nil, nil
	}

	var valid []int
	for _, sStr := range severities {
		s, err := strconv.Atoi(sStr)
		if err != nil {
			return nil, models.NewAPIError(models.ErrCodeInvalidSeverity,
				fmt.Sprintf("'%s' is not a valid integer", sStr)).
				WithField("Severity")
		}
		if !models.IsValidSeverity(s) {
			return nil, models.NewAPIError(models.ErrCodeInvalidSeverity,
				fmt.Sprintf("value %d is out of range (must be 0-7)", s)).
				WithField("Severity").
				WithDetails("See RFC-5424 for valid severity levels")
		}
		valid = append(valid, s)
	}

	return valid, nil
}

// ValidatePriorities is a deprecated alias for ValidateSeverities.
// The query parameter ?Priority= is accepted as an alias for ?Severity=.
func ValidatePriorities(priorities []string) ([]int, error) {
	return ValidateSeverities(priorities)
}

// ValidateFacilities validates facility values and returns them as integers.
func ValidateFacilities(facilities []string) ([]int, error) {
	if len(facilities) == 0 {
		return nil, nil
	}

	var validFacilities []int
	for _, fStr := range facilities {
		f, err := strconv.Atoi(fStr)
		if err != nil {
			return nil, models.NewAPIError(models.ErrCodeInvalidFacility,
				fmt.Sprintf("'%s' is not a valid integer", fStr)).
				WithField("Facility")
		}
		if !models.IsValidFacility(f) {
			return nil, models.NewAPIError(models.ErrCodeInvalidFacility,
				fmt.Sprintf("value %d is out of range (must be 0-23)", f)).
				WithField("Facility").
				WithDetails("See RFC-5424 for valid facility codes")
		}
		validFacilities = append(validFacilities, f)
	}

	return validFacilities, nil
}

// ValidateMessages validates message search terms.
func ValidateMessages(messages []string) ([]string, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	var validMessages []string
	for _, msg := range messages {
		if len(msg) < 2 {
			return nil, models.NewAPIError(models.ErrCodeInvalidParameter,
				"search term must be at least 2 characters long").
				WithField("Message").
				WithDetails(fmt.Sprintf("Term '%s' is too short", msg))
		}
		validMessages = append(validMessages, msg)
	}

	return validMessages, nil
}

// ValidateDateRange validates and parses date range parameters.
func ValidateDateRange(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	var startDate, endDate time.Time
	var err error

	// Parse start date
	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, models.NewAPIError(models.ErrCodeInvalidParameter,
				"invalid format").
				WithField("start_date").
				WithDetails("Expected ISO 8601/RFC3339 format (e.g., 2025-02-15T10:00:00Z)")
		}
	} else {
		startDate = time.Now().Add(-24 * time.Hour)
	}

	// Parse end date
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, models.NewAPIError(models.ErrCodeInvalidParameter,
				"invalid format").
				WithField("end_date").
				WithDetails("Expected ISO 8601/RFC3339 format (e.g., 2025-02-15T10:00:00Z)")
		}
	} else {
		endDate = time.Now()
	}

	// Validate range
	if startDate.After(endDate) {
		return time.Time{}, time.Time{}, models.NewAPIError(
			models.ErrCodeInvalidDateRange,
			"start_date cannot be after end_date")
	}

	if endDate.Sub(startDate) > 90*24*time.Hour {
		return time.Time{}, time.Time{}, models.NewAPIError(
			models.ErrCodeInvalidDateRange,
			"date range cannot exceed 90 days").
			WithDetails(fmt.Sprintf("Requested range: %.1f days", endDate.Sub(startDate).Hours()/24))
	}

	return startDate, endDate, nil
}

// ValidatePagination validates limit and offset parameters.
func ValidatePagination(limitStr, offsetStr string) (int, int, error) {
	const (
		defaultLimit = 10
		maxLimit     = 1000
	)

	// Parse offset
	offset := 0
	if offsetStr != "" {
		val, err := strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, models.NewAPIError(models.ErrCodeInvalidParameter,
				fmt.Sprintf("'%s' is not a valid integer", offsetStr)).
				WithField("offset")
		}
		if val < 0 {
			return 0, 0, models.NewAPIError(models.ErrCodeInvalidParameter,
				"must be non-negative").
				WithField("offset")
		}
		offset = val
	}

	// Parse limit
	limit := defaultLimit
	if limitStr != "" {
		val, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, models.NewAPIError(models.ErrCodeInvalidParameter,
				fmt.Sprintf("'%s' is not a valid integer", limitStr)).
				WithField("limit")
		}
		if val <= 0 {
			return 0, 0, models.NewAPIError(models.ErrCodeInvalidParameter,
				"must be greater than 0").
				WithField("limit")
		}
		if val > maxLimit {
			return 0, 0, models.NewAPIError(models.ErrCodeInvalidParameter,
				fmt.Sprintf("cannot exceed %d", maxLimit)).
				WithField("limit").
				WithDetails(fmt.Sprintf("Requested: %d", val))
		}
		limit = val
	}

	return limit, offset, nil
}

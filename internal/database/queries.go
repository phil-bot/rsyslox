package database

import (
	"fmt"

	"github.com/phil-bot/rsyslox/internal/models"
)

// QueryLogs executes a logs query with filters and pagination
func (db *DB) QueryLogs(whereClause string, args []interface{}, limit, offset int) ([]models.LogEntry, error) {
	query := fmt.Sprintf(`
		SELECT ID, CustomerID, ReceivedAt, DeviceReportedTime, Facility, Priority,
		       FromHost, Message, NTSeverity, Importance, EventSource, EventUser,
		       EventCategory, EventID, EventBinaryData, MaxAvailable, CurrUsage,
		       MinUsage, MaxUsage, InfoUnitID, SysLogTag, EventLogType,
		       GenericFileName, SystemID
		FROM SystemEvents
		WHERE %s
		ORDER BY ReceivedAt DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, limit, offset)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var entries []models.LogEntry
	for rows.Next() {
		var entry models.LogEntry
		if err := entry.ScanFromRows(rows); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	if entries == nil {
		entries = []models.LogEntry{}
	}

	return entries, nil
}

// CountLogs counts total matching entries for a query
func (db *DB) CountLogs(whereClause string, args []interface{}) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM SystemEvents WHERE %s", whereClause)

	var total int
	if err := db.QueryRow(query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count query failed: %v", err)
	}

	return total, nil
}

// QueryDistinctValues queries distinct values for a column with optional filters.
//
// "Severity" is a virtual column — it is computed from the Priority column using
// MOD 8, which works correctly for both legacy and modern rsyslog formats as well
// as mixed datasets.
func (db *DB) QueryDistinctValues(column, whereClause string, args []interface{}) (interface{}, error) {
	// "Severity" is a virtual column derived from Priority MOD 8
	if column == "Severity" {
		return db.queryDistinctSeverity(whereClause, args)
	}

	query := fmt.Sprintf(
		"SELECT DISTINCT %s FROM SystemEvents WHERE %s AND %s IS NOT NULL ORDER BY %s ASC",
		column, whereClause, column, column,
	)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("meta query failed: %v", err)
	}
	defer rows.Close()

	// Facility column returns values with RFC labels
	if column == "Facility" {
		return scanMetaFacilityValues(rows)
	}

	// Check if column is an integer type
	if db.isIntegerColumn(column) {
		return scanIntValues(rows)
	}

	// String columns
	return scanStringValues(rows)
}

// queryDistinctSeverity returns distinct Severity values with RFC labels.
// Uses Priority MOD 8 to work across legacy, modern, and mixed datasets.
func (db *DB) queryDistinctSeverity(whereClause string, args []interface{}) ([]models.MetaValue, error) {
	query := fmt.Sprintf(
		"SELECT DISTINCT Priority MOD 8 AS Severity FROM SystemEvents WHERE %s ORDER BY Severity ASC",
		whereClause,
	)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("severity meta query failed: %v", err)
	}
	defer rows.Close()

	var values []models.MetaValue
	for rows.Next() {
		var val int
		if err := rows.Scan(&val); err != nil {
			continue
		}
		values = append(values, models.MetaValue{
			Val:   val,
			Label: models.GetSeverityLabel(val),
		})
	}

	if values == nil {
		values = []models.MetaValue{}
	}
	return values, nil
}

// isIntegerColumn checks if a column is an integer type
func (db *DB) isIntegerColumn(column string) bool {
	intColumns := []string{
		"ID", "CustomerID", "NTSeverity", "Importance", "EventCategory", "EventID",
		"MaxAvailable", "CurrUsage", "MinUsage", "MaxUsage", "InfoUnitID", "SystemID",
		"Priority",
	}
	for _, col := range intColumns {
		if col == column {
			return true
		}
	}
	return false
}

// Helper functions for scanning different value types

func scanMetaFacilityValues(rows interface {
	Next() bool
	Scan(...interface{}) error
}) ([]models.MetaValue, error) {
	var values []models.MetaValue
	for rows.Next() {
		var val int
		if err := rows.Scan(&val); err != nil {
			continue
		}
		values = append(values, models.MetaValue{
			Val:   val,
			Label: models.GetFacilityLabel(val),
		})
	}

	if values == nil {
		values = []models.MetaValue{}
	}
	return values, nil
}

func scanIntValues(rows interface {
	Next() bool
	Scan(...interface{}) error
}) ([]int, error) {
	var values []int
	for rows.Next() {
		var val int
		if err := rows.Scan(&val); err != nil {
			continue
		}
		values = append(values, val)
	}

	if values == nil {
		values = []int{}
	}
	return values, nil
}

func scanStringValues(rows interface {
	Next() bool
	Scan(...interface{}) error
}) ([]string, error) {
	var values []string
	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			continue
		}
		if val != "" {
			values = append(values, val)
		}
	}

	if values == nil {
		values = []string{}
	}
	return values, nil
}

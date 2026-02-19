package models

import (
	"database/sql"
	"time"
)

// LogEntry represents a single log entry from the SystemEvents table.
//
// Priority/Severity handling:
// Older rsyslog versions (< 8.2204.0) stored the Severity (0-7) in the
// Priority column. Newer versions store the RFC PRI value (Facility*8 + Severity).
// ScanFromRows detects the format per row and always exposes:
//   - Priority: RFC PRI value (Facility*8 + Severity)
//   - Severity: RFC Severity value (0-7)
//   - Severity_Label: human-readable severity label
//   - Facility / Facility_Label: unchanged
type LogEntry struct {
	ID                 int        `json:"ID"`
	CustomerID         *int64     `json:"CustomerID"`
	ReceivedAt         time.Time  `json:"ReceivedAt"`
	DeviceReportedTime *time.Time `json:"DeviceReportedTime"`
	Facility           int        `json:"Facility"`
	FacilityLabel      string     `json:"Facility_Label"`
	Priority           int        `json:"Priority"`
	Severity           int        `json:"Severity"`
	SeverityLabel      string     `json:"Severity_Label"`
	FromHost           string     `json:"FromHost"`
	Message            string     `json:"Message"`
	NTSeverity         *int       `json:"NTSeverity"`
	Importance         *int       `json:"Importance"`
	EventSource        *string    `json:"EventSource"`
	EventUser          *string    `json:"EventUser"`
	EventCategory      *int       `json:"EventCategory"`
	EventID            *int       `json:"EventID"`
	EventBinaryData    *string    `json:"EventBinaryData"`
	MaxAvailable       *int       `json:"MaxAvailable"`
	CurrUsage          *int       `json:"CurrUsage"`
	MinUsage           *int       `json:"MinUsage"`
	MaxUsage           *int       `json:"MaxUsage"`
	InfoUnitID         *int       `json:"InfoUnitID"`
	SysLogTag          *string    `json:"SysLogTag"`
	EventLogType       *string    `json:"EventLogType"`
	GenericFileName    *string    `json:"GenericFileName"`
	SystemID           *int       `json:"SystemID"`
}

// ScanFromRows scans a database row into a LogEntry.
//
// The raw Priority column value is inspected per row to determine the storage format:
//   - rawPriority > 7  → modern format: Priority = PRI, Severity = Priority % 8
//   - rawPriority <= 7 → legacy format: Severity = Priority, Priority = Facility*8 + Severity
//
// This handles mixed datasets (after a rsyslog version upgrade) correctly.
func (e *LogEntry) ScanFromRows(rows *sql.Rows) error {
	var customerID, ntSeverity, importance, eventCategory, eventID sql.NullInt64
	var maxAvail, currUsage, minUsage, maxUsage, infoUnitID, systemID sql.NullInt64
	var deviceTime sql.NullTime
	var eventSource, eventUser, eventBinary, sysLogTag, eventLogType, genericFile sql.NullString

	var rawPriority int

	err := rows.Scan(
		&e.ID, &customerID, &e.ReceivedAt, &deviceTime,
		&e.Facility, &rawPriority, &e.FromHost, &e.Message,
		&ntSeverity, &importance, &eventSource, &eventUser,
		&eventCategory, &eventID, &eventBinary, &maxAvail, &currUsage,
		&minUsage, &maxUsage, &infoUnitID, &sysLogTag, &eventLogType,
		&genericFile, &systemID,
	)
	if err != nil {
		return err
	}

	// Per-row format detection: derive RFC-compliant Priority and Severity
	if rawPriority > 7 {
		// Modern rsyslog (>= 8.2204.0): Priority column already contains RFC PRI
		e.Priority = rawPriority
		e.Severity = rawPriority % 8
	} else {
		// Legacy rsyslog (< 8.2204.0): Priority column contains Severity (0-7)
		e.Severity = rawPriority
		e.Priority = e.Facility*8 + rawPriority
	}

	// RFC labels
	e.SeverityLabel = GetSeverityLabel(e.Severity)
	e.FacilityLabel = GetFacilityLabel(e.Facility)

	// Map nullable fields
	if customerID.Valid {
		e.CustomerID = &customerID.Int64
	}
	if deviceTime.Valid {
		e.DeviceReportedTime = &deviceTime.Time
	}
	if ntSeverity.Valid {
		val := int(ntSeverity.Int64)
		e.NTSeverity = &val
	}
	if importance.Valid {
		val := int(importance.Int64)
		e.Importance = &val
	}
	if eventSource.Valid {
		e.EventSource = &eventSource.String
	}
	if eventUser.Valid {
		e.EventUser = &eventUser.String
	}
	if eventCategory.Valid {
		val := int(eventCategory.Int64)
		e.EventCategory = &val
	}
	if eventID.Valid {
		val := int(eventID.Int64)
		e.EventID = &val
	}
	if eventBinary.Valid {
		e.EventBinaryData = &eventBinary.String
	}
	if maxAvail.Valid {
		val := int(maxAvail.Int64)
		e.MaxAvailable = &val
	}
	if currUsage.Valid {
		val := int(currUsage.Int64)
		e.CurrUsage = &val
	}
	if minUsage.Valid {
		val := int(minUsage.Int64)
		e.MinUsage = &val
	}
	if maxUsage.Valid {
		val := int(maxUsage.Int64)
		e.MaxUsage = &val
	}
	if infoUnitID.Valid {
		val := int(infoUnitID.Int64)
		e.InfoUnitID = &val
	}
	if sysLogTag.Valid {
		e.SysLogTag = &sysLogTag.String
	}
	if eventLogType.Valid {
		e.EventLogType = &eventLogType.String
	}
	if genericFile.Valid {
		e.GenericFileName = &genericFile.String
	}
	if systemID.Valid {
		val := int(systemID.Int64)
		e.SystemID = &val
	}

	return nil
}

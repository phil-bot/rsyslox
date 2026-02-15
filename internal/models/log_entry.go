package models

import (
	"database/sql"
	"time"
)

// LogEntry represents a single log entry from the SystemEvents table
type LogEntry struct {
	ID                 int        `json:"ID"`
	CustomerID         *int64     `json:"CustomerID"`
	ReceivedAt         time.Time  `json:"ReceivedAt"`
	DeviceReportedTime *time.Time `json:"DeviceReportedTime"`
	Facility           int        `json:"Facility"`
	FacilityLabel      string     `json:"Facility_Label"`
	Priority           int        `json:"Priority"`
	PriorityLabel      string     `json:"Priority_Label"`
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

// ScanFromRows scans a database row into a LogEntry
func (e *LogEntry) ScanFromRows(rows *sql.Rows) error {
	var customerID, ntSeverity, importance, eventCategory, eventID sql.NullInt64
	var maxAvail, currUsage, minUsage, maxUsage, infoUnitID, systemID sql.NullInt64
	var deviceTime sql.NullTime
	var eventSource, eventUser, eventBinary, sysLogTag, eventLogType, genericFile sql.NullString

	err := rows.Scan(
		&e.ID, &customerID, &e.ReceivedAt, &deviceTime,
		&e.Facility, &e.Priority, &e.FromHost, &e.Message,
		&ntSeverity, &importance, &eventSource, &eventUser,
		&eventCategory, &eventID, &eventBinary, &maxAvail, &currUsage,
		&minUsage, &maxUsage, &infoUnitID, &sysLogTag, &eventLogType,
		&genericFile, &systemID,
	)
	if err != nil {
		return err
	}

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

	// Add RFC labels
	e.PriorityLabel = GetPriorityLabel(e.Priority)
	e.FacilityLabel = GetFacilityLabel(e.Facility)

	return nil
}

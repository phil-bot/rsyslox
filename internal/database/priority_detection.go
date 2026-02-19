package database

import "log"

// PriorityMode describes how the Priority column is stored in the database.
// This differs between rsyslog versions:
//   - Legacy (< 8.2204.0): Priority column = Severity (0-7)
//   - Modern (>= 8.2204.0): Priority column = RFC PRI (Facility*8 + Severity)
type PriorityMode int

const (
	// PriorityModeLegacy: Priority column contains Severity (0-7).
	PriorityModeLegacy PriorityMode = iota

	// PriorityModeModern: Priority column contains RFC PRI (Facility*8 + Severity).
	PriorityModeModern

	// PriorityModeMixed: dataset contains both legacy and modern entries.
	PriorityModeMixed
)

// String returns a human-readable description of the priority mode.
func (m PriorityMode) String() string {
	switch m {
	case PriorityModeLegacy:
		return "legacy (Priority = Severity 0-7)"
	case PriorityModeModern:
		return "modern (Priority = Facility*8 + Severity)"
	case PriorityModeMixed:
		return "mixed (legacy + modern entries present)"
	default:
		return "unknown"
	}
}

// detectPriorityMode examines the database to determine the Priority column format.
//
// Strategy: sample the oldest and newest non-kernel entries (Facility > 0).
// Kernel entries (Facility = 0) are excluded because their Priority values are
// ambiguous — in modern format a low-severity kernel message still has Priority <= 7,
// making it indistinguishable from a legacy entry.
//
// Decision table:
//
//	oldest <= 7 AND newest <= 7  → Legacy
//	oldest > 7  AND newest > 7  → Modern
//	oldest <= 7 AND newest > 7  → Mixed (rsyslog was updated at some point)
//
// Fallback: if no non-kernel entries exist, legacy mode is assumed and a warning
// is logged.
func (db *DB) detectPriorityMode() PriorityMode {
	var oldest, newest int
	var oldestFound, newestFound bool

	row := db.QueryRow(
		"SELECT Priority FROM SystemEvents WHERE Facility > 0 ORDER BY ReceivedAt ASC LIMIT 1",
	)
	if err := row.Scan(&oldest); err == nil {
		oldestFound = true
	}

	row = db.QueryRow(
		"SELECT Priority FROM SystemEvents WHERE Facility > 0 ORDER BY ReceivedAt DESC LIMIT 1",
	)
	if err := row.Scan(&newest); err == nil {
		newestFound = true
	}

	if !oldestFound && !newestFound {
		log.Println("⚠ Priority mode detection: no non-kernel entries found, assuming legacy mode")
		return PriorityModeLegacy
	}

	oldestIsModern := oldestFound && oldest > 7
	newestIsModern := newestFound && newest > 7

	var mode PriorityMode
	switch {
	case !oldestIsModern && !newestIsModern:
		mode = PriorityModeLegacy
	case oldestIsModern && newestIsModern:
		mode = PriorityModeModern
	default:
		// Oldest entry is legacy-format, newest is modern-format → rsyslog was updated
		mode = PriorityModeMixed
	}

	log.Printf("✓ Priority mode detected: %s (oldest non-kernel Priority=%d, newest=%d)",
		mode, oldest, newest)
	return mode
}

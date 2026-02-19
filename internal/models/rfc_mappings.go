package models

import "fmt"

// RFC-5424 Severity Levels
const (
	SeverityEmergency     = 0
	SeverityAlert         = 1
	SeverityCritical      = 2
	SeverityError         = 3
	SeverityWarning       = 4
	SeverityNotice        = 5
	SeverityInformational = 6
	SeverityDebug         = 7
)

// RFC-5424 Facility Codes
const (
	FacilityKern     = 0
	FacilityUser     = 1
	FacilityMail     = 2
	FacilityDaemon   = 3
	FacilityAuth     = 4
	FacilitySyslog   = 5
	FacilityLpr      = 6
	FacilityNews     = 7
	FacilityUucp     = 8
	FacilityCron     = 9
	FacilityAuthPriv = 10
	FacilityFtp      = 11
	FacilityNtp      = 12
	FacilityLogAudit = 13
	FacilityLogAlert = 14
	FacilityClock    = 15
	FacilityLocal0   = 16
	FacilityLocal1   = 17
	FacilityLocal2   = 18
	FacilityLocal3   = 19
	FacilityLocal4   = 20
	FacilityLocal5   = 21
	FacilityLocal6   = 22
	FacilityLocal7   = 23
)

// RFCSeverity maps severity values (0-7) to RFC-5424 labels
var RFCSeverity = map[int]string{
	0: "Emergency",
	1: "Alert",
	2: "Critical",
	3: "Error",
	4: "Warning",
	5: "Notice",
	6: "Informational",
	7: "Debug",
}

// RFCFacility maps facility values to RFC-5424 facility labels
var RFCFacility = map[int]string{
	0:  "kern",
	1:  "user",
	2:  "mail",
	3:  "daemon",
	4:  "auth",
	5:  "syslog",
	6:  "lpr",
	7:  "news",
	8:  "uucp",
	9:  "cron",
	10: "authpriv",
	11: "ftp",
	12: "ntp",
	13: "logaudit",
	14: "logalert",
	15: "clock",
	16: "local0",
	17: "local1",
	18: "local2",
	19: "local3",
	20: "local4",
	21: "local5",
	22: "local6",
	23: "local7",
}

// GetSeverityLabel returns the RFC-5424 label for a severity value (0-7).
func GetSeverityLabel(severity int) string {
	if label, ok := RFCSeverity[severity]; ok {
		return label
	}
	return fmt.Sprintf("Unknown(%d)", severity)
}

// GetFacilityLabel returns the RFC-5424 label for a facility value.
func GetFacilityLabel(facility int) string {
	if label, ok := RFCFacility[facility]; ok {
		return label
	}
	return fmt.Sprintf("Unknown(%d)", facility)
}

// GetPriorityLabel is kept for internal compatibility.
// Prefer GetSeverityLabel for new code.
func GetPriorityLabel(severity int) string {
	return GetSeverityLabel(severity)
}

// IsValidSeverity checks if a severity value is valid (0-7).
func IsValidSeverity(severity int) bool {
	return severity >= 0 && severity <= 7
}

// IsValidPriority is kept for internal compatibility.
// Prefer IsValidSeverity for new code.
func IsValidPriority(priority int) bool {
	return IsValidSeverity(priority)
}

// IsValidFacility checks if a facility value is valid (0-23).
func IsValidFacility(facility int) bool {
	return facility >= 0 && facility <= 23
}

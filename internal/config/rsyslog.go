package config

import (
	"fmt"
	"os"
	"regexp"
)

// ParseRsyslogConfig reads MySQL connection details from rsyslog config file
func ParseRsyslogConfig(configPath string) (string, string, string, string, error) {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", "", "", "", fmt.Errorf("rsyslog config file not found: %s", configPath)
	}

	// Read file
	content, err := os.ReadFile(configPath)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to read rsyslog config: %v", err)
	}

	// Parse configuration
	// Pattern: action(type="ommysql" server="host" db="dbname" uid="user" pwd="password")
	pattern := `action\(type="ommysql"\s+server="([^"]+)"\s+db="([^"]+)"\s+uid="([^"]+)"\s+pwd="([^"]+)"\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(string(content))

	if len(matches) != 5 {
		return "", "", "", "", fmt.Errorf("failed to parse MySQL connection parameters from rsyslog config")
	}

	dbHost := matches[1]
	dbName := matches[2]
	dbUser := matches[3]
	dbPass := matches[4]

	return dbUser, dbPass, dbName, dbHost, nil
}

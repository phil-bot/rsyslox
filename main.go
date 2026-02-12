package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Configuration holds all app configuration
type Configuration struct {
	DBUser          string
	DBPass          string
	DBName          string
	DBHost          string
	ServerHost      string
	ServerPort      string
	APIKey          string
	AllowedOrigins  []string
	UseSSL          bool
	SSLCertFile     string
	SSLKeyFile      string
	RsyslogConfPath string
	InstallPath     string
}

// Global variables
var (
	config          *Configuration
	db              *sql.DB
	availableColumns []string // Available columns from SystemEvents
)

// RFC Mappings
var rfcSeverity = map[int]string{
	0: "Emergency",
	1: "Alert",
	2: "Critical",
	3: "Error",
	4: "Warning",
	5: "Notice",
	6: "Informational",
	7: "Debug",
}

var rfcFacility = map[int]string{
	0: "kern", 1: "user", 2: "mail", 3: "daemon", 4: "auth", 5: "syslog",
	6: "lpr", 7: "news", 8: "uucp", 9: "cron", 10: "authpriv", 11: "ftp",
	12: "ntp", 13: "logaudit", 14: "logalert", 15: "clock",
	16: "local0", 17: "local1", 18: "local2", 19: "local3",
	20: "local4", 21: "local5", 22: "local6", 23: "local7",
}

// LogEntry represents a single log entry from the database
type LogEntry struct {
	ID            int       `json:"ID"`
	ReceivedAt    time.Time `json:"ReceivedAt"`
	FromHost      string    `json:"FromHost"`
	Priority      int       `json:"Priority"`
	PriorityLabel string    `json:"Priority_Label"`
	Facility      int       `json:"Facility"`
	FacilityLabel string    `json:"Facility_Label"`
	Message       string    `json:"Message"`
}

// LogsResponse represents the response structure for /logs endpoint
type LogsResponse struct {
	Total  int        `json:"total"`
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	Rows   []LogEntry `json:"rows"`
}

// MetaValue represents a meta value with optional label
type MetaValue struct {
	Val   int    `json:"val"`
	Label string `json:"label"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// loadConfiguration loads configuration from environment and rsyslog config
func loadConfiguration() (*Configuration, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}
	installPath := filepath.Dir(execPath)

	config := &Configuration{
		ServerHost:      getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:      getEnv("SERVER_PORT", "8000"),
		APIKey:          getEnv("API_KEY", ""),
		UseSSL:          getEnv("USE_SSL", "false") == "true",
		SSLCertFile:     getEnv("SSL_CERTFILE", filepath.Join(installPath, "certs", "cert.pem")),
		SSLKeyFile:      getEnv("SSL_KEYFILE", filepath.Join(installPath, "certs", "key.pem")),
		RsyslogConfPath: getEnv("RSYSLOG_CONFIG_PATH", "/etc/rsyslog.d/mysql.conf"),
		InstallPath:     installPath,
	}

	// Parse allowed origins
	originsStr := getEnv("ALLOWED_ORIGINS", "*")
	config.AllowedOrigins = strings.Split(originsStr, ",")

	// Read rsyslog configuration
	dbUser, dbPass, dbName, dbHost, err := readRsyslogConfig(config.RsyslogConfPath)
	if err != nil {
		return nil, err
	}

	config.DBUser = dbUser
	config.DBPass = dbPass
	config.DBName = dbName
	config.DBHost = dbHost

	return config, nil
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// readRsyslogConfig reads MySQL connection details from rsyslog config
func readRsyslogConfig(configPath string) (string, string, string, string, error) {
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

	log.Printf("Database connection loaded: %s@%s/%s", dbUser, dbHost, dbName)

	return dbUser, dbPass, dbName, dbHost, nil
}

// initDatabase initializes the database connection and creates indexes
func initDatabase(cfg *Configuration) error {
	// Build connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName)

	// Open database connection
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established")

	// Create indexes
	createIndexes()

	return nil
}

// createIndexes creates necessary database indexes
func createIndexes() {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_receivedat ON SystemEvents (ReceivedAt)",
		"CREATE INDEX IF NOT EXISTS idx_host_time ON SystemEvents (FromHost, ReceivedAt)",
		"CREATE INDEX IF NOT EXISTS idx_priority ON SystemEvents (Priority)",
		"CREATE INDEX IF NOT EXISTS idx_facility ON SystemEvents (Facility)",
	}

	for _, query := range indexes {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Index creation info: %v", err)
		}
	}

	// Try to create fulltext index (may fail if already exists)
	if _, err := db.Exec("ALTER TABLE SystemEvents ADD FULLTEXT(Message)"); err != nil {
		log.Printf("Fulltext index info: %v", err)
	}

	log.Println("Database indexes created/verified")
}

// loadAvailableColumns loads all column names from SystemEvents table
func loadAvailableColumns() error {
	query := "SHOW COLUMNS FROM SystemEvents"
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query columns: %v", err)
	}
	defer rows.Close()

	availableColumns = []string{}
	for rows.Next() {
		var field, colType, null, key, def, extra sql.NullString
		if err := rows.Scan(&field, &colType, &null, &key, &def, &extra); err != nil {
			log.Printf("Warning: failed to scan column info: %v", err)
			continue
		}
		if field.Valid {
			availableColumns = append(availableColumns, field.String)
		}
	}

	log.Printf("Loaded %d columns from SystemEvents table", len(availableColumns))
	return nil
}

// Middleware: CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range config.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if config.AllowedOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Middleware: Logging
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// Middleware: Authentication
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.APIKey == "" {
			next(w, r)
			return
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey != config.APIKey {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Invalid or missing API key"})
			return
		}

		next(w, r)
	}
}

// Helper: Respond with JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// Handler: Root
func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Try to serve index.html
	indexPath := filepath.Join(config.InstallPath, "index.html")
	if _, err := os.Stat(indexPath); err == nil {
		http.ServeFile(w, r, indexPath)
		return
	}

	// Otherwise show API info
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"name":    "rsyslog REST API",
		"version": "0.2.0",
		"endpoints": map[string]string{
			"logs":   "/logs?limit=10&Priority=3",
			"meta":   "/meta or /meta/{column}",
			"health": "/health",
		},
	})
}

// Handler: Health Check
func handleHealth(w http.ResponseWriter, r *http.Request) {
	// Test database connection
	dbStatus := "connected"
	if err := db.Ping(); err != nil {
		dbStatus := "disconnected"
		respondJSON(w, http.StatusServiceUnavailable, HealthResponse{
			Status:    "unhealthy",
			Database:  dbStatus,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	respondJSON(w, http.StatusOK, HealthResponse{
		Status:    "healthy",
		Database:  dbStatus,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// Handler: Get Logs
func handleLogs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	
	offsetStr := query.Get("offset")
	limitStr := query.Get("limit")

	offset := 0
	if offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}

	limit := 10
	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 && val <= 1000 {
			limit = val
		}
	}

	// Parse filters
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")
	fromHost := query.Get("FromHost")
	priorityStr := query.Get("Priority")
	facilityStr := query.Get("Facility")
	message := query.Get("Message")

	where, args, err := buildFilters(startDate, endDate, fromHost, priorityStr, facilityStr, message)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Count total matching entries
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM SystemEvents WHERE %s", where)
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("Count query error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		return
	}

	// Query entries with pagination
	sqlQuery := fmt.Sprintf(`
		SELECT ID, ReceivedAt, FromHost, Priority, Facility, Message 
		FROM SystemEvents 
		WHERE %s 
		ORDER BY ReceivedAt DESC 
		LIMIT ? OFFSET ?
	`, where)

	args = append(args, limit, offset)
	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		return
	}
	defer rows.Close()

	var entries []LogEntry
	for rows.Next() {
		var entry LogEntry
		if err := rows.Scan(&entry.ID, &entry.ReceivedAt, &entry.FromHost,
			&entry.Priority, &entry.Facility, &entry.Message); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		entry.PriorityLabel = rfcSeverity[entry.Priority]
		if entry.PriorityLabel == "" {
			entry.PriorityLabel = fmt.Sprintf("Unknown(%d)", entry.Priority)
		}

		entry.FacilityLabel = rfcFacility[entry.Facility]
		if entry.FacilityLabel == "" {
			entry.FacilityLabel = fmt.Sprintf("Unknown(%d)", entry.Facility)
		}

		entries = append(entries, entry)
	}

	if entries == nil {
		entries = []LogEntry{}
	}

	respondJSON(w, http.StatusOK, LogsResponse{
		Total:  total,
		Offset: offset,
		Limit:  limit,
		Rows:   entries,
	})
}

// Handler: Get Meta
func handleMeta(w http.ResponseWriter, r *http.Request) {
	// Special case: /meta without column -> list all available columns
	if r.URL.Path == "/meta" {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"available_columns": availableColumns,
			"usage":            "GET /meta/{column} to get distinct values for a column",
		})
		return
	}

	// Extract column from path
	column := strings.TrimPrefix(r.URL.Path, "/meta/")
	column = strings.TrimSpace(column)

	// Validate column exists in database
	columnExists := false
	for _, col := range availableColumns {
		if col == column {
			columnExists = true
			break
		}
	}

	if !columnExists {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("Invalid column '%s'. Available columns: %s", column, strings.Join(availableColumns, ", ")),
		})
		return
	}

	// Parse filters (all filters are supported)
	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")
	fromHost := query.Get("FromHost")
	priorityStr := query.Get("Priority")
	facilityStr := query.Get("Facility")
	message := query.Get("Message")

	where, args, err := buildFilters(startDate, endDate, fromHost, priorityStr, facilityStr, message)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Query distinct values with NULL handling
	sqlQuery := fmt.Sprintf("SELECT DISTINCT %s FROM SystemEvents WHERE %s AND %s IS NOT NULL ORDER BY %s ASC", column, where, column, column)
	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Meta query error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		return
	}
	defer rows.Close()

	// Handle different column types
	if column == "Priority" || column == "Facility" {
		// Integer columns with RFC labels
		var values []MetaValue
		for rows.Next() {
			var val int
			if err := rows.Scan(&val); err != nil {
				continue
			}

			var label string
			if column == "Priority" {
				label = rfcSeverity[val]
				if label == "" {
					label = fmt.Sprintf("Unknown(%d)", val)
				}
			} else {
				label = rfcFacility[val]
				if label == "" {
					label = fmt.Sprintf("Unknown(%d)", val)
				}
			}

			values = append(values, MetaValue{Val: val, Label: label})
		}

		if values == nil {
			values = []MetaValue{}
		}
		respondJSON(w, http.StatusOK, values)
	} else if strings.Contains(strings.ToLower(column), "id") || 
	          column == "NTSeverity" || column == "Importance" || 
	          column == "EventCategory" || column == "EventID" || 
	          column == "MaxAvailable" || column == "CurrUsage" || 
	          column == "MinUsage" || column == "MaxUsage" || 
	          column == "InfoUnitID" || column == "SystemID" || column == "CustomerID" {
		// Integer columns without special labels
		var values []int
		for rows.Next() {
			var val sql.NullInt64
			if err := rows.Scan(&val); err != nil {
				continue
			}
			if val.Valid {
				values = append(values, int(val.Int64))
			}
		}

		if values == nil {
			values = []int{}
		}
		respondJSON(w, http.StatusOK, values)
	} else {
		// String/Text columns
		var values []string
		for rows.Next() {
			var val sql.NullString
			if err := rows.Scan(&val); err != nil {
				continue
			}
			if val.Valid && val.String != "" {
				values = append(values, val.String)
			}
		}

		if values == nil {
			values = []string{}
		}
		respondJSON(w, http.StatusOK, values)
	}
}

// buildFilters constructs WHERE clause and arguments for SQL query
func buildFilters(startDate, endDate, fromHost, priorityStr, facilityStr, message string) (string, []interface{}, error) {
	var conditions []string
	var args []interface{}

	// Parse dates
	var startDt, endDt time.Time
	var err error

	if startDate != "" {
		startDt, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			return "", nil, fmt.Errorf("invalid start_date format. Expected ISO 8601/RFC3339")
		}
	} else {
		startDt = time.Now().Add(-24 * time.Hour)
	}

	if endDate != "" {
		endDt, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			return "", nil, fmt.Errorf("invalid end_date format. Expected ISO 8601/RFC3339")
		}
	} else {
		endDt = time.Now()
	}

	// Validate date range
	if startDt.After(endDt) {
		return "", nil, fmt.Errorf("start_date cannot be after end_date")
	}

	if endDt.Sub(startDt) > 90*24*time.Hour {
		return "", nil, fmt.Errorf("date range cannot exceed 90 days")
	}

	// Date range filter
	conditions = append(conditions, "ReceivedAt BETWEEN ? AND ?")
	args = append(args, startDt, endDt)

	// FromHost filter
	if fromHost != "" {
		conditions = append(conditions, "FromHost = ?")
		args = append(args, fromHost)
	}

	// Priority filter
	if priorityStr != "" {
		priority, err := strconv.Atoi(priorityStr)
		if err != nil || priority < 0 || priority > 7 {
			return "", nil, fmt.Errorf("Priority must be between 0 and 7")
		}
		conditions = append(conditions, "Priority = ?")
		args = append(args, priority)
	}

	// Facility filter
	if facilityStr != "" {
		facility, err := strconv.Atoi(facilityStr)
		if err != nil || facility < 0 || facility > 23 {
			return "", nil, fmt.Errorf("Facility must be between 0 and 23")
		}
		conditions = append(conditions, "Facility = ?")
		args = append(args, facility)
	}

	// Message filter
	if message != "" {
		if len(strings.TrimSpace(message)) < 2 {
			return "", nil, fmt.Errorf("Message search must be at least 2 characters")
		}
		conditions = append(conditions, "Message LIKE ?")
		args = append(args, "%"+message+"%")
	}

	where := strings.Join(conditions, " AND ")
	return where, args, nil
}

// Handler: Static files
func handleStatic(w http.ResponseWriter, r *http.Request) {
	staticPath := filepath.Join(config.InstallPath, "static")
	http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))).ServeHTTP(w, r)
}

func main() {
	// Load configuration
	var err error
	config, err = loadConfiguration()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := initDatabase(config); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Load available columns from SystemEvents table
	if err := loadAvailableColumns(); err != nil {
		log.Fatalf("Failed to load available columns: %v", err)
	}

	// Log configuration
	if config.APIKey == "" {
		log.Println("WARNING: Running without API key authentication! Set API_KEY environment variable for production.")
	}
	log.Printf("CORS allowed origins: %v", config.AllowedOrigins)

	// Setup routes
	http.HandleFunc("/", corsMiddleware(loggingMiddleware(handleRoot)))
	http.HandleFunc("/health", corsMiddleware(loggingMiddleware(handleHealth)))
	http.HandleFunc("/logs", corsMiddleware(loggingMiddleware(authMiddleware(handleLogs))))
	http.HandleFunc("/meta", corsMiddleware(loggingMiddleware(authMiddleware(handleMeta))))   // List all columns
	http.HandleFunc("/meta/", corsMiddleware(loggingMiddleware(authMiddleware(handleMeta)))) // Specific column
	http.HandleFunc("/static/", corsMiddleware(loggingMiddleware(handleStatic)))

	// Start server
	addr := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)

	if config.UseSSL {
		// Check if SSL files exist
		if _, err := os.Stat(config.SSLCertFile); os.IsNotExist(err) {
			log.Fatalf("SSL certificate not found: %s", config.SSLCertFile)
		}
		if _, err := os.Stat(config.SSLKeyFile); os.IsNotExist(err) {
			log.Fatalf("SSL key not found: %s", config.SSLKeyFile)
		}

		log.Printf("Starting HTTPS server on https://%s", addr)
		if err := http.ListenAndServeTLS(addr, config.SSLCertFile, config.SSLKeyFile, nil); err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		log.Printf("WARNING: Running without SSL! Enable USE_SSL=true for production.")
		log.Printf("Starting HTTP server on http://%s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}

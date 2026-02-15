package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config holds all application configuration
type Config struct {
	// Database
	DBHost             string
	DBName             string
	DBUser             string
	DBPass             string
	DBConnectionString string // NEW: Optional full connection string

	// Server
	ServerHost string
	ServerPort string

	// Security
	APIKey      string
	UseSSL      bool
	SSLCertFile string
	SSLKeyFile  string

	// CORS
	AllowedOrigins []string

	// Paths
	InstallPath     string
	RsyslogConfPath string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}
	installPath := filepath.Dir(execPath)

	cfg := &Config{
		ServerHost:      getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:      getEnv("SERVER_PORT", "8000"),
		APIKey:          getEnv("API_KEY", ""),
		UseSSL:          getEnv("USE_SSL", "false") == "true",
		SSLCertFile:     getEnv("SSL_CERTFILE", filepath.Join(installPath, "certs", "cert.pem")),
		SSLKeyFile:      getEnv("SSL_KEYFILE", filepath.Join(installPath, "certs", "key.pem")),
		RsyslogConfPath: getEnv("RSYSLOG_CONFIG_PATH", "/etc/rsyslog.d/mysql.conf"),
		InstallPath:     installPath,
	}

	// Parse CORS origins
	originsStr := getEnv("ALLOWED_ORIGINS", "*")
	cfg.AllowedOrigins = strings.Split(originsStr, ",")
	for i := range cfg.AllowedOrigins {
		cfg.AllowedOrigins[i] = strings.TrimSpace(cfg.AllowedOrigins[i])
	}

	// Database configuration - multiple options
	if err := cfg.loadDatabaseConfig(); err != nil {
		return nil, err
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	cfg.logConfiguration()
	return cfg, nil
}

// loadDatabaseConfig loads database configuration with fallback options
func (c *Config) loadDatabaseConfig() error {
	// Option 1: Full connection string (NEW in v0.2.3+)
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if connStr != "" {
		c.DBConnectionString = connStr
		log.Println("Database connection from DB_CONNECTION_STRING")
		return nil
	}

	// Option 2: Individual environment variables (RECOMMENDED)
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	if dbHost != "" && dbName != "" && dbUser != "" && dbPass != "" {
		c.DBHost = dbHost
		c.DBName = dbName
		c.DBUser = dbUser
		c.DBPass = dbPass
		log.Printf("Database connection from environment: %s@%s/%s", dbUser, dbHost, dbName)
		return nil
	}

	// Option 3: Fallback to rsyslog config file
	log.Println("DB_* environment variables not set, trying rsyslog config file...")
	dbUser, dbPass, dbName, dbHost, err := ParseRsyslogConfig(c.RsyslogConfPath)
	if err != nil {
		return fmt.Errorf("failed to load database config: %v\n\n"+
			"Please set one of:\n"+
			"  1. DB_CONNECTION_STRING=user:pass@tcp(host)/dbname\n"+
			"  2. DB_HOST, DB_NAME, DB_USER, DB_PASS individually\n"+
			"  3. Readable rsyslog config at %s", err, c.RsyslogConfPath)
	}

	c.DBUser = dbUser
	c.DBPass = dbPass
	c.DBName = dbName
	c.DBHost = dbHost
	log.Printf("Database connection from rsyslog config: %s@%s/%s", dbUser, dbHost, dbName)

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Database validation
	if c.DBConnectionString == "" {
		if c.DBHost == "" || c.DBName == "" || c.DBUser == "" || c.DBPass == "" {
			return fmt.Errorf("incomplete database configuration")
		}
	}

	// Server validation
	if c.ServerPort == "" {
		return fmt.Errorf("SERVER_PORT cannot be empty")
	}

	// SSL validation
	if c.UseSSL {
		if _, err := os.Stat(c.SSLCertFile); os.IsNotExist(err) {
			return fmt.Errorf("SSL certificate not found: %s", c.SSLCertFile)
		}
		if _, err := os.Stat(c.SSLKeyFile); os.IsNotExist(err) {
			return fmt.Errorf("SSL key not found: %s", c.SSLKeyFile)
		}
	}

	return nil
}

// GetDSN returns the MySQL DSN (Data Source Name) for database connection
func (c *Config) GetDSN() string {
	if c.DBConnectionString != "" {
		// Parse and ensure parseTime=true
		if !strings.Contains(c.DBConnectionString, "parseTime=") {
			if strings.Contains(c.DBConnectionString, "?") {
				return c.DBConnectionString + "&parseTime=true&charset=utf8mb4"
			}
			return c.DBConnectionString + "?parseTime=true&charset=utf8mb4"
		}
		return c.DBConnectionString
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4",
		c.DBUser, c.DBPass, c.DBHost, c.DBName)
}

// logConfiguration logs the loaded configuration (without sensitive data)
func (c *Config) logConfiguration() {
	log.Println("Configuration loaded:")
	log.Printf("  Server: %s:%s", c.ServerHost, c.ServerPort)
	log.Printf("  SSL: %v", c.UseSSL)
	
	if c.DBConnectionString != "" {
		log.Printf("  Database: [connection string]")
	} else {
		log.Printf("  Database: %s@%s/%s", c.DBUser, c.DBHost, c.DBName)
	}
	
	if c.APIKey == "" {
		log.Println("  ⚠️  WARNING: Running without API key authentication!")
		log.Println("     Set API_KEY environment variable for production.")
	} else {
		log.Printf("  API Key: %s...%s", c.APIKey[:4], c.APIKey[len(c.APIKey)-4:])
	}
	
	log.Printf("  CORS Origins: %v", c.AllowedOrigins)
	log.Printf("  Install Path: %s", c.InstallPath)
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

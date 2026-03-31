package config

import "time"

// Config is the root configuration structure, mapped 1:1 to config.toml.
type Config struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Auth     AuthConfig     `toml:"auth"`
	Cleanup  CleanupConfig  `toml:"cleanup"`

	// Runtime-only fields (not persisted to TOML)
	InstallPath string `toml:"-"`
	ConfigPath  string `toml:"-"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host           string   `toml:"host"`
	Port           int      `toml:"port"`
	UseSSL         bool     `toml:"use_ssl"`
	SSLCertFile    string   `toml:"ssl_cert"`
	SSLKeyFile     string   `toml:"ssl_key"`
	AllowedOrigins []string `toml:"allowed_origins"`

	// Server-side defaults — applied to new browser sessions that have not
	// yet stored their own preferences in localStorage.
	AutoRefreshInterval int    `toml:"auto_refresh_interval"` // seconds
	DefaultTimeRange    string `toml:"default_time_range"`    // e.g. "24h"
	DefaultLanguage     string `toml:"default_language"`      // "en" | "de"
	DefaultFontSize     string `toml:"default_font_size"`     // "small" | "medium" | "large"
	DefaultTimeFormat   string `toml:"default_time_format"`   // "24h" | "12h"
}

// DatabaseConfig holds database connection settings.
// Password is stored AES-GCM encrypted with prefix "enc:".
type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"name"`
	User     string `toml:"user"`
	Password string `toml:"password"` // may be "enc:<base64>" or plaintext during setup
}

// ReadOnlyKey is a named API key for read-only access.
// The actual key is stored as a SHA-256 hex hash.
type ReadOnlyKey struct {
	Name    string `toml:"name"`
	KeyHash string `toml:"key_hash"` // sha256 hex
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	AdminPasswordHash string        `toml:"admin_password_hash"` // bcrypt
	ReadOnlyKeys      []ReadOnlyKey `toml:"read_only_keys"`
}

// CleanupConfig holds the log cleanup / housekeeping settings.
type CleanupConfig struct {
	Enabled          bool          `toml:"enabled"`
	DiskPath         string        `toml:"disk_path"`
	ThresholdPercent float64       `toml:"threshold_percent"`
	BatchSize        int           `toml:"batch_size"`
	Interval         time.Duration `toml:"interval"`
}

// defaults returns a Config pre-filled with sensible defaults.
func defaults() *Config {
	return &Config{
		Server: ServerConfig{
			Host:                "0.0.0.0",
			Port:                8000,
			UseSSL:              false,
			SSLCertFile:         "/etc/rsyslox/certs/cert.pem",
			SSLKeyFile:          "/etc/rsyslox/certs/key.pem",
			AllowedOrigins:      []string{"*"},
			AutoRefreshInterval: 30,
			DefaultTimeRange:    "24h",
			DefaultLanguage:     "en",
			DefaultFontSize:     "medium",
			DefaultTimeFormat:   "24h",
		},
		Database: DatabaseConfig{
			Host: "localhost",
			Port: 3306,
			Name: "Syslog",
		},
		Auth: AuthConfig{
			ReadOnlyKeys: []ReadOnlyKey{},
		},
		Cleanup: CleanupConfig{
			Enabled:          false,
			DiskPath:         "/var/lib/mysql",
			ThresholdPercent: 85.0,
			BatchSize:        1000,
			Interval:         15 * time.Minute,
		},
	}
}

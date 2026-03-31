package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/phil-bot/rsyslox/internal/cleanup"
	"github.com/phil-bot/rsyslox/internal/config"
	"github.com/phil-bot/rsyslox/internal/models"
)

// ConfigView is the safe, sanitized view of the config returned to the frontend.
type ConfigView struct {
	Server   ServerView   `json:"server"`
	Database DatabaseView `json:"database"`
	Cleanup  CleanupView  `json:"cleanup"`
}

type ServerView struct {
	Host                string   `json:"host"`
	Port                int      `json:"port"`
	UseSSL              bool     `json:"use_ssl"`
	SSLCertFile         string   `json:"ssl_cert"`
	SSLKeyFile          string   `json:"ssl_key"`
	AllowedOrigins      []string `json:"allowed_origins"`
	AutoRefreshInterval int      `json:"auto_refresh_interval"`
	DefaultTimeRange    string   `json:"default_time_range"`
	DefaultLanguage     string   `json:"default_language"`
	DefaultFontSize     string   `json:"default_font_size"`
	DefaultTimeFormat   string   `json:"default_time_format"`
}

type DatabaseView struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Name string `json:"name"`
	User string `json:"user"`
}

type CleanupView struct {
	Enabled          bool    `json:"enabled"`
	DiskPath         string  `json:"disk_path"`
	ThresholdPercent float64 `json:"threshold_percent"`
	BatchSize        int     `json:"batch_size"`
	IntervalSeconds  int     `json:"interval_seconds"`
}

type ConfigUpdateRequest struct {
	Server   *ServerUpdateRequest   `json:"server,omitempty"`
	Database *DatabaseUpdateRequest `json:"database,omitempty"`
	Cleanup  *CleanupUpdateRequest  `json:"cleanup,omitempty"`
}

type ServerUpdateRequest struct {
	Host                string   `json:"host,omitempty"`
	Port                *int     `json:"port,omitempty"`
	AllowedOrigins      []string `json:"allowed_origins,omitempty"`
	AutoRefreshInterval *int     `json:"auto_refresh_interval,omitempty"`
	UseSSL              *bool    `json:"use_ssl,omitempty"`
	DefaultTimeRange    string   `json:"default_time_range,omitempty"`
	DefaultLanguage     string   `json:"default_language,omitempty"`
	DefaultFontSize     string   `json:"default_font_size,omitempty"`
	DefaultTimeFormat   string   `json:"default_time_format,omitempty"`
}

type DatabaseUpdateRequest struct {
	Host     string `json:"host,omitempty"`
	Port     *int   `json:"port,omitempty"`
	Name     string `json:"name,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type CleanupUpdateRequest struct {
	Enabled          *bool    `json:"enabled,omitempty"`
	DiskPath         string   `json:"disk_path,omitempty"`
	ThresholdPercent *float64 `json:"threshold_percent,omitempty"`
	BatchSize        *int     `json:"batch_size,omitempty"`
	IntervalSeconds  *int     `json:"interval_seconds,omitempty"`
}

// ConfigHandler handles GET and PATCH /api/admin/config.
type ConfigHandler struct {
	cfg     *config.Config
	cleaner *cleanup.Cleaner
}

func NewConfigHandler(cfg *config.Config, cleaner *cleanup.Cleaner) *ConfigHandler {
	return &ConfigHandler{cfg: cfg, cleaner: cleaner}
}

func (h *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPatch:
		h.handlePatch(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed,
			models.NewAPIError("METHOD_NOT_ALLOWED", "Only GET and PATCH are allowed"))
	}
}

func (h *ConfigHandler) handleGet(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, toConfigView(h.cfg))
}

func (h *ConfigHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	var req ConfigUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest,
			models.NewAPIError(models.ErrCodeInvalidParameter, "Invalid JSON body"))
		return
	}

	if s := req.Server; s != nil {
		if s.Host != "" {
			h.cfg.Server.Host = s.Host
		}
		if s.Port != nil {
			if *s.Port < 1 || *s.Port > 65535 {
				respondError(w, http.StatusBadRequest,
					models.NewValidationError("port", "Must be between 1 and 65535"))
				return
			}
			h.cfg.Server.Port = *s.Port
		}
		if len(s.AllowedOrigins) > 0 {
			h.cfg.Server.AllowedOrigins = s.AllowedOrigins
		}
		if s.AutoRefreshInterval != nil {
			if *s.AutoRefreshInterval < 5 {
				respondError(w, http.StatusBadRequest,
					models.NewValidationError("auto_refresh_interval", "Minimum interval is 5 seconds"))
				return
			}
			h.cfg.Server.AutoRefreshInterval = *s.AutoRefreshInterval
		}
		if s.UseSSL != nil {
			h.cfg.Server.UseSSL = *s.UseSSL
		}
		if s.DefaultTimeRange != "" {
			h.cfg.Server.DefaultTimeRange = s.DefaultTimeRange
		}
		if s.DefaultLanguage != "" {
			h.cfg.Server.DefaultLanguage = s.DefaultLanguage
		}
		if s.DefaultFontSize != "" {
			h.cfg.Server.DefaultFontSize = s.DefaultFontSize
		}
		if s.DefaultTimeFormat != "" {
			h.cfg.Server.DefaultTimeFormat = s.DefaultTimeFormat
		}
	}

	if d := req.Database; d != nil {
		if d.Host != "" {
			h.cfg.Database.Host = d.Host
		}
		if d.Port != nil {
			if *d.Port < 1 || *d.Port > 65535 {
				respondError(w, http.StatusBadRequest,
					models.NewValidationError("database.port", "Must be between 1 and 65535"))
				return
			}
			h.cfg.Database.Port = *d.Port
		}
		if d.Name != "" {
			h.cfg.Database.Name = d.Name
		}
		if d.User != "" {
			h.cfg.Database.User = d.User
		}
		if d.Password != "" {
			encrypted, err := config.EncryptPassword(d.Password)
			if err != nil {
				respondError(w, http.StatusInternalServerError,
					models.NewAPIError("INTERNAL_ERROR", "Failed to encrypt password"))
				return
			}
			h.cfg.Database.Password = encrypted
		}
	}

	if c := req.Cleanup; c != nil {
		if c.Enabled != nil {
			h.cfg.Cleanup.Enabled = *c.Enabled
		}
		if c.DiskPath != "" {
			h.cfg.Cleanup.DiskPath = c.DiskPath
		}
		if c.ThresholdPercent != nil {
			if *c.ThresholdPercent <= 0 || *c.ThresholdPercent > 100 {
				respondError(w, http.StatusBadRequest,
					models.NewValidationError("threshold_percent", "Must be between 1 and 100"))
				return
			}
			h.cfg.Cleanup.ThresholdPercent = *c.ThresholdPercent
		}
		if c.BatchSize != nil {
			if *c.BatchSize <= 0 {
				respondError(w, http.StatusBadRequest,
					models.NewValidationError("batch_size", "Must be greater than 0"))
				return
			}
			h.cfg.Cleanup.BatchSize = *c.BatchSize
		}
		if c.IntervalSeconds != nil {
			if *c.IntervalSeconds < 60 {
				respondError(w, http.StatusBadRequest,
					models.NewValidationError("interval_seconds", "Minimum interval is 60 seconds"))
				return
			}
			h.cfg.Cleanup.Interval = time.Duration(*c.IntervalSeconds) * time.Second
		}
	}

	if err := config.Save(h.cfg); err != nil {
		log.Printf("Config update: failed to save: %v", err)
		respondError(w, http.StatusInternalServerError,
			models.NewAPIError("INTERNAL_ERROR", "Failed to save configuration"))
		return
	}

	// Propagate cleanup changes to the running goroutine without restart.
	if req.Cleanup != nil && h.cleaner != nil {
		h.cleaner.UpdateConfig(cleanup.Config{
			Enabled:          h.cfg.Cleanup.Enabled,
			DiskPath:         h.cfg.Cleanup.DiskPath,
			ThresholdPercent: h.cfg.Cleanup.ThresholdPercent,
			BatchSize:        h.cfg.Cleanup.BatchSize,
			Interval:         h.cfg.Cleanup.Interval,
		})
		log.Printf("Cleanup: config updated live (enabled=%v, threshold=%.1f%%)",
			h.cfg.Cleanup.Enabled, h.cfg.Cleanup.ThresholdPercent)
	}

	log.Println("Admin: configuration updated")
	respondJSON(w, http.StatusOK, toConfigView(h.cfg))
}

func toConfigView(cfg *config.Config) ConfigView {
	return ConfigView{
		Server: ServerView{
			Host:                cfg.Server.Host,
			Port:                cfg.Server.Port,
			UseSSL:              cfg.Server.UseSSL,
			SSLCertFile:         cfg.Server.SSLCertFile,
			SSLKeyFile:          cfg.Server.SSLKeyFile,
			AllowedOrigins:      cfg.Server.AllowedOrigins,
			AutoRefreshInterval: cfg.Server.AutoRefreshInterval,
			DefaultTimeRange:    cfg.Server.DefaultTimeRange,
			DefaultLanguage:     cfg.Server.DefaultLanguage,
			DefaultFontSize:     cfg.Server.DefaultFontSize,
			DefaultTimeFormat:   cfg.Server.DefaultTimeFormat,
		},
		Database: DatabaseView{
			Host: cfg.Database.Host,
			Port: cfg.Database.Port,
			Name: cfg.Database.Name,
			User: cfg.Database.User,
		},
		Cleanup: CleanupView{
			Enabled:          cfg.Cleanup.Enabled,
			DiskPath:         cfg.Cleanup.DiskPath,
			ThresholdPercent: cfg.Cleanup.ThresholdPercent,
			BatchSize:        cfg.Cleanup.BatchSize,
			IntervalSeconds:  int(cfg.Cleanup.Interval.Seconds()),
		},
	}
}

package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/phil-bot/rsyslog-rest-api/internal/config"
	"github.com/phil-bot/rsyslog-rest-api/internal/database"
	"github.com/phil-bot/rsyslog-rest-api/internal/handlers"
	"github.com/phil-bot/rsyslog-rest-api/internal/middleware"
)

// Server represents the HTTP server
type Server struct {
	config  *config.Config
	db      *database.DB
	router  *http.ServeMux
	version string
}

// New creates a new server instance
func New(cfg *config.Config, db *database.DB, version string) *Server {
	return &Server{
		config:  cfg,
		db:      db,
		router:  http.NewServeMux(),
		version: version,
	}
}

// SetupRoutes configures all HTTP routes and middleware
func (s *Server) SetupRoutes() {
	// Create handlers
	healthHandler := handlers.NewHealthHandler(s.db, s.version)
	rootHandler := handlers.NewRootHandler(s.config.InstallPath, s.version)
	logsHandler := handlers.NewLogsHandler(s.db)
	metaHandler := handlers.NewMetaHandler(s.db)

	// Create middleware chain
	corsMiddleware := middleware.CORS(s.config.AllowedOrigins)
	loggingMiddleware := middleware.Logging()
	authMiddleware := middleware.Auth(s.config.APIKey)

	// Public routes (no auth)
	s.router.Handle("/", 
		corsMiddleware(loggingMiddleware(rootHandler)))
	s.router.Handle("/health", 
		corsMiddleware(loggingMiddleware(healthHandler)))

	// Protected routes (with auth)
	s.router.Handle("/logs", 
		corsMiddleware(loggingMiddleware(authMiddleware(logsHandler))))
	s.router.Handle("/meta", 
		corsMiddleware(loggingMiddleware(authMiddleware(metaHandler))))
	s.router.Handle("/meta/", 
		corsMiddleware(loggingMiddleware(authMiddleware(metaHandler))))

	log.Println("✓ Routes configured")
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.ServerHost, s.config.ServerPort)

	if s.config.UseSSL {
		log.Printf("Starting HTTPS server on https://%s", addr)
		log.Printf("Version: %s", s.version)
		return http.ListenAndServeTLS(addr, s.config.SSLCertFile, s.config.SSLKeyFile, s.router)
	}

	log.Printf("⚠️  WARNING: Running without SSL! Enable USE_SSL=true for production.")
	log.Printf("Starting HTTP server on http://%s", addr)
	log.Printf("Version: %s", s.version)
	return http.ListenAndServe(addr, s.router)
}

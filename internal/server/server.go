// Package server wires together all HTTP handlers, middleware and the embedded
// frontend. The route layout is:
//
//	/                         → embedded Vue frontend (or setup wizard redirect)
//	/docs                     → embedded Redoc API documentation
//	/health                   → health check (public)
//	/api/setup                → first-run wizard (localhost only, no config)
//	/api/admin/login          → admin login (public)
//	/api/admin/logout         → admin logout (admin token)
//	/api/admin/config         → configuration (admin token)
//	/api/admin/keys           → read-only key management (admin token)
//	/api/admin/cleanup/status → partition status + list (admin token)
//	/api/logs                 → log entries (read-only key or admin token)
//	/api/meta                 → metadata (read-only key or admin token)
//	/api/meta/                → metadata column values (read-only key or admin token)
package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/phil-bot/rsyslox/internal/auth"
	"github.com/phil-bot/rsyslox/internal/cleanup"
	"github.com/phil-bot/rsyslox/internal/config"
	"github.com/phil-bot/rsyslox/internal/database"
	"github.com/phil-bot/rsyslox/internal/handlers"
	"github.com/phil-bot/rsyslox/internal/handlers/admin"
	"github.com/phil-bot/rsyslox/internal/handlers/setup"
	"github.com/phil-bot/rsyslox/internal/middleware"
)

// Server represents the HTTP server.
type Server struct {
	cfg          *config.Config
	db           *database.DB
	cleaner      *cleanup.Cleaner
	router       *http.ServeMux
	version      string
	setupMode    bool
	authMgr      *auth.Manager
	sessionStore *auth.SessionStore
}

// New creates a new Server instance.
func New(cfg *config.Config, db *database.DB, cleaner *cleanup.Cleaner, version string, setupMode bool) *Server {
	return &Server{
		cfg:          cfg,
		db:           db,
		cleaner:      cleaner,
		router:       http.NewServeMux(),
		version:      version,
		setupMode:    setupMode,
		authMgr:      auth.New(cfg),
		sessionStore: auth.NewSessionStore(),
	}
}

// SetupRoutes configures all HTTP routes and middleware.
func (s *Server) SetupRoutes() {
	cors          := middleware.CORS(s.cfg.Server.AllowedOrigins)
	logging       := middleware.Logging()
	authRO        := middleware.AuthReadOnly(s.authMgr, s.sessionStore)
	authAdmin     := middleware.AuthAdmin(s.sessionStore)
	localhostOnly := middleware.LocalhostOnly()

	// --- Frontend ---
	s.router.Handle("/", cors(logging(s.frontendHandler())))

	// --- Docs ---
	docsHandler := http.StripPrefix("/docs", s.docsHandler())
	s.router.Handle("/docs", http.RedirectHandler("/docs/", http.StatusMovedPermanently))
	s.router.Handle("/docs/", cors(logging(docsHandler)))

	// --- Health (public) ---
	s.router.Handle("/health", cors(logging(handlers.NewHealthHandler(s.db, s.version, s.cfg))))

	// --- Setup wizard ---
	setupHandler := setup.New(s.cfg, s.sessionStore)
	if s.setupMode {
		s.router.Handle("/api/setup", cors(logging(setupHandler)))
		log.Println("⚠️  Running in setup mode — open the web UI to complete setup")
		return
	}
	s.router.Handle("/api/setup", cors(logging(localhostOnly(setupHandler))))

	// --- Admin: auth ---
	s.router.Handle("/api/admin/login",  cors(logging(admin.NewLoginHandler(s.authMgr, s.sessionStore))))
	s.router.Handle("/api/admin/logout", cors(logging(authAdmin(admin.NewLogoutHandler(s.sessionStore)))))

	// --- Admin: config, keys, ssl, restart, disk, cleanup status ---
	s.router.Handle("/api/admin/config",          cors(logging(authAdmin(admin.NewConfigHandler(s.cfg, s.cleaner)))))
	s.router.Handle("/api/admin/keys",             cors(logging(authAdmin(admin.NewKeysHandler(s.cfg)))))
	s.router.Handle("/api/admin/keys/",            cors(logging(authAdmin(admin.NewKeysHandler(s.cfg)))))
	s.router.Handle("/api/admin/ssl/",             cors(logging(authAdmin(admin.NewSSLHandler(s.cfg)))))
	s.router.Handle("/api/admin/restart",          cors(logging(authAdmin(admin.NewRestartHandler()))))
	s.router.Handle("/api/admin/disk",             cors(logging(authAdmin(admin.NewDiskHandler(s.cfg)))))
	s.router.Handle("/api/admin/cleanup/status",   cors(logging(authAdmin(admin.NewCleanupStatusHandler(s.cleaner)))))

	// --- API: logs and meta ---
	s.router.Handle("/api/logs",  cors(logging(authRO(handlers.NewLogsHandler(s.db)))))
	s.router.Handle("/api/meta",  cors(logging(authRO(handlers.NewMetaHandler(s.db)))))
	s.router.Handle("/api/meta/", cors(logging(authRO(handlers.NewMetaHandler(s.db)))))

	log.Println("✓ Routes configured")
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)

	if s.cfg.Server.UseSSL {
		if err := config.EnsureSSLCerts(&s.cfg.Server); err != nil {
			return fmt.Errorf("SSL setup failed: %w", err)
		}
		log.Printf("Starting HTTPS server on https://%s", addr)
		return http.ListenAndServeTLS(addr,
			s.cfg.Server.SSLCertFile,
			s.cfg.Server.SSLKeyFile,
			s.router)
	}

	if !s.setupMode {
		log.Printf("⚠️  WARNING: Running without SSL! Enable use_ssl=true for production.")
	}
	log.Printf("Starting HTTP server on http://%s", addr)
	return http.ListenAndServe(addr, s.router)
}

// frontendHandler serves the embedded Vue app.
func (s *Server) frontendHandler() http.Handler {
	sub, err := fs.Sub(FrontendFS, "frontend/dist")
	if err != nil {
		log.Println("⚠️  No embedded frontend found. Run 'make frontend' first.")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<!DOCTYPE html><html><body>
				<h2>rsyslox</h2>
				<p>Frontend not built. Run <code>make frontend</code> first.</p>
				<p><a href="/health">Health check</a></p>
			</body></html>`))
		})
	}

	fileServer := http.FileServer(http.FS(sub))
	indexHTML, indexErr := fs.ReadFile(sub, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fsPath := strings.TrimPrefix(r.URL.Path, "/")
		if fsPath == "" {
			fsPath = "index.html"
		}
		if fsPath != "index.html" {
			if _, openErr := sub.Open(fsPath); openErr == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
		}
		if indexErr != nil {
			http.Error(w, "Frontend not available", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(indexHTML)
	})
}

// docsHandler serves the embedded Redoc documentation.
func (s *Server) docsHandler() http.Handler {
	sub, err := fs.Sub(DocsFS, "docs/api-ui")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "API documentation not available", http.StatusNotFound)
		})
	}
	return http.FileServer(http.FS(sub))
}

// Package server wires together all HTTP handlers, middleware and the embedded
// frontend. The route layout is:
//
//	/                  → embedded Vue frontend (or setup wizard redirect)
//	/docs              → embedded Redoc API documentation
//	/health            → health check (public)
//	/api/setup         → first-run wizard (localhost only, no config)
//	/api/admin/login   → admin login (public)
//	/api/admin/logout  → admin logout (admin token)
//	/api/admin/config  → configuration (admin token)
//	/api/admin/keys    → read-only key management (admin token)
//	/api/logs          → log entries (read-only key or admin token)
//	/api/meta          → metadata (read-only key or admin token)
//	/api/meta/         → metadata column values (read-only key or admin token)
package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/phil-bot/rsyslox/internal/auth"
	"github.com/phil-bot/rsyslox/internal/config"
	"github.com/phil-bot/rsyslox/internal/database"
	"github.com/phil-bot/rsyslox/internal/handlers"
	"github.com/phil-bot/rsyslox/internal/handlers/admin"
	"github.com/phil-bot/rsyslox/internal/handlers/setup"
	"github.com/phil-bot/rsyslox/internal/middleware"
)

// frontendFS and docsFS are injected at build time via go:embed in embed.go.
var (
	FrontendFS embed.FS
	DocsFS     embed.FS
)

// Server represents the HTTP server.
type Server struct {
	cfg          *config.Config
	db           *database.DB
	router       *http.ServeMux
	version      string
	setupMode    bool
	authMgr      *auth.Manager
	sessionStore *auth.SessionStore
}

// New creates a new Server instance.
// setupMode=true means no config file was found; only the setup wizard is enabled.
func New(cfg *config.Config, db *database.DB, version string, setupMode bool) *Server {
	return &Server{
		cfg:          cfg,
		db:           db,
		router:       http.NewServeMux(),
		version:      version,
		setupMode:    setupMode,
		authMgr:      auth.New(cfg),
		sessionStore: auth.NewSessionStore(),
	}
}

// SetupRoutes configures all HTTP routes and middleware.
func (s *Server) SetupRoutes() {
	cors := middleware.CORS(s.cfg.Server.AllowedOrigins)
	logging := middleware.Logging()
	authRO := middleware.AuthReadOnly(s.authMgr, s.sessionStore)
	authAdmin := middleware.AuthAdmin(s.sessionStore)
	localhostOnly := middleware.LocalhostOnly()

	// --- Frontend ---
	frontendHandler := s.frontendHandler()
	s.router.Handle("/", cors(logging(frontendHandler)))

	// --- Docs (Redoc, offline) ---
	docsHandler := s.docsHandler()
	s.router.Handle("/docs", cors(logging(docsHandler)))
	s.router.Handle("/docs/", cors(logging(docsHandler)))

	// --- Health (public) ---
	healthHandler := handlers.NewHealthHandler(s.db, s.version, s.cfg)
	s.router.Handle("/health", cors(logging(healthHandler)))

	// --- Setup wizard (localhost only, only in setup mode) ---
	setupHandler := setup.New(s.cfg, s.sessionStore)
	s.router.Handle("/api/setup", cors(logging(localhostOnly(setupHandler))))

	if s.setupMode {
		log.Println("⚠️  Running in setup mode — only /api/setup and /health are active")
		log.Println("   Open http://localhost:" + fmt.Sprint(s.cfg.Server.Port) + " to complete setup")
		return
	}

	// --- Admin: login / logout (public, rate-limited by bcrypt cost) ---
	loginHandler := admin.NewLoginHandler(s.authMgr, s.sessionStore)
	logoutHandler := admin.NewLogoutHandler(s.sessionStore)
	s.router.Handle("/api/admin/login", cors(logging(loginHandler)))
	s.router.Handle("/api/admin/logout", cors(logging(authAdmin(logoutHandler))))

	// --- Admin: config and key management (admin token required) ---
	configHandler := admin.NewConfigHandler(s.cfg)
	keysHandler   := admin.NewKeysHandler(s.cfg)
	sslHandler    := admin.NewSSLHandler(s.cfg)
	s.router.Handle("/api/admin/config", cors(logging(authAdmin(configHandler))))
	s.router.Handle("/api/admin/keys",   cors(logging(authAdmin(keysHandler))))
	s.router.Handle("/api/admin/keys/",  cors(logging(authAdmin(keysHandler))))
	s.router.Handle("/api/admin/ssl/",   cors(logging(authAdmin(sslHandler))))

	// --- API: logs and meta (read-only key or admin token) ---
	logsHandler := handlers.NewLogsHandler(s.db)
	metaHandler := handlers.NewMetaHandler(s.db)
	s.router.Handle("/api/logs", cors(logging(authRO(logsHandler))))
	s.router.Handle("/api/meta", cors(logging(authRO(metaHandler))))
	s.router.Handle("/api/meta/", cors(logging(authRO(metaHandler))))

	log.Println("✓ Routes configured")
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)

	if s.cfg.Server.UseSSL {
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
// In setup mode it serves the same SPA — the Vue router shows the wizard.
func (s *Server) frontendHandler() http.Handler {
	sub, err := fs.Sub(FrontendFS, "frontend/dist")
	if err != nil {
		// No embedded frontend (dev mode without npm build) — serve placeholder
		log.Println("⚠️  No embedded frontend found. Run 'make frontend' first.")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<!DOCTYPE html><html><body>
				<h2>rsyslox</h2>
				<p>Frontend not built. Run <code>make frontend</code> first.</p>
				<p><a href="/health">Health check</a></p>
			</body></html>`))
		})
	}

	fileServer := http.FileServer(http.FS(sub))

	// SPA fallback: serve index.html for unknown paths so Vue Router can handle them
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || path == "" {
			r.URL.Path = "/index.html"
		}
		// Try to serve the file; if not found, fall back to index.html
		_, err := sub.Open(path)
		if err != nil {
			r.URL.Path = "/index.html"
		}
		fileServer.ServeHTTP(w, r)
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

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/phil-bot/rsyslox/internal/auth"
	"github.com/phil-bot/rsyslox/internal/cleanup"
	"github.com/phil-bot/rsyslox/internal/config"
	"github.com/phil-bot/rsyslox/internal/database"
	"github.com/phil-bot/rsyslox/internal/server"
)

// Version is set at build time via ldflags.
var Version = "dev"

func main() {
	// Subcommand: rsyslox hash-password <plaintext>
	if len(os.Args) == 3 && os.Args[1] == "hash-password" {
		hash, err := auth.HashAdminPassword(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(hash)
		return
	}

	log.Println("========================================")
	log.Println("rsyslox", Version)
	log.Println("========================================")

	server.FrontendFS = frontendFS
	server.DocsFS = docsFS

	cfg, setupMode, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	if setupMode {
		log.Println("⚠️  No configuration found — starting in setup mode")
		log.Printf("   Setup wizard available at http://<this-host>:%d", cfg.Server.Port)
		srv := server.New(cfg, nil, nil, Version, true)
		srv.SetupRoutes()
		if err := srv.Start(); err != nil {
			log.Fatalf("❌ Server error: %v", err)
		}
		return
	}

	// Connect to database.
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Start cleanup service.
	cleaner := cleanup.New(db.DB, cleanup.Config{
		Enabled:          cfg.Cleanup.Enabled,
		DiskPath:         cfg.Cleanup.DiskPath,
		ThresholdPercent: cfg.Cleanup.ThresholdPercent,
		Interval:         cfg.Cleanup.Interval,
	})
	cleaner.Start()
	defer cleaner.Stop()

	// Start server — pass cleaner so admin handlers can query its status.
	srv := server.New(cfg, db, cleaner, Version, false)
	srv.SetupRoutes()

	log.Println("========================================")
	log.Println("✓ Ready to accept connections")
	log.Println("========================================")

	if err := srv.Start(); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}

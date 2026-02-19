package main

import (
	"log"

	"github.com/phil-bot/rsyslox/internal/cleanup"
	"github.com/phil-bot/rsyslox/internal/config"
	"github.com/phil-bot/rsyslox/internal/database"
	"github.com/phil-bot/rsyslox/internal/server"
)

// Version is set at build time via ldflags
var Version = "dev"

func main() {
	log.Println("========================================")
	log.Println("rsyslox")
	log.Println("========================================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Start cleanup service
	cleaner := cleanup.New(db.DB, cleanup.Config{
		Enabled:             cfg.CleanupEnabled,
		DiskPath:            cfg.CleanupDiskPath,
		ThresholdPercent:    cfg.CleanupThresholdPercent,
		BatchSize:           cfg.CleanupBatchSize,
		Interval:            cfg.CleanupInterval,
	})
	cleaner.Start()
	defer cleaner.Stop()

	// Create and configure server
	srv := server.New(cfg, db, Version)
	srv.SetupRoutes()

	log.Println("========================================")
	log.Println("✓ Ready to accept connections")
	log.Println("========================================")

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}

package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/phil-bot/rsyslog-rest-api/internal/config"
)

// DB wraps the database connection and provides helper methods
type DB struct {
	*sql.DB
	AvailableColumns []string
}

// Connect establishes a connection to the database
func Connect(cfg *config.Config) (*DB, error) {
	dsn := cfg.GetDSN()

	// Open database connection
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("✓ Database connection established")

	db := &DB{DB: sqlDB}

	// Initialize database
	if err := db.initialize(); err != nil {
		return nil, err
	}

	return db, nil
}

// initialize performs initial database setup
func (db *DB) initialize() error {
	// Create indexes
	if err := db.createIndexes(); err != nil {
		return err
	}

	// Load available columns
	if err := db.loadColumns(); err != nil {
		return err
	}

	return nil
}

// loadColumns loads all column names from SystemEvents table
func (db *DB) loadColumns() error {
	query := "SHOW COLUMNS FROM SystemEvents"
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query columns: %v", err)
	}
	defer rows.Close()

	db.AvailableColumns = []string{}
	for rows.Next() {
		var field, colType, null, key, def, extra sql.NullString
		if err := rows.Scan(&field, &colType, &null, &key, &def, &extra); err != nil {
			log.Printf("Warning: failed to scan column info: %v", err)
			continue
		}
		if field.Valid {
			db.AvailableColumns = append(db.AvailableColumns, field.String)
		}
	}

	log.Printf("✓ Loaded %d columns from SystemEvents table", len(db.AvailableColumns))
	return nil
}

// IsValidColumn checks if a column name exists in the SystemEvents table
func (db *DB) IsValidColumn(column string) bool {
	for _, col := range db.AvailableColumns {
		if col == column {
			return true
		}
	}
	return false
}

// Health checks the database connection health
func (db *DB) Health() error {
	return db.Ping()
}

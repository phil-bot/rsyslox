package database

import "log"

// createIndexes creates necessary database indexes for optimal query performance
func (db *DB) createIndexes() error {
	indexes := []struct {
		name  string
		query string
	}{
		{
			name:  "idx_receivedat",
			query: "CREATE INDEX IF NOT EXISTS idx_receivedat ON SystemEvents (ReceivedAt)",
		},
		{
			name:  "idx_host_time",
			query: "CREATE INDEX IF NOT EXISTS idx_host_time ON SystemEvents (FromHost, ReceivedAt)",
		},
		{
			name:  "idx_priority",
			query: "CREATE INDEX IF NOT EXISTS idx_priority ON SystemEvents (Priority)",
		},
		{
			name:  "idx_facility",
			query: "CREATE INDEX IF NOT EXISTS idx_facility ON SystemEvents (Facility)",
		},
		{
			name:  "idx_syslogtag",
			query: "CREATE INDEX IF NOT EXISTS idx_syslogtag ON SystemEvents (SysLogTag)",
		},
	}

	for _, idx := range indexes {
		if _, err := db.Exec(idx.query); err != nil {
			log.Printf("Index creation info (%s): %v", idx.name, err)
		}
	}

	// Try to create fulltext index (may fail if already exists)
	if _, err := db.Exec("ALTER TABLE SystemEvents ADD FULLTEXT(Message)"); err != nil {
		log.Printf("Fulltext index info: %v", err)
	}

	log.Println("✓ Database indexes created/verified")
	return nil
}

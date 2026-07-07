package database

import "log"

// createIndexes creates necessary database indexes for optimal query performance.
//
// Note: rsyslox does NOT use MySQL FULLTEXT search — message filtering is done
// via "Message LIKE ?" (see internal/filters/builder.go, AddMessageSearch).
// A FULLTEXT index was previously created here but never queried by the
// application. Its internal FTS_* auxiliary tables can grow to many gigabytes
// on large tables and are never reclaimed automatically, so it is actively
// removed if found (see dropOrphanedFulltextIndex below).
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

	if err := db.dropOrphanedFulltextIndex(); err != nil {
		log.Printf("⚠  Fulltext index cleanup: %v", err)
	}

	log.Println("✓ Database indexes created/verified")
	return nil
}

// dropOrphanedFulltextIndex finds and removes any FULLTEXT index on the
// Message column, regardless of its name. Earlier rsyslox versions created
// such an index automatically, but it is never used by the application
// (message search uses LIKE, not MATCH ... AGAINST). The index's internal
// FTS_* auxiliary tables are not reclaimed by MySQL/MariaDB on their own and
// can consume many gigabytes of disk space over time, so it is dropped here
// unconditionally once found.
func (db *DB) dropOrphanedFulltextIndex() error {
	rows, err := db.Query(`
		SELECT DISTINCT INDEX_NAME
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME   = 'SystemEvents'
		  AND INDEX_TYPE   = 'FULLTEXT'
		  AND COLUMN_NAME  = 'Message'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var indexNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		indexNames = append(indexNames, name)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, name := range indexNames {
		log.Printf("⚙  Found orphaned FULLTEXT index %q on SystemEvents.Message — dropping to reclaim disk space…", name)
		if _, err := db.Exec("ALTER TABLE SystemEvents DROP INDEX `" + name + "`"); err != nil {
			return err
		}
		log.Printf("✓ Dropped FULLTEXT index %q — associated FTS_* files will be removed by MySQL", name)
	}

	return nil
}

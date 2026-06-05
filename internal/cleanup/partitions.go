package cleanup

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	weeksAhead      = 6
	maxDropsPerRun  = 3
	futurePartition = "p_future"
	partitionTable  = "SystemEvents"
)

// IsPartitioned returns true if SystemEvents is already using RANGE partitioning.
func IsPartitioned(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_NAME   = ?
		  AND TABLE_SCHEMA = DATABASE()
		  AND PARTITION_NAME IS NOT NULL
	`, partitionTable).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("partition check failed: %w", err)
	}
	return count > 0, nil
}

// MigrateToPartitions migrates an unpartitioned SystemEvents table to weekly
// RANGE COLUMNS partitioning.
func MigrateToPartitions(db *sql.DB) error {
	log.Println("⚙  Partition migration: starting — this may take a while on large tables…")

	if _, err := db.Exec(`
		ALTER TABLE SystemEvents
		  DROP PRIMARY KEY,
		  ADD  PRIMARY KEY (ID, ReceivedAt)
	`); err != nil {
		log.Printf("⚠  Partition migration: PK alteration note: %v", err)
	}

	oldest, err := oldestReceivedAt(db)
	if err != nil {
		return fmt.Errorf("could not determine oldest entry: %w", err)
	}

	clause, err := buildPartitionClause(oldest, weeksAhead)
	if err != nil {
		return fmt.Errorf("could not build partition clause: %w", err)
	}

	query := fmt.Sprintf(
		"ALTER TABLE %s PARTITION BY RANGE COLUMNS(ReceivedAt) (%s)",
		partitionTable, clause,
	)

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("ALTER TABLE failed: %w", err)
	}

	log.Println("✓ Partition migration complete")
	return nil
}

// EnsureFuturePartitions creates weekly partitions for the next weeksAhead
// weeks if they do not exist yet.
func EnsureFuturePartitions(db *sql.DB) error {
	now := time.Now()
	for i := 1; i <= weeksAhead; i++ {
		weekStart := startOfWeek(now.AddDate(0, 0, 7*i))
		weekEnd := weekStart.AddDate(0, 0, 7)
		name := partitionName(weekStart)

		exists, err := partitionExists(db, name)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		query := fmt.Sprintf(`
			ALTER TABLE %s REORGANIZE PARTITION %s INTO (
				PARTITION %s VALUES LESS THAN ('%s'),
				PARTITION %s VALUES LESS THAN (MAXVALUE)
			)`,
			partitionTable,
			futurePartition,
			name, weekEnd.Format("2006-01-02"),
			futurePartition,
		)

		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create partition %s: %w", name, err)
		}
		log.Printf("✓ Partition created: %s (< %s)", name, weekEnd.Format("2006-01-02"))
	}
	return nil
}

// DropOldestPartitions drops up to maxDropsPerRun of the oldest data partitions.
// Returns the number of partitions actually dropped.
func DropOldestPartitions(db *sql.DB) (int, error) {
	dropped := 0
	for dropped < maxDropsPerRun {
		name, err := oldestPartitionName(db)
		if err == sql.ErrNoRows {
			break
		}
		if err != nil {
			return dropped, fmt.Errorf("query failed: %w", err)
		}

		query := fmt.Sprintf("ALTER TABLE %s DROP PARTITION %s", partitionTable, name)
		if _, err := db.Exec(query); err != nil {
			return dropped, fmt.Errorf("DROP PARTITION %s failed: %w", name, err)
		}

		log.Printf("✓ Partition dropped: %s", name)
		dropped++
	}
	return dropped, nil
}

// HasAlterPrivilege checks whether the current DB user has ALTER privilege.
func HasAlterPrivilege(db *sql.DB) (bool, error) {
	_, err := db.Exec(fmt.Sprintf("ALTER TABLE %s COMMENT = ''", partitionTable))
	if err != nil {
		return false, err
	}
	return true, nil
}

// listPartitions returns all partitions of SystemEvents with row count and size.
func listPartitions(db *sql.DB) ([]Partition, error) {
	rows, err := db.Query(`
		SELECT
			PARTITION_NAME,
			COALESCE(TABLE_ROWS, 0),
			COALESCE(ROUND((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024, 2), 0)
		FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_NAME   = ?
		  AND TABLE_SCHEMA = DATABASE()
		  AND PARTITION_NAME IS NOT NULL
		ORDER BY PARTITION_DESCRIPTION ASC
	`, partitionTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Partition
	for rows.Next() {
		var p Partition
		if err := rows.Scan(&p.Name, &p.Rows, &p.SizeMB); err != nil {
			continue
		}
		result = append(result, p)
	}
	if result == nil {
		result = []Partition{}
	}
	return result, nil
}

// --------------------------------------------------------------------------
// Internal helpers
// --------------------------------------------------------------------------

func oldestReceivedAt(db *sql.DB) (time.Time, error) {
	var t sql.NullTime
	if err := db.QueryRow("SELECT MIN(ReceivedAt) FROM SystemEvents").Scan(&t); err != nil {
		return time.Time{}, err
	}
	if !t.Valid {
		return startOfWeek(time.Now()), nil
	}
	return t.Time, nil
}

func oldestPartitionName(db *sql.DB) (string, error) {
	var name string
	err := db.QueryRow(`
		SELECT PARTITION_NAME
		FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_NAME   = ?
		  AND TABLE_SCHEMA = DATABASE()
		  AND PARTITION_NAME != ?
		ORDER BY PARTITION_DESCRIPTION ASC
		LIMIT 1
	`, partitionTable, futurePartition).Scan(&name)
	return name, err
}

func partitionExists(db *sql.DB, name string) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_NAME    = ?
		  AND TABLE_SCHEMA  = DATABASE()
		  AND PARTITION_NAME = ?
	`, partitionTable, name).Scan(&count)
	return count > 0, err
}

func buildPartitionClause(oldest time.Time, ahead int) (string, error) {
	now := time.Now()
	end := now.AddDate(0, 0, 7*ahead)

	type entry struct {
		name    string
		weekEnd time.Time
	}
	var entries []entry

	cur := startOfWeek(oldest)
	for !cur.After(end) {
		weekEnd := cur.AddDate(0, 0, 7)
		entries = append(entries, entry{
			name:    partitionName(cur),
			weekEnd: weekEnd,
		})
		cur = weekEnd
	}

	if len(entries) == 0 {
		return "", fmt.Errorf("no partitions generated")
	}

	sql := ""
	for _, e := range entries {
		sql += fmt.Sprintf(
			"\n\tPARTITION %s VALUES LESS THAN ('%s'),",
			e.name, e.weekEnd.Format("2006-01-02"),
		)
	}
	sql += fmt.Sprintf("\n\tPARTITION %s VALUES LESS THAN (MAXVALUE)", futurePartition)

	return sql, nil
}

func startOfWeek(t time.Time) time.Time {
	t = t.UTC()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := t.AddDate(0, 0, -(weekday - 1))
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
}

func partitionName(weekStart time.Time) string {
	year, week := weekStart.ISOWeek()
	return fmt.Sprintf("p%d_w%02d", year, week)
}

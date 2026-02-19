package cleanup

import (
	"database/sql"
	"log"
	"syscall"
	"time"
)

// Cleaner periodically removes old database entries when disk usage exceeds a threshold.
type Cleaner struct {
	db           *sql.DB
	cfg          Config
	stopCh       chan struct{}
}

// Config holds the cleanup configuration.
type Config struct {
	// Enabled enables or disables the cleanup service.
	Enabled bool

	// DiskPath is the filesystem path to monitor for disk usage (e.g. /var/lib/mysql).
	DiskPath string

	// ThresholdPercent is the maximum allowed disk usage in percent (e.g. 85).
	// When usage exceeds this value, old records will be deleted.
	ThresholdPercent float64

	// BatchSize is the number of records to delete per cleanup run.
	BatchSize int

	// Interval is how often the cleanup check runs.
	Interval time.Duration
}

// New creates a new Cleaner instance.
func New(db *sql.DB, cfg Config) *Cleaner {
	return &Cleaner{
		db:     db,
		cfg:    cfg,
		stopCh: make(chan struct{}),
	}
}

// Start launches the cleanup loop in a background goroutine.
func (c *Cleaner) Start() {
	if !c.cfg.Enabled {
		log.Println("⏭  Cleanup service disabled")
		return
	}

	log.Printf("✓ Cleanup service started (threshold: %.1f%%, interval: %s, batch: %d)",
		c.cfg.ThresholdPercent, c.cfg.Interval, c.cfg.BatchSize)

	go c.run()
}

// Stop signals the cleanup loop to stop.
func (c *Cleaner) Stop() {
	close(c.stopCh)
}

// run is the main cleanup loop.
func (c *Cleaner) run() {
	ticker := time.NewTicker(c.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.check()
		case <-c.stopCh:
			log.Println("Cleanup service stopped")
			return
		}
	}
}

// check evaluates the current disk usage and deletes records if necessary.
func (c *Cleaner) check() {
	usedPercent, err := diskUsagePercent(c.cfg.DiskPath)
	if err != nil {
		log.Printf("⚠️  Cleanup: failed to get disk usage for %s: %v", c.cfg.DiskPath, err)
		return
	}

	log.Printf("Cleanup: disk usage at %.1f%% (threshold: %.1f%%)", usedPercent, c.cfg.ThresholdPercent)

	if usedPercent < c.cfg.ThresholdPercent {
		return
	}

	log.Printf("⚠️  Cleanup: disk usage %.1f%% exceeds threshold %.1f%% — deleting %d old records",
		usedPercent, c.cfg.ThresholdPercent, c.cfg.BatchSize)

	deleted, err := c.deleteOldestRecords(c.cfg.BatchSize)
	if err != nil {
		log.Printf("❌ Cleanup: failed to delete records: %v", err)
		return
	}

	log.Printf("✓ Cleanup: deleted %d records", deleted)
}

// deleteOldestRecords removes the oldest N records from SystemEvents.
// Returns the number of actually deleted rows.
func (c *Cleaner) deleteOldestRecords(n int) (int64, error) {
	// Use a subquery with a derived table to work around MySQL's limitation
	// of not being able to reference the target table in a DELETE subquery directly.
	query := `
		DELETE FROM SystemEvents
		WHERE ID IN (
			SELECT id FROM (
				SELECT ID as id FROM SystemEvents
				ORDER BY ReceivedAt ASC, ID ASC
				LIMIT ?
			) AS oldest
		)
	`

	result, err := c.db.Exec(query, n)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// diskUsagePercent returns the used disk space as a percentage for the given path.
func diskUsagePercent(path string) (float64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)

	if total == 0 {
		return 0, nil
	}

	used := total - free
	return float64(used) / float64(total) * 100.0, nil
}

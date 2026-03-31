package cleanup

import (
	"database/sql"
	"log"
	"sync"
	"syscall"
	"time"
)

// Cleaner periodically removes old database entries when disk usage exceeds a threshold.
// Config can be updated at runtime via UpdateConfig without a process restart.
type Cleaner struct {
	db      *sql.DB
	cfg     Config
	mu      sync.RWMutex
	stopCh  chan struct{}
	resetCh chan struct{} // signals the run loop to re-read config
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
		db:      db,
		cfg:     cfg,
		stopCh:  make(chan struct{}),
		resetCh: make(chan struct{}, 1),
	}
}

// Start launches the cleanup loop in a background goroutine.
func (c *Cleaner) Start() {
	c.mu.RLock()
	cfg := c.cfg
	c.mu.RUnlock()

	if !cfg.Enabled {
		log.Println("⏭  Cleanup service disabled (can be enabled in admin without restart)")
	} else {
		log.Printf("✓ Cleanup service started (threshold: %.1f%%, interval: %s, batch: %d)",
			cfg.ThresholdPercent, cfg.Interval, cfg.BatchSize)
	}
	go c.run()
}

// Stop signals the cleanup loop to stop.
func (c *Cleaner) Stop() {
	close(c.stopCh)
}

// UpdateConfig updates the cleanup configuration at runtime.
// Changes take effect on the next tick or immediately if the service
// was disabled and is now enabled.
func (c *Cleaner) UpdateConfig(cfg Config) {
	c.mu.Lock()
	c.cfg = cfg
	c.mu.Unlock()

	select {
	case c.resetCh <- struct{}{}:
	default:
	}
}

// run is the main cleanup loop.
func (c *Cleaner) run() {
	var ticker *time.Ticker
	var tickerInterval time.Duration

	defer func() {
		if ticker != nil {
			ticker.Stop()
		}
		log.Println("Cleanup service stopped")
	}()

	for {
		c.mu.RLock()
		cfg := c.cfg
		c.mu.RUnlock()

		if !cfg.Enabled {
			if ticker != nil {
				ticker.Stop()
				ticker = nil
				tickerInterval = 0
			}
			select {
			case <-time.After(5 * time.Second):
			case <-c.resetCh:
			case <-c.stopCh:
				return
			}
			continue
		}

		interval := cfg.Interval
		if interval <= 0 {
			interval = 15 * time.Minute
		}
		if ticker == nil || tickerInterval != interval {
			if ticker != nil {
				ticker.Stop()
			}
			tickerInterval = interval
			ticker = time.NewTicker(tickerInterval)
		}

		select {
		case <-ticker.C:
			c.check()
		case <-c.resetCh:
			ticker.Stop()
			ticker = nil
			tickerInterval = 0
		case <-c.stopCh:
			return
		}
	}
}

// check evaluates the current disk usage and deletes records if necessary.
func (c *Cleaner) check() {
	c.mu.RLock()
	cfg := c.cfg
	c.mu.RUnlock()

	usedPercent, err := diskUsagePercent(cfg.DiskPath)
	if err != nil {
		log.Printf("⚠️  Cleanup: failed to get disk usage for %s: %v", cfg.DiskPath, err)
		return
	}

	log.Printf("Cleanup: disk usage at %.1f%% (threshold: %.1f%%)", usedPercent, cfg.ThresholdPercent)

	if usedPercent < cfg.ThresholdPercent {
		return
	}

	log.Printf("⚠️  Cleanup: disk usage %.1f%% exceeds threshold %.1f%% — deleting %d old records",
		usedPercent, cfg.ThresholdPercent, cfg.BatchSize)

	deleted, err := c.deleteOldestRecords(cfg.BatchSize)
	if err != nil {
		log.Printf("❌ Cleanup: failed to delete records: %v", err)
		return
	}

	log.Printf("✓ Cleanup: deleted %d records", deleted)
}

// deleteOldestRecords removes the oldest N records from SystemEvents.
func (c *Cleaner) deleteOldestRecords(n int) (int64, error) {
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
//
// Uses stat.Bavail (blocks available to unprivileged users) — consistent with
// what the disk widget endpoint reports. stat.Bfree includes blocks reserved
// for root and would show a lower usage than what is actually visible.
func diskUsagePercent(path string) (float64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	avail := stat.Bavail * uint64(stat.Bsize) // available to unprivileged users

	if total == 0 {
		return 0, nil
	}

	used := total - avail
	return float64(used) / float64(total) * 100.0, nil
}

package cleanup

import (
	"database/sql"
	"log"
	"sync"
	"syscall"
	"time"
)

// Mode describes the current operating mode of the Cleaner.
type Mode string

const (
	ModeUnknown    Mode = "unknown"
	ModeMigrating  Mode = "migrating"
	ModePartition  Mode = "partition"
	ModeFailed     Mode = "failed"
)

// Status holds the current runtime state of the Cleaner, returned by Status().
type Status struct {
	Mode       Mode        `json:"mode"`
	Healthy    bool        `json:"healthy"`
	Error      string      `json:"error,omitempty"`
	Partitions []Partition `json:"partitions"`
}

// Partition represents a single MySQL partition with usage stats.
type Partition struct {
	Name   string  `json:"name"`
	Rows   int64   `json:"rows"`
	SizeMB float64 `json:"size_mb"`
}

// Cleaner periodically frees disk space by dropping the oldest weekly partition
// when disk usage exceeds the configured threshold.
//
// The Cleaner requires ALTER privilege on SystemEvents. If the table is not yet
// partitioned it migrates automatically on startup. If partitioning fails the
// Cleaner marks itself as failed and does NOT fall back to DELETE.
type Cleaner struct {
	db              *sql.DB
	mu              sync.RWMutex
	cfg             Config
	stopCh          chan struct{}
	restartCh       chan struct{}
	mode            Mode
	statusErr       string
	partitioned     bool
	alterPrivileged bool
	lastWeeklyMaint int
}

// Config holds the cleanup configuration.
type Config struct {
	Enabled          bool
	DiskPath         string
	ThresholdPercent float64
	Interval         time.Duration
}

// New creates a new Cleaner instance.
func New(db *sql.DB, cfg Config) *Cleaner {
	return &Cleaner{
		db:        db,
		cfg:       cfg,
		stopCh:    make(chan struct{}),
		restartCh: make(chan struct{}, 1),
		mode:      ModeUnknown,
	}
}

// Start launches the cleanup loop in a background goroutine.
func (c *Cleaner) Start() {
	c.mu.RLock()
	enabled := c.cfg.Enabled
	threshold := c.cfg.ThresholdPercent
	interval := c.cfg.Interval
	c.mu.RUnlock()

	if !enabled {
		log.Println("⏭  Cleanup service disabled")
		return
	}

	log.Printf("✓ Cleanup service starting (threshold: %.1f%%, interval: %s)",
		threshold, interval)

	c.initPartitions()

	go c.run()
}

// Stop signals the cleanup loop to stop.
func (c *Cleaner) Stop() {
	close(c.stopCh)
}

// UpdateConfig applies a new configuration live without restarting the process.
func (c *Cleaner) UpdateConfig(newCfg Config) {
	c.mu.Lock()
	oldInterval := c.cfg.Interval
	c.cfg = newCfg
	c.mu.Unlock()

	if newCfg.Interval != oldInterval {
		select {
		case c.restartCh <- struct{}{}:
		default:
		}
	}

	log.Printf("Cleanup: config updated live (enabled=%v, threshold=%.1f%%, interval=%s)",
		newCfg.Enabled, newCfg.ThresholdPercent, newCfg.Interval)
}

// Status returns the current runtime state including partition list.
func (c *Cleaner) Status() Status {
	c.mu.RLock()
	mode := c.mode
	statusErr := c.statusErr
	c.mu.RUnlock()

	s := Status{
		Mode:       mode,
		Healthy:    mode == ModePartition,
		Error:      statusErr,
		Partitions: []Partition{},
	}

	if mode == ModePartition {
		parts, err := listPartitions(c.db)
		if err != nil {
			log.Printf("Cleanup status: failed to list partitions: %v", err)
		} else {
			s.Partitions = parts
		}
	}

	return s
}

// --------------------------------------------------------------------------
// Internal
// --------------------------------------------------------------------------

func (c *Cleaner) setMode(m Mode, errMsg string) {
	c.mu.Lock()
	c.mode = m
	c.statusErr = errMsg
	c.mu.Unlock()
}

// initPartitions checks partition status and migrates if necessary.
func (c *Cleaner) initPartitions() {
	ok, err := HasAlterPrivilege(c.db)
	if err != nil || !ok {
		msg := "DB user lacks ALTER privilege on SystemEvents. " +
			"Run: GRANT ALTER ON Syslog.SystemEvents TO 'rsyslox'@'localhost'; FLUSH PRIVILEGES;"
		log.Printf("❌ Cleanup: %s", msg)
		c.setMode(ModeFailed, msg)
		c.alterPrivileged = false
		return
	}
	c.alterPrivileged = true

	partitioned, err := IsPartitioned(c.db)
	if err != nil {
		msg := "Could not check partition status: " + err.Error()
		log.Printf("❌ Cleanup: %s", msg)
		c.setMode(ModeFailed, msg)
		return
	}

	if !partitioned {
		c.setMode(ModeMigrating, "")
		log.Println("⚙  Cleanup: SystemEvents is not partitioned — starting automatic migration…")
		if err := MigrateToPartitions(c.db); err != nil {
			msg := "Migration failed: " + err.Error()
			log.Printf("❌ Cleanup: %s", msg)
			c.setMode(ModeFailed, msg)
			return
		}
	}

	if err := EnsureFuturePartitions(c.db); err != nil {
		log.Printf("⚠  Cleanup: EnsureFuturePartitions on init: %v", err)
	}

	c.partitioned = true
	c.setMode(ModePartition, "")
	log.Println("✓ Cleanup: partition mode active")
}

// run is the main cleanup loop.
func (c *Cleaner) run() {
	c.mu.RLock()
	interval := c.cfg.Interval
	c.mu.RUnlock()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.RLock()
			enabled := c.cfg.Enabled
			c.mu.RUnlock()
			if enabled {
				c.check()
				c.weeklyMaintenance()
			}

		case <-c.restartCh:
			ticker.Stop()
			c.mu.RLock()
			interval = c.cfg.Interval
			c.mu.RUnlock()
			ticker = time.NewTicker(interval)
			log.Printf("Cleanup: ticker restarted with new interval %s", interval)

		case <-c.stopCh:
			log.Println("Cleanup service stopped")
			return
		}
	}
}

// check evaluates disk usage and drops partitions if necessary.
func (c *Cleaner) check() {
	c.mu.RLock()
	mode := c.mode
	diskPath := c.cfg.DiskPath
	threshold := c.cfg.ThresholdPercent
	c.mu.RUnlock()

	if mode != ModePartition {
		log.Printf("⚠  Cleanup: skipping check — mode is %s", mode)
		return
	}

	usedPercent, err := diskUsagePercent(diskPath)
	if err != nil {
		log.Printf("⚠  Cleanup: failed to get disk usage for %s: %v", diskPath, err)
		return
	}

	log.Printf("Cleanup: disk usage at %.1f%% (threshold: %.1f%%)", usedPercent, threshold)

	if usedPercent < threshold {
		return
	}

	log.Printf("⚠  Cleanup: disk usage %.1f%% exceeds threshold %.1f%%", usedPercent, threshold)

	dropped, err := DropOldestPartitions(c.db)
	if err != nil {
		log.Printf("❌ Cleanup: partition drop error: %v", err)
		return
	}
	if dropped == 0 {
		log.Println("⚠  Cleanup: no droppable partitions found — cannot free space")
		return
	}
	log.Printf("✓ Cleanup: dropped %d partition(s)", dropped)

	if after, err := diskUsagePercent(diskPath); err == nil {
		log.Printf("Cleanup: disk usage after drop: %.1f%%", after)
	}
}

// weeklyMaintenance ensures future partitions are created once per calendar week.
func (c *Cleaner) weeklyMaintenance() {
	c.mu.RLock()
	mode := c.mode
	c.mu.RUnlock()

	if mode != ModePartition {
		return
	}

	_, currentWeek := time.Now().ISOWeek()
	if currentWeek == c.lastWeeklyMaint {
		return
	}
	if err := EnsureFuturePartitions(c.db); err != nil {
		log.Printf("⚠  Cleanup: weekly EnsureFuturePartitions: %v", err)
		return
	}
	c.lastWeeklyMaint = currentWeek
	log.Println("✓ Cleanup: weekly partition maintenance done")
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

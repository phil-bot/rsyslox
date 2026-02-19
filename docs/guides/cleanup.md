# Cleanup / Housekeeping

Automatic deletion of old log entries to prevent disk overflow.

## Overview

The cleanup service monitors the disk usage of a configured path (e.g. `/var/lib/mysql`) and automatically deletes the oldest entries from the `SystemEvents` table when usage exceeds a configurable threshold.

This prevents the operating system from crashing due to a full disk while keeping the database operational at all times — without requiring manual intervention or fixed retention periods.

## How It Works

1. Every `CLEANUP_INTERVAL` (default: `15m`) the service checks the current disk usage of `CLEANUP_DISK_PATH`.
2. If usage exceeds `CLEANUP_THRESHOLD_PERCENT` (default: `85`), a cleanup run is triggered.
3. The `CLEANUP_BATCH_SIZE` oldest records (ordered by `ReceivedAt ASC`) are deleted.
4. The result is logged and the process repeats on the next interval tick.

```
Every CLEANUP_INTERVAL
        │
        ▼
  Disk usage > CLEANUP_THRESHOLD_PERCENT?
        │                   │
       No                  Yes
        │                   │
      Skip        Delete CLEANUP_BATCH_SIZE oldest records
                            │
                            ▼
                       Log result
```

## Configuration

Enable and configure the cleanup service in your `.env` file:

```bash
# Enable the cleanup service
CLEANUP_ENABLED=true

# Filesystem path to monitor
CLEANUP_DISK_PATH=/var/lib/mysql

# Disk usage threshold in percent
CLEANUP_THRESHOLD_PERCENT=85

# Records to delete per cleanup run
CLEANUP_BATCH_SIZE=1000

# Check interval (Go duration format)
CLEANUP_INTERVAL=15m
```

### Parameters

| Variable | Default | Description |
|----------|---------|-------------|
| `CLEANUP_ENABLED` | `false` | Enables or disables the service |
| `CLEANUP_DISK_PATH` | `/var/lib/mysql` | Filesystem path/mount point to monitor |
| `CLEANUP_THRESHOLD_PERCENT` | `85` | Max allowed disk usage in percent (1–99) |
| `CLEANUP_BATCH_SIZE` | `1000` | Number of records deleted per run |
| `CLEANUP_INTERVAL` | `15m` | Check interval (e.g. `5m`, `15m`, `1h`) |

### CLEANUP_DISK_PATH

This must point to the **mount point of the partition** where MySQL/MariaDB stores its data files.

```bash
# Find the correct path
df -h /var/lib/mysql

# Example output:
# Filesystem   Size  Used Avail Use% Mounted on
# /dev/sdb1     50G   40G   10G  80% /var/lib/mysql
#
# → Use: CLEANUP_DISK_PATH=/var/lib/mysql
```

If MySQL data is stored on the root partition:

```bash
CLEANUP_DISK_PATH=/
```

### CLEANUP_THRESHOLD_PERCENT

Choose a value that gives enough headroom before the disk fills completely.

```bash
# React early — recommended for production
CLEANUP_THRESHOLD_PERCENT=80

# Default — balanced
CLEANUP_THRESHOLD_PERCENT=85

# More permissive — keeps more history
CLEANUP_THRESHOLD_PERCENT=90
```

!> **Important:** Do not set this too close to `100`. MySQL needs free space for transaction logs and temporary files.

### CLEANUP_BATCH_SIZE

Controls how many records are removed per cleanup run.

```bash
# Low volume — delete slowly, keep more history
CLEANUP_BATCH_SIZE=500

# Default — balanced
CLEANUP_BATCH_SIZE=1000

# High volume — free space quickly
CLEANUP_BATCH_SIZE=5000
```

Larger values free up space faster but create more database load per cycle.

### CLEANUP_INTERVAL

How often the disk is checked. Uses [Go duration format](https://pkg.go.dev/time#ParseDuration).

```bash
# High volume systems
CLEANUP_INTERVAL=5m

# Default
CLEANUP_INTERVAL=15m

# Low volume systems
CLEANUP_INTERVAL=1h
```

## Log Output

When the cleanup service is active, you will see entries like:

```
✓ Cleanup service started (threshold: 85.0%, interval: 15m0s, batch: 1000)
Cleanup: disk usage at 72.3% (threshold: 85.0%)
Cleanup: disk usage at 86.1% (threshold: 85.0%)
⚠️  Cleanup: disk usage 86.1% exceeds threshold 85.0% — deleting 1000 old records
✓ Cleanup: deleted 1000 records
```

If the service is disabled:

```
⏭  Cleanup service disabled
```

## Recommended Configurations

### Production (default)

```bash
CLEANUP_ENABLED=true
CLEANUP_DISK_PATH=/var/lib/mysql
CLEANUP_THRESHOLD_PERCENT=85
CLEANUP_BATCH_SIZE=1000
CLEANUP_INTERVAL=15m
```

### High Volume Systems

```bash
CLEANUP_ENABLED=true
CLEANUP_DISK_PATH=/var/lib/mysql
CLEANUP_THRESHOLD_PERCENT=75
CLEANUP_BATCH_SIZE=10000
CLEANUP_INTERVAL=5m
```

### Low Volume / Long Retention

```bash
CLEANUP_ENABLED=true
CLEANUP_DISK_PATH=/var/lib/mysql
CLEANUP_THRESHOLD_PERCENT=90
CLEANUP_BATCH_SIZE=500
CLEANUP_INTERVAL=1h
```

## Monitoring

### Check Disk Usage

```bash
# Check the monitored partition
df -h /var/lib/mysql

# Check table size in MySQL
mysql -u rsyslog -p Syslog -e "
SELECT
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
    table_rows AS 'Rows (approx)'
FROM information_schema.TABLES
WHERE table_schema = 'Syslog' AND table_name = 'SystemEvents';
"

# Check oldest and newest record
mysql -u rsyslog -p Syslog -e "
SELECT
    MIN(ReceivedAt) AS oldest,
    MAX(ReceivedAt) AS newest,
    COUNT(*) AS total
FROM SystemEvents;
"
```

### Watch Cleanup in Real Time

```bash
# Follow service logs and filter cleanup messages
sudo journalctl -u rsyslox -f | grep -i cleanup
```

## Troubleshooting

### Cleanup Not Triggering

```bash
# Verify CLEANUP_ENABLED is exactly "true" (not "1" or "yes")
grep CLEANUP_ENABLED /opt/rsyslox/.env

# Check actual disk usage
df -h /var/lib/mysql

# Make sure the path is on the correct partition
df -h | grep mysql
```

### Disk Still Fills Up

If disk usage grows faster than the cleanup can handle, increase the aggressiveness:

```bash
CLEANUP_THRESHOLD_PERCENT=70
CLEANUP_BATCH_SIZE=5000
CLEANUP_INTERVAL=5m
```

### Records Not Being Deleted

```
❌ Cleanup: failed to delete records: ...
```

If the database user has only `SELECT` permissions (read-only user), the cleanup service needs `DELETE` as well:

```sql
GRANT DELETE ON Syslog.SystemEvents TO 'rsyslog'@'localhost';
FLUSH PRIVILEGES;
```

?> **Note:** If you are using a dedicated read-only API user for security reasons, consider creating a separate cleanup user with only `DELETE` permission, or grant `DELETE` only on `SystemEvents`.

## More Resources

- [Configuration Reference](configuration.md) - All configuration options
- [Performance Guide](performance.md) - Database optimization
- [Deployment Guide](deployment.md) - Production setup
- [Troubleshooting](troubleshooting.md) - General troubleshooting

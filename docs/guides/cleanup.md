# Cleanup / Housekeeping

Automatic deletion of old log entries to prevent disk overflow. Configured in **Admin → Database → Log Cleanup**.

## Overview

The cleanup service monitors disk usage at a configured path. When usage exceeds the threshold it frees space by **dropping the oldest weekly partition** from `SystemEvents`. This gives back disk space to the OS immediately — no `OPTIMIZE TABLE` required.

## How It Works

```
Startup
  │
  ├─ Is SystemEvents partitioned?
  │     No  → Migrate automatically (weekly partitions)
  │     Yes → Continue
  │
  ├─ Does the DB user have ALTER privilege?
  │     No  → Fall back to legacy DELETE (disk space NOT reclaimed)
  │     Yes → Partition mode active
  │
Every <interval>
  │
  ├─ Disk usage > threshold?
  │     No  → Skip
  │     Yes → Drop oldest partition(s) (max 3 per run)
  │            └─ Log new usage after drop
  │
  └─ Once per week → Ensure next 6 weeks of partitions exist
```

## Partition Strategy

- **Granularity:** one partition per ISO calendar week (`p2026_w23`)
- **Pre-creation:** 6 weeks into the future, maintained automatically
- **Catch-all:** `p_future` is always kept and never dropped
- **Max drops per run:** 3 partitions (to avoid accidental bulk deletion)

## Database Requirements

The cleanup service needs **ALTER** privilege in addition to the standard SELECT:

```sql
-- Read-only access (minimum for rsyslox without cleanup)
GRANT SELECT ON Syslog.SystemEvents TO 'rsyslox'@'localhost';

-- Full cleanup support (required for partition mode)
GRANT SELECT, DELETE, ALTER ON Syslog.SystemEvents TO 'rsyslox'@'localhost';
FLUSH PRIVILEGES;
```

> If ALTER is not granted, rsyslox falls back to the legacy `DELETE`-based cleanup.
> Old rows are removed but **disk space is not reclaimed** until MySQL reclaims it
> internally. A warning is logged at startup.

## Automatic Migration

On first start after enabling cleanup, rsyslox automatically migrates an unpartitioned
`SystemEvents` table:

1. Extends the PRIMARY KEY to `(ID, ReceivedAt)`
2. Creates weekly partitions from the oldest existing entry up to 6 weeks ahead
3. Adds a `p_future` catch-all partition

> **Note:** Migration can be slow on large tables (millions of rows).
> rsyslox remains fully operational during migration — it runs in the background.

## Configuration

Configure via **Admin panel → Database → Log Cleanup**. Changes take effect immediately — no restart needed.

| Setting | Description | Default |
|---|---|---|
| Enabled | Toggle the cleanup service | off |
| Disk path | Mount point to monitor | `/var/lib/mysql` |
| Threshold % | Trigger cleanup above this disk usage | 85 % |
| Batch size | Rows deleted per run (legacy fallback only) | 1 000 |
| Interval | Seconds between disk checks | 900 |

### Disk Path

This must be the **mount point of the partition** where MySQL/MariaDB stores its data:

```bash
# Find the correct value
df -h /var/lib/mysql

# Example output:
# Filesystem   Size  Used Avail Use% Mounted on
# /dev/sdb1     50G   40G   10G  80% /var/lib/mysql
#
# → Use: /var/lib/mysql
```

!> The cleanup service checks disk usage on the **local filesystem**. It only works
correctly if the database runs on the same host as rsyslox.

## Log Output

**Partition mode (normal):**
```
✓ Cleanup service starting (threshold: 85.0%, interval: 15m0s)
⚙  Cleanup: SystemEvents is not partitioned — starting automatic migration…
✓ Partition migration complete
✓ Cleanup: partition mode active
Cleanup: disk usage at 72.3% (threshold: 85.0%)
Cleanup: disk usage at 87.1% (threshold: 85.0%)
⚠  Cleanup: disk usage 87.1% exceeds threshold 85.0%
✓ Partition dropped: p2026_w18
Cleanup: disk usage after drop: 81.4%
✓ Cleanup: weekly partition maintenance done
```

**Legacy fallback (ALTER not granted):**
```
⚠  Cleanup: DB user lacks ALTER privilege
   Falling back to legacy DELETE cleanup — disk space will NOT be reclaimed after deletion.
   Grant ALTER on Syslog.SystemEvents to the rsyslox DB user to enable partition mode.
```

```bash
# Watch cleanup messages in real time
sudo journalctl -u rsyslox -f | grep -i cleanup
```

## Manual Partition Management

### View current partitions

```sql
SELECT PARTITION_NAME, TABLE_ROWS,
       PARTITION_DESCRIPTION AS less_than
FROM INFORMATION_SCHEMA.PARTITIONS
WHERE TABLE_NAME   = 'SystemEvents'
  AND TABLE_SCHEMA = 'Syslog'
ORDER BY PARTITION_DESCRIPTION;
```

### Manually drop a partition

```sql
ALTER TABLE SystemEvents DROP PARTITION p2026_w18;
```

### Manually add a partition for a specific week

```sql
-- Example: add week 30 of 2026 (Mon 2026-07-20 – Sun 2026-07-26)
ALTER TABLE SystemEvents REORGANIZE PARTITION p_future INTO (
    PARTITION p2026_w30 VALUES LESS THAN ('2026-07-27'),
    PARTITION p_future  VALUES LESS THAN (MAXVALUE)
);
```

## Recommended Configurations

**Production (default)**

| Setting | Value |
|---|---|
| Enabled | true |
| Disk path | `/var/lib/mysql` |
| Threshold % | 85 |
| Interval | 900 s |

**High Volume / Tight Disk**

| Setting | Value |
|---|---|
| Threshold % | 75 |
| Interval | 300 s |

**Low Volume / Long Retention**

| Setting | Value |
|---|---|
| Threshold % | 90 |
| Interval | 3 600 s |

## Monitoring

```bash
# Check the monitored partition
df -h /var/lib/mysql

# Check table size and row count in MySQL
mysql -u rsyslox -p Syslog -e "
SELECT
  PARTITION_NAME,
  TABLE_ROWS,
  ROUND((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024, 1) AS size_mb
FROM INFORMATION_SCHEMA.PARTITIONS
WHERE TABLE_NAME   = 'SystemEvents'
  AND TABLE_SCHEMA = 'Syslog'
ORDER BY PARTITION_DESCRIPTION;
"

# Oldest and newest record
mysql -u rsyslox -p Syslog -e "
SELECT MIN(ReceivedAt) AS oldest, MAX(ReceivedAt) AS newest, COUNT(*) AS total
FROM SystemEvents;
"
```

## Troubleshooting

**Cleanup not triggering**

Check actual disk usage — it may genuinely be below the threshold:
```bash
df -h /var/lib/mysql
```

Verify the service is enabled in **Admin → Database → Log Cleanup**.

**"DB user lacks ALTER privilege" warning**

Grant the required privilege and restart rsyslox:
```sql
GRANT ALTER ON Syslog.SystemEvents TO 'rsyslox'@'localhost';
FLUSH PRIVILEGES;
```

**Migration is slow**

This is expected on large tables. rsyslox continues to serve requests normally
during migration. Check progress with:
```sql
SHOW PROCESSLIST;
```

**Disk still fills up after cleanup**

The partition was dropped but the OS did not reclaim space immediately — this is
unusual in partition mode but can happen on some filesystems. Check:
```bash
# Verify the inode is actually released
lsof | grep -i mysql | grep deleted
```

## More Resources

- [Configuration Reference](../getting-started/configuration.md)
- [Security Guide](../guides/security.md) — DB user permissions
- [Performance Guide](performance.md)
- [Troubleshooting](troubleshooting.md)

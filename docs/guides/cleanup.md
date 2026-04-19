# Cleanup / Housekeeping

Automatic deletion of old log entries to prevent disk overflow. Configured in **Admin → Database → Log Cleanup**.

## Overview

The cleanup service monitors disk usage at a configured path and automatically deletes the oldest entries from `SystemEvents` when usage exceeds a threshold. This prevents disk-full crashes without requiring fixed retention periods or manual intervention.

## How It Works

```
Every <interval>
       │
       ▼
 Disk usage > threshold?
       │             │
      No            Yes
       │             │
     Skip   Delete <batch_size> oldest records (ordered by ReceivedAt ASC)
                     │
                     ▼
                Log result, repeat next tick
```

## Configuration

Configure via **Admin panel → Database → Log Cleanup**. Changes take effect immediately — no restart needed.

| Setting | Description | Default |
|---|---|---|
| Enabled | Toggle the cleanup service | off |
| Disk path | Mount point to monitor | `/var/lib/mysql` |
| Threshold % | Trigger cleanup above this disk usage | 85 % |
| Batch size | Rows deleted per cleanup run | 1 000 |
| Interval | Seconds between checks | 900 |

### Disk Path

This must be the **mount point of the partition** where MySQL/MariaDB stores its data files:

```bash
# Find the correct value
df -h /var/lib/mysql

# Example output:
# Filesystem   Size  Used Avail Use% Mounted on
# /dev/sdb1     50G   40G   10G  80% /var/lib/mysql
#
# → Use: /var/lib/mysql
```

If MySQL data is on the root partition:
```
/
```

### Threshold Percent

Choose a value that gives enough headroom before the disk fills completely.

```
80 % — react early, recommended for production
85 % — default, balanced
90 % — more permissive, keeps more history
```

!> Do not set this close to `100`. MySQL needs free space for transaction logs and temp files.

### Batch Size

```
500   — low volume, delete slowly
1000  — default
5000  — high volume, free space quickly
```

Larger values free space faster but create more database load per cycle.

### Interval

How often the disk is checked (in seconds). Examples: `300` (5 min), `900` (15 min, default), `3600` (1 h).

## Database Permissions

The cleanup service needs `DELETE` on `SystemEvents`. If you use a read-only database user, grant `DELETE` as well:

```sql
GRANT DELETE ON Syslog.SystemEvents TO 'rsyslox'@'localhost';
FLUSH PRIVILEGES;
```

## Log Output

When active, the service logs to systemd journal:

```
✓ Cleanup service started (threshold: 85.0%, interval: 15m0s, batch: 1000)
Cleanup: disk usage at 72.3% (threshold: 85.0%)
Cleanup: disk usage at 86.1% (threshold: 85.0%)
⚠️  Cleanup: disk usage 86.1% exceeds threshold 85.0% — deleting 1000 old records
✓ Cleanup: deleted 1000 records
```

```bash
# Watch cleanup messages in real time
sudo journalctl -u rsyslox -f | grep -i cleanup
```

## Recommended Configurations

**Production (default)**

| Setting | Value |
|---|---|
| Enabled | true |
| Disk path | `/var/lib/mysql` |
| Threshold % | 85 |
| Batch size | 1 000 |
| Interval | 900 s |

**High Volume Systems**

| Setting | Value |
|---|---|
| Threshold % | 75 |
| Batch size | 10 000 |
| Interval | 300 s |

**Low Volume / Long Retention**

| Setting | Value |
|---|---|
| Threshold % | 90 |
| Batch size | 500 |
| Interval | 3 600 s |

## Monitoring

```bash
# Check the monitored partition
df -h /var/lib/mysql

# Check table size and row count in MySQL
mysql -u rsyslox -p Syslog -e "
SELECT
  ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
  table_rows AS 'Rows (approx)'
FROM information_schema.TABLES
WHERE table_schema = 'Syslog' AND table_name = 'SystemEvents';
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

**Disk still fills up**

Increase aggressiveness:
```
Threshold: 70 %
Batch size: 5 000
Interval: 300 s
```

**"Failed to delete records" error**

The database user lacks `DELETE` permission — see [Database Permissions](#database-permissions) above.

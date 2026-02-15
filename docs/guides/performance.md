# Performance Guide

Optimization and tuning for rsyslog REST API.

## Performance Metrics

Typical performance with 1M+ logs:

| Operation | Response Time | Throughput |
|-----------|---------------|------------|
| Health check | <5ms | 10,000+ req/s |
| Simple query (limit 10) | 10-50ms | 1,000 req/s |
| Complex filter | 50-200ms | 500 req/s |
| Aggregation (meta) | 100-500ms | 200 req/s |

## Database Optimization

### Indexes

Automatically created:
```sql
CREATE INDEX idx_receivedat ON SystemEvents (ReceivedAt);
CREATE INDEX idx_host_time ON SystemEvents (FromHost, ReceivedAt);
CREATE INDEX idx_priority ON SystemEvents (Priority);
CREATE INDEX idx_facility ON SystemEvents (Facility);
CREATE INDEX IF NOT EXISTS idx_message USING FULLTEXT (Message);
```

Additional indexes:
```sql
CREATE INDEX idx_syslogtag ON SystemEvents (SysLogTag);
CREATE INDEX idx_host_priority ON SystemEvents (FromHost, Priority);
```

### Query Optimization

```bash
# Good: Uses indexes
?Priority=3&start_date=...&limit=100

# Slow: Full table scan
?Message=keyword  # Full-text search

# Best: Combine indexed fields
?FromHost=web01&Priority=3&start_date=...
```

### Connection Pool

```go
// Increase for high traffic
db.SetMaxOpenConns(50)    // Default: 25
db.SetMaxIdleConns(10)    // Default: 5
db.SetConnMaxLifetime(5 * time.Minute)
```

### MySQL Configuration

```ini
# my.cnf optimizations
[mysqld]
innodb_buffer_pool_size = 2G
innodb_log_file_size = 512M
innodb_flush_log_at_trx_commit = 2
query_cache_size = 256M
query_cache_type = 1
max_connections = 200
```

## API Optimization

### Pagination

```bash
# Bad: Large results
?limit=10000  # Slow!

# Good: Reasonable pages
?limit=100&offset=0  # Fast

# Best: Time-based pagination
?start_date=...&end_date=...&limit=100
```

### Filtering

```bash
# Fast: Indexed fields
?Priority=3
?FromHost=webserver01
?start_date=...

# Slower: Full-text search
?Message=keyword

# Optimize: Combine filters
?FromHost=web01&Priority=3  # Uses multiple indexes
```

## Monitoring

### Application Metrics

```bash
# Query performance
time curl -H "X-API-Key: $KEY" "http://localhost:8000/logs?limit=1000"

# Concurrent requests
ab -n 1000 -c 100 http://localhost:8000/health
```

### Database Stats

```sql
-- Slow queries
SHOW PROCESSLIST;

-- Index usage
SHOW INDEX FROM SystemEvents;

-- Table size
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS "Size (MB)"
FROM information_schema.TABLES
WHERE table_schema = "Syslog";
```

## Scaling

### Vertical Scaling

1. **CPU:** More cores for concurrent queries
2. **RAM:** Larger buffer pool for MySQL
3. **Storage:** SSD for faster I/O

### Horizontal Scaling

```
Load Balancer (nginx)
├── API Instance 1
├── API Instance 2
└── API Instance 3
     ↓
MySQL Master/Slave Replication
```

See [Deployment Guide](deployment.md#scaling) for details.

## Benchmarking

### ApacheBench

```bash
# Health endpoint
ab -n 10000 -c 100 http://localhost:8000/health

# API endpoint
ab -n 1000 -c 10 \
   -H "X-API-Key: $KEY" \
   "http://localhost:8000/logs?limit=10"
```

### wrk

```bash
# Install wrk
sudo apt-get install wrk

# Benchmark
wrk -t4 -c100 -d30s \
    -H "X-API-Key: $KEY" \
    http://localhost:8000/logs?limit=10
```

## Troubleshooting Performance

**Slow queries:**
- Check database indexes
- Reduce time window
- Use smaller limits

**High memory usage:**
- Reduce connection pool size
- Lower result limits
- Optimize MySQL buffer pool

**High CPU:**
- Enable query cache
- Add more indexes
- Scale horizontally

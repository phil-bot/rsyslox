# Performance Guide

[← Back to overview](index.md)

Performance optimization and benchmarking for rsyslog REST API.

## 📊 Performance characteristics

### Typical performance

**Hardware: 2 CPU, 4GB RAM, SSD**

| Metric | Value |
|--------|------|
| Requests/second | 500-1000 |
| Response Time (avg) | 10-50ms |
| Response Time (p95) | 100-200ms |
| Database Queries | 1-2 per request |
| Memory Usage | 50-100MB |
| CPU Usage | 5-15% |

**With 1 million logs in DB:**
- Simple Query (`limit=10`): ~10-20ms
- Filtered Query: ~20-50ms
- Complex Multi-Filter: ~50-100ms

---

## ⚡ Optimizations

### 1. database indexing

**Auto-Created Indexes:**

The API automatically creates the following indexes at startup:

```sql
CREATE INDEX idx_receivedat ON SystemEvents (ReceivedAt);
CREATE INDEX idx_host_time ON SystemEvents (FromHost, ReceivedAt);
CREATE INDEX idx_priority ON SystemEvents (Priority);
CREATE INDEX idx_facility ON SystemEvents (Facility);
CREATE FULLTEXT INDEX ON SystemEvents (Message);
```

**Check manually:**

```sql
-- Show indexes
SHOW INDEX FROM SystemEvents;

-- Check index usage
EXPLAIN SELECT * FROM SystemEvents
WHERE ReceivedAt > '2025-01-01'
AND FromHost = 'webserver01'
LIMIT 10;
```

**Additional indexes:**

```sql
-- If a lot is filtered by SysLogTag
CREATE INDEX idx_syslogtag ON SystemEvents (SysLogTag);

-- Composite index for frequent combinations
CREATE INDEX idx_host_priority_time
  ON SystemEvents (FromHost, Priority, ReceivedAt);
```

### 2. query optimization

**Best Practices:**

```bash
# ✅ FAST - Small time windows
?start_date=2025-02-09T10:00:00Z&end_date=2025-02-09T11:00:00Z&limit=100

# ❌ SLOW - Large time windows
?start_date=2025-01-01T00:00:00Z&end_date=2025-02-09T23:59:59Z&limit=1000

# ✅ FAST - Indexed fields
?FromHost=webserver01&Priority=3

# ⚠️ SLOW - Text search (full text, but still slower)
?Message=error
```

**Use pagination correctly:**

```bash
# ✅ CORRECT - Small pages
?limit=100&offset=0
?limit=100&offset=100
?limit=100&offset=200

# ❌ WRONG - Large offset (becomes slower!)
?limit=100&offset=50000
```

### 3. connection pooling

**In main.go (already configured):**

```go
db.SetMaxOpenConns(25) // Max parallel connections
db.SetMaxIdleConns(5) // Idle connections in the pool
db.SetConnMaxLifetime(5 * time.Minute)
```

**Adapt for more load:**

```go
// For high traffic:
db.SetMaxOpenConns(50)
db.SetMaxIdleConns(10)
```

### 4. database configuration

**MySQL/MariaDB Tuning:**

```ini
# /etc/mysql/mariadb.conf.d/50-server.cnf

[mysqld]
# InnoDB Buffer Pool (50-70% of RAM)
innodb_buffer_pool_size = 2G

# Query Cache (MySQL 5.7, deprecated in 8.0)
query_cache_size = 64M
query_cache_type = 1

# Connections
max_connections = 100

# Logging (disable in production for performance)
slow_query_log = 0
general_log = 0
```

**After changes:**
```bash
sudo systemctl restart mysql
```

---

## 🔍 Monitoring

### Application Metrics

**Simple Stats Script:**

```bash
#!/bin/bash
# /usr/local/bin/api-stats.sh

echo "=== API Performance Stats ==="

# Response Times (nginx)
echo "Average Response Time (last 1000 requests):"
tail -1000 /var/log/nginx/rsyslog-api-access.log | \
  awk '{print $NF}' | \
  awk '{sum+=$1; count++} END {print sum/count " seconds"}'

# Requests per minute
echo ""
echo "Requests per minute (last hour):"
tail -60000 /var/log/nginx/rsyslog-api-access.log | \
  awk '{print $4}' | cut -d: -f2 | sort | uniq -c | \
  awk '{sum+=$1; count++} END {print sum/count " req/min"}'

# Status codes
echo ""
echo "Status codes (last 1000 requests):"
tail -1000 /var/log/nginx/rsyslog-api-access.log | \
  awk '{print $9}' | sort | uniq -c
```

### Database Monitoring

```bash
#!/bin/bash
# db-stats.sh

mysql -u root -p <<EOF
-- Connection Stats
SHOW STATUS LIKE 'Threads_connected';
SHOW STATUS LIKE 'Max_used_connections';

-- Query Stats
SHOW STATUS LIKE 'Queries';
SHOW STATUS LIKE 'Slow_queries';

-- InnoDB Buffer Pool
SHOW STATUS LIKE 'Innodb_buffer_pool%';

-- Table Stats
SELECT
  table_name,
  table_rows,
  ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
FROM information_schema.TABLES
WHERE table_schema = 'Syslog';
EOF
```

### Prometheus + Grafana (Advanced)

**Prometheus nginx-exporter:**

```bash
# Installation
sudo apt-get install prometheus-nginx-exporter

# Config
sudo nano /etc/default/prometheus-nginx-exporter
```

```bash
ARGS="--nginx.scrape-uri=http://localhost/nginx_status"
```

**nginx config:**

```nginx
server {
    listen 127.0.0.1:80;
    location /nginx_status {
        stub_status on;
        access_log off;
    }
}
```

---

## 🧪 Benchmarking

### Simple Benchmark (ApacheBench)

```bash
# Install
sudo apt-get install apache2-utils

# Health Endpoint (Baseline)
ab -n 10000 -c 100 http://localhost:8000/health

# Interpret results:
# Requests per second: X [#/sec] (mean)
# Time per request: X [ms] (mean)
```

### API Endpoint Benchmark

```bash
# With API Key
API_KEY="your-api-key"

# Simple query
ab -n 1000 -c 10 \
  -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?limit=10"

# Filtered Query
from -n 500 -c 10 \
  -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=50"
```

### Load Testing (wrk)

```bash
# Install
sudo apt-get install wrk

# 30 seconds, 10 threads, 100 connections
wrk -t10 -c100 -d30s \
  -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?limit=10"

# Output:
# Requests/sec: X
# Transfer/sec: X
# Latency Distribution
```

### Database Query Performance

```sql
-- Activate slow queries
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1; -- > 1 second

-- View slow query log
-- /var/log/mysql/mysql-slow.log

-- Analyze query performance
EXPLAIN SELECT * FROM SystemEvents
WHERE ReceivedAt > '2025-01-01'
AND FromHost = 'webserver01'
LIMIT 10;
```

---

## 📈 Scaling Strategies

### Vertical Scaling

**CPU:**
- API is mostly I/O-bound (waits for database)
- More CPU helps with high number of requests
- 2-4 cores recommended for production

**Memory:**
- API itself: 50-100 MB
- Database Buffer Pool: 50-70% of available RAM
- Minimum: 1GB total, Recommended: 4GB+

**Storage:**
- SSD recommended (Database)
- Logs grow: approx. 1MB per 10,000 entries
- Index size: approx. 30-50% of the data size

### Horizontal scaling

**Load Balancer Setup:**

```nginx
# nginx load balancer
upstream rsyslog_api_cluster {
    least_conn; # Load Balancing Method
    
    server 192.168.1.10:8000 weight=1;
    server 192.168.1.11:8000 weight=1;
    server 192.168.1.12:8000 weight=1;
    
    # Health Check (nginx Plus)
    # health_check interval=5s;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;
    
    location / {
        proxy_pass http://rsyslog_api_cluster;
        # ...
    }
}
```

**Shared Database:**
- All API instances → same database
- Database becomes bottleneck
- Consider: Read Replicas

### Database Scaling

**Read Replicas:**

```bash
# Master-Slave Replication
# Master: Writes (rsyslog)
# Slaves: Reads (API)

# In .env of different API instances:
# API-1: DB_HOST=slave1.db.local
# API-2: DB_HOST=slave2.db.local
```

**Partitioning:**

```sql
-- Partitioning by date
ALTER TABLE SystemEvents
PARTITION BY RANGE (YEAR(ReceivedAt)) (
    PARTITION p2023 VALUES LESS THAN (2024),
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION pmax VALUES LESS THAN MAXVALUE
);
```

---

## 🎯 Performance Tuning Checklist

### Application Level

- [ ] Database indexes created
- [ ] Connection Pool configured
- [ ] Small time windows for queries
- [ ] Pagination instead of large limits
- [ ] Logging optimized (no debug in production)

### Database level

- [ ] InnoDB Buffer Pool optimized
- [ ] Slow Query Log activated
- [ ] Indexes for frequent queries
- [ ] Query Cache (MySQL 5.7)
- [ ] Old logs archived/deleted

### System Level

- [ ] SSD for database
- [ ] Sufficient RAM
- [ ] Firewall optimized
- [ ] nginx/Apache configured
- [ ] Rate limiting activated

---

## 🔧 Troubleshooting Performance Issues

### Symptom: Slow responses

**Diagnosis:**

```bash
#1 Check database
mysql Syslog -e "SHOW PROCESSLIST"

# 2. slow queries
tail -100 /var/log/mysql/mysql-slow.log

# 3. API logs
sudo journalctl -u rsyslog-rest-api -n 100

# 4. system resources
top
htop
iostat -x 1
```

**Solutions:**

- Indexes missing? → Create
- Large time windows? → Restrict
- Many connections? → Enlarge pool
- Disk I/O high? → Use SSD

### Symptom: High memory usage

**Diagnosis:**

```bash
# API Memory
ps aux | grep rsyslog-rest-api

# Database Memory
sudo mysqladmin -u root -p memory
```

**Solutions:**

- Memory leak? → Update to latest version
- DB buffer too large? → Reduce
- Too many connections? → Set limit

### Symptom: High CPU Usage

**Diagnosis:**

```bash
# CPU per process
top -p $(pgrep rsyslog-rest-api)

# Database CPU
mysqladmin -u root -p processlist
```

**Solutions:**

- Inefficient queries? → Use EXPLAIN
- Many requests? → Rate limiting
- Missing indexes? → Create

---

## 📊 Benchmarks

### Test Setup

- **Hardware:** 4 CPU, 8GB RAM, SSD
- **Database:** 1M logs
- **Network:** Localhost

### Results

| Query Type | Requests/s | Avg Response Time |
|------------|-----------|-------------------|
| Health Check | 5000 | 2ms |
| Simple (limit=10) | 800 | 12ms |
| Filtered (1 filter) | 600 | 16ms |
| Multi-Filter (3 filters) | 400 | 25ms |
| Large Result (limit=1000) | 50 | 200ms |

---

## 💡 Best practices

1. **Use small time windows** - Max. 7-30 days
2. **Limit limit** - Max. 100-500, not 1000
3. **Use pagination** - Instead of large limits
4. **Filter indexed fields** - FromHost, Priority, ReceivedAt
5. **Archive old logs** - Clean up regularly
6. **Set up monitoring** - Detect problems at an early stage
7. **Load testing** - Before production deployment

---

[← Back to overview](index.md)

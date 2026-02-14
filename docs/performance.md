# Performance Guide

[← Zurück zur Übersicht](index.md)

Performance-Optimierung und Benchmarking für rsyslog REST API.

## 📊 Performance Charakteristiken

### Typische Performance

**Hardware: 2 CPU, 4GB RAM, SSD**

| Metrik | Wert |
|--------|------|
| Requests/Sekunde | 500-1000 |
| Response Time (avg) | 10-50ms |
| Response Time (p95) | 100-200ms |
| Database Queries | 1-2 pro Request |
| Memory Usage | 50-100MB |
| CPU Usage | 5-15% |

**Mit 1 Million Logs in DB:**
- Simple Query (`limit=10`): ~10-20ms
- Filtered Query: ~20-50ms
- Complex Multi-Filter: ~50-100ms

---

## ⚡ Optimierungen

### 1. Database Indexing

**Auto-Created Indexes:**

Die API erstellt automatisch folgende Indexes beim Start:

```sql
CREATE INDEX idx_receivedat ON SystemEvents (ReceivedAt);
CREATE INDEX idx_host_time ON SystemEvents (FromHost, ReceivedAt);
CREATE INDEX idx_priority ON SystemEvents (Priority);
CREATE INDEX idx_facility ON SystemEvents (Facility);
CREATE FULLTEXT INDEX ON SystemEvents (Message);
```

**Manuell prüfen:**

```sql
-- Indexes anzeigen
SHOW INDEX FROM SystemEvents;

-- Index-Nutzung prüfen
EXPLAIN SELECT * FROM SystemEvents 
WHERE ReceivedAt > '2025-01-01' 
AND FromHost = 'webserver01' 
LIMIT 10;
```

**Zusätzliche Indexes:**

```sql
-- Falls viel nach SysLogTag gefiltert wird
CREATE INDEX idx_syslogtag ON SystemEvents (SysLogTag);

-- Composite Index für häufige Kombinationen
CREATE INDEX idx_host_priority_time 
  ON SystemEvents (FromHost, Priority, ReceivedAt);
```

### 2. Query Optimization

**Best Practices:**

```bash
# ✅ SCHNELL - Kleine Zeitfenster
?start_date=2025-02-09T10:00:00Z&end_date=2025-02-09T11:00:00Z&limit=100

# ❌ LANGSAM - Große Zeitfenster
?start_date=2025-01-01T00:00:00Z&end_date=2025-02-09T23:59:59Z&limit=1000

# ✅ SCHNELL - Indexierte Felder
?FromHost=webserver01&Priority=3

# ⚠️ LANGSAM - Text-Suche (Fulltext, aber trotzdem langsamer)
?Message=error
```

**Pagination richtig nutzen:**

```bash
# ✅ RICHTIG - Kleine Pages
?limit=100&offset=0
?limit=100&offset=100
?limit=100&offset=200

# ❌ FALSCH - Große Offset (wird langsamer!)
?limit=100&offset=50000
```

### 3. Connection Pooling

**In main.go (bereits konfiguriert):**

```go
db.SetMaxOpenConns(25)     // Max parallele Connections
db.SetMaxIdleConns(5)      // Idle Connections im Pool
db.SetConnMaxLifetime(5 * time.Minute)
```

**Anpassen für mehr Last:**

```go
// Für High-Traffic:
db.SetMaxOpenConns(50)
db.SetMaxIdleConns(10)
```

### 4. Database Configuration

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

# Logging (disable in production für Performance)
slow_query_log = 0
general_log = 0
```

**Nach Änderungen:**
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

# Status Codes
echo ""
echo "Status Codes (last 1000 requests):"
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

# Results interpretieren:
# Requests per second: X [#/sec] (mean)
# Time per request: X [ms] (mean)
```

### API Endpoint Benchmark

```bash
# Mit API Key
API_KEY="your-api-key"

# Simple Query
ab -n 1000 -c 10 \
  -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?limit=10"

# Filtered Query
ab -n 500 -c 10 \
  -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=50"
```

### Load Testing (wrk)

```bash
# Install
sudo apt-get install wrk

# 30 Sekunden, 10 Threads, 100 Connections
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
-- Slow Queries aktivieren
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;  -- > 1 Sekunde

-- Slow Query Log ansehen
-- /var/log/mysql/mysql-slow.log

-- Query Performance analysieren
EXPLAIN SELECT * FROM SystemEvents 
WHERE ReceivedAt > '2025-01-01' 
AND FromHost = 'webserver01' 
LIMIT 10;
```

---

## 📈 Scaling Strategies

### Vertical Scaling

**CPU:**
- API ist meist I/O-bound (wartet auf Database)
- Mehr CPU hilft bei hoher Request-Anzahl
- 2-4 Cores empfohlen für Production

**Memory:**
- API selbst: 50-100 MB
- Database Buffer Pool: 50-70% of available RAM
- Minimum: 1GB total, Empfohlen: 4GB+

**Storage:**
- SSD empfohlen (Database)
- Logs wachsen: ca. 1MB pro 10.000 Einträge
- Index-Size: ca. 30-50% der Datengröße

### Horizontal Scaling

**Load Balancer Setup:**

```nginx
# nginx load balancer
upstream rsyslog_api_cluster {
    least_conn;  # Load Balancing Method
    
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
- Alle API-Instanzen → gleiche Datenbank
- Database wird Bottleneck
- Consider: Read Replicas

### Database Scaling

**Read Replicas:**

```bash
# Master-Slave Replication
# Master: Writes (rsyslog)
# Slaves: Reads (API)

# In .env verschiedener API-Instanzen:
# API-1: DB_HOST=slave1.db.local
# API-2: DB_HOST=slave2.db.local
```

**Partitioning:**

```sql
-- Partitionierung nach Datum
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

- [ ] Database Indexes erstellt
- [ ] Connection Pool konfiguriert
- [ ] Kleine Zeitfenster bei Queries
- [ ] Pagination statt große Limits
- [ ] Logging optimiert (kein Debug in Production)

### Database Level

- [ ] InnoDB Buffer Pool optimiert
- [ ] Slow Query Log aktiviert
- [ ] Indexes für häufige Queries
- [ ] Query Cache (MySQL 5.7)
- [ ] Alte Logs archiviert/gelöscht

### System Level

- [ ] SSD für Database
- [ ] Ausreichend RAM
- [ ] Firewall optimiert
- [ ] nginx/Apache konfiguriert
- [ ] Rate Limiting aktiviert

---

## 🔧 Troubleshooting Performance Issues

### Symptom: Langsame Responses

**Diagnose:**

```bash
# 1. Database prüfen
mysql Syslog -e "SHOW PROCESSLIST"

# 2. Slow Queries
tail -100 /var/log/mysql/mysql-slow.log

# 3. API Logs
sudo journalctl -u rsyslog-rest-api -n 100

# 4. System Resources
top
htop
iostat -x 1
```

**Lösungen:**

- Indexes fehlen? → Erstellen
- Große Zeitfenster? → Einschränken
- Viele Connections? → Pool vergrößern
- Disk I/O hoch? → SSD nutzen

### Symptom: Hohe Memory Usage

**Diagnose:**

```bash
# API Memory
ps aux | grep rsyslog-rest-api

# Database Memory
sudo mysqladmin -u root -p memory
```

**Lösungen:**

- Memory Leak? → Update auf neueste Version
- DB Buffer zu groß? → Reduzieren
- Zu viele Connections? → Limit setzen

### Symptom: High CPU Usage

**Diagnose:**

```bash
# CPU per Process
top -p $(pgrep rsyslog-rest-api)

# Database CPU
mysqladmin -u root -p processlist
```

**Lösungen:**

- Ineffiziente Queries? → EXPLAIN nutzen
- Viele Requests? → Rate Limiting
- Fehlende Indexes? → Erstellen

---

## 📊 Benchmarks

### Test Setup

- **Hardware:** 4 CPU, 8GB RAM, SSD
- **Database:** 1M Logs
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

## 💡 Best Practices

1. **Kleine Zeitfenster verwenden** - Max. 7-30 Tage
2. **Limit begrenzen** - Max. 100-500, nicht 1000
3. **Pagination nutzen** - Statt großer Limits
4. **Indexierte Felder filtern** - FromHost, Priority, ReceivedAt
5. **Alte Logs archivieren** - Regelmäßig aufräumen
6. **Monitoring einrichten** - Frühzeitig Probleme erkennen
7. **Load Testing** - Vor Production Deployment

---

[← Zurück zur Übersicht](index.md)

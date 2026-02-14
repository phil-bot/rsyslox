# Troubleshooting Guide

[← Back to overview](index.md)

Solutions for common problems and troubleshooting.

## 🔍 Quick diagnosis

### Check service status

```bash
# Status
sudo systemctl status rsyslog-rest-api

# Active logs
sudo journalctl -u rsyslog-rest-api -f

# Last 50 lines
sudo journalctl -u rsyslog-rest-api -n 50
```

### Health Check

```bash
# Simple health check
curl http://localhost:8000/health

# With details
curl -v http://localhost:8000/health
```

---

## ❌ Common problems

### 1. service does not start

#### Symptom
```bash
sudo systemctl start rsyslog-rest-api
# Job for rsyslog-rest-api.service failed
```

#### Causes & solutions

**A) Binary not found**

```bash
# Check if binary exists
ls -la /opt/rsyslog-rest-api/rsyslog-rest-api

# If not available: Reinstall
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-amd64
sudo cp rsyslog-rest-api-linux-amd64 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api
```

**B) .env file missing**

```bash
# Check
ls -la /opt/rsyslog-rest-api/.env

# If not available
sudo nano /opt/rsyslog-rest-api/.env
# Enter minimum config
```

**C) Permissions problem**

```bash
# Set correct permissions
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api
```

**D) Port already in use**

```bash
# Check which process is using port 8000
sudo lsof -i :8000
sudo netstat -tlnp | grep 8000

# Set another port in .env
sudo nano /opt/rsyslog-rest-api/.env
# SERVER_PORT=8080
```

**E) Check logs**

```bash
# Detailed errors
sudo journalctl -u rsyslog-rest-api -n 100 --no-pager
```

---

### 2. database connection failed

#### Symptom
```json
{
  "status": "unhealthy",
  "database": "disconnected"
}
```

#### Solutions

**A) Check credentials**

```bash
# Check .env file
sudo grep "^DB_" /opt/rsyslog-rest-api/.env

# Test manually
mysql -h DB_HOST -u DB_USER -pDB_PASS DB_NAME
```

**B) MySQL is not running**

```bash
# Check status
sudo systemctl status mysql
# or
sudo systemctl status mariadb

# Start if stopped
sudo systemctl start mysql
```

**C) User/rights missing**

```bash
# MySQL as root
sudo mysql

# Create user
CREATE USER 'rsyslog'@'localhost' IDENTIFIED BY 'password';
GRANT SELECT ON Syslog.* TO 'rsyslog'@'localhost';
FLUSH PRIVILEGES;
```

**D) Database does not exist**

```bash
sudo mysql

# Create database
CREATE DATABASE IF NOT EXISTS Syslog;

# Create table (if not existing)
USE Syslog;
CREATE TABLE IF NOT EXISTS SystemEvents (
    ID int unsigned not null auto_increment primary key,
    ReceivedAt datetime NULL,
    FromHost varchar(60) NULL,
    Priority smallint NULL,
    Facility smallint NULL,
    Message text,
    SysLogTag varchar(60)
);
```

**E) Network problem (Remote DB)**

```bash
# Test connection
ping DB_HOST
telnet DB_HOST 3306

# Check firewall
sudo ufw status
sudo iptables -L
```

---

### 3. API key does not work

#### Symptom
```json
{
  "error": "Invalid or missing API key"
}
```

#### Solutions

**A) Use key correctly**

```bash
# Read key from .env
API_KEY=$(sudo grep "^API_KEY=" /opt/rsyslog-rest-api/.env | cut -d'=' -f2)

# Test with exact key
curl -H "X-API-Key: $API_KEY" "http://localhost:8000/logs?limit=1"
```

**B) Whitespace/Formatting**

```bash
# Check .env (no spaces!)
sudo cat /opt/rsyslog-rest-api/.env | grep API_KEY

# Correct: API_KEY=abc123
# Incorrect: API_KEY = abc123
# Incorrect: API_KEY=abc123 (space at the end)
```

**C) Restart service after change**

```bash
sudo systemctl restart rsyslog-rest-api
```

---

### 4. no logs / empty result

#### Symptom
```json
{
  "total": 0,
  "offset": 0,
  "limit": 10,
  "rows": []
}
```

#### Solutions

**A) rsyslog does not write to DB**

```bash
# rsyslog status
sudo systemctl status rsyslog

# check rsyslog config
cat /etc/rsyslog.d/mysql.conf

# restart rsyslog
sudo systemctl restart rsyslog

# Send test log
logger -t test "Test message from logger"

# Check in DB
sudo mysql Syslog -e "SELECT COUNT(*) FROM SystemEvents"
```

**B) Filter to strict**

```bash
# Test without filter
curl -H "X-API-Key: YOUR_KEY" "http://localhost:8000/logs?limit=10"

# Extend time window
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=2025-01-01T00:00:00Z&limit=10"
```

**C) Incorrect table**

```bash
# Check correct table
sudo mysql Syslog -e "SHOW TABLES"
# Should: SystemEvents

# Count entries
sudo mysql Syslog -e "SELECT COUNT(*) FROM SystemEvents"
```

---

### 5. SSL/TLS problems

#### Symptom
```
curl: (60) SSL certificate problem
```

#### Solutions

**A) Self-Signed Certificate (Development)**

```bash
# Ignore warning (Development only!)
curl -k https://localhost:8000/health
```

**B) Certificate Path incorrect**

```bash
# Check paths in .env
sudo grep "^SSL_" /opt/rsyslog-rest-api/.env

# Files exist?
ls -la /opt/rsyslog-rest-api/certs/

# Permissions
sudo chmod 600 /opt/rsyslog-rest-api/certs/*.pem
```

**C) Let's Encrypt path**

```bash
# Correct path for Let's Encrypt
SSL_CERTFILE=/etc/letsencrypt/live/yourdomain.com/fullchain.pem
SSL_KEYFILE=/etc/letsencrypt/live/yourdomain.com/privkey.pem

# API must have access
sudo chmod +r /etc/letsencrypt/live/yourdomain.com/*.pem
# OR
sudo usermod -a -G ssl-cert root
```

---

### 6. CORS Errors (Browser)

#### Symptom
```
Access to fetch at 'http://api.example.com' from origin
'http://dashboard.example.com' has been blocked by CORS policy
```

#### Solutions

**A) Set ALLOWED_ORIGINS correctly**

```bash
sudo nano /opt/rsyslog-rest-api/.env

# Development (allow all)
ALLOWED_ORIGINS=*

# Production (specific)
ALLOWED_ORIGINS=https://dashboard.example.com,https://app.example.com

# IMPORTANT: https:// prefix!
```

**B) Restart service**

```bash
sudo systemctl restart rsyslog-rest-api
```

**C) Empty browser cache**

```bash
# Chrome: Ctrl+Shift+Delete
# Firefox: Ctrl+Shift+Delete
# Or: Test incognito mode
```

---

### 7. performance problems

#### Symptom
Slow response times, timeouts

#### Solutions

**A) Restrict time window**

```bash
# Instead of 90 days...
?start_date=2025-02-09T00:00:00Z&end_date=2025-02-09T23:59:59Z

# ...use only 1 day
```

**B) Use limit**

```bash
# Not: ?limit=10000
# Better: ?limit=100 with pagination
```

**C) Check indexes**

```bash
sudo mysql Syslog -e "SHOW INDEX FROM SystemEvents"

# If missing, create a new one
sudo mysql Syslog <<EOF
CREATE INDEX idx_receivedat ON SystemEvents (ReceivedAt);
CREATE INDEX idx_host_time ON SystemEvents (FromHost, ReceivedAt);
CREATE INDEX idx_priority ON SystemEvents (Priority);
EOF
```

**D) Database Stats**

```bash
# Table size
sudo mysql Syslog -e "
SELECT
  COUNT(*) as rows,
  ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
FROM information_schema.TABLES
WHERE table_schema = 'Syslog' AND table_name = 'SystemEvents'
"
```

---

## 🐞 Debug Mode

### Activate verbose logging

For detailed logs (development):

```bash
# Stop service
sudo systemctl stop rsyslog-rest-api

# Start manually in the foreground
cd /opt/rsyslog-rest-api
sudo -E ./rsyslog-rest-api

# Test in another terminal
curl http://localhost:8000/health

# View output
```

### Log SQL queries

In `main.go` (for developers):

```go
// Before db.Query()
log.Printf("SQL: %s", sqlQuery)
log.Printf("Args: %v", args)
```

---

## 🔬 Extended diagnosis

### Network debugging

```bash
# Check port lists
sudo ss -tlnp | grep 8000
sudo netstat -tlnp | grep 8000

# Firewall
sudo ufw status verbose
sudo iptables -L -n -v

# Test from other host
telnet API_HOST 8000
```

### MySQL Connection Debugging

```bash
# MySQL Connections
sudo mysql -e "SHOW PROCESSLIST"

# MySQL Errors
sudo tail -f /var/log/mysql/error.log

# MySQL Slow Queries
sudo tail -f /var/log/mysql/slow.log
```

### System Resources

```bash
# CPU/Memory
top
htop

# Disk Space
df -h

# Disk I/O
iostat -x 1
```

---

## 📋 FAQ

### Can I use the API without an API key?

**Yes**, set `API_KEY=` (empty) in `.env`. Only recommended for development!

### How big can limit be?

**Maximum 1000**. Use pagination for larger amounts of data.

### What time period is allowed?

**Maximum 90 days** between `start_date` and `end_date`.

### Does the API work with PostgreSQL?

**No**, only MySQL/MariaDB. PostgreSQL support is not planned.

### Can I use multiple API keys?

**No**, currently only one global API key. Use the same key for multiple clients (or create a feature request).

### How many requests/second are possible?

Currently **no** rate limiting. Performance depends on hardware and database size. For production: rate limiting via reverse proxy.

### Can I add my own fields?

**Yes**, all columns in `SystemEvents` are automatically available. Simply add column in DB, restart API.

### Is HTTP/2 supported?

**No**, currently only HTTP/1.1. HTTP/2 possible via reverse proxy (nginx).

### Can I run the API in Docker?

**Yes**, see [Docker Guide](docker.md). Recommended for testing, not for production.

---

## 🆘 Further help

### Collect logs

For bug reports:

```bash
# System Info
uname -a
cat /etc/os-release

# API version
/opt/rsyslog-rest-api/rsyslog-rest-api --version

# Service Status
sudo systemctl status rsyslog-rest-api

# Logs
sudo journalctl -u rsyslog-rest-api -n 100 --no-pager > api-logs.txt

# Config (blacken passwords!)
sudo cat /opt/rsyslog-rest-api/.env | sed 's/PASS=.*/PASS=REDACTED/' > config.txt
```

### Create GitHub issue

→ [GitHub Issues](https://github.com/phil-bot/rsyslog-rest-api/issues)

**Template:**
```markdown
**Environment:**
- OS: Ubuntu 22.04
- API Version: v0.2.2
- Installation: Binary/Source/Docker

**Problem:**
[Description]

**Steps to Reproduce:**
1. ...
2. ...

**Expected:**
[What should happen]

**Actual:**
[What actually happens]

**Logs:**
```
[Insert logs here]
```

**Config:**
[Config (without passwords!)]
```

---

[← Back to overview](index.md) | [Forward to Docker →](docker.md)

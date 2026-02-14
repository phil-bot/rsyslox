# Docker Test Environment

[← Back to overview](index.md)

Complete Docker test environment with live generated logs.

## 🐳 About the test environment

The Docker environment offers a fully functional setup for testing and development:

✅ **Ubuntu 24.04** container
✅ **rsyslog + MariaDB** pre-installed
✅ **10 initial test logs** at startup
✅ **Live Log Generator** - 3 new logs every 10 seconds
✅ **All 25+ fields** filled with realistic data
✅ **No API key authentication** (can be activated optionally)

**Perfect for:**
- API testing
- Feature development
- Demos
- Testing multi-value filters
- Testing Extended Columns

---

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose installed
- Go 1.21+ (for build)
- Make

### Setup

```bash
# Build 1st binary (on host system)
make build-static

# 2. start docker container
cd docker
docker-compose up -d

# 3. wait until ready (approx. 30 seconds)
docker-compose logs -f
# Wait for: "✓ Environment Ready!"

# 4. testing
curl http://localhost:8000/health
curl "http://localhost:8000/logs?limit=5"
```

**Done!** The API is now running on `http://localhost:8000`.

---

## 📦 What is included?

### Container setup

```
rsyslog-rest-api-test/
├── Ubuntu 24.04
├── MariaDB Server
├── rsyslog + rsyslog-mysql
├── rsyslog-rest-api (your binary)
├── log-generator.sh (live data)
└── test-v0.2.2.sh (test suite)
```

### Initial data

At container start **10 test logs** are created:

```sql
INSERT INTO SystemEvents VALUES
  ('webserver01', Priority=6, 'User login successful'),
  ('webserver01', Priority=3, 'Failed login attempt'),
  ('dbserver01', Priority=4, 'Database connection timeout'),
  ('dbserver01', Priority=6, 'Query executed successfully'),
  ('appserver01', Priority=5, 'Application started'),
  ('appserver01', Priority=3, 'Critical error in module'),
  ('webserver02', Priority=6, 'HTTP request: GET /api/users'),
  ('webserver02', Priority=4, 'Slow response time detected'),
  ('mailserver01', Priority=2, 'Mail queue growing rapidly'),
  ('mailserver01', Priority=6, 'Email sent successfully')
```

### Live Log Generator 🔥

**Every 10 seconds** **3 new realistic logs** are generated:

**hosts:**
- webserver01, webserver02
- dbserver01
- appserver01
- mailserver01
- firewall01

**SysLogTags:**
- sshd, nginx, mysqld, node, postfix, iptables, systemd, docker

**Priorities (weighted):**
- 6 (Info) - 50%
- 5 (Notice) - 25%
- 4 (Warning) - 12.5%
- 3 (Error) - 10%
- 2 (Critical) - 2.5%

**Extended Fields:**
- Realistic event IDs (e.g. 4624 for SSH login, 200-500 for HTTP)
- EventSource, EventUser, NTSeverity
- Importance, EventCategory, SystemID
- All 25+ columns filled!

**Example logs:**
```
[2025-02-09 10:30:15] [webserver01] [INFO] [sshd] User login successful
[2025-02-09 10:30:16] [dbserver01] [WARNING] [mysqld] Slow query detected: 2500ms
[2025-02-09 10:30:17] [appserver01] [ERROR] [node] Connection refused
```

---

## 🧪 Testing

### Manual tests

```bash
# Health Check
curl http://localhost:8000/health

# Latest 5 logs
curl "http://localhost:8000/logs?limit=5"

# Errors only
curl "http://localhost:8000/logs?Priority=3&limit=10"

# Multi-Value: Multiple hosts
curl "http://localhost:8000/logs?FromHost=webserver01&FromHost=webserver02"

# See extended columns
curl "http://localhost:8000/logs?limit=1" | jq .rows[0]

# Meta: All hosts
curl "http://localhost:8000/meta/FromHost"

# Meta: Hosts with errors
curl "http://localhost:8000/meta/FromHost?Priority=3&Priority=4"
```

### Execute test suite

The Docker environment contains a comprehensive test suite:

```bash
# In docker/ directory
./test-v0.2.2.sh
```

**Tests:**
- ✅ Health Check
- ✅ Basic Log Retrieval
- ✅ Multi-Value Filters
- ✅ Extended Columns
- ✅ Meta endpoints
- ✅ Backward Compatibility

**Expected output:**
```
==========================================
rsyslog REST API v0.2.2 - Test Suite
==========================================
...
Waiting for API... ready

Running tests...

[1] Health Check... ✓ OK (HTTP 200)
[2] Get Logs (default)... ✓ OK (HTTP 200)
[3] Multi-value FromHost... ✓ OK (HTTP 200)
...
==========================================
Test Summary
==========================================
Passed: 22
Failed: 0
Total: 22

✓ All tests passed!
```

### Watch live data

```bash
# View log generator output
docker exec rsyslog-rest-api-test tail -f /var/log/log-generator.log

# Watch database count live
watch -n 5 'docker exec rsyslog-rest-api-test mysql -N Syslog -e "SELECT COUNT(*) FROM SystemEvents"'

# Container logs
docker-compose logs -f
```

---

## ⚙️ configuration

### Activate API key (optional)

By default, the API runs **without** authentication. To activate it:

**1. edit docker-compose.yml:**

```yaml
# docker/docker-compose.yml
environment:
  - SERVER_PORT=8000
  - ALLOWED_ORIGINS=*
  - API_KEY=test123456789 # <-- Activate
```

**2. restart container:**

```bash
docker-compose down
docker-compose up -d
```

**3. test with API key:**

```bash
curl -H "X-API-Key: test123456789" "http://localhost:8000/logs?limit=5"
```

### Change port

**docker-compose.yml:**

```yaml
ports:
  - "8080:8000" # <-- Host:Container
```

```bash
docker-compose down
docker-compose up -d
```

API then to: `http://localhost:8080`

### Customize log generator

*edit *docker/log-generator.sh:**

```bash
# Configuration
INTERVAL=10 # Seconds between bursts (default: 10)
LOGS_PER_BURST=3 # Logs per burst (default: 3)
```

**Rebuild container:**

```bash
docker-compose down
docker-compose up -d --build
```

---

## 📊 Monitoring

### Container status

```bash
# Status
docker-compose ps

# Logs live
docker-compose logs -f

# API logs
docker exec rsyslog-rest-api-test tail -f /var/log/rsyslog-rest-api.log

# Generator logs
docker exec rsyslog-rest-api-test tail -f /var/log/log-generator.log
```

### Database

```bash
# MySQL Shell
docker exec -it rsyslog-rest-api-test mysql -u rsyslog -ppassword Syslog

# Count logs
docker exec rsyslog-rest-api-test mysql -N Syslog -e "SELECT COUNT(*) FROM SystemEvents"

# Latest logs
docker exec rsyslog-rest-api-test mysql Syslog -e "
  SELECT ReceivedAt, FromHost, Priority, Message
  FROM SystemEvents
  ORDER BY ReceivedAt DESC
  LIMIT 5
"

# Statistics
docker exec rsyslog-rest-api-test mysql Syslog -e "
  SELECT
    FromHost,
    COUNT(*) as count,
    AVG(Priority) as avg_priority
  FROM SystemEvents
  GROUP BY FromHost
"
```

### Processes in the container

```bash
# Running processes
docker exec rsyslog-rest-api-test ps aux

# API running?
docker exec rsyslog-rest-api-test ps aux | grep rsyslog-rest-api

# Generator running?
docker exec rsyslog-rest-api-test ps aux | grep log-generator
```

---

## 🔧 Advanced usage

### Container Shell

```bash
# Bash shell in the container
docker exec -it rsyslog-rest-api-test bash

# In the container:
cd /opt/rsyslog-rest-api
cat .env
./rsyslog-rest-api --help
```

### redeploy binary

After code changes:

```bash
#1 Rebuild
cd .. && make build-static

# 2. restart container (copies automatically)
cd docker
docker-compose restart
```

### Reset database

```bash
# In the container
docker exec -it rsyslog-rest-api-test mysql Syslog -e "TRUNCATE TABLE SystemEvents"

# Reload initial data
docker-compose restart
```

---

## 🛠️ Troubleshooting

### Container does not start

```bash
# View logs
docker-compose logs

# Binary missing?
ls -la ../build/rsyslog-rest-api
# If not available: make build-static

# Rebuild
docker-compose down
docker-compose up -d --build
```

### API does not respond

```bash
# API is running?
docker exec rsyslog-rest-api-test ps aux | grep rsyslog-rest-api

# Logs
docker exec rsyslog-rest-api-test cat /var/log/rsyslog-rest-api.log

# Port correct?
docker-compose ps
# Should show: 0.0.0.0:8000->8000/tcp

# Test from host
curl -v http://localhost:8000/health
```

### Log generator is not running

```bash
# Check process
docker exec rsyslog-rest-api-test ps aux | grep log-generator

# View logs
docker exec rsyslog-rest-api-test cat /var/log/log-generator.log

# Start manually
docker exec -it rsyslog-rest-api-test bash
/opt/rsyslog-rest-api/log-generator.sh
```

### Database problems

```bash
# MySQL is running?
docker exec rsyslog-rest-api-test systemctl status mysql

# Test connection
docker exec rsyslog-rest-api-test mysql -u rsyslog -ppassword Syslog -e "SELECT 1"

# Table exists?
docker exec rsyslog-rest-api-test mysql Syslog -e "SHOW TABLES"
```

---

## 🧹 Cleanup

### Stop container

```bash
# Stop (data remains)
docker-compose stop

# Start
docker-compose start
```

### Remove completely

```bash
# Remove container and network
docker-compose down

# With volumes (DELETE DATA!)
docker-compose down -v
```

### Release disk space

```bash
# Unused containers/images/volumes
docker system prune

# All dangling images
docker image prune
```

---

## 📚 Files

### Structure

```
docker/
├── Dockerfile # Container image
├── docker-compose.yml # Compose configuration
├── entrypoint.sh # Startup script
├── log-generator.sh # Live log generator
├── test.sh # Basic test suite
└── test-v0.2.2.sh # Extended test suite
```

### Important paths in the container

```
/opt/rsyslog-rest-api/ # API installation
  ├── rsyslog-rest-api # Binary (from host)
  ├── .env # Config
  └── log-generator.sh # Generator

/etc/rsyslog.d/mysql.conf # rsyslog Config

/var/log/
  ├── rsyslog-rest-api.log # API logs
  └── log-generator.log # Generator logs
```

---

## 🎯 Use Cases

### Feature development

```bash
#1 Change code
vim main.go

# 2. rebuild
make build-static

# 3. restart container
cd docker && docker-compose restart

# 4. test
curl "http://localhost:8000/new-endpoint"
```

### Prepare demo

```bash
# 1. start container
docker-compose up -d

# 2. wait until data has been generated (2-3 minutes)
sleep 180

# 3. demo
curl "http://localhost:8000/logs?Priority=3&Priority=4&limit=10" | jq
```

### Performance Testing

```bash
# Load Testing (with ApacheBench)
ab -n 1000 -c 10 "http://localhost:8000/logs?limit=100"

# With filters
from -n 500 -c 5 "http://localhost:8000/logs?FromHost=webserver01&Priority=3&limit=50"
```

---

## 🆚 Docker vs. production

### Differences

| Feature | Docker | Production |
|---------|--------|------------|
| API Key | Optional | **Required** |
| SSL/TLS | No | **Yes** |
| Database | In Container | External |
| Persistence | Volatile | Persistent |
| Performance | Limited | Optimized |

### Not for Production!

Docker setup is **only for testing/development**!

For Production see:
- [Deployment Guide](deployment.md)
- [Security Guide](security.md)

---

## 💡 Tips

1. **Long run** - After 1 hour you have ~1000 logs to test
2. **Test multi-value** - Perfect for new filter features
3. **Extended Columns** - All fields have realistic values
4. **Performance** - Test with different `limit` values
5. **Pagination** - Test with `offset` and `limit`

---

[← Back to overview](index.md) | [Next to Development →](development.md)

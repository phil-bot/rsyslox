# Docker Test Environment

Simple test environment: **Build on host, runtime in container**.

**NEW in v0.2.3:** 🎯 **Unified test script** - use `/build-and-test.sh` from project root!

## Quick Start

```bash
# Build binary
make build-static

# Run complete test suite (from project root)
./build-and-test.sh
```

**That's it!** The script handles everything: Docker, tests, cleanup.

---

## Manual Docker Usage

If you want to work with Docker directly:

### Start Environment

```bash
# Start Docker
cd docker
docker-compose up -d

# Check logs
docker-compose logs -f
```

### Test API

```bash
# Health check
curl http://localhost:8000/health

# Get logs
curl "http://localhost:8000/logs?limit=5"

# With API key (if configured)
curl -H "X-API-Key: test123456789" "http://localhost:8000/logs?limit=5"
```

### Stop Environment

```bash
cd docker
docker-compose down

# With cleanup
docker-compose down -v
```

---

## What's Included

- ✅ Ubuntu 24.04
- ✅ rsyslog + MariaDB
- ✅ 10 initial test entries
- ✅ **Live log generator** (3 new entries every 10 seconds)
- ✅ No API key required (optional)

---

## Live Logs Feature 🔥

The container generates realistic syslog messages continuously!

### What it does:
- **Every 10 seconds:** 3 new log entries
- **6 Hosts:** webserver01, webserver02, dbserver01, appserver01, mailserver01, firewall01
- **Realistic messages:** Login attempts, database queries, errors, warnings
- **Various priorities:** Info (most), Notice, Warning, Error, Critical (rare)

### Watch it in action:
```bash
# Initial count
curl -s "http://localhost:8000/logs" | jq .total
# → 10 entries

# Wait 30 seconds (3 bursts × 3 logs = 9 new logs)
sleep 30

# New count
curl -s "http://localhost:8000/logs" | jq .total
# → 19 entries (growing live!)
```

---

## Testing

### Automated Testing (Recommended)

```bash
# From project root - runs everything!
./build-and-test.sh

# Options
./build-and-test.sh --skip-build  # Skip build, just test
./build-and-test.sh --cleanup     # Stop Docker and cleanup
./build-and-test.sh --help        # Show help
```

### Manual Testing

```bash
# Start environment
cd docker
docker-compose up -d

# Run individual tests
curl http://localhost:8000/health
curl "http://localhost:8000/logs?limit=5"
curl "http://localhost:8000/meta/FromHost"
```

---

## Testing v0.2.3 Features

### Multi-Value Filters

```bash
# Multiple hosts
curl "http://localhost:8000/logs?FromHost=webserver01&FromHost=webserver02&limit=5"

# Multiple priorities
curl "http://localhost:8000/logs?Priority=3&Priority=4&limit=5"

# Combined
curl "http://localhost:8000/logs?FromHost=webserver01&FromHost=dbserver01&Priority=3&Priority=6"
```

### Structured Errors

```bash
# Test invalid priority
curl "http://localhost:8000/logs?Priority=99"

# Expected response (v0.2.3+):
{
  "code": "INVALID_PRIORITY",
  "message": "value 99 is out of range (must be 0-7)",
  "details": "See RFC-5424 for valid priority levels",
  "field": "Priority"
}
```

### Extended Columns

```bash
# See all 25+ columns
curl "http://localhost:8000/logs?limit=1" | jq .rows[0]

# Filter by SysLogTag
curl "http://localhost:8000/logs?SysLogTag=nginx&limit=5"
```

---

## Enable API Key (Optional)

Edit `docker-compose.yml`:
```yaml
environment:
  - API_KEY=test123456789
```

Restart:
```bash
docker-compose down && docker-compose up -d
```

Then use:
```bash
curl -H "X-API-Key: test123456789" "http://localhost:8000/logs?limit=5"
```

---

## Monitor Live Logs

```bash
# Watch log generator
docker exec rsyslog-rest-api-test tail -f /var/log/log-generator.log

# Watch database grow
watch -n 5 'docker exec rsyslog-rest-api-test mysql -N Syslog -e "SELECT COUNT(*) FROM SystemEvents"'

# Container logs
docker-compose logs -f
```

---

## Configuration

### Adjust Log Generation Rate

Edit `docker/log-generator.sh`:
```bash
INTERVAL=10  # Seconds between bursts (default: 10)
LOGS_PER_BURST=3  # Logs per burst (default: 3)
```

Rebuild:
```bash
docker-compose down
docker-compose up -d --build
```

### Change Port

Edit `docker-compose.yml`:
```yaml
ports:
  - "8080:8000"
```

---

## Database Access

```bash
# Connect
docker exec -it rsyslog-rest-api-test mysql -u rsyslog -ppassword Syslog

# Check latest logs
SELECT ReceivedAt, FromHost, Priority, Message 
FROM SystemEvents 
ORDER BY ReceivedAt DESC 
LIMIT 5;

# Watch count grow
SELECT COUNT(*) FROM SystemEvents;
```

---

## Troubleshooting

### Binary not found
```bash
cd .. && make build-static
cd docker && docker-compose up -d
```

### Live logs not working
```bash
# Check generator
docker exec rsyslog-rest-api-test ps aux | grep log-generator

# View generator logs
docker exec rsyslog-rest-api-test cat /var/log/log-generator.log

# Restart
docker-compose restart
```

### API not responding
```bash
# Check API logs
docker logs rsyslog-rest-api-test

# Check API process
docker exec rsyslog-rest-api-test ps aux | grep rsyslog-rest-api

# Restart
docker-compose restart
```

---

## Deprecated Test Scripts

⚠️ **As of v0.2.3, the following scripts are deprecated:**

- ❌ `/docker/test.sh`
- ❌ `/docker/test-v0.2.2.sh`

**Use instead:** `/build-and-test.sh` from project root

These old scripts will be removed in v0.3.0.

---

## Advantages

✅ **Simple** - One command to test everything  
✅ **Fast** - Binary pre-built  
✅ **Realistic** - Live log generation  
✅ **Flexible** - Easy to restart  
✅ **No buildx** - Standard Docker  
✅ **Unified** - Single test script for all versions

---

## What's Next?

**v0.3.0 Preview:**
- Negation filters (`exclude`, `not`)
- More complex filtering
- Statistics endpoint
- Ready to test with live data!

Full documentation: https://github.com/phil-bot/rsyslog-rest-api

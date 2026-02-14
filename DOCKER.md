# Docker Test Environment

Simple test environment: **Build on host, runtime in container**.

**NEW in v0.2.2:** 🔥 **Live log generation** - realistic syslog messages appear every 10 seconds!

## Quick Start

```bash
# 1. Build binary
make build-static

# 2. Start Docker
cd docker
docker-compose up -d

# 3. Test (no API key needed by default!)
curl http://localhost:8000/health
curl "http://localhost:8000/logs?limit=5"
```

## Live Logs Feature 🔥

The container now generates realistic syslog messages continuously!

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

## What's Included

- ✅ Ubuntu 24.04
- ✅ rsyslog + MariaDB
- ✅ 10 initial test entries
- ✅ **Live log generator** (NEW!)
- ✅ No API key required (optional)

## Testing v0.2.2 Features

### Multi-Value Filters (NEW!)

```bash
# Multiple hosts
curl "http://localhost:8000/logs?FromHost=webserver01&FromHost=webserver02&limit=5"

# Multiple priorities
curl "http://localhost:8000/logs?Priority=3&Priority=4&limit=5"

# Combined
curl "http://localhost:8000/logs?FromHost=webserver01&FromHost=dbserver01&Priority=3&Priority=6"
```

### Extended Columns (NEW!)

```bash
# See all 25 columns
curl "http://localhost:8000/logs?limit=1" | jq .rows[0]

# Filter by SysLogTag
curl "http://localhost:8000/logs?SysLogTag=nginx&limit=5"
```

### Run Test Suite

```bash
# Extended test suite for v0.2.2
./test-v0.2.2.sh
# Tests multi-value filters, extended columns, live data
```

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

## Monitor Live Logs

```bash
# Watch log generator
docker exec rsyslog-rest-api-test tail -f /var/log/log-generator.log

# Watch database grow
watch -n 5 'docker exec rsyslog-rest-api-test mysql -N Syslog -e "SELECT COUNT(*) FROM SystemEvents"'

# Container logs
docker-compose logs -f
```

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

## Stop & Clean

```bash
# Stop (data persists)
docker-compose stop

# Remove (data persists)
docker-compose down

# Delete everything
docker-compose down -v
```

## Advantages

✅ **Simple** - No complex build  
✅ **Fast** - Binary pre-built  
✅ **Realistic** - Live log generation  
✅ **Flexible** - Easy to restart  
✅ **No buildx** - Standard Docker

## What's Next?

**v0.3.0 Preview:**
- Negation filters (`exclude`, `not`)
- More complex filtering
- Ready to test with live data!

Full documentation: https://github.com/phil-bot/rsyslog-rest-api

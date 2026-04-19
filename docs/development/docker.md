# Docker Testing Environment

A self-contained test environment with live log generation for development and testing.

## What's Included

- Ubuntu 24.04 + rsyslog + MariaDB
- 10 initial test log entries
- Live log generator — 3 logs every 10 seconds
- All 25+ SystemEvents fields populated
- Authentication optional (disabled by default for easy testing)

## Quick Start

```bash
# 1. Build static binary
make build-static

# 2. Start the test container
cd docker
docker-compose up -d

# 3. Verify
curl "http://localhost:8000/health"
curl "http://localhost:8000/api/logs?limit=5"
```

Open `http://localhost:8000` in your browser — the full web UI is available.

## Setup Wizard in Docker

On first start, the container runs the setup wizard. It is pre-configured with test credentials via environment variables — the wizard completes automatically. If you want to redo it:

```bash
docker-compose down -v
docker-compose up -d
```

## Live Log Generator

The container generates realistic log entries every 10 seconds:

- **Hosts:** webserver01, webserver02, dbserver01, appserver01, mailserver01, firewall01
- **Tags:** sshd, nginx, mysqld, node, postfix, iptables
- **Severities:** Weighted distribution — more Informational, fewer Critical
- **All fields populated:** EventID, EventSource, EventUser, etc.

```bash
# Watch logs being generated
docker exec rsyslox-test tail -f /var/log/log-generator.log

# Watch row count grow
watch -n 5 'docker exec rsyslox-test mysql -N Syslog -e "SELECT COUNT(*) FROM SystemEvents"'
```

## Testing the API

```bash
# Multi-value filters
curl "http://localhost:8000/api/logs?FromHost=webserver01&FromHost=webserver02&limit=5"
curl "http://localhost:8000/api/logs?Severity=3&Severity=4&limit=5"

# Metadata queries
curl "http://localhost:8000/api/meta/FromHost"
curl "http://localhost:8000/api/meta/Severity"

# Error validation
curl "http://localhost:8000/api/logs?Severity=99"
```

## Enable API Key Authentication

The `API_KEY` environment variable was removed in v0.4.0. API keys are now
created through the Admin panel:

1. Open `http://localhost:8000/admin` and log in with the admin password
   (default: set during Docker container startup via `ADMIN_PASSWORD` env var)
2. Navigate to **API Keys → Create**
3. Enter a name (e.g. `testing`) and copy the key shown once after creation

```bash
API_KEY="your-key-here"
curl -H "X-API-Key: $API_KEY" "http://localhost:8000/api/logs"
```

## Configuration

### Change Port

```yaml
ports:
  - "9000:8000"  # Host:Container
```

### Adjust Log Generation

Edit `docker/log-generator.sh`:

```bash
INTERVAL=10        # seconds between bursts
LOGS_PER_BURST=3   # logs per burst
```

Rebuild:
```bash
docker-compose down
docker-compose up -d --build
```

## Monitoring the Container

```bash
# All container logs
docker-compose logs -f

# rsyslox application logs
docker exec rsyslox-test journalctl -u rsyslox -n 50

# Generator logs
docker exec rsyslox-test cat /var/log/log-generator.log

# Database row count
docker exec rsyslox-test mysql Syslog -e "SELECT COUNT(*) FROM SystemEvents"
```

## Database Access

```bash
# Connect directly to MariaDB
docker exec -it rsyslox-test mysql -u rsyslog -ppassword Syslog

# Recent entries
SELECT ReceivedAt, FromHost, Priority, Message
FROM SystemEvents
ORDER BY ReceivedAt DESC
LIMIT 10;
```

## Cleanup

```bash
# Stop (data persists in volume)
docker-compose stop

# Remove containers (data persists)
docker-compose down

# Remove everything including data
docker-compose down -v
```

## Troubleshooting

**Binary not found at container start:**
```bash
cd ..
make build-static
cd docker && docker-compose up -d
```

**Live logs not generating:**
```bash
docker exec rsyslox-test ps aux | grep log-generator
# If not running:
docker-compose restart
```

**Port 8000 already in use:**
```bash
sudo lsof -i :8000
# Change port in docker-compose.yml ports section
```

## Performance Testing

```bash
# Install wrk
sudo apt-get install wrk

# Benchmark
wrk -t4 -c50 -d30s "http://localhost:8000/api/logs?limit=10"
```

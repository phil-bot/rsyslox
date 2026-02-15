# Docker Testing Environment

Complete test environment with live log generation.

## Features

- ✅ Ubuntu 24.04 + rsyslog + MariaDB
- ✅ 10 initial test entries
- ✅ **Live log generator** - 3 logs every 10 seconds
- ✅ All 25+ fields populated
- ✅ No API key required (optional)

## Quick Start

```bash
# 1. Build binary
make build-static

# 2. Start Docker
cd docker
docker-compose up -d

# 3. Test
curl "http://localhost:8000/health"
curl "http://localhost:8000/logs?limit=5"
```

## What's Included

```
docker/
├── Dockerfile
├── docker-compose.yml
├── entrypoint.sh
├── log-generator.sh    # Live log generation
└── test.sh            # Test suite
```

## Live Log Generator

Generates realistic logs every 10 seconds:

- **Hosts:** web01, web02, db01, app01, mail01, firewall01
- **Tags:** sshd, nginx, mysqld, node, postfix, iptables
- **Priorities:** Weighted (more INFO, less CRITICAL)
- **All fields:** EventID, EventSource, EventUser, etc.

```bash
# Watch logs being generated
docker exec rsyslog-rest-api-test tail -f /var/log/log-generator.log

# Check database growth
watch -n 5 'docker exec rsyslog-rest-api-test mysql -N Syslog -e "SELECT COUNT(*) FROM SystemEvents"'
```

## Testing v0.2.3 Features

### Multi-Value Filters

```bash
# Multiple hosts
curl "http://localhost:8000/logs?FromHost=web01&FromHost=web02&limit=5"

# Multiple priorities
curl "http://localhost:8000/logs?Priority=3&Priority=4&limit=5"

# Combined
curl "http://localhost:8000/logs?FromHost=web01&FromHost=db01&Priority=3&Priority=6"
```

### Run Test Suite

```bash
cd docker
./test.sh
```

## Configuration

### Enable API Key

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
curl -H "X-API-Key: test123456789" "http://localhost:8000/logs"
```

### Change Port

Edit `docker-compose.yml`:

```yaml
ports:
  - "9000:8000"  # Host:Container
```

### Adjust Log Generation

Edit `docker/log-generator.sh`:

```bash
INTERVAL=10        # Seconds between bursts
LOGS_PER_BURST=3   # Logs per burst
```

Rebuild:
```bash
docker-compose down
docker-compose up -d --build
```

## Monitoring

```bash
# Container logs
docker-compose logs -f

# API logs
docker exec rsyslog-rest-api-test cat /var/log/rsyslog-rest-api.log

# Generator logs
docker exec rsyslog-rest-api-test cat /var/log/log-generator.log

# Database count
docker exec rsyslog-rest-api-test mysql Syslog -e "SELECT COUNT(*) FROM SystemEvents"
```

## Database Access

```bash
# Connect to MySQL
docker exec -it rsyslog-rest-api-test mysql -u rsyslog -ppassword Syslog

# Example queries
SELECT ReceivedAt, FromHost, Priority, Message 
FROM SystemEvents 
ORDER BY ReceivedAt DESC 
LIMIT 10;
```

## Cleanup

```bash
# Stop (data persists)
docker-compose stop

# Remove (data persists)
docker-compose down

# Delete everything including data
docker-compose down -v
```

## Troubleshooting

**Binary not found:**
```bash
cd .. && make build-static
cd docker && docker-compose up -d
```

**Live logs not generating:**
```bash
# Check generator status
docker exec rsyslog-rest-api-test ps aux | grep log-generator

# Restart container
docker-compose restart
```

**Container won't start:**
```bash
# Check logs
docker-compose logs

# Rebuild
docker-compose down
docker-compose up -d --build
```

## Advanced Usage

### Custom Test Data

```bash
# Connect to container
docker exec -it rsyslog-rest-api-test bash

# Insert custom log
mysql Syslog <<EOF
INSERT INTO SystemEvents (ReceivedAt, FromHost, Priority, Message, SysLogTag)
VALUES (NOW(), 'testhost', 3, 'Custom test message', 'test');

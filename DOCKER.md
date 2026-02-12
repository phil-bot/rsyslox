# Docker Test Environment

Simple test environment: **Build on host, runtime in container**.

## Workflow

```bash
# 1. Build on host
make build-static
# Creates: build/rsyslog-rest-api

# 2. Start Docker
cd docker
docker-compose up -d

# 3. Wait (about 20 seconds)
docker-compose logs -f
# When you see "Environment Ready!" → Ctrl+C

# 4. Test
curl http://localhost:8000/health
curl -H "X-API-Key: test123456789" "http://localhost:8000/logs?limit=5"
```

That's it! No buildx, no complicated build process.

## What happens?

### On the host (your machine):
```bash
make build-static
```
- Builds the binary with Go
- Result: `build/rsyslog-rest-api` (8 MB)

### In the container:
```bash
docker-compose up -d
```
1. ✅ Mounts `../build` to `/host-build` (read-only)
2. ✅ Copies binary to `/opt/rsyslog-rest-api/`
3. ✅ Starts MariaDB
4. ✅ Creates database with 10 test entries
5. ✅ Configures rsyslog
6. ✅ Creates `.env` with API_KEY
7. ✅ Starts the API

## Prerequisites

- Docker (no buildx!)
- Docker Compose
- Make
- Go 1.21+ (for `make build-static`)

## Testing the API

```bash
# Health Check
curl http://localhost:8000/health

# Get logs
curl -H "X-API-Key: test123456789" "http://localhost:8000/logs?limit=5"

# Only errors
curl -H "X-API-Key: test123456789" "http://localhost:8000/logs?Priority=3"

# All hosts
curl -H "X-API-Key: test123456789" "http://localhost:8000/meta/FromHost"

# Test suite
./test.sh
```

## Troubleshooting

### Binary not found
```
✗ ERROR: Binary not found!
```

**Solution:**
```bash
cd .. && make build-static
cd docker && docker-compose up -d
```

### API won't start
```bash
# View logs
docker-compose logs

# Inside container
docker exec -it rsyslog-rest-api-test cat /var/log/rsyslog-rest-api.log
```

## Restart

```bash
# After new build
cd .. && make build-static
cd docker && docker-compose restart

# Complete rebuild
cd docker
docker-compose down
docker-compose up -d --build
```

## What's included?

- Ubuntu 24.04
- rsyslog with MySQL support
- MariaDB with test database "Syslog"
- 10 sample log entries
- Automatic setup on start

**Container layout:**
```
/opt/rsyslog-rest-api/
├── rsyslog-rest-api          # Binary (copied from host)
└── .env                      # Auto-generated config

/etc/rsyslog.d/mysql.conf     # rsyslog MySQL output
/var/log/rsyslog-rest-api.log # API logs
```

## Database Access

```bash
# Inside container
docker exec -it rsyslog-rest-api-test bash
mysql -u rsyslog -ppassword Syslog

# Check data
SELECT ReceivedAt, FromHost, Priority, Message 
FROM SystemEvents 
ORDER BY ReceivedAt DESC 
LIMIT 5;
```

### Add custom test data
```bash
docker exec -it rsyslog-rest-api-test mysql -u rsyslog -ppassword Syslog <<'EOF'
INSERT INTO SystemEvents (ReceivedAt, FromHost, Priority, Facility, Message, SysLogTag)
VALUES (NOW(), 'testhost', 6, 1, 'Custom test message', 'test');
EOF

# Verify
curl -H "X-API-Key: test123456789" \
  "http://localhost:8000/logs?FromHost=testhost"
```

## Configuration

### Change API Key

Edit `docker-compose.yml`:
```yaml
environment:
  - API_KEY=your-new-key
```

Then restart:
```bash
docker-compose down
docker-compose up -d
```

### Change Port

Edit `docker-compose.yml`:
```yaml
ports:
  - "8080:8000"  # Host:Container
```

Test with:
```bash
curl http://localhost:8080/health
```

## Stop & Clean

```bash
# Stop (data persists)
docker-compose stop

# Stop and remove (data persists)
docker-compose down

# Delete EVERYTHING (including data!)
docker-compose down -v
```

## Performance Test

```bash
# Apache Bench (if installed)
ab -n 1000 -c 10 -H "X-API-Key: test123456789" \
  http://localhost:8000/logs?limit=10

# Or simple timing
time curl -H "X-API-Key: test123456789" \
  "http://localhost:8000/logs?limit=100"
```

## Logs

```bash
# Container logs
docker-compose logs -f

# API logs
docker exec -it rsyslog-rest-api-test \
  tail -f /var/log/rsyslog-rest-api.log

# MariaDB logs
docker exec -it rsyslog-rest-api-test \
  tail -f /var/log/mysql/error.log
```

## Advantages

### 1. Simple
- No complex build steps in container
- Clear separation: build vs runtime
- Easy to understand

### 2. Fast
- Binary already built when container starts
- No waiting for Go build
- Faster restart

### 3. Flexible
```bash
# Change code
nano main.go

# Rebuild (on host)
make build-static

# Restart container (uses new binary)
cd docker
docker-compose restart
```

### 4. No buildx issues
- Works with standard Docker
- No plugin required
- Fewer dependencies

## vs. Real Server

| Feature | Docker | Real Server |
|---------|--------|-------------|
| Installation | Automatic | `make install` |
| Binary | Auto-built on host | Copied from host |
| API Start | Automatic | systemd |
| Database | Auto-setup | Manual |
| Test Data | Included | None |

The container simulates the **complete installation** - perfect for testing!

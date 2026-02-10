[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg)](https://golang.org)

# rsyslog REST API

High-performance REST API for querying rsyslog data stored in MariaDB/MySQL. Written in Go for maximum performance and minimal resource usage.

## Features

- 🚀 Fast - ~15k requests/second
- 💾 Low memory - ~20 MB RAM
- 📦 Single binary - no dependencies
- 🔒 Secure - API key auth, input validation
- 🎯 Complete - all filters, RFC-5424 labels
- 🔄 Dynamic - auto-discovers database columns

## Quick Start

### Prerequisites

- rsyslog with ommysql plugin
- MariaDB/MySQL with rsyslog data
- Go 1.21+ (for building)

### Build

```bash
git clone https://github.com/phil-bot/rsyslog-rest-api.git
cd rsyslog-rest-api
make build
```

### Install

```bash
sudo make install
sudo nano /opt/rsyslog-rest-api/.env  # Set API_KEY
sudo systemctl enable --now rsyslog-rest-api
```

### Test

```bash
curl http://localhost:8000/health
curl -H "X-API-Key: YOUR_KEY" "http://localhost:8000/logs?limit=10"
```

## Configuration

Edit `/opt/rsyslog-rest-api/.env`:

```bash
API_KEY=                  # Generate: openssl rand -hex 32
SERVER_PORT=8000
ALLOWED_ORIGINS=*         # Restrict in production!
```

## API Endpoints

### GET /health
Health check (no authentication required)

```bash
curl http://localhost:8000/health
```

### GET /logs
Query log entries

**Parameters:**
- `start_date` - ISO 8601 datetime
- `end_date` - ISO 8601 datetime  
- `FromHost` - Hostname filter
- `Priority` - Priority 0-7 (0=Emergency, 7=Debug)
- `Facility` - Facility 0-23
- `Message` - Text search
- `offset` - Pagination offset
- `limit` - Results limit (max 1000)

**Example:**
```bash
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=20"
```

**Response:**
```json
{
  "total": 152,
  "offset": 0,
  "limit": 20,
  "rows": [
    {
      "ID": 12345,
      "ReceivedAt": "2026-02-08T10:30:15Z",
      "FromHost": "webserver01",
      "Priority": 3,
      "Priority_Label": "Error",
      "Facility": 1,
      "Facility_Label": "user",
      "Message": "Authentication failed"
    }
  ]
}
```

### GET /meta
List all available database columns

```bash
curl -H "X-API-Key: YOUR_KEY" http://localhost:8000/meta
```

### GET /meta/{column}
Get distinct values for any column. Supports all filters from `/logs`.

```bash
# All hosts
curl -H "X-API-Key: YOUR_KEY" http://localhost:8000/meta/FromHost

# Syslog tags from specific host
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=webserver01"
```

## Priority Levels (RFC-5424)

| Value | Label | Description |
|-------|-------|-------------|
| 0 | Emergency | System unusable |
| 1 | Alert | Action must be taken immediately |
| 2 | Critical | Critical conditions |
| 3 | Error | Error conditions |
| 4 | Warning | Warning conditions |
| 5 | Notice | Normal but significant |
| 6 | Informational | Informational messages |
| 7 | Debug | Debug messages |

## Examples

### Get all errors from last hour
```bash
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&start_date=$(date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%SZ)&limit=100"
```

### Search for failed logins
```bash
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=failed+login&limit=50"
```

### Get logs from specific host and facility
```bash
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&Facility=4&limit=20"
```

## Performance

Benchmarks on Ubuntu 24.04, 4 cores, 8GB RAM:

| Metric | Value |
|--------|-------|
| Requests/sec | ~15,000 |
| Avg Latency | 0.6 ms |
| Memory Usage | 20 MB |
| Binary Size | 8 MB |

## Service Management

```bash
# Start
sudo systemctl start rsyslog-rest-api

# Stop
sudo systemctl stop rsyslog-rest-api

# Status
sudo systemctl status rsyslog-rest-api

# Logs
sudo journalctl -u rsyslog-rest-api -f
```

## Development

```bash
# Run locally
cp .env.example .env
nano .env
go run .

# Build
make build

# Clean
make clean
```

## Docker Test Environment

Simple test setup - build on host, run in Docker:

```bash
# 1. Build binary
make build-static

# 2. Start Docker
cd docker
docker-compose up -d

# 3. Watch startup
docker-compose logs -f
# Wait for "Environment Ready!" then Ctrl+C

# 4. Test
curl http://localhost:8000/health
curl -H "X-API-Key: test123456789" "http://localhost:8000/logs?limit=5"
./test.sh
```

**Includes:**
- rsyslog + MariaDB
- 10 test log entries
- Automatic setup

See [DOCKER.md](DOCKER.md) for details.

## Troubleshooting

### API won't start
```bash
# Check logs
sudo journalctl -u rsyslog-rest-api -n 50

# Test manually
cd /opt/rsyslog-rest-api
export $(cat .env | xargs)
./rsyslog-rest-api
```

### Database connection fails
```bash
# Verify rsyslog config
sudo cat /etc/rsyslog.d/mysql.conf

# Test MySQL connection
mysql -h localhost -u rsyslog -p Syslog -e "SELECT COUNT(*) FROM SystemEvents;"
```

### No logs returned
```bash
# Check if data exists
mysql -u rsyslog -p Syslog -e "SELECT COUNT(*) FROM SystemEvents;"

# Try wider time range
curl "http://localhost:8000/logs?start_date=2020-01-01T00:00:00Z"
```

## License

MIT License - Copyright (c) 2026 Phillip Grothues

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions welcome! Please open an issue or pull request.

---

**Made with ❤️ for high-performance log analysis**

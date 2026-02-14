# rsyslog REST API

High-performance REST API for rsyslog/MySQL written in Go.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/phil-bot/rsyslog-rest-api)](https://github.com/phil-bot/rsyslog-rest-api/releases)

## Features

- 🚀 **High Performance** - Written in Go for maximum speed
- 🔍 **Advanced Filtering** - Multi-value filters for complex queries
- 📊 **Extended Columns** - All 25+ SystemEvents fields accessible
- 🔐 **Secure** - API key authentication, SSL/TLS support
- 🐳 **Docker Ready** - Complete test environment with live data generation
- 📝 **REST API** - Clean JSON responses
- 🎯 **RFC-5424 Compliant** - Proper syslog severity and facility labels

## What's New in v0.2.2

### Multi-Value Filters
All filter parameters now support multiple values:

```bash
# Multiple hosts
curl "http://localhost:8000/logs?FromHost=web01&FromHost=web02&FromHost=db01"

# Multiple priorities (Error + Critical)
curl "http://localhost:8000/logs?Priority=2&Priority=3"

# Combinations
curl "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=3&Priority=4"
```

### Extended Columns
Response now includes all 25 SystemEvents columns:
- Core: ID, ReceivedAt, FromHost, Priority, Facility, Message
- Extended: DeviceReportedTime, SysLogTag, EventSource, EventUser, EventID, EventCategory, NTSeverity, Importance, CustomerID, SystemID, and more

### Live Log Generator (Docker)
Test environment now generates realistic logs continuously with all fields populated!

## Quick Start

### Prerequisites

- Go 1.21 or higher
- rsyslog with MySQL/MariaDB support
- MySQL/MariaDB database

### Installation

#### Option 1: Download Binary (Recommended)

```bash
# Download latest release
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-amd64

# Make executable
chmod +x rsyslog-rest-api-linux-amd64
sudo mv rsyslog-rest-api-linux-amd64 /usr/local/bin/rsyslog-rest-api

# Run
rsyslog-rest-api
```

#### Option 2: Install from Source

```bash
# Clone repository
git clone https://github.com/phil-bot/rsyslog-rest-api.git
cd rsyslog-rest-api

# Install
sudo make install

# Start service
sudo systemctl enable --now rsyslog-rest-api
```

### Configuration

Create or edit `/opt/rsyslog-rest-api/.env`:

```bash
# API Security (REQUIRED for production!)
API_KEY=your-secret-key-here

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8000

# Database Connection (RECOMMENDED)
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=your-database-password

# SSL/TLS (optional)
USE_SSL=false
SSL_CERTFILE=/opt/rsyslog-rest-api/certs/cert.pem
SSL_KEYFILE=/opt/rsyslog-rest-api/certs/key.pem

# CORS
ALLOWED_ORIGINS=*
```

**Generate API Key:**
```bash
openssl rand -hex 32
```

## API Endpoints

### GET /health

Health check endpoint.

**Example:**
```bash
curl http://localhost:8000/health
```

**Response:**
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-02-09T10:30:00Z"
}
```

### GET /logs

Retrieve log entries with filtering and pagination.

**Parameters:**

**Pagination:**
- `offset` - Results offset (default: 0)
- `limit` - Results limit (default: 10, max: 1000)

**Date Range:**
- `start_date` - ISO 8601 datetime (default: 24h ago)
- `end_date` - ISO 8601 datetime (default: now)

**Filters (all support multi-value):**
- `FromHost` - Filter by hostname(s)
- `Priority` - Filter by severity (0-7)
- `Facility` - Filter by facility (0-23)
- `Message` - Text search (supports multiple terms with OR logic)
- `SysLogTag` - Filter by syslog tag(s)

**Multi-Value Support:**

All filters accept multiple values by repeating the parameter:

```bash
# Single value (legacy)
?FromHost=web01

# Multiple values (NEW in v0.2.2!)
?FromHost=web01&FromHost=web02&FromHost=web03
```

**Examples:**

```bash
# Get latest 10 logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=10"

# Get errors from last hour
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&start_date=2025-02-09T09:00:00Z"

# Multiple hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&limit=20"

# Errors AND warnings from specific hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=3&Priority=4"

# Search multiple terms (OR logic)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=error&Message=failed&Message=timeout"
```

**Response:**
```json
{
  "total": 1234,
  "offset": 0,
  "limit": 10,
  "rows": [
    {
      "ID": 12345,
      "CustomerID": 42,
      "ReceivedAt": "2025-02-09T10:30:15Z",
      "DeviceReportedTime": "2025-02-09T10:30:13Z",
      "Facility": 1,
      "Facility_Label": "user",
      "Priority": 3,
      "Priority_Label": "Error",
      "FromHost": "webserver01",
      "Message": "Connection timeout to database",
      "SysLogTag": "nginx",
      "EventSource": "web-service",
      "EventUser": "www-data",
      "EventID": 504,
      "EventCategory": 5,
      "NTSeverity": 3000,
      "Importance": 4,
      "SystemID": 1,
      "InfoUnitID": 2
    }
  ]
}
```

**Note:** Extended fields use `omitempty` - they only appear if the database has a value (not NULL).

### GET /meta

List all available columns for filtering.

**Example:**
```bash
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta"
```

**Response:**
```json
{
  "available_columns": [
    "ID", "CustomerID", "ReceivedAt", "DeviceReportedTime",
    "Facility", "Priority", "FromHost", "Message", "NTSeverity",
    "Importance", "EventSource", "EventUser", "EventCategory",
    "EventID", "SysLogTag", "InfoUnitID", "SystemID"
  ],
  "usage": "GET /meta/{column} to get distinct values for a column"
}
```

### GET /meta/{column}

Get distinct values for a specific column.

**Supports the same filters as /logs** (including multi-value!)

**Examples:**

```bash
# Get all unique hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost"

# Get all priorities with labels
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/Priority"

# Get all SysLogTags
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag"

# Get hosts that logged errors (multi-value filter)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3&Priority=4"

# Get all event sources from web servers
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/EventSource?FromHost=web01&FromHost=web02"
```

**Response (for Priority/Facility with labels):**
```json
[
  {"val": 0, "label": "Emergency"},
  {"val": 1, "label": "Alert"},
  {"val": 3, "label": "Error"},
  {"val": 6, "label": "Informational"}
]
```

**Response (for other columns):**
```json
["webserver01", "webserver02", "dbserver01", "appserver01"]
```

## Available Columns

### Core Fields (always present)
- `ID` - Log entry ID
- `ReceivedAt` - Time received by rsyslog
- `FromHost` - Source hostname
- `Priority` - Severity (0-7)
- `Priority_Label` - RFC-5424 label (Emergency, Alert, Critical, Error, Warning, Notice, Informational, Debug)
- `Facility` - Facility (0-23)
- `Facility_Label` - RFC-5424 label (kern, user, mail, daemon, auth, syslog, etc.)
- `Message` - Log message text

### Extended Fields (when available)
- `CustomerID` - Customer identifier
- `DeviceReportedTime` - Original device timestamp
- `SysLogTag` - Syslog tag/program name
- `NTSeverity` - Windows NT severity level
- `Importance` - Event importance rating (1-5)
- `EventSource` - Event source identifier
- `EventUser` - Associated user
- `EventCategory` - Event category code
- `EventID` - Event ID number
- `EventBinaryData` - Binary event data
- `MaxAvailable` - Maximum available resource
- `CurrUsage` - Current resource usage
- `MinUsage` - Minimum resource usage
- `MaxUsage` - Maximum resource usage
- `InfoUnitID` - Information unit identifier
- `EventLogType` - Event log type
- `GenericFileName` - Associated filename
- `SystemID` - System identifier

## Priority Levels (RFC-5424)

| Value | Label | Description |
|-------|-------|-------------|
| 0 | Emergency | System is unusable |
| 1 | Alert | Action must be taken immediately |
| 2 | Critical | Critical conditions |
| 3 | Error | Error conditions |
| 4 | Warning | Warning conditions |
| 5 | Notice | Normal but significant |
| 6 | Informational | Informational messages |
| 7 | Debug | Debug messages |

## Facility Codes (RFC-5424)

| Value | Label | Description |
|-------|-------|-------------|
| 0 | kern | Kernel messages |
| 1 | user | User-level messages |
| 2 | mail | Mail system |
| 3 | daemon | System daemons |
| 4 | auth | Security/authorization |
| 16-23 | local0-7 | Local use |

See [RFC-5424](https://tools.ietf.org/html/rfc5424) for complete list.

## Docker Test Environment

Perfect for testing and development with **realistic live data**!

### Features
- ✅ Complete rsyslog + MariaDB setup
- ✅ 10 initial test entries
- ✅ **Live log generator** - 3 new entries every 10 seconds
- ✅ **All 25+ fields filled** with realistic values
- ✅ No API key required (optional)

### Quick Start

```bash
# Build binary
make build-static

# Start Docker
cd docker
docker-compose up -d

# Test
curl "http://localhost:8000/health"
curl "http://localhost:8000/logs?limit=5"

# Run test suite
./test-v0.2.2.sh
```

### What the Live Generator Does

**Every 10 seconds, 3 new realistic log entries with:**
- Realistic Event IDs (4624 for SSH login, 200-500 for nginx, etc.)
- Varied hosts (webserver01, dbserver01, appserver01, etc.)
- Multiple syslog tags (sshd, nginx, mysqld, postfix, etc.)
- All extended fields populated (EventSource, EventUser, NTSeverity, Importance, etc.)
- Weighted priorities (more INFO, less CRITICAL)

**Perfect for:**
- Testing multi-value filters
- Seeing extended columns in action
- Demos with growing data
- Performance testing

See [DOCKER.md](DOCKER.md) for detailed documentation.

## Development

### Build

```bash
# Standard build
make build

# Static build (no libc dependency)
make build-static

# Install locally
sudo make install
```

### Test

```bash
# Start Docker test environment
cd docker
docker-compose up -d

# Run test suite
./test-v0.2.2.sh

# Manual tests
curl "http://localhost:8000/health"
curl "http://localhost:8000/logs?limit=5"
```

### Project Structure

```
rsyslog-rest-api/
├── main.go                    # Main application
├── Makefile                   # Build automation
├── .env.example               # Configuration template
├── rsyslog-rest-api.service   # systemd service file
├── docker/                    # Test environment
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── entrypoint.sh
│   ├── log-generator.sh       # Realistic log generator
│   └── test-v0.2.2.sh         # Test suite
├── DOCKER.md                  # Docker documentation
└── README.md                  # This file
```

## Security

### API Key

**Always use an API key in production:**

```bash
# Generate secure key
openssl rand -hex 32

# Set in .env
API_KEY=your-generated-key-here
```

### SSL/TLS

**Enable SSL for production:**

```bash
# .env configuration
USE_SSL=true
SSL_CERTFILE=/path/to/cert.pem
SSL_KEYFILE=/path/to/key.pem
```

### Database Credentials

**Store securely in .env (not in rsyslog config):**

```bash
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=secure-password-here
```

**Set proper file permissions:**
```bash
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
```

## Troubleshooting

### API won't start

```bash
# Check logs
sudo journalctl -u rsyslog-rest-api -n 50

# Check config
sudo cat /opt/rsyslog-rest-api/.env

# Test database connection
mysql -u rsyslog -p Syslog
```

### No data returned

```bash
# Check database
mysql -u rsyslog -p Syslog -e "SELECT COUNT(*) FROM SystemEvents"

# Check filters
curl "http://localhost:8000/logs"  # No filters

# Check API key
curl -H "X-API-Key: YOUR_KEY" "http://localhost:8000/logs"
```

### Multi-value filters not working

```bash
# Correct syntax (repeat parameter)
?FromHost=web01&FromHost=web02

# INCORRECT (comma-separated doesn't work)
?FromHost=web01,web02
```

## Performance

- **Written in Go** - Compiled, high performance
- **Connection pooling** - Efficient database usage
- **Indexes created automatically** - Fast queries
- **Pagination** - Handles large datasets
- **Tested** - 10,000+ logs/second query performance

## Roadmap

### v0.3.0 (Planned)
- Negation filters (`exclude`, `not`)
- Advanced filter combinations
- Complex query support

### v0.4.0 (Planned)
- Statistics endpoint (`/stats`)
- Aggregations by column
- Timeline/histogram features
- Performance optimizations

See [GitHub Issues](https://github.com/phil-bot/rsyslog-rest-api/issues) for more.

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- **Issues:** [GitHub Issues](https://github.com/phil-bot/rsyslog-rest-api/issues)
- **Discussions:** [GitHub Discussions](https://github.com/phil-bot/rsyslog-rest-api/discussions)
- **Documentation:** [DOCKER.md](DOCKER.md) for Docker setup

## Credits

Created with ❤️ for the syslog community.

**Built with:**
- [Go](https://go.dev/)
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- [rsyslog](https://www.rsyslog.com/)
- [MariaDB](https://mariadb.org/)

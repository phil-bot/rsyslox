# Changelog

[← Back to overview](index.md)

All important changes to the project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/de/1.0.0/),
and this project follows [Semantic Versioning](https://semver.org/lang/de/).

---

## [Unreleased]

### Planned for v0.3.0
- Negation filters (`exclude`, `not`)
- Complex filter combinations
- Unit tests
- GitHub Actions CI/CD

### Planned for v0.4.0
- Statistics endpoint (`/stats`)
- Aggregation support
- Timeline/Histogram features
- WebSocket support for live logs

---

## [0.2.2] - 2025-02-09

### ✨ Added
- **Multi-Value Filter Support** - All filter parameters now accept multiple values
  ```bash
  # NEW: Multiple hosts
  ?FromHost=web01&FromHost=web02&FromHost=web03
  
  # NEW: Multiple Priorities
  ?Priority=3&Priority=4
  ```
  
- **Extended Columns** - All 25+ SystemEvents columns in Response
  - CustomerID, DeviceReportedTime, SysLogTag
  - EventSource, EventUser, EventID, EventCategory
  - NTSeverity, Importance, SystemID
  - And many more (see API Reference)
  
- **Enhanced Meta Endpoint** - `/meta/{column}` now supports filters
  ```bash
  # Hosts that had errors
  /meta/FromHost?Priority=3&Priority=4
  ```

- **Live Log Generator** (Docker) - Continuous test data
  - 3 new logs every 10 seconds
  - Realistic messages, event IDs, extended fields
  - 6 different hosts, 8 different tags

### 🔄 Changed
- Response structure extended - extended fields with `omitempty`
- Meta endpoint now provides structured responses for priority/facility

### 🐛 Fixed
- Filter validation improved
- NULL handling for extended columns

### 📚 Documentation
- New documentation structure
- Extended API Reference
- Docker Guide updated
- New examples for multi-value filters

---

## [0.2.1] - 2025-01-15

### 🐛 Fixed
- Database connection pooling optimized
- Fixed memory leak in long-running instances

### 🔄 Changed
- Default `limit` reduced from 100 to 10
- Logging improved

---

## [0.2.0] - 2024-12-20

### ✨ Added
- **RFC-5424 Labels** - Priority and facility with correct labels
  ```json
  {
    "Priority": 3,
    "Priority_Label": "Error",
    "Facility": 1,
    "Facility_Label": "user"
  }
  ```

- **Meta Endpoint** - `/meta` and `/meta/{column}`
  ```bash
  # Available columns
  GET /meta
  
  # Unique values
  GET /meta/FromHost
  GET /meta/Priority
  ```

- **CORS Support** - Configurable via `ALLOWED_ORIGINS`

- **SSL/TLS Support** - Optionally activatable
  ```bash
  USE_SSL=true
  SSL_CERTFILE=/path/to/cert.pem
  SSL_KEYFILE=/path/to/key.pem
  ```

### 🔄 Changed
- API-Key Authentication optional (developed for production)
- Database Configuration via Environment Variables (recommended)
- systemd service with security hardening

### 📚 Documentation
- README completely revised
- API Reference created
- Docker Guide added

---

## [0.1.5] - 2024-11-10

### 🐛 Fixed
- Pagination Edge Cases
- Date Range Validation
- SQL Injection Prevention improved

---

## [0.1.0] - 2024-10-01

### ✨ Initial Release

**Features:**
- REST API for rsyslog/MySQL
- `/health` endpoint
- `/logs` endpoint with filtering
  - Priority, Facility, FromHost, Message
  - Date Range (start_date, end_date)
  - Pagination (limit, offset)
- API-Key Authentication
- Connection Pooling
- Auto-Indexing

**Endpoints:**
- `GET /health` - Health Check
- `GET /logs` - Log Retrieval

**Configuration:**
- Environment Variables via `.env`
- rsyslog Config File Parsing (Fallback)

**Deployment:**
- systemd Service File
- Makefile for Build
- Static Binary Support

---

## Migration Guides

### Upgrade from v0.2.1 to v0.2.2

**No breaking changes!

✅ **Compatible:**
- Single-value filters continue to work
- All existing endpoints unchanged
- Response format extended (but backward-compatible)

✅ **Newly available:**
- Multi-value filter (optional use)
- Extended Columns (appear automatically if available)
- Meta-Filtering

**Update Steps:**
```bash
# 1. update binary
sudo systemctl stop rsyslog-rest-api
sudo cp rsyslog-rest-api-v0.2.2 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo systemctl start rsyslog-rest-api

# 2. testing
curl http://localhost:8000/health
curl -H "X-API-Key: KEY" "http://localhost:8000/logs?limit=1"

# Done! No config changes necessary.
```

---

### Upgrade from v0.1.x to v0.2.0

**Breaking Change: Response Format**

**v0.1.x Response:**
```json
{
  "Priority": 3,
  "Facility": 1
}
```

**v0.2.0+ Response:**
```json
{
  "Priority": 3,
  "Priority_Label": "Error",
  "Facility": 1,
  "Facility_Label": "user"
}
```

**Migration:**
- Clients must be able to ignore new fields (should not be a problem)
- No code change necessary if only `Priority`/`Facility` is used numerically

**Update Steps:**
```bash
# 1. Adapt .env (new variables)
sudo nano /opt/rsyslog-rest-api/.env

# Add new:
ALLOWED_ORIGINS=* # or specific domains

# Optional:
USE_SSL=false
SSL_CERTFILE=
SSL_KEYFILE=

# 2nd binary update
sudo systemctl stop rsyslog-rest-api
sudo cp rsyslog-rest-api-v0.2.0 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo systemctl start rsyslog-rest-api

# 3. testing
curl http://localhost:8000/health
```

---

## Deprecated Features

### v0.2.0
- **rsyslog config file parsing** (fallback, still supported)
  - Please use `DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASS` in `.env`
  - `RSYSLOG_CONFIG_PATH` will be removed in the future

---

## Security Fixes

### v0.2.1
- Fixed potential SQL injection in custom queries
- Improved input validation

### v0.2.0
- API key authentication added
- SSL/TLS support

---

## Performance Improvements

### v0.2.2
- Database Connection Pooling optimized
- Query performance for multi-value filter

### v0.2.0
- Automatic indexing on startup
- Connection Pooling implemented

---

## Known Issues

### v0.2.2
- No known issues

### v0.2.1
- No known issues

### v0.2.0
- ~~Memory leak on long running instances~~ (Fixed in v0.2.1)

---

## Upcoming Features (Planned)

### Short Term (v0.3.0)
- [ ] Negation filters
- [ ] Unit Tests
- [ ] CI/CD Pipeline

### Medium Term (v0.4.0)
- [ ] Statistics Endpoint
- [ ] Aggregations
- [ ] WebSocket Support

### Long Term
- [ ] PostgreSQL Support
- [ ] GraphQL API
- [ ] Web UI

---

## Version History Summary

| Version | Release Date | Highlights |
|---------|--------------|------------|
| 0.2.2 | 2025-02-09 | Multi-Value Filter, Extended Columns, Live Generator |
| 0.2.1 | 2025-01-15 | Bug Fixes, Performance |
| 0.2.0 | 2024-12-20 | RFC Labels, Meta Endpoint, SSL/TLS, CORS |
| 0.1.5 | 2024-11-10 | Bug Fixes |
| 0.1.0 | 2024-10-01 | Initial Release |

---

## Contributing

Found a bug? Want a feature?

→ [GitHub Issues](https://github.com/phil-bot/rsyslog-rest-api/issues)
→ [GitHub Discussions](https://github.com/phil-bot/rsyslog-rest-api/discussions)

---

[← Back to overview](index.md)

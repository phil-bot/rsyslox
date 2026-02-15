# API Reference

[← Back to overview](index.md)

Complete API documentation for all endpoints.

## 🔐 Authentication

### API Key Authentication

All protected endpoints require an API key in the header:

```http
X-API-Key: your-api-key-here
```

**Example:**
```bash
curl -H "X-API-Key: a3d7f8c9e2b4a6d8..." "http://localhost:8000/logs"
```

**Note:** If `API_KEY` in `.env` is empty, **no** authentication is required (only for development!).

---

## 📍 Endpoints overview

| Endpoint | Method | Auth | Description |
|----------|---------|------|--------------|
| `/health` | GET | ❌ | Health Check |
| `/logs` | GET | ✅ | Logs with filtering and pagination |
| `/meta` | GET | ✅ | List available columns |
| `/meta/{column}` | GET | ✅ | Unique values of a column |

---

## GET /health

Health Check endpoint without authentication.

### Request

```http
GET /health HTTP/1.1
Host: localhost:8000
```

```bash
curl http://localhost:8000/health
```

### Response

**Success (200 OK):**
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-02-09T10:30:00Z"
}
```

**Error (503 Service Unavailable):**
```json
{
  "status": "unhealthy",
  "database": "disconnected",
  "timestamp": "2025-02-09T10:30:00Z"
}
```

### Status codes

| Code | Meaning |
|------|-----------|
| 200 | API and database working |
| 503 | Database not accessible |

---

## GET /logs

Retrieves log entries with filtering and pagination.

### Request

```http
GET /logs?limit=10&Priority=3 HTTP/1.1
Host: localhost:8000
X-API-Key: your-api-key
```

### Query Parameter

#### Pagination

| Parameter | Type | Default | Description |
|-----------|------|---------|--------------|
| `offset` | Integer | 0 | Start position (skipping entries) |
| `limit` | Integer | 10 | Maximum number of results (max: 1000) |

#### Time filter

| Parameter | Type | Default | Format | Description |
|-----------|------|---------|--------|--------------|
| `start_date` | DateTime | -24h | ISO 8601 | Start date/time |
| `end_date` | DateTime | now | ISO 8601 | End date/time |

**ISO 8601 Format:** `2025-02-09T10:30:00Z` or `2025-02-09T10:30:00+01:00`

**Max. Time span:** 90 days

#### Content-Filter (Multi-Value!)

All filters support **multiple values** by repeating the parameter:

| Parameter | Type | Multi | Description |
|-----------|------|-------|--------------|
| `FromHost` | String | ✅ | Filter host name(s) |
| `Priority` | Integer | ✅ | Filter Severity (0-7) |
| `Facility` | Integer | ✅ | Filter Facility (0-23) |
| `Message` | String | ✅ | Text search (OR logic) |
| `SysLogTag` | String | ✅ | Filter syslog tag |

**Multi-Value Syntax:**
```bash
# Repeat multiple values = parameter
?FromHost=web01&FromHost=web02&FromHost=db01

# NOT: Comma-separated (does NOT work!)
?FromHost=web01,web02,db01 # ❌ FALSE
```

### Priority Values (RFC-5424)

| Value | Label | Description |
|-------|-------|--------------|
| 0 | Emergency | System unusable |
| 1 | Alert | Immediate action required |
| 2 | Critical | Critical condition |
| 3 | Error | Error conditions |
| 4 | Warning | Warnings |
| 5 | Notice | Normal but significant |
| 6 | Informational | Informational Messages |
| 7 | Debug | Debug Messages |

### Facility Values (RFC-5424)

| Value | Label | Description |
|-------|-------|--------------|
| 0 | kern | kernel messages |
| 1 | user | user-level messages |
| 2 | mail | mail system |
| 3 | daemon | system daemons |
| 4 | auth | security/authorization |
| 5 | syslog | syslog internal |
| 16-23 | local0-7 | Local use |

[Full list: RFC-5424](https://tools.ietf.org/html/rfc5424)

### Request examples

#### Simple query

```bash
# Latest 10 logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=10"
```

#### Time filter

```bash
# Logs of the last hour
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=2025-02-09T09:00:00Z&end_date=2025-02-09T10:00:00Z"

# Logs from yesterday
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=2025-02-08T00:00:00Z&end_date=2025-02-08T23:59:59Z"
```

#### Single-Value Filter

```bash
# Errors only (Priority 3)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3"

# From a specific host
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01"

# Text search
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=login"
```

#### Multi-Value Filter (NEW in v0.2.2!)

```bash
# Multiple hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&FromHost=db01"

# Errors AND warnings
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&Priority=4"

# Multiple facilities
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Facility=1&Facility=4&Facility=16"

# Multiple search terms (OR logic)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=error&Message=failed&Message=timeout"
```

#### Combined filters

```bash
# Errors from multiple hosts in the last hour
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=3&start_date=2025-02-09T09:00:00Z&limit=20"

# All priorities from specific host
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=dbserver01&Priority=2&Priority=3&Priority=4"
```

#### Pagination

```bash
# First 10 entries
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=10&offset=0"

# Next 10 entries
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=10&offset=10"

# Maximum (1000)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=1000"
```

### Response

**Success (200 OK):**

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

### Response fields

#### Mandatory fields (always present)

| Field | Type | Description |
|------|------|--------------|
| `ID` | Integer | Log entry ID |
| `ReceivedAt` | DateTime | Received time at rsyslog |
| `FromHost` | String | Source Hostname |
| `Priority` | Integer | Severity (0-7) |
| `Priority_Label` | String | RFC label (e.g. "Error") |
| `Facility` | Integer | Facility (0-23) |
| `Facility_Label` | String | RFC-Label (e.g. "user") |
| `Message` | String | Log message |

#### Extended fields (optional, if available)

| Field | Type | Description |
|------|------|--------------|
| `CustomerID` | Integer | Customer ID |
| `DeviceReportedTime` | DateTime | Original timestamp from device |
| `SysLogTag` | String | Syslog tag/program name |
| `EventSource` | String | EventSource |
| `EventUser` | String | Associated user |
| `EventID` | Integer | Event-ID |
| `EventCategory` | Integer | Event-Category |
| `NTSeverity` | Integer | Windows NT Severity |
| `Importance` | Integer | Importance Rating (1-5) |
| `EventBinaryData` | String | Binary event data |
| `MaxAvailable` | Integer | Max. available resources |
| `CurrUsage` | Integer | Current resource usage |
| `MinUsage` | Integer | Minimum utilization |
| `MaxUsage` | Integer | Maximum usage |
| `InfoUnitID` | Integer | Info-Unit-ID |
| `EventLogType` | String | Event-Log-Type |
| `GenericFileName` | String | Associated file name |
| `SystemID` | Integer | System-ID |

**Note:** Extended fields use `omitempty` - they only appear if the database has a value (not NULL).

### Error Responses

**401 Unauthorized:**
```json
{
  "error": "Invalid or missing API key"
}
```

**400 Bad Request:**
```json
{
  "error": "Priority must be between 0 and 7"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Database error"
}
```

---

## GET /meta

Lists all available columns for filtering.

### Request

```http
GET /meta HTTP/1.1
Host: localhost:8000
X-API-Key: your-api-key
```

```bash
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta"
```

### Response

**Success (200 OK):**

```json
{
  "available_columns": [
    "ID",
    "CustomerID",
    "ReceivedAt",
    "DeviceReportedTime",
    "Facility",
    "Priority",
    "FromHost",
    "Message",
    "NTSeverity",
    "Importance",
    "EventSource",
    "EventUser",
    "EventCategory",
    "EventID",
    "SysLogTag",
    "InfoUnitID",
    "SystemID"
  ],
  "usage": "GET /meta/{column} to get distinct values for a column"
}
```

---

## GET /meta/{column}

Retrieves unique values of a specific column.

### Request

```http
GET /meta/FromHost HTTP/1.1
Host: localhost:8000
X-API-Key: your-api-key
```

### Path parameter

| Parameter | Type | Description |
|-----------|------|--------------|
| `column` | String | Column name (from `/meta`) |

### Query Parameter

**All filters from `/logs` are supported (multi-value too!)

This enables filtered meta queries:

```bash
# All hosts that had errors
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3&Priority=4"

# All SysLogTags from specific hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=web01&FromHost=web02"
```

### Request examples

#### Simple meta queries

```bash
# All hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost"

# All priorities used
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/Priority"

# All SysLogTags
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag"

# All event sources
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/EventSource"
```

#### Filtered meta queries

```bash
# Hosts that had errors
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3"

# SysLogTags from specific hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=webserver01&FromHost=webserver02"

# Priorities in the last hour
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/Priority?start_date=2025-02-09T09:00:00Z"
```

### Response

#### For Priority/Facility (with labels)

```json
[
  { "val": 0, "label": "Emergency" },
  { "val": 1, "label": "Alert" },
  { "val": 3, "label": "Error" },
  { "val": 6, "label": "Informational" }
]
```

#### For integer columns (IDs)

```json
[1, 2, 5, 10, 42, 100]
```

#### For string columns

```json
[
  "webserver01",
  "webserver02",
  "dbserver01",
  "appserver01"
]
```

### Error Response

**400 Bad Request (invalid column):**

```json
{
  "error": "Invalid column 'InvalidCol'. Available columns: ID, CustomerID, ReceivedAt, ..."
}
```

---

## 🔢 HTTP status codes

| Code | Meaning | Usage |
|------|-----------|------------|
| 200 | OK | Successful request |
| 400 | Bad Request | Invalid Parameters |
| 401 | Unauthorized | API key missing or invalid |
| 500 | Internal Server Error | Server/Database Error |
| 503 | Service Unavailable | Database Unavailable |

---

## 📊 Rate Limiting

Currently **no** rate limiting implemented.

**Recommendation:** Use a reverse proxy (nginx/Apache) for rate limiting in production.

→ [Deployment: Reverse Proxy](deployment.md#reverse-proxy)

---

## 🔒 CORS

CORS is configured via `ALLOWED_ORIGINS` in `.env`:

```bash
# Development (all origins)
ALLOWED_ORIGINS=*

# Production (specific domains)
ALLOWED_ORIGINS=https://dashboard.example.com,https://app.example.com
```

→ [Configuration: CORS](configuration.md#cors-configuration)

---

## 💡 Best practices

### Performance

1. **Use limit:** Always set `limit` (Default: 10, Max: 1000)
2. **Limit time window:** Smaller time periods = faster queries
3. **Use pagination:** Retrieve large results in chunks
4. **Filter indexed fields:** `Priority`, `Facility`, `FromHost`, `ReceivedAt`

### Security

1. **Rotate API key:** Generate new key regularly
2. **Use HTTPS:** Always use SSL/TLS in Production
3. **Restrict CORS:** Only allow necessary origins
4. **Rate limiting:** Via reverse proxy

### Error handling

1. **Check HTTP status:** Do not accept only 200
2. **Retry-Logic:** At 500/503 with backoff
3. **Set timeout:** Configure client-side timeout

---

## 📖 Further resources

- [Installation](installation.md) - Setup and deployment
- [Configuration](configuration.md) - Configuration
- [Examples](examples.md) - Practical examples
- [Troubleshooting](troubleshooting.md) - Troubleshooting

---

[← Back to overview](index.md) | [Forward to Examples →](examples.md)

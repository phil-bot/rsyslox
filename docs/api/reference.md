# API Reference

Complete API documentation for rsyslog REST API v0.2.3.

## Authentication

All protected endpoints require an API key:

```http
X-API-Key: your-api-key-here
```

**Example:**
```bash
curl -H "X-API-Key: a3d7f8c9..." "http://localhost:8000/logs"
```

?> **Development:** If `API_KEY` in `.env` is empty, authentication is disabled.

## Base URL

```
http://localhost:8000
```

Or with custom host/port from configuration.

## Endpoints

### GET /health

Health check endpoint (no authentication required).

**Request:**
```bash
curl http://localhost:8000/health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-02-15T10:30:00Z"
}
```

**Response (503 Service Unavailable):**
```json
{
  "status": "unhealthy",
  "database": "disconnected",
  "timestamp": "2025-02-15T10:30:00Z"
}
```

---

### GET /logs

Retrieve log entries with filtering and pagination.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `offset` | Integer | 0 | Skip N entries |
| `limit` | Integer | 10 | Max results (max: 1000) |
| `start_date` | DateTime | -24h | Start datetime (ISO 8601) |
| `end_date` | DateTime | now | End datetime (ISO 8601) |
| `FromHost` | String | - | Filter by hostname (multi-value) |
| `Priority` | Integer | - | Filter by severity 0-7 (multi-value) |
| `Facility` | Integer | - | Filter by facility 0-23 (multi-value) |
| `Message` | String | - | Text search (multi-value, OR) |
| `SysLogTag` | String | - | Filter by syslog tag (multi-value) |

**Multi-Value Support (v0.2.2+):**

Repeat parameter for multiple values:

```bash
# Multiple hosts
?FromHost=web01&FromHost=web02&FromHost=db01

# Multiple priorities
?Priority=3&Priority=4

# NOT this (won't work):
?FromHost=web01,web02  # ❌ Wrong
```

**Priority Values (RFC-5424):**

| Value | Label | Description |
|-------|-------|-------------|
| 0 | Emergency | System unusable |
| 1 | Alert | Action required immediately |
| 2 | Critical | Critical conditions |
| 3 | Error | Error conditions |
| 4 | Warning | Warning conditions |
| 5 | Notice | Normal but significant |
| 6 | Informational | Informational |
| 7 | Debug | Debug messages |

**Examples:**

```bash
# Latest 10 logs
curl -H "X-API-Key: $KEY" "http://localhost:8000/logs?limit=10"

# Errors from last hour
START=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Priority=3&start_date=$START"

# Multiple hosts (NEW in v0.2.2!)
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&FromHost=db01"

# Errors AND warnings
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Priority=3&Priority=4"

# Combined filters
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=3&limit=20"
```

**Response (200 OK):**

```json
{
  "total": 1234,
  "offset": 0,
  "limit": 10,
  "rows": [
    {
      "ID": 12345,
      "CustomerID": 42,
      "ReceivedAt": "2025-02-15T10:30:15Z",
      "DeviceReportedTime": "2025-02-15T10:30:13Z",
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

**Response Fields:**

Core fields (always present):
- `ID` - Log entry ID
- `ReceivedAt` - Time received by rsyslog
- `FromHost` - Source hostname
- `Priority` / `Priority_Label` - Severity
- `Facility` / `Facility_Label` - Facility
- `Message` - Log message

Extended fields (when available):
- `CustomerID`, `DeviceReportedTime`, `SysLogTag`
- `EventSource`, `EventUser`, `EventID`, `EventCategory`
- `NTSeverity`, `Importance`, `SystemID`, `InfoUnitID`
- More... (25+ total fields)

---

### GET /meta

List all available columns.

**Request:**
```bash
curl -H "X-API-Key: $KEY" "http://localhost:8000/meta"
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

---

### GET /meta/{column}

Get distinct values for a column.

**Supports all filters from /logs!**

**Examples:**

```bash
# All hosts
curl -H "X-API-Key: $KEY" "http://localhost:8000/meta/FromHost"

# All priorities with labels
curl -H "X-API-Key: $KEY" "http://localhost:8000/meta/Priority"

# Hosts that logged errors (filtered!)
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3&Priority=4"

# SysLogTags from specific hosts
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=web01&FromHost=web02"
```

**Response (Priority/Facility with labels):**
```json
[
  { "val": 0, "label": "Emergency" },
  { "val": 3, "label": "Error" },
  { "val": 6, "label": "Informational" }
]
```

**Response (Integer columns):**
```json
[1, 2, 5, 10, 42]
```

**Response (String columns):**
```json
["webserver01", "webserver02", "dbserver01"]
```

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Success |
| 400 | Bad Request - Invalid parameters |
| 401 | Unauthorized - API key missing/invalid |
| 500 | Internal Server Error |
| 503 | Service Unavailable - Database down |

## Rate Limiting

Currently **no** built-in rate limiting.

!> **Production:** Use reverse proxy (nginx/Apache) for rate limiting.

## Best Practices

**Performance:**
- Always set `limit` parameter
- Use smaller time windows
- Paginate large results
- Filter on indexed fields

**Security:**
- Use HTTPS in production
- Rotate API keys regularly
- Restrict CORS origins
- Implement rate limiting

**Reliability:**
- Check `/health` before queries
- Handle errors gracefully
- Implement retry logic
- Monitor API availability

## What's New in v0.2.3

- ✅ Improved multi-value filter performance
- ✅ Better error validation messages
- ✅ Enhanced meta endpoint filtering
- ✅ Bug fixes and stability improvements

[View Full Changelog](../development/changelog.md)

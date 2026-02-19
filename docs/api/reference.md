# API Reference

Complete API documentation for rsyslog REST API v0.3.0.

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

---

## Priority vs. Severity

rsyslog changed how it populates the `Priority` database column depending on version:

| rsyslog version | `Priority` column contains |
|---|---|
| < 8.2204.0 (legacy) | Severity only (0–7) |
| ≥ 8.2204.0 (modern) | RFC PRI = `Facility × 8 + Severity` |

The API detects the storage format automatically at startup by sampling the oldest
and newest non-kernel entries. Mixed datasets (produced by a rsyslog upgrade) are
handled correctly on a per-row basis.

**API response fields are always RFC-5424 compliant**, regardless of the rsyslog version:

| Field | Description | Example |
|---|---|---|
| `Priority` | RFC PRI value (`Facility × 8 + Severity`) | `25` |
| `Severity` | Severity value 0–7 | `1` |
| `Severity_Label` | Human-readable severity | `"Alert"` |
| `Facility` | Facility value 0–23 | `3` |
| `Facility_Label` | Human-readable facility | `"daemon"` |

---

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
| `Severity` | Integer | - | Filter by severity 0-7 (multi-value) |
| `Priority` | Integer | - | Deprecated alias for `Severity` |
| `Facility` | Integer | - | Filter by facility 0-23 (multi-value) |
| `Message` | String | - | Text search (multi-value, OR) |
| `SysLogTag` | String | - | Filter by syslog tag (multi-value) |

**Multi-Value Support:**

Repeat parameter for multiple values:

```bash
?Severity=3&Severity=4
?FromHost=web01&FromHost=web02
```

**Severity Values (RFC-5424):**

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
  "http://localhost:8000/logs?Severity=3&start_date=$START"

# Multiple hosts
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&FromHost=db01"

# Errors AND warnings
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Severity=3&Severity=4"

# Combined filters
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?FromHost=web01&Severity=3&limit=20"
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
      "Facility": 3,
      "Facility_Label": "daemon",
      "Priority": 25,
      "Severity": 1,
      "Severity_Label": "Alert",
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
- `Priority` - RFC PRI value (`Facility × 8 + Severity`)
- `Severity` / `Severity_Label` - RFC severity (0–7) with label
- `Facility` / `Facility_Label` - RFC facility with label
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
    "EventID", "SysLogTag", "InfoUnitID", "SystemID",
    "Severity"
  ],
  "usage": "GET /meta/{column} to get distinct values for a column"
}
```

?> **Note:** `Severity` is a virtual column — it is computed from the `Priority` column
at query time and is not a physical database column.

---

### GET /meta/{column}

Get distinct values for a column across **all data** (no default time filter).

**Key behavior:**
- Without any filters: returns all distinct values from the entire dataset.
- Filters are **optional** and narrow the result set when provided.
- Unlike `/logs`, **no default date range is applied**.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `start_date` | DateTime | **none** | Optional start datetime (ISO 8601). No default. |
| `end_date` | DateTime | **none** | Optional end datetime (ISO 8601). No default. |
| `FromHost` | String | - | Filter by hostname (multi-value) |
| `Severity` | Integer | - | Filter by severity 0-7 (multi-value) |
| `Priority` | Integer | - | Deprecated alias for `Severity` |
| `Facility` | Integer | - | Filter by facility 0-23 (multi-value) |
| `Message` | String | - | Text search (multi-value, OR) |
| `SysLogTag` | String | - | Filter by syslog tag (multi-value) |

**Examples:**

```bash
# All distinct hosts (entire dataset)
curl -H "X-API-Key: $KEY" "http://localhost:8000/meta/FromHost"

# All severity values with labels (virtual column)
curl -H "X-API-Key: $KEY" "http://localhost:8000/meta/Severity"

# Hosts that logged errors
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/FromHost?Severity=3&Severity=4"

# SysLogTags from specific hosts
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=web01&FromHost=web02"

# Hosts with errors in a specific time window
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/FromHost?Severity=3&start_date=2025-02-01T00:00:00Z&end_date=2025-02-15T23:59:59Z"
```

**Response (`Severity` or `Facility` — with labels):**
```json
[
  { "val": 1, "label": "Alert" },
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

---

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
- Always set `limit` parameter on `/logs`
- Use smaller time windows on `/logs`
- Paginate large results
- Filter on indexed fields
- Cache `/meta` responses — they cover the full dataset and change slowly

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

## What's New in v0.3.0

- ✅ `Severity` is now the correct RFC-5424 field (0–7), exposed in all log responses
- ✅ `Priority` in responses now contains the true RFC PRI value (`Facility × 8 + Severity`)
- ✅ Automatic detection of rsyslog Priority column format (legacy / modern / mixed)
- ✅ Mixed datasets (before/after rsyslog upgrade) are handled correctly per row
- ✅ `?Severity=` filter parameter introduced; `?Priority=` kept as deprecated alias
- ✅ `/meta/Severity` virtual column returns distinct severity values with labels

## What's New in v0.2.4

- ✅ `/meta/{column}` returns all distinct values without a default time filter

## What's New in v0.2.3

- ✅ Improved multi-value filter performance
- ✅ Better error validation messages
- ✅ Enhanced meta endpoint filtering
- ✅ Bug fixes and stability improvements

[View Full Changelog](../development/changelog.md)

# Quick Start

Get up and running with rsyslox in under 5 minutes.

## Prerequisites

- rsyslox installed and running — see [Installation](installation.md)
- rsyslog writing to MySQL/MariaDB

## Step 1: Open the Log Viewer

Navigate to `http://<your-host>:8000` in your browser. The log viewer loads immediately after setup.

The left sidebar contains all filter controls:

- **Time range** — select a duration (15m / 1h / 6h / 24h / 7d / 30d) or set custom dates
- **Severity** — click one or more severity levels to filter
- **Facility**, **Host**, **Tag** — multi-select filter pills
- **Message search** — free-text search across the message field

## Step 2: Browse Logs

The table shows log entries matching the active filters. Click any row to open the detail panel with the full message and all fields.

**Useful actions:**

| Action | How |
|---|---|
| Select rows | Checkbox on the left |
| Export selection | Select rows → Export CSV / Export JSON |
| Toggle auto-refresh | Refresh button in the toolbar (shows countdown) |
| Navigate pages | Pagination bar at the bottom |
| Load all entries | Toggle the view mode button in the toolbar |

## Step 3: Use the API

rsyslox exposes a REST API for external tools. You need a read-only API key — create one in **Admin → API Keys**.

**Get your key, then test it:**

```bash
API_KEY="your-key-here"
curl -H "X-API-Key: $API_KEY" "http://localhost:8000/api/logs?limit=5"
```

**Retrieve errors from the last hour:**

```bash
START=$(date -u -d "1 hour ago" "+%Y-%m-%dT%H:%M:%SZ")
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/api/logs?Severity=3&start_date=$START"
```

**Filter by multiple hosts:**

```bash
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/api/logs?FromHost=web01&FromHost=web02"
```

## Step 4: Explore the API Docs

Interactive API documentation is available directly in the app:

```
http://<your-host>:8000/docs
```

It covers all endpoints, parameters, and response formats — no separate reference page needed.

## Admin Panel

Navigate to `http://<your-host>:8000/admin`. Log in with your admin password to manage:

- Server settings and SSL
- Log cleanup configuration
- API key creation and revocation
- Browser preferences (language, font size, time format, auto-refresh interval)

## Next Steps

- [Deploy to Production](../guides/deployment.md)
- [Troubleshooting](../guides/troubleshooting.md)
- [Security Guide](../guides/security.md)

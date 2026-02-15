# Quick Start

Get up and running with rsyslog REST API in 5 minutes!

## Prerequisites

- rsyslog REST API installed
- rsyslog configured with MySQL
- MySQL/MariaDB running

If not installed yet: [Installation Guide](installation.md)

## Step 1: Start the API

```bash
# Navigate to installation directory
cd /opt/rsyslog-rest-api

# Start API (foreground for testing)
rsyslog-rest-api
```

You should see:
```
Database connection established
Starting HTTP server on http://0.0.0.0:8000
```

## Step 2: Health Check

Open a new terminal and test:

```bash
curl http://localhost:8000/health
```

Response:
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-02-15T10:30:00Z"
}
```

✅ API is running!

## Step 3: Your First API Call

Get your API key:

```bash
API_KEY=$(grep "^API_KEY=" /opt/rsyslog-rest-api/.env | cut -d'=' -f2)
echo $API_KEY
```

Retrieve latest logs:

```bash
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?limit=5"
```

## Step 4: Try Filters

### Filter by Priority

Get only errors:

```bash
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=10"
```

### Filter by Host

Logs from specific server:

```bash
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&limit=10"
```

### Multiple Filters (v0.2.3!)

Errors from multiple hosts:

```bash
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=3"
```

## Step 5: Query Metadata

Get all hosts:

```bash
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/meta/FromHost"
```

Get all priorities:

```bash
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/meta/Priority"
```

## Common Use Cases

### Last Hour Errors

```bash
START=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?Priority=3&start_date=$START"
```

### Search for Keyword

```bash
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?Message=login"
```

### Pagination

```bash
# First page
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?limit=10&offset=0"

# Second page
curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?limit=10&offset=10"
```

## Production Setup

Ready for production?

```bash
# Stop foreground process (Ctrl+C)

# Install systemd service
sudo cp rsyslog-rest-api.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now rsyslog-rest-api

# Check status
sudo systemctl status rsyslog-rest-api
```

## Next Steps

- [Full API Reference](../api/reference.md)
- [More Examples](../api/examples.md)
- [Deploy to Production](../guides/deployment.md)

## Troubleshooting

**API won't start:**
- Check logs: `sudo journalctl -u rsyslog-rest-api -n 50`
- Verify database: `mysql -u rsyslog -p Syslog`

**No logs returned:**
- Check database has data: `mysql -u rsyslog -p Syslog -e "SELECT COUNT(*) FROM SystemEvents"`
- Try without filters first

**Authentication failed:**
- Verify API key: `grep API_KEY /opt/rsyslog-rest-api/.env`
- Check header format: `X-API-Key: your-key`

More help: [Troubleshooting Guide](../guides/troubleshooting.md)

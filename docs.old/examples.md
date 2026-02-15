# Examples & Use Cases

[← Back to overview](index.md)

Practical examples of common use cases.

## 🎯 Basic queries

### Retrieve latest logs

```bash
# Last 10 logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=10"

# Last 50 logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=50"
```

### Logs from a specific time period

```bash
# Last hour
START=$(date -u -d '1 hour ago' '+%Y-%m-%dD%H:%M:%SZ')
END=$(date -u '+%Y-%m-%dD%H:%M:%SZ')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=$START&end_date=$END"

# Today (00:00 until now)
START=$(date -u '+%Y-%m-%dT00:00:00Z')
END=$(date -u '+%Y-%m-%dD%H:%M:%SZ')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=$START&end_date=$END"

# Specific tag
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=2025-02-09T00:00:00Z&end_date=2025-02-09T23:59:59Z"

# Last 24 hours
START=$(date -u -d '24 hours ago' '+%Y-%m-%dD%H:%M:%SZ')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=$START&limit=100"
```

---

## 🔍 Filter by severity

### Single priority

```bash
# Errors only (priority 3)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=20"

# Only warnings (priority 4)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=4&limit=20"

# Critical problems (priority 0-2)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=0&Priority=1&Priority=2&limit=50"
```

### Multiple priorities

```bash
# Errors AND warnings
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&Priority=4&limit=50"

# All problems (critical, error, warning)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=2&Priority=3&Priority=4&limit=100"

# Only informative logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=6&Priority=7&limit=50"
```

---

## 🖥️ Filter by host

### Single host

```bash
# Logs from webserver01
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&limit=20"

# Logs from dbserver01
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=dbserver01&limit=20"
```

### Multiple hosts (NEW in v0.2.2!)

```bash
# All web servers
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&FromHost=webserver02&FromHost=webserver03&limit=50"

# Web + Database Server
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&FromHost=dbserver01&limit=50"

# All servers in a group
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=app01&FromHost=app02&FromHost=app03&FromHost=app04&limit=100"
```

---

## 🏷️ Filter by SysLogTag

```bash
# SSH logs only
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?SysLogTag=sshd&limit=20"

# Only nginx logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?SysLogTag=nginx&limit=20"

# Multiple tags
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?SysLogTag=sshd&SysLogTag=sudo&SysLogTag=systemd&limit=50"
```

---

## 🔎 Text search

### Single search term

```bash
# Search for "login"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=login&limit=20"

# Search for "error"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=error&limit=50"

# Search for "failed"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=failed&limit=30"
```

### Multiple search terms (OR logic!)

```bash
# "error" OR "failed" OR "timeout"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=error&Message=failed&Message=timeout&limit=50"

# "login" OR "logout" OR "authentication"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=login&Message=logout&Message=authentication&limit=50"
```

---

## 🎨 Combined filters

### Host + Priority

```bash
# Errors from webserver01
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&Priority=3&limit=20"

# Warnings from multiple hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=4&limit=50"
```

### Host + Priority + Time

```bash
# Errors of the last hour from dbserver01
START=$(date -u -d '1 hour ago' '+%Y-%m-%dD%H:%M:%SZ')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=dbserver01&Priority=3&start_date=$START&limit=50"

# All problems from multiple hosts today
START=$(date -u '+%Y-%m-%dT00:00:00Z')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=2&Priority=3&Priority=4&start_date=$START&limit=100"
```

### Host + Message

```bash
# Login attempts on webserver01
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&Message=login&limit=30"

# Error logs with specific keywords
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=appserver01&Message=error&Message=exception&Priority=3&limit=50"
```

### Complex filter combination

```bash
# Errors AND warnings from multiple web servers with "timeout" in message
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&FromHost=web03&Priority=3&Priority=4&Message=timeout&limit=50"

# Critical problems of DB servers today
START=$(date -u '+%Y-%m-%dT00:00:00Z')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=db01&FromHost=db02&Priority=0&Priority=1&Priority=2&start_date=$START&limit=100"
```

---

## 📄 Pagination

### Basic paging

```bash
# First 50 entries
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=50&offset=0"

# Second page (51-100)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=50&offset=50"

# Third page (101-150)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=50&offset=100"
```

### Pagination with filter

```bash
# Page 1 of Errors
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=50&offset=0"

# Page 2 of Errors
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=50&offset=50"
```

### Iterate all data (bash script)

```bash
#!/bin/bash
API_KEY="your-api-key"
LIMIT=100
OFFSET=0

while true; do
    RESPONSE=$(curl -s -H "X-API-Key: $API_KEY" \
      "http://localhost:8000/logs?limit=$LIMIT&offset=$OFFSET&Priority=3")
    
    COUNT=$(echo "$RESPONSE" | jq '.rows | length')
    
    if [ "$COUNT" -eq 0 ]; then
        echo "Done!"
        break
    fi
    
    echo "Process $COUNT entries (offset: $OFFSET)..."
    echo "$RESPONSE" | jq '.rows[] | .Message'
    
    OFFSET=$((OFFSET + LIMIT))
    sleep 1
done
```

---

## 📊 Query metadata

### Available hosts

```bash
# All hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost"

# Hosts that had errors today
START=$(date -u '+%Y-%m-%dT00:00:00Z')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3&start_date=$START"

# Hosts with problems (priority 2-4)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=2&Priority=3&Priority=4"
```

### Available tags

```bash
# All SysLogTags
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag"

# Tags from specific hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=webserver01&FromHost=webserver02"

# Tags from error logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag?Priority=3"
```

### Available priorities

```bash
# All used priorities
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/Priority"

# Priorities from specific host
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/Priority?FromHost=dbserver01"
```

---

## 🛠️ Practical use cases

### Monitoring dashboard

```bash
#!/bin/bash
# dashboard.sh - Simple Monitoring Dashboard

API_KEY="your-api-key"
API_URL="http://localhost:8000"

echo "=== Syslog Dashboard ==="
echo ""

# Total logs today
START=$(date -u '+%Y-%m-%dT00:00:00Z')
TOTAL=$(curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?start_date=$START&limit=1" | jq .total)
echo "Total logs today: $TOTAL"

# Errors today
ERRORS=$(curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?Priority=3&start_date=$START&limit=1" | jq .total)
echo "Errors today: $ERRORS"

# Top 5 hosts with the most errors
echo ""
echo "Top 5 hosts with errors:"
curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/meta/FromHost?Priority=3&start_date=$START" | jq -r '.[]' | head -5

# Last 5 Critical/Errors
echo ""
echo "Last 5 critical logs:"
curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?Priority=2&Priority=3&limit=5" | \
  jq -r '.rows[] | "\(.ReceivedAt) [\(.FromHost)] \(.Message)"'
```

### Login monitoring

```bash
#!/bin/bash
# login-monitor.sh - SSH login monitoring

API_KEY="your-api-key"
API_URL="http://localhost:8000"
START=$(date -u -d '1 hour ago' '+%Y-%m-%dD%H:%M:%SZ')

# Successful logins
echo "=== Successful SSH logins (last hour) ==="
curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?SysLogTag=sshd&Message=Accepted&start_date=$START&limit=50" | \
  jq -r '.rows[] | "\(.ReceivedAt) - \(.FromHost) - \(.Message)"'

echo ""

# Failed logins
echo "=== Failed SSH logins (last hour) ==="
curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?SysLogTag=sshd&Message=Failed&start_date=$START&limit=50" | \
  jq -r '.rows[] | "\(.ReceivedAt) - \(.FromHost) - \(.Message)"'
```

### Error report

```bash
#!/bin/bash
# error-report.sh - Daily error report

API_KEY="your-api-key"
API_URL="http://localhost:8000"
START=$(date -u '+%Y-%m-%dT00:00:00Z')
DATE=$(date '+%Y-%m-%d')

OUTPUT="error-report-$DATE.txt"

{
  echo "====================================="
  echo "Error report for $DATE"
  echo "====================================="
  echo ""
  
  # Statistics
  TOTAL_ERRORS=$(curl -s -H "X-API-Key: $API_KEY" \
    "$API_URL/logs?Priority=3&start_date=$START&limit=1" | jq .total)
  
  TOTAL_WARNINGS=$(curl -s -H "X-API-Key: $API_KEY" \
    "$API_URL/logs?Priority=4&start_date=$START&limit=1" | jq .total)
  
  echo "Errors: $TOTAL_ERRORS"
  echo "Warnings: $TOTAL_WARNINGS"
  echo ""
  
  # Hosts with the most errors
  echo "Hosts with most errors:"
  curl -s -H "X-API-Key: $API_KEY" \
    "$API_URL/meta/FromHost?Priority=3&start_date=$START" | \
    jq -r '.[]' | nl
  
  echo ""
  echo "Top 20 Error Messages:"
  curl -s -H "X-API-Key: $API_KEY" \
    "$API_URL/logs?Priority=3&start_date=$START&limit=20" | \
    jq -r '.rows[] | "[\(.FromHost)] \(.Message)"'
  
} > "$OUTPUT"

echo "Report saved: $OUTPUT"
```

### Alert system (simple)

```bash
#!/bin/bash
# alert-check.sh - Simple alert system
# Cron: */5 * * * * * /path/to/alert-check.sh

API_KEY="your-api-key"
API_URL="http://localhost:8000"
START=$(date -u -d '5 minutes ago' '+%Y-%m-%dD%H:%M:%SZ')
ALERT_EMAIL="admin@example.com"

# Check critical logs
CRITICAL=$(curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?Priority=0&Priority=1&Priority=2&start_date=$START&limit=1" | jq .total)

if [ "$CRITICAL" -gt 0 ]; then
  echo "ALERT: $CRITICAL critical logs in the last 5 minutes!" | \
    mail -s "Critical Syslog Alert" "$ALERT_EMAIL"
fi

# Too many errors?
ERRORS=$(curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?Priority=3&start_date=$START&limit=1" | jq .total)

if [ "$ERRORS" -gt 100 ]; then
  echo "WARNING: $ERRORS Errors in the last 5 minutes!" | \
    mail -s "High Error Rate Alert" "$ALERT_EMAIL"
fi
```

---

## 🐍 Python examples

### Simple client

```python
#!/usr/bin/env python3
import requests
import json
from datetime import datetime, timedelta

API_KEY = "your-api-key"
API_URL = "http://localhost:8000"

def get_logs(priority=None, from_host=None, limit=10):
    """Retrieve logs with filters"""
    headers = {"X-API-Key": API_KEY}
    params = {"limit": limit}
    
    if priority:
        params["Priority"] = priority
    if from_host:
        params["FromHost"] = from_host
    
    response = requests.get(f"{API_URL}/logs", headers=headers, params=params)
    response.raise_for_status()
    return response.json()

# Example usage
if __name__ == "__main__":
    # Last 10 errors
    data = get_logs(priority=3, limit=10)
    
    print(f "Total Errors: {data['total']}")
    print(f"\nLast {len(data['rows'])} Errors:")
    
    for log in data['rows']:
        print(f"[{log['ReceivedAt']}] {log['FromHost']}: {log['Message']}")
```

### Monitoring with Multi-Value

```python
#!/usr/bin/env python3
import requests
from datetime import datetime, timedelta

API_KEY = "your-api-key"
API_URL = "http://localhost:8000"

def get_errors_from_hosts(hosts, hours=1):
    """Retrieve errors from multiple hosts"""
    headers = {"X-API-Key": API_KEY}
    
    # Calculate start time
    start = (datetime.utcnow() - timedelta(hours=hours)).isoformat() + 'Z'
    
    # Multi-Value Parameter
    params = [
        ("Priority", "3"),
        ("start_date", start),
        ("limit", "100")
    ]
    
    # Add hosts (Multi-Value!)
    for host in hosts:
        params.append(("FromHost", host))
    
    response = requests.get(f"{API_URL}/logs", headers=headers, params=params)
    response.raise_for_status()
    return response.json()

# Example
hosts = ["webserver01", "webserver02", "webserver03"]
data = get_errors_from_hosts(hosts, hours=1)

print(f "Errors from {len(hosts)} hosts: {data['total']}")
for log in data['rows']:
    print(f"{log['FromHost']}: {log['Message']}")
```

---

## 🌐 JavaScript/Node.js examples

### Fetch API (Browser/Node)

```javascript
const API_KEY = 'your-api-key';
const API_URL = 'http://localhost:8000';

async function getLogs(options = {}) {
  const params = new URLSearchParams();
  
  if (options.priority) params.append('Priority', options.priority);
  if (options.fromHost) params.append('FromHost', options.fromHost);
  if (options.limit) params.append('limit', options.limit);
  
  const response = await fetch(`${API_URL}/logs?${params}`, {
    headers: {
      'X-API-Key': API_KEY
    }
  });
  
  if (!response.ok) {
    throw new Error(`API Error: ${response.status}`);
  }
  
  return response.json();
}

// Usage
getLogs({ priority: 3, limit: 20 })
  .then(data => {
    console.log(`Total: ${data.total}`);
    data.rows.forEach(log => {
      console.log(`${log.FromHost}: ${log.Message}`);
    });
  })
  .catch(console.error);
```

### Multi-Value Filter (JavaScript)

```javascript
async function getLogsFromMultipleHosts(hosts, priority) {
  const params = new URLSearchParams();
  
  // Multi-Value: Hosts
  hosts.forEach(host => params.append('FromHost', host));
  
  // Priority
  if (Array.isArray(priority)) {
    priority.forEach(p => params.append('Priority', p));
  } else if (priority) {
    params.append('Priority', priority);
  }
  
  params.append('limit', '100');
  
  const response = await fetch(`${API_URL}/logs?${params}`, {
    headers: { 'X-API-Key': API_KEY }
  });
  
  return response.json();
}

// Usage
const hosts = ['web01', 'web02', 'db01'];
const priorities = [3, 4]; // Errors & Warnings

getLogsFromMultipleHosts(hosts, priorities)
  .then(data => console.log(data))
  .catch(console.error);
```

---

## 💡 Tips & Best Practices

### Performance

1. **Always use limit** - default is only 10, max 1000
2. **Limit time frame** - Shorter time frames = faster
3. **Use pagination** - Not all data at once
4. **Filter indexed fields** - `FromHost`, `Priority`, `ReceivedAt`

### Error handling

```bash
# Check HTTP status
RESPONSE=$(curl -s -w "\n%{http_code}" -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?limit=10")

HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
BODY=$(echo "$RESPONSE" | head -n -1)

if [ "$HTTP_CODE" -eq 200 ]; then
  echo "$BODY" | jq
else
  echo "Error: HTTP $HTTP_CODE"
  echo "$BODY"
fi
```

### Logging

```bash
# Log requests
LOG_FILE="api-requests.log"

curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=10" | \
  tee -a "$LOG_FILE" | jq
```

---

[← Back to overview](index.md) | [Next to Deployment →](deployment.md)

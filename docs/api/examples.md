# API Examples

Practical examples for common use cases.

## Basic Queries

### Latest Logs

```bash
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?limit=10"
```

### Specific Time Range

```bash
# Last 24 hours
START=$(date -u -d '24 hours ago' '+%Y-%m-%dT%H:%M:%SZ')
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?start_date=$START&limit=50"

# Yesterday
START=$(date -u -d 'yesterday 00:00' '+%Y-%m-%dT%H:%M:%SZ')
END=$(date -u -d 'yesterday 23:59' '+%Y-%m-%dT%H:%M:%SZ')
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?start_date=$START&end_date=$END"
```

## Filtering Examples

### By Priority

```bash
# Only errors
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Priority=3"

# Errors and warnings
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Priority=3&Priority=4"

# Critical and above (0-2)
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Priority=0&Priority=1&Priority=2"
```

### By Host

```bash
# Single host
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?FromHost=webserver01"

# Multiple hosts (v0.2.2+)
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&FromHost=db01"
```

### By Message Content

```bash
# Search for keyword
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Message=login"

# Multiple keywords (OR logic)
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Message=error&Message=failed&Message=timeout"
```

### Combined Filters

```bash
# Errors from specific hosts
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=3"

# Complex query
START=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?FromHost=web01&Priority=3&Priority=4&start_date=$START&limit=100"
```

## Pagination

### Basic Pagination

```bash
# Page 1 (first 10)
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?limit=10&offset=0"

# Page 2 (next 10)
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?limit=10&offset=10"

# Page 3 (next 10)
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?limit=10&offset=20"
```

### Iterate Through All Results

```bash
#!/bin/bash
LIMIT=100
OFFSET=0

while true; do
    RESPONSE=$(curl -s -H "X-API-Key: $KEY" \
      "http://localhost:8000/logs?limit=$LIMIT&offset=$OFFSET")
    
    COUNT=$(echo "$RESPONSE" | jq '.rows | length')
    
    if [ "$COUNT" -eq 0 ]; then
        break
    fi
    
    echo "Processing $COUNT logs at offset $OFFSET"
    # Process logs here
    
    OFFSET=$((OFFSET + LIMIT))
done
```

## Metadata Queries

### List Available Values

```bash
# All hosts
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/FromHost"

# All priorities
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/Priority"

# All syslog tags
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/SysLogTag"
```

### Filtered Metadata

```bash
# Hosts that logged errors
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3"

# SysLogTags from web servers
curl -H "X-API-Key: $KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=web01&FromHost=web02"
```

## Practical Use Cases

### Monitor Specific Application

```bash
#!/bin/bash
# Monitor nginx errors

while true; do
    curl -s -H "X-API-Key: $KEY" \
      "http://localhost:8000/logs?SysLogTag=nginx&Priority=3&limit=5" \
      | jq -r '.rows[] | "\(.ReceivedAt) \(.FromHost) \(.Message)"'
    
    sleep 60
done
```

### Generate Error Report

```bash
#!/bin/bash
# Daily error report

DATE=$(date '+%Y-%m-%d')
START="${DATE}T00:00:00Z"
END="${DATE}T23:59:59Z"

echo "Error Report for $DATE"
echo "===================="
echo ""

# Get all errors
ERRORS=$(curl -s -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Priority=3&start_date=$START&end_date=$END&limit=1000")

# Count by host
echo "$ERRORS" | jq -r '.rows[] | .FromHost' | sort | uniq -c | sort -rn

echo ""
echo "Total errors: $(echo "$ERRORS" | jq '.total')"
```

### Simple Alerting

```bash
#!/bin/bash
# Alert on critical errors

CHECK_MINUTES=5
START=$(date -u -d "$CHECK_MINUTES minutes ago" '+%Y-%m-%dT%H:%M:%SZ')

CRITICAL=$(curl -s -H "X-API-Key: $KEY" \
  "http://localhost:8000/logs?Priority=2&start_date=$START" \
  | jq '.total')

if [ "$CRITICAL" -gt 0 ]; then
    echo "ALERT: $CRITICAL critical errors in last $CHECK_MINUTES minutes!"
    # Send notification (email, Slack, etc.)
fi
```

## Language Examples

### Python

```python
import requests
import json
from datetime import datetime, timedelta

API_KEY = "your-api-key"
BASE_URL = "http://localhost:8000"

headers = {"X-API-Key": API_KEY}

# Get latest logs
response = requests.get(
    f"{BASE_URL}/logs",
    headers=headers,
    params={"limit": 10}
)

logs = response.json()
print(f"Total logs: {logs['total']}")

for log in logs['rows']:
    print(f"{log['ReceivedAt']} [{log['Priority_Label']}] {log['Message']}")

# Multi-value filter (v0.2.2+)
response = requests.get(
    f"{BASE_URL}/logs",
    headers=headers,
    params={
        "FromHost": ["web01", "web02"],
        "Priority": [3, 4],
        "limit": 20
    }
)
```

### JavaScript/Node.js

```javascript
const API_KEY = 'your-api-key';
const BASE_URL = 'http://localhost:8000';

// Get latest logs
async function getLogs() {
    const response = await fetch(`${BASE_URL}/logs?limit=10`, {
        headers: {
            'X-API-Key': API_KEY
        }
    });
    
    const logs = await response.json();
    console.log(`Total logs: ${logs.total}`);
    
    logs.rows.forEach(log => {
        console.log(`${log.ReceivedAt} [${log.Priority_Label}] ${log.Message}`);
    });
}

// Multi-value filter
async function getErrorsFromMultipleHosts() {
    const params = new URLSearchParams();
    params.append('FromHost', 'web01');
    params.append('FromHost', 'web02');
    params.append('Priority', '3');
    params.append('Priority', '4');
    
    const response = await fetch(`${BASE_URL}/logs?${params}`, {
        headers: {
            'X-API-Key': API_KEY
        }
    });
    
    return await response.json();
}
```

## Tips & Tricks

### Performance

```bash
# Use smaller time windows
START=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')
# Better than querying last 30 days

# Limit results
?limit=100  # Don't retrieve more than needed

# Use indexed fields for filtering
# Fast: Priority, Facility, FromHost, ReceivedAt
# Slower: Message (full-text search)
```

### Output Formatting

```bash
# Pretty JSON
curl ... | jq

# Extract specific field
curl ... | jq -r '.rows[] | .Message'

# CSV format
curl ... | jq -r '.rows[] | [.ReceivedAt, .FromHost, .Message] | @csv'

# Table format
curl ... | jq -r '.rows[] | "\(.ReceivedAt)\t\(.FromHost)\t\(.Message)"' | column -t
```

### Debugging

```bash
# Verbose output
curl -v -H "X-API-Key: $KEY" "http://localhost:8000/logs"

# Check response headers
curl -I -H "X-API-Key: $KEY" "http://localhost:8000/logs"

# Time the request
time curl -H "X-API-Key: $KEY" "http://localhost:8000/logs?limit=1000"
```

## More Examples

Need help with a specific use case? Check:

- [Deployment Guide](../guides/deployment.md) - Production examples
- [Troubleshooting](../guides/troubleshooting.md) - Debug examples
- [GitHub Discussions](https://github.com/phil-bot/rsyslog-rest-api/discussions) - Community examples

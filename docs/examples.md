# Examples & Use Cases

[← Zurück zur Übersicht](index.md)

Praktische Beispiele für häufige Anwendungsfälle.

## 🎯 Grundlegende Abfragen

### Neueste Logs abrufen

```bash
# Letzte 10 Logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=10"

# Letzte 50 Logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=50"
```

### Logs aus bestimmtem Zeitraum

```bash
# Letzte Stunde
START=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')
END=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=$START&end_date=$END"

# Heute (00:00 bis jetzt)
START=$(date -u '+%Y-%m-%dT00:00:00Z')
END=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=$START&end_date=$END"

# Bestimmter Tag
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=2025-02-09T00:00:00Z&end_date=2025-02-09T23:59:59Z"

# Letzte 24 Stunden
START=$(date -u -d '24 hours ago' '+%Y-%m-%dT%H:%M:%SZ')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=$START&limit=100"
```

---

## 🔍 Nach Severity filtern

### Einzelne Priority

```bash
# Nur Errors (Priority 3)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=20"

# Nur Warnings (Priority 4)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=4&limit=20"

# Kritische Probleme (Priority 0-2)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=0&Priority=1&Priority=2&limit=50"
```

### Mehrere Priorities

```bash
# Errors UND Warnings
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&Priority=4&limit=50"

# Alle Probleme (Critical, Error, Warning)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=2&Priority=3&Priority=4&limit=100"

# Nur informative Logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=6&Priority=7&limit=50"
```

---

## 🖥️ Nach Host filtern

### Einzelner Host

```bash
# Logs von webserver01
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&limit=20"

# Logs von dbserver01
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=dbserver01&limit=20"
```

### Mehrere Hosts (NEU in v0.2.2!)

```bash
# Alle Webserver
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&FromHost=webserver02&FromHost=webserver03&limit=50"

# Web + Database Server
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&FromHost=dbserver01&limit=50"

# Alle Server einer Gruppe
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=app01&FromHost=app02&FromHost=app03&FromHost=app04&limit=100"
```

---

## 🏷️ Nach SysLogTag filtern

```bash
# Nur SSH-Logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?SysLogTag=sshd&limit=20"

# Nur nginx Logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?SysLogTag=nginx&limit=20"

# Mehrere Tags
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?SysLogTag=sshd&SysLogTag=sudo&SysLogTag=systemd&limit=50"
```

---

## 🔎 Text-Suche

### Einzelner Suchbegriff

```bash
# Suche nach "login"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=login&limit=20"

# Suche nach "error"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=error&limit=50"

# Suche nach "failed"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=failed&limit=30"
```

### Mehrere Suchbegriffe (OR-Logik!)

```bash
# "error" ODER "failed" ODER "timeout"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=error&Message=failed&Message=timeout&limit=50"

# "login" ODER "logout" ODER "authentication"
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=login&Message=logout&Message=authentication&limit=50"
```

---

## 🎨 Kombinierte Filter

### Host + Priority

```bash
# Errors von webserver01
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&Priority=3&limit=20"

# Warnings von mehreren Hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=4&limit=50"
```

### Host + Priority + Zeit

```bash
# Errors der letzten Stunde von dbserver01
START=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=dbserver01&Priority=3&start_date=$START&limit=50"

# Alle Probleme von mehreren Hosts heute
START=$(date -u '+%Y-%m-%dT00:00:00Z')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=2&Priority=3&Priority=4&start_date=$START&limit=100"
```

### Host + Message

```bash
# Login-Versuche auf webserver01
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01&Message=login&limit=30"

# Fehler-Logs mit bestimmten Keywords
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=appserver01&Message=error&Message=exception&Priority=3&limit=50"
```

### Komplexe Filter-Kombination

```bash
# Errors UND Warnings von mehreren Web-Servern mit "timeout" in Message
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&FromHost=web03&Priority=3&Priority=4&Message=timeout&limit=50"

# Kritische Probleme von DB-Servern heute
START=$(date -u '+%Y-%m-%dT00:00:00Z')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=db01&FromHost=db02&Priority=0&Priority=1&Priority=2&start_date=$START&limit=100"
```

---

## 📄 Pagination

### Grundlegendes Paging

```bash
# Erste 50 Einträge
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=50&offset=0"

# Zweite Seite (51-100)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=50&offset=50"

# Dritte Seite (101-150)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=50&offset=100"
```

### Pagination mit Filter

```bash
# Seite 1 von Errors
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=50&offset=0"

# Seite 2 von Errors
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=50&offset=50"
```

### Alle Daten iterieren (Bash-Script)

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
        echo "Fertig!"
        break
    fi
    
    echo "Verarbeite $COUNT Einträge (offset: $OFFSET)..."
    echo "$RESPONSE" | jq '.rows[] | .Message'
    
    OFFSET=$((OFFSET + LIMIT))
    sleep 1
done
```

---

## 📊 Metadaten abfragen

### Verfügbare Hosts

```bash
# Alle Hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost"

# Hosts die heute Errors hatten
START=$(date -u '+%Y-%m-%dT00:00:00Z')
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3&start_date=$START"

# Hosts mit Problemen (Priority 2-4)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=2&Priority=3&Priority=4"
```

### Verfügbare Tags

```bash
# Alle SysLogTags
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag"

# Tags von bestimmten Hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=webserver01&FromHost=webserver02"

# Tags von Error-Logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag?Priority=3"
```

### Verfügbare Priorities

```bash
# Alle verwendeten Priorities
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/Priority"

# Priorities von bestimmtem Host
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/Priority?FromHost=dbserver01"
```

---

## 🛠️ Praktische Use Cases

### Monitoring-Dashboard

```bash
#!/bin/bash
# dashboard.sh - Simple Monitoring Dashboard

API_KEY="your-api-key"
API_URL="http://localhost:8000"

echo "=== Syslog Dashboard ==="
echo ""

# Total Logs heute
START=$(date -u '+%Y-%m-%dT00:00:00Z')
TOTAL=$(curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?start_date=$START&limit=1" | jq .total)
echo "Total Logs heute: $TOTAL"

# Errors heute
ERRORS=$(curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?Priority=3&start_date=$START&limit=1" | jq .total)
echo "Errors heute: $ERRORS"

# Top 5 Hosts mit meisten Errors
echo ""
echo "Top 5 Hosts mit Errors:"
curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/meta/FromHost?Priority=3&start_date=$START" | jq -r '.[]' | head -5

# Letzte 5 Critical/Errors
echo ""
echo "Letzte 5 kritische Logs:"
curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?Priority=2&Priority=3&limit=5" | \
  jq -r '.rows[] | "\(.ReceivedAt) [\(.FromHost)] \(.Message)"'
```

### Login-Überwachung

```bash
#!/bin/bash
# login-monitor.sh - SSH Login Überwachung

API_KEY="your-api-key"
API_URL="http://localhost:8000"
START=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')

# Erfolgreiche Logins
echo "=== Erfolgreiche SSH Logins (letzte Stunde) ==="
curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?SysLogTag=sshd&Message=Accepted&start_date=$START&limit=50" | \
  jq -r '.rows[] | "\(.ReceivedAt) - \(.FromHost) - \(.Message)"'

echo ""

# Fehlgeschlagene Logins
echo "=== Fehlgeschlagene SSH Logins (letzte Stunde) ==="
curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?SysLogTag=sshd&Message=Failed&start_date=$START&limit=50" | \
  jq -r '.rows[] | "\(.ReceivedAt) - \(.FromHost) - \(.Message)"'
```

### Error-Report

```bash
#!/bin/bash
# error-report.sh - Täglicher Error Report

API_KEY="your-api-key"
API_URL="http://localhost:8000"
START=$(date -u '+%Y-%m-%dT00:00:00Z')
DATE=$(date '+%Y-%m-%d')

OUTPUT="error-report-$DATE.txt"

{
  echo "====================================="
  echo "Error Report für $DATE"
  echo "====================================="
  echo ""
  
  # Statistik
  TOTAL_ERRORS=$(curl -s -H "X-API-Key: $API_KEY" \
    "$API_URL/logs?Priority=3&start_date=$START&limit=1" | jq .total)
  
  TOTAL_WARNINGS=$(curl -s -H "X-API-Key: $API_KEY" \
    "$API_URL/logs?Priority=4&start_date=$START&limit=1" | jq .total)
  
  echo "Errors:   $TOTAL_ERRORS"
  echo "Warnings: $TOTAL_WARNINGS"
  echo ""
  
  # Hosts mit meisten Errors
  echo "Hosts mit meisten Errors:"
  curl -s -H "X-API-Key: $API_KEY" \
    "$API_URL/meta/FromHost?Priority=3&start_date=$START" | \
    jq -r '.[]' | nl
  
  echo ""
  echo "Top 20 Error Messages:"
  curl -s -H "X-API-Key: $API_KEY" \
    "$API_URL/logs?Priority=3&start_date=$START&limit=20" | \
    jq -r '.rows[] | "[\(.FromHost)] \(.Message)"'
  
} > "$OUTPUT"

echo "Report gespeichert: $OUTPUT"
```

### Alert-System (einfach)

```bash
#!/bin/bash
# alert-check.sh - Einfaches Alert-System
# Cron: */5 * * * * /path/to/alert-check.sh

API_KEY="your-api-key"
API_URL="http://localhost:8000"
START=$(date -u -d '5 minutes ago' '+%Y-%m-%dT%H:%M:%SZ')
ALERT_EMAIL="admin@example.com"

# Kritische Logs prüfen
CRITICAL=$(curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?Priority=0&Priority=1&Priority=2&start_date=$START&limit=1" | jq .total)

if [ "$CRITICAL" -gt 0 ]; then
  echo "ALERT: $CRITICAL kritische Logs in letzten 5 Minuten!" | \
    mail -s "Critical Syslog Alert" "$ALERT_EMAIL"
fi

# Zu viele Errors?
ERRORS=$(curl -s -H "X-API-Key: $API_KEY" \
  "$API_URL/logs?Priority=3&start_date=$START&limit=1" | jq .total)

if [ "$ERRORS" -gt 100 ]; then
  echo "WARNING: $ERRORS Errors in letzten 5 Minuten!" | \
    mail -s "High Error Rate Alert" "$ALERT_EMAIL"
fi
```

---

## 🐍 Python Beispiele

### Einfacher Client

```python
#!/usr/bin/env python3
import requests
import json
from datetime import datetime, timedelta

API_KEY = "your-api-key"
API_URL = "http://localhost:8000"

def get_logs(priority=None, from_host=None, limit=10):
    """Logs abrufen mit Filtern"""
    headers = {"X-API-Key": API_KEY}
    params = {"limit": limit}
    
    if priority:
        params["Priority"] = priority
    if from_host:
        params["FromHost"] = from_host
    
    response = requests.get(f"{API_URL}/logs", headers=headers, params=params)
    response.raise_for_status()
    return response.json()

# Beispiel-Nutzung
if __name__ == "__main__":
    # Letzte 10 Errors
    data = get_logs(priority=3, limit=10)
    
    print(f"Total Errors: {data['total']}")
    print(f"\nLetzte {len(data['rows'])} Errors:")
    
    for log in data['rows']:
        print(f"[{log['ReceivedAt']}] {log['FromHost']}: {log['Message']}")
```

### Monitoring mit Multi-Value

```python
#!/usr/bin/env python3
import requests
from datetime import datetime, timedelta

API_KEY = "your-api-key"
API_URL = "http://localhost:8000"

def get_errors_from_hosts(hosts, hours=1):
    """Errors von mehreren Hosts abrufen"""
    headers = {"X-API-Key": API_KEY}
    
    # Start-Zeit berechnen
    start = (datetime.utcnow() - timedelta(hours=hours)).isoformat() + 'Z'
    
    # Multi-Value Parameter
    params = [
        ("Priority", "3"),
        ("start_date", start),
        ("limit", "100")
    ]
    
    # Hosts hinzufügen (Multi-Value!)
    for host in hosts:
        params.append(("FromHost", host))
    
    response = requests.get(f"{API_URL}/logs", headers=headers, params=params)
    response.raise_for_status()
    return response.json()

# Beispiel
hosts = ["webserver01", "webserver02", "webserver03"]
data = get_errors_from_hosts(hosts, hours=1)

print(f"Errors von {len(hosts)} Hosts: {data['total']}")
for log in data['rows']:
    print(f"{log['FromHost']}: {log['Message']}")
```

---

## 🌐 JavaScript/Node.js Beispiele

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

// Verwendung
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

// Verwendung
const hosts = ['web01', 'web02', 'db01'];
const priorities = [3, 4]; // Errors & Warnings

getLogsFromMultipleHosts(hosts, priorities)
  .then(data => console.log(data))
  .catch(console.error);
```

---

## 💡 Tipps & Best Practices

### Performance

1. **Immer limit verwenden** - Standard ist nur 10, max 1000
2. **Zeitfenster einschränken** - Kürzere Zeiträume = schneller
3. **Pagination nutzen** - Nicht alle Daten auf einmal
4. **Indexierte Felder filtern** - `FromHost`, `Priority`, `ReceivedAt`

### Fehlerbehandlung

```bash
# HTTP Status prüfen
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
# Requests loggen
LOG_FILE="api-requests.log"

curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8000/logs?Priority=3&limit=10" | \
  tee -a "$LOG_FILE" | jq
```

---

[← Zurück zur Übersicht](index.md) | [Weiter zu Deployment →](deployment.md)

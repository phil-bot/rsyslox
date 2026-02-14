# API Reference

[← Zurück zur Übersicht](index.md)

Vollständige API-Dokumentation für alle Endpunkte.

## 🔐 Authentifizierung

### API Key Authentication

Alle geschützten Endpunkte erfordern einen API-Key im Header:

```http
X-API-Key: your-api-key-here
```

**Beispiel:**
```bash
curl -H "X-API-Key: a3d7f8c9e2b4a6d8..." "http://localhost:8000/logs"
```

**Hinweis:** Wenn `API_KEY` in `.env` leer ist, ist **keine** Authentifizierung erforderlich (nur für Development!).

---

## 📍 Endpunkte Übersicht

| Endpunkt | Methode | Auth | Beschreibung |
|----------|---------|------|--------------|
| `/health` | GET | ❌ | Health Check |
| `/logs` | GET | ✅ | Logs mit Filterung und Pagination |
| `/meta` | GET | ✅ | Verfügbare Spalten auflisten |
| `/meta/{column}` | GET | ✅ | Eindeutige Werte einer Spalte |

---

## GET /health

Health Check Endpunkt ohne Authentifizierung.

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

### Status Codes

| Code | Bedeutung |
|------|-----------|
| 200 | API und Datenbank funktionieren |
| 503 | Datenbank nicht erreichbar |

---

## GET /logs

Ruft Log-Einträge mit Filterung und Pagination ab.

### Request

```http
GET /logs?limit=10&Priority=3 HTTP/1.1
Host: localhost:8000
X-API-Key: your-api-key
```

### Query Parameter

#### Pagination

| Parameter | Type | Default | Beschreibung |
|-----------|------|---------|--------------|
| `offset` | Integer | 0 | Startposition (Überspringen von Einträgen) |
| `limit` | Integer | 10 | Maximale Anzahl Ergebnisse (max: 1000) |

#### Zeitfilter

| Parameter | Type | Default | Format | Beschreibung |
|-----------|------|---------|--------|--------------|
| `start_date` | DateTime | -24h | ISO 8601 | Startdatum/Zeit |
| `end_date` | DateTime | now | ISO 8601 | Enddatum/Zeit |

**ISO 8601 Format:** `2025-02-09T10:30:00Z` oder `2025-02-09T10:30:00+01:00`

**Max. Zeitspanne:** 90 Tage

#### Content-Filter (Multi-Value!)

Alle Filter unterstützen **mehrere Werte** durch Wiederholung des Parameters:

| Parameter | Type | Multi | Beschreibung |
|-----------|------|-------|--------------|
| `FromHost` | String | ✅ | Hostname(s) filtern |
| `Priority` | Integer | ✅ | Severity filtern (0-7) |
| `Facility` | Integer | ✅ | Facility filtern (0-23) |
| `Message` | String | ✅ | Text-Suche (OR-Logik) |
| `SysLogTag` | String | ✅ | Syslog-Tag filtern |

**Multi-Value Syntax:**
```bash
# Mehrere Werte = Parameter wiederholen
?FromHost=web01&FromHost=web02&FromHost=db01

# NICHT: Komma-getrennt (funktioniert NICHT!)
?FromHost=web01,web02,db01  # ❌ FALSCH
```

### Priority Values (RFC-5424)

| Value | Label | Beschreibung |
|-------|-------|--------------|
| 0 | Emergency | System unbrauchbar |
| 1 | Alert | Sofortige Maßnahmen erforderlich |
| 2 | Critical | Kritischer Zustand |
| 3 | Error | Fehlerbedingungen |
| 4 | Warning | Warnungen |
| 5 | Notice | Normal aber signifikant |
| 6 | Informational | Informationsmeldungen |
| 7 | Debug | Debug-Nachrichten |

### Facility Values (RFC-5424)

| Value | Label | Beschreibung |
|-------|-------|--------------|
| 0 | kern | Kernel-Nachrichten |
| 1 | user | User-Level-Nachrichten |
| 2 | mail | Mail-System |
| 3 | daemon | System-Daemons |
| 4 | auth | Sicherheit/Autorisierung |
| 5 | syslog | Syslog intern |
| 16-23 | local0-7 | Lokale Verwendung |

[Vollständige Liste: RFC-5424](https://tools.ietf.org/html/rfc5424)

### Request Beispiele

#### Einfache Abfrage

```bash
# Neueste 10 Logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=10"
```

#### Zeitfilter

```bash
# Logs der letzten Stunde
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=2025-02-09T09:00:00Z&end_date=2025-02-09T10:00:00Z"

# Logs von gestern
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=2025-02-08T00:00:00Z&end_date=2025-02-08T23:59:59Z"
```

#### Single-Value Filter

```bash
# Nur Errors (Priority 3)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3"

# Von einem bestimmten Host
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=webserver01"

# Text-Suche
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=login"
```

#### Multi-Value Filter (NEU in v0.2.2!)

```bash
# Mehrere Hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&FromHost=db01"

# Errors UND Warnings
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&Priority=4"

# Mehrere Facilities
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Facility=1&Facility=4&Facility=16"

# Mehrere Such-Begriffe (OR-Logik)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Message=error&Message=failed&Message=timeout"
```

#### Kombinierte Filter

```bash
# Errors von mehreren Hosts in letzter Stunde
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=3&start_date=2025-02-09T09:00:00Z&limit=20"

# Alle Priorities von spezifischem Host
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=dbserver01&Priority=2&Priority=3&Priority=4"
```

#### Pagination

```bash
# Erste 10 Einträge
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=10&offset=0"

# Nächste 10 Einträge
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

### Response Felder

#### Pflicht-Felder (immer vorhanden)

| Feld | Type | Beschreibung |
|------|------|--------------|
| `ID` | Integer | Log-Eintrags-ID |
| `ReceivedAt` | DateTime | Empfangszeit bei rsyslog |
| `FromHost` | String | Quell-Hostname |
| `Priority` | Integer | Severity (0-7) |
| `Priority_Label` | String | RFC-Label (z.B. "Error") |
| `Facility` | Integer | Facility (0-23) |
| `Facility_Label` | String | RFC-Label (z.B. "user") |
| `Message` | String | Log-Nachricht |

#### Erweiterte Felder (optional, wenn vorhanden)

| Feld | Type | Beschreibung |
|------|------|--------------|
| `CustomerID` | Integer | Kunden-ID |
| `DeviceReportedTime` | DateTime | Original-Zeitstempel vom Gerät |
| `SysLogTag` | String | Syslog-Tag/Programmname |
| `EventSource` | String | Event-Quelle |
| `EventUser` | String | Zugehöriger Benutzer |
| `EventID` | Integer | Event-ID |
| `EventCategory` | Integer | Event-Kategorie |
| `NTSeverity` | Integer | Windows NT Severity |
| `Importance` | Integer | Wichtigkeits-Rating (1-5) |
| `EventBinaryData` | String | Binäre Event-Daten |
| `MaxAvailable` | Integer | Max. verfügbare Ressourcen |
| `CurrUsage` | Integer | Aktuelle Ressourcen-Nutzung |
| `MinUsage` | Integer | Minimale Nutzung |
| `MaxUsage` | Integer | Maximale Nutzung |
| `InfoUnitID` | Integer | Info-Unit-ID |
| `EventLogType` | String | Event-Log-Typ |
| `GenericFileName` | String | Zugehöriger Dateiname |
| `SystemID` | Integer | System-ID |

**Hinweis:** Erweiterte Felder verwenden `omitempty` - sie erscheinen nur, wenn die Datenbank einen Wert hat (nicht NULL).

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

Listet alle verfügbaren Spalten für Filterung auf.

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

Ruft eindeutige Werte einer bestimmten Spalte ab.

### Request

```http
GET /meta/FromHost HTTP/1.1
Host: localhost:8000
X-API-Key: your-api-key
```

### Path Parameter

| Parameter | Type | Beschreibung |
|-----------|------|--------------|
| `column` | String | Spaltenname (aus `/meta`) |

### Query Parameter

**Alle Filter von `/logs` werden unterstützt!** (Multi-Value auch!)

Dies ermöglicht gefilterte Meta-Abfragen:

```bash
# Alle Hosts die Errors hatten
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3&Priority=4"

# Alle SysLogTags von bestimmten Hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=web01&FromHost=web02"
```

### Request Beispiele

#### Einfache Meta-Abfragen

```bash
# Alle Hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost"

# Alle verwendeten Priorities
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/Priority"

# Alle SysLogTags
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag"

# Alle Event-Quellen
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/EventSource"
```

#### Gefilterte Meta-Abfragen

```bash
# Hosts die Errors hatten
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3"

# SysLogTags von spezifischen Hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/SysLogTag?FromHost=webserver01&FromHost=webserver02"

# Priorities in letzter Stunde
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/Priority?start_date=2025-02-09T09:00:00Z"
```

### Response

#### Für Priority/Facility (mit Labels)

```json
[
  { "val": 0, "label": "Emergency" },
  { "val": 1, "label": "Alert" },
  { "val": 3, "label": "Error" },
  { "val": 6, "label": "Informational" }
]
```

#### Für Integer-Spalten (IDs)

```json
[1, 2, 5, 10, 42, 100]
```

#### Für String-Spalten

```json
[
  "webserver01",
  "webserver02",
  "dbserver01",
  "appserver01"
]
```

### Error Response

**400 Bad Request (ungültige Spalte):**

```json
{
  "error": "Invalid column 'InvalidCol'. Available columns: ID, CustomerID, ReceivedAt, ..."
}
```

---

## 🔢 HTTP Status Codes

| Code | Bedeutung | Verwendung |
|------|-----------|------------|
| 200 | OK | Erfolgreiche Anfrage |
| 400 | Bad Request | Ungültige Parameter |
| 401 | Unauthorized | API-Key fehlt oder ungültig |
| 500 | Internal Server Error | Server-/Datenbank-Fehler |
| 503 | Service Unavailable | Datenbank nicht erreichbar |

---

## 📊 Rate Limiting

Aktuell **kein** Rate Limiting implementiert.

**Empfehlung:** Verwende einen Reverse Proxy (nginx/Apache) für Rate Limiting in Production.

→ [Deployment: Reverse Proxy](deployment.md#reverse-proxy)

---

## 🔒 CORS

CORS wird über `ALLOWED_ORIGINS` in `.env` konfiguriert:

```bash
# Development (alle Origins)
ALLOWED_ORIGINS=*

# Production (spezifische Domains)
ALLOWED_ORIGINS=https://dashboard.example.com,https://app.example.com
```

→ [Configuration: CORS](configuration.md#cors-configuration)

---

## 💡 Best Practices

### Performance

1. **Limit verwenden:** Immer `limit` setzen (Default: 10, Max: 1000)
2. **Zeitfenster begrenzen:** Kleinere Zeiträume = schnellere Queries
3. **Pagination nutzen:** Große Ergebnisse in Chunks abrufen
4. **Indexierte Felder filtern:** `Priority`, `Facility`, `FromHost`, `ReceivedAt`

### Sicherheit

1. **API-Key rotieren:** Regelmäßig neuen Key generieren
2. **HTTPS verwenden:** In Production immer SSL/TLS
3. **CORS einschränken:** Nur notwendige Origins erlauben
4. **Rate Limiting:** Über Reverse Proxy

### Fehlerbehandlung

1. **HTTP-Status prüfen:** Nicht nur 200 annehmen
2. **Retry-Logic:** Bei 500/503 mit Backoff
3. **Timeout setzen:** Client-seitig Timeout konfigurieren

---

## 📖 Weitere Ressourcen

- [Installation](installation.md) - Setup und Deployment
- [Configuration](configuration.md) - Konfiguration
- [Examples](examples.md) - Praktische Beispiele
- [Troubleshooting](troubleshooting.md) - Fehlersuche

---

[← Zurück zur Übersicht](index.md) | [Weiter zu Examples →](examples.md)

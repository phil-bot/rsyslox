# Changelog

[← Zurück zur Übersicht](index.md)

Alle wichtigen Änderungen am Projekt werden in dieser Datei dokumentiert.

Das Format basiert auf [Keep a Changelog](https://keepachangelog.com/de/1.0.0/),
und dieses Projekt folgt [Semantic Versioning](https://semver.org/lang/de/).

---

## [Unreleased]

### Geplant für v0.3.0
- Negation filters (`exclude`, `not`)
- Complex filter combinations
- Unit Tests
- GitHub Actions CI/CD

### Geplant für v0.4.0
- Statistics endpoint (`/stats`)
- Aggregation support
- Timeline/Histogram features
- WebSocket support für Live-Logs

---

## [0.2.2] - 2025-02-09

### ✨ Added
- **Multi-Value Filter Support** - Alle Filter-Parameter akzeptieren nun mehrere Werte
  ```bash
  # NEU: Mehrere Hosts
  ?FromHost=web01&FromHost=web02&FromHost=web03
  
  # NEU: Mehrere Priorities
  ?Priority=3&Priority=4
  ```
  
- **Extended Columns** - Alle 25+ SystemEvents-Spalten in Response
  - CustomerID, DeviceReportedTime, SysLogTag
  - EventSource, EventUser, EventID, EventCategory
  - NTSeverity, Importance, SystemID
  - Und viele mehr (siehe API Reference)
  
- **Enhanced Meta Endpoint** - `/meta/{column}` unterstützt nun Filter
  ```bash
  # Hosts die Errors hatten
  /meta/FromHost?Priority=3&Priority=4
  ```

- **Live Log Generator** (Docker) - Kontinuierliche Test-Daten
  - 3 neue Logs alle 10 Sekunden
  - Realistische Messages, Event IDs, Extended Fields
  - 6 verschiedene Hosts, 8 verschiedene Tags

### 🔄 Changed
- Response-Struktur erweitert - Extended Fields mit `omitempty`
- Meta-Endpoint gibt nun strukturierte Antworten für Priority/Facility

### 🐛 Fixed
- Filter-Validierung verbessert
- NULL-Handling für Extended Columns

### 📚 Documentation
- Neue Dokumentations-Struktur
- Erweiterte API Reference
- Docker Guide aktualisiert
- Neue Beispiele für Multi-Value Filter

---

## [0.2.1] - 2025-01-15

### 🐛 Fixed
- Database connection pooling optimiert
- Memory leak in long-running instances behoben

### 🔄 Changed
- Default `limit` von 100 auf 10 reduziert
- Logging verbessert

---

## [0.2.0] - 2024-12-20

### ✨ Added
- **RFC-5424 Labels** - Priority und Facility mit korrekten Labels
  ```json
  {
    "Priority": 3,
    "Priority_Label": "Error",
    "Facility": 1,
    "Facility_Label": "user"
  }
  ```

- **Meta Endpoint** - `/meta` und `/meta/{column}`
  ```bash
  # Verfügbare Spalten
  GET /meta
  
  # Eindeutige Werte
  GET /meta/FromHost
  GET /meta/Priority
  ```

- **CORS Support** - Konfigurierbar via `ALLOWED_ORIGINS`

- **SSL/TLS Support** - Optional aktivierbar
  ```bash
  USE_SSL=true
  SSL_CERTFILE=/path/to/cert.pem
  SSL_KEYFILE=/path/to/key.pem
  ```

### 🔄 Changed
- API-Key Authentication optional (entwickelt für Production)
- Database Configuration via Environment Variables (empfohlen)
- systemd Service mit Security Hardening

### 📚 Documentation
- README komplett überarbeitet
- API Reference erstellt
- Docker Guide hinzugefügt

---

## [0.1.5] - 2024-11-10

### 🐛 Fixed
- Pagination Edge Cases
- Date Range Validation
- SQL Injection Prevention verbessert

---

## [0.1.0] - 2024-10-01

### ✨ Initial Release

**Features:**
- REST API für rsyslog/MySQL
- `/health` Endpoint
- `/logs` Endpoint mit Filtering
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
- Makefile für Build
- Static Binary Support

---

## Migration Guides

### Upgrade von v0.2.1 auf v0.2.2

**Keine Breaking Changes!** 

✅ **Kompatibel:**
- Single-Value Filter funktionieren weiterhin
- Alle bestehenden Endpoints unverändert
- Response-Format erweitert (aber backward-compatible)

✅ **Neu verfügbar:**
- Multi-Value Filter (optional nutzen)
- Extended Columns (erscheinen automatisch wenn vorhanden)
- Meta-Filtering

**Update Steps:**
```bash
# 1. Binary aktualisieren
sudo systemctl stop rsyslog-rest-api
sudo cp rsyslog-rest-api-v0.2.2 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo systemctl start rsyslog-rest-api

# 2. Testen
curl http://localhost:8000/health
curl -H "X-API-Key: KEY" "http://localhost:8000/logs?limit=1"

# Fertig! Keine Config-Änderungen nötig.
```

---

### Upgrade von v0.1.x auf v0.2.0

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
- Clients müssen neue Felder ignorieren können (sollte kein Problem sein)
- Keine Code-Änderung nötig wenn nur `Priority`/`Facility` numerisch genutzt

**Update Steps:**
```bash
# 1. .env anpassen (neue Variablen)
sudo nano /opt/rsyslog-rest-api/.env

# Neu hinzufügen:
ALLOWED_ORIGINS=*  # oder spezifische Domains

# Optional:
USE_SSL=false
SSL_CERTFILE=
SSL_KEYFILE=

# 2. Binary aktualisieren
sudo systemctl stop rsyslog-rest-api
sudo cp rsyslog-rest-api-v0.2.0 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo systemctl start rsyslog-rest-api

# 3. Testen
curl http://localhost:8000/health
```

---

## Deprecated Features

### v0.2.0
- **rsyslog Config File Parsing** (Fallback, noch unterstützt)
  - Bitte nutze `DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASS` in `.env`
  - `RSYSLOG_CONFIG_PATH` wird in Zukunft entfernt

---

## Security Fixes

### v0.2.1
- Fixed potential SQL injection in custom queries
- Improved input validation

### v0.2.0
- API-Key Authentication hinzugefügt
- SSL/TLS Support

---

## Performance Improvements

### v0.2.2
- Database Connection Pooling optimiert
- Query Performance für Multi-Value Filter

### v0.2.0
- Automatic Indexing bei Startup
- Connection Pooling implementiert

---

## Known Issues

### v0.2.2
- Keine bekannten Issues

### v0.2.1
- Keine bekannten Issues

### v0.2.0
- ~~Memory leak bei lang laufenden Instanzen~~ (Fixed in v0.2.1)

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

[← Zurück zur Übersicht](index.md)

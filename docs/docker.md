# Docker Test Environment

[← Zurück zur Übersicht](index.md)

Komplette Docker-Testumgebung mit live generierten Logs.

## 🐳 Über die Testumgebung

Die Docker-Umgebung bietet ein voll funktionsfähiges Setup zum Testen und Entwickeln:

✅ **Ubuntu 24.04** Container  
✅ **rsyslog + MariaDB** vorinstalliert  
✅ **10 initiale Test-Logs** bei Start  
✅ **Live Log Generator** - 3 neue Logs alle 10 Sekunden  
✅ **Alle 25+ Felder** mit realistischen Daten gefüllt  
✅ **Keine API-Key-Authentifizierung** (optional aktivierbar)  

**Perfekt für:**
- API Testing
- Feature-Entwicklung
- Demos
- Multi-Value-Filter testen
- Extended Columns testen

---

## 🚀 Quick Start

### Voraussetzungen

- Docker & Docker Compose installiert
- Go 1.21+ (für Build)
- Make

### Setup

```bash
# 1. Binary bauen (auf Host-System)
make build-static

# 2. Docker Container starten
cd docker
docker-compose up -d

# 3. Warten bis bereit (ca. 30 Sekunden)
docker-compose logs -f
# Warte auf: "✓ Environment Ready!"

# 4. Testen
curl http://localhost:8000/health
curl "http://localhost:8000/logs?limit=5"
```

**Fertig!** Die API läuft nun auf `http://localhost:8000`.

---

## 📦 Was ist enthalten?

### Container-Setup

```
rsyslog-rest-api-test/
├── Ubuntu 24.04
├── MariaDB Server
├── rsyslog + rsyslog-mysql
├── rsyslog-rest-api (dein Binary)
├── log-generator.sh (Live-Daten)
└── test-v0.2.2.sh (Test-Suite)
```

### Initiale Daten

Bei Container-Start werden **10 Test-Logs** erstellt:

```sql
INSERT INTO SystemEvents VALUES
  ('webserver01', Priority=6, 'User login successful'),
  ('webserver01', Priority=3, 'Failed login attempt'),
  ('dbserver01', Priority=4, 'Database connection timeout'),
  ('dbserver01', Priority=6, 'Query executed successfully'),
  ('appserver01', Priority=5, 'Application started'),
  ('appserver01', Priority=3, 'Critical error in module'),
  ('webserver02', Priority=6, 'HTTP request: GET /api/users'),
  ('webserver02', Priority=4, 'Slow response time detected'),
  ('mailserver01', Priority=2, 'Mail queue growing rapidly'),
  ('mailserver01', Priority=6, 'Email sent successfully')
```

### Live Log Generator 🔥

**Alle 10 Sekunden** werden **3 neue realistische Logs** generiert:

**Hosts:**
- webserver01, webserver02
- dbserver01
- appserver01
- mailserver01
- firewall01

**SysLogTags:**
- sshd, nginx, mysqld, node, postfix, iptables, systemd, docker

**Priorities (gewichtet):**
- 6 (Info) - 50%
- 5 (Notice) - 25%
- 4 (Warning) - 12.5%
- 3 (Error) - 10%
- 2 (Critical) - 2.5%

**Extended Fields:**
- Realistische Event IDs (z.B. 4624 für SSH Login, 200-500 für HTTP)
- EventSource, EventUser, NTSeverity
- Importance, EventCategory, SystemID
- Alle 25+ Spalten gefüllt!

**Beispiel-Logs:**
```
[2025-02-09 10:30:15] [webserver01] [INFO] [sshd] User login successful
[2025-02-09 10:30:16] [dbserver01] [WARNING] [mysqld] Slow query detected: 2500ms
[2025-02-09 10:30:17] [appserver01] [ERROR] [node] Connection refused
```

---

## 🧪 Testing

### Manuelle Tests

```bash
# Health Check
curl http://localhost:8000/health

# Neueste 5 Logs
curl "http://localhost:8000/logs?limit=5"

# Nur Errors
curl "http://localhost:8000/logs?Priority=3&limit=10"

# Multi-Value: Mehrere Hosts
curl "http://localhost:8000/logs?FromHost=webserver01&FromHost=webserver02"

# Extended Columns sehen
curl "http://localhost:8000/logs?limit=1" | jq .rows[0]

# Meta: Alle Hosts
curl "http://localhost:8000/meta/FromHost"

# Meta: Hosts mit Errors
curl "http://localhost:8000/meta/FromHost?Priority=3&Priority=4"
```

### Test-Suite ausführen

Die Docker-Umgebung enthält eine umfassende Test-Suite:

```bash
# In docker/ Verzeichnis
./test-v0.2.2.sh
```

**Tests:**
- ✅ Health Check
- ✅ Basic Log Retrieval
- ✅ Multi-Value Filters
- ✅ Extended Columns
- ✅ Meta Endpoints
- ✅ Backward Compatibility

**Erwartete Ausgabe:**
```
==========================================
rsyslog REST API v0.2.2 - Test Suite
==========================================
...
Waiting for API... ready

Running tests...

[1] Health Check... ✓ OK (HTTP 200)
[2] Get Logs (default)... ✓ OK (HTTP 200)
[3] Multi-value FromHost... ✓ OK (HTTP 200)
...
==========================================
Test Summary
==========================================
Passed: 22
Failed: 0
Total:  22

✓ All tests passed!
```

### Live-Daten beobachten

```bash
# Log-Generator-Output ansehen
docker exec rsyslog-rest-api-test tail -f /var/log/log-generator.log

# Datenbank-Count live beobachten
watch -n 5 'docker exec rsyslog-rest-api-test mysql -N Syslog -e "SELECT COUNT(*) FROM SystemEvents"'

# Container Logs
docker-compose logs -f
```

---

## ⚙️ Konfiguration

### API-Key aktivieren (Optional)

Standardmäßig läuft die API **ohne** Authentifizierung. Zum Aktivieren:

**1. docker-compose.yml bearbeiten:**

```yaml
# docker/docker-compose.yml
environment:
  - SERVER_PORT=8000
  - ALLOWED_ORIGINS=*
  - API_KEY=test123456789  # <-- Aktivieren
```

**2. Container neu starten:**

```bash
docker-compose down
docker-compose up -d
```

**3. Mit API-Key testen:**

```bash
curl -H "X-API-Key: test123456789" "http://localhost:8000/logs?limit=5"
```

### Port ändern

**docker-compose.yml:**

```yaml
ports:
  - "8080:8000"  # <-- Host:Container
```

```bash
docker-compose down
docker-compose up -d
```

API dann auf: `http://localhost:8080`

### Log-Generator anpassen

**docker/log-generator.sh bearbeiten:**

```bash
# Configuration
INTERVAL=10        # Sekunden zwischen Bursts (default: 10)
LOGS_PER_BURST=3   # Logs pro Burst (default: 3)
```

**Container neu bauen:**

```bash
docker-compose down
docker-compose up -d --build
```

---

## 📊 Monitoring

### Container Status

```bash
# Status
docker-compose ps

# Logs live
docker-compose logs -f

# API Logs
docker exec rsyslog-rest-api-test tail -f /var/log/rsyslog-rest-api.log

# Generator Logs
docker exec rsyslog-rest-api-test tail -f /var/log/log-generator.log
```

### Datenbank

```bash
# MySQL Shell
docker exec -it rsyslog-rest-api-test mysql -u rsyslog -ppassword Syslog

# Count Logs
docker exec rsyslog-rest-api-test mysql -N Syslog -e "SELECT COUNT(*) FROM SystemEvents"

# Neueste Logs
docker exec rsyslog-rest-api-test mysql Syslog -e "
  SELECT ReceivedAt, FromHost, Priority, Message 
  FROM SystemEvents 
  ORDER BY ReceivedAt DESC 
  LIMIT 5
"

# Statistics
docker exec rsyslog-rest-api-test mysql Syslog -e "
  SELECT 
    FromHost, 
    COUNT(*) as count,
    AVG(Priority) as avg_priority
  FROM SystemEvents 
  GROUP BY FromHost
"
```

### Prozesse im Container

```bash
# Laufende Prozesse
docker exec rsyslog-rest-api-test ps aux

# API läuft?
docker exec rsyslog-rest-api-test ps aux | grep rsyslog-rest-api

# Generator läuft?
docker exec rsyslog-rest-api-test ps aux | grep log-generator
```

---

## 🔧 Erweiterte Nutzung

### Container Shell

```bash
# Bash Shell im Container
docker exec -it rsyslog-rest-api-test bash

# Im Container:
cd /opt/rsyslog-rest-api
cat .env
./rsyslog-rest-api --help
```

### Binary neu deployen

Nach Code-Änderungen:

```bash
# 1. Neu bauen
cd .. && make build-static

# 2. Container neustarten (kopiert automatisch)
cd docker
docker-compose restart
```

### Datenbank zurücksetzen

```bash
# Im Container
docker exec -it rsyslog-rest-api-test mysql Syslog -e "TRUNCATE TABLE SystemEvents"

# Initial-Daten neu laden
docker-compose restart
```

---

## 🛠️ Troubleshooting

### Container startet nicht

```bash
# Logs ansehen
docker-compose logs

# Binary fehlt?
ls -la ../build/rsyslog-rest-api
# Falls nicht vorhanden: make build-static

# Neu bauen
docker-compose down
docker-compose up -d --build
```

### API antwortet nicht

```bash
# API läuft?
docker exec rsyslog-rest-api-test ps aux | grep rsyslog-rest-api

# Logs
docker exec rsyslog-rest-api-test cat /var/log/rsyslog-rest-api.log

# Port richtig?
docker-compose ps
# Sollte zeigen: 0.0.0.0:8000->8000/tcp

# Von Host testen
curl -v http://localhost:8000/health
```

### Log-Generator läuft nicht

```bash
# Prozess prüfen
docker exec rsyslog-rest-api-test ps aux | grep log-generator

# Logs ansehen
docker exec rsyslog-rest-api-test cat /var/log/log-generator.log

# Manuell starten
docker exec -it rsyslog-rest-api-test bash
/opt/rsyslog-rest-api/log-generator.sh
```

### Datenbank-Probleme

```bash
# MySQL läuft?
docker exec rsyslog-rest-api-test systemctl status mysql

# Connection testen
docker exec rsyslog-rest-api-test mysql -u rsyslog -ppassword Syslog -e "SELECT 1"

# Tabelle existiert?
docker exec rsyslog-rest-api-test mysql Syslog -e "SHOW TABLES"
```

---

## 🧹 Cleanup

### Container stoppen

```bash
# Stoppen (Daten bleiben)
docker-compose stop

# Starten
docker-compose start
```

### Komplett entfernen

```bash
# Container und Netzwerk entfernen
docker-compose down

# Mit Volumes (DATEN LÖSCHEN!)
docker-compose down -v
```

### Disk Space freigeben

```bash
# Ungenutzte Container/Images/Volumes
docker system prune

# Alle dangling images
docker image prune
```

---

## 📚 Dateien

### Struktur

```
docker/
├── Dockerfile              # Container-Image
├── docker-compose.yml      # Compose-Konfiguration
├── entrypoint.sh           # Startup-Script
├── log-generator.sh        # Live-Log-Generator
├── test.sh                 # Basic Test-Suite
└── test-v0.2.2.sh          # Extended Test-Suite
```

### Wichtige Pfade im Container

```
/opt/rsyslog-rest-api/          # API Installation
  ├── rsyslog-rest-api          # Binary (von Host)
  ├── .env                      # Config
  └── log-generator.sh          # Generator

/etc/rsyslog.d/mysql.conf       # rsyslog Config

/var/log/
  ├── rsyslog-rest-api.log      # API Logs
  └── log-generator.log         # Generator Logs
```

---

## 🎯 Use Cases

### Feature-Entwicklung

```bash
# 1. Code ändern
vim main.go

# 2. Neu bauen
make build-static

# 3. Container neustarten
cd docker && docker-compose restart

# 4. Testen
curl "http://localhost:8000/new-endpoint"
```

### Demo vorbereiten

```bash
# 1. Container starten
docker-compose up -d

# 2. Warten bis Daten generiert wurden (2-3 Minuten)
sleep 180

# 3. Demo
curl "http://localhost:8000/logs?Priority=3&Priority=4&limit=10" | jq
```

### Performance Testing

```bash
# Load Testing (mit ApacheBench)
ab -n 1000 -c 10 "http://localhost:8000/logs?limit=100"

# Mit Filtern
ab -n 500 -c 5 "http://localhost:8000/logs?FromHost=webserver01&Priority=3&limit=50"
```

---

## 🆚 Docker vs. Production

### Unterschiede

| Feature | Docker | Production |
|---------|--------|------------|
| API Key | Optional | **Erforderlich** |
| SSL/TLS | Nein | **Ja** |
| Database | In Container | Extern |
| Persistence | Flüchtig | Persistent |
| Performance | Begrenzt | Optimiert |

### Nicht für Production!

Docker-Setup ist **nur für Testing/Development**!

Für Production siehe:
- [Deployment Guide](deployment.md)
- [Security Guide](security.md)

---

## 💡 Tipps

1. **Lange laufen lassen** - Nach 1 Stunde hast du ~1000 Logs zum Testen
2. **Multi-Value testen** - Perfekt für neue Filter-Features
3. **Extended Columns** - Alle Felder haben realistische Werte
4. **Performance** - Test mit verschiedenen `limit` Werten
5. **Pagination** - Test mit `offset` und `limit`

---

[← Zurück zur Übersicht](index.md) | [Weiter zu Development →](development.md)

# Development Guide

[← Zurück zur Übersicht](index.md)

Entwicklungs-Guide und Contributing-Richtlinien.

## 🏗️ Projekt-Architektur

### Überblick

```
rsyslog-rest-api
├── main.go                  # Haupt-Applikation (Single File!)
├── go.mod                   # Go Dependencies
├── go.sum                   # Dependency Checksums
├── Makefile                 # Build Automation
├── .env.example             # Config Template
├── rsyslog-rest-api.service # systemd Service
├── docker/                  # Docker Testumgebung
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── entrypoint.sh
│   ├── log-generator.sh
│   └── test-v0.2.2.sh
└── docs/                    # Dokumentation
    └── ...
```

### Technologie-Stack

- **Sprache:** Go 1.21+
- **Database Driver:** go-sql-driver/mysql v1.7.1
- **Datenbank:** MySQL 5.7+ / MariaDB 10.3+
- **Testing:** Docker + Bash Scripts

---

## 🚀 Setup Development Environment

### Voraussetzungen

```bash
# Go installieren (1.21+)
# Ubuntu/Debian
sudo apt-get install golang-1.21

# Oder: https://go.dev/dl/

# Verifizieren
go version
# go version go1.21.x linux/amd64

# Git
sudo apt-get install git

# Make
sudo apt-get install make

# Docker (für Testing)
sudo apt-get install docker.io docker-compose
```

### Repository klonen

```bash
git clone https://github.com/phil-bot/rsyslog-rest-api.git
cd rsyslog-rest-api
```

### Dependencies installieren

```bash
# Go Modules Download
go mod download

# Verifizieren
go mod verify
```

---

## 🔨 Build

### Development Build

```bash
# Standard Build
make build

# Binary ist in: ./build/rsyslog-rest-api
./build/rsyslog-rest-api
```

### Static Build (Production)

```bash
# Static Binary (keine libc dependency)
make build-static

# Verifizieren (keine Dependencies)
ldd ./build/rsyslog-rest-api
# not a dynamic executable
```

### Mit Version

```bash
# Version aus Git Tag
VERSION=v0.2.2 make build

# Oder manuell
go build -ldflags "-s -w -X main.Version=v0.2.2" -o build/rsyslog-rest-api .
```

### Clean

```bash
make clean
# Entfernt build/ Verzeichnis
```

---

## 🧪 Testing

### Docker Testumgebung

**Setup:**
```bash
# Binary bauen
make build-static

# Container starten
cd docker
docker-compose up -d

# Logs verfolgen
docker-compose logs -f
```

**Manuelle Tests:**
```bash
# Health Check
curl http://localhost:8000/health

# Logs abrufen
curl "http://localhost:8000/logs?limit=5"

# Multi-Value Filter
curl "http://localhost:8000/logs?FromHost=web01&FromHost=web02"
```

**Test-Suite:**
```bash
cd docker
./test-v0.2.2.sh

# Erwartete Ausgabe:
# Passed: 22
# Failed: 0
```

**Cleanup:**
```bash
docker-compose down
```

### Unit Tests (TODO)

Aktuell keine Unit Tests vorhanden. Geplant für v0.3.0.

```bash
# Future
go test ./...
go test -cover ./...
```

---

## 📝 Code-Struktur

### main.go Aufbau

```go
// Global Variables
var (
    config          *Configuration
    db              *sql.DB
    availableColumns []string
)

// RFC Mappings
var rfcSeverity = map[int]string{...}
var rfcFacility = map[int]string{...}

// Structs
type Configuration struct {...}
type LogEntry struct {...}
type LogsResponse struct {...}

// Functions
func loadConfiguration() (*Configuration, error)
func initDatabase(cfg *Configuration) error
func loadAvailableColumns() error

// Middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc
func authMiddleware(next http.HandlerFunc) http.HandlerFunc

// Handlers
func handleRoot(w http.ResponseWriter, r *http.Request)
func handleHealth(w http.ResponseWriter, r *http.Request)
func handleLogs(w http.ResponseWriter, r *http.Request)
func handleMeta(w http.ResponseWriter, r *http.Request)

// Helpers
func buildFilters(...) (string, []interface{}, error)
func respondJSON(w http.ResponseWriter, status int, data interface{})
```

### Wichtige Funktionen

**buildFilters()** - Kritischste Funktion!

```go
func buildFilters(
    startDate, endDate string,
    fromHosts, priorities, facilities, messages, sysLogTags []string
) (string, []interface{}, error)
```

- Erstellt WHERE clause für SQL
- Validiert alle Parameter
- Unterstützt Multi-Value Filter
- Returns: (whereClause, args, error)

**handleLogs()** - Haupt-Endpoint

```go
func handleLogs(w http.ResponseWriter, r *http.Request)
```

- Parse query parameters
- Build filters
- Count total
- Query database
- Map results to LogEntry structs
- Return JSON

---

## 🎨 Coding Guidelines

### Go Best Practices

**Code-Style:**
```bash
# gofmt (automatisch)
go fmt ./...

# golint (empfohlen)
go install golang.org/x/lint/golint@latest
golint ./...

# go vet (Fehler finden)
go vet ./...
```

**Naming:**
- Exportierte Funktionen: `PascalCase`
- Private Funktionen: `camelCase`
- Konstanten: `ALL_CAPS` oder `PascalCase`

**Error Handling:**
```go
// ✅ RICHTIG
if err != nil {
    return fmt.Errorf("failed to do X: %v", err)
}

// ❌ FALSCH
if err != nil {
    panic(err)  // Nur in main() bei Startup
}
```

### SQL Best Practices

**Prepared Statements:**
```go
// ✅ RICHTIG - Verhindert SQL Injection
query := "SELECT * FROM SystemEvents WHERE FromHost = ?"
db.Query(query, fromHost)

// ❌ FALSCH - SQL Injection möglich!
query := fmt.Sprintf("SELECT * FROM SystemEvents WHERE FromHost = '%s'", fromHost)
db.Query(query)
```

**Connection Pooling:**
```go
// In initDatabase()
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

---

## 🔧 Adding Features

### Neuer Filter-Parameter

**1. URL Parameter parsen:**
```go
// In handleLogs()
newParam := query["NewParam"]  // Multi-Value Array
```

**2. buildFilters() erweitern:**
```go
func buildFilters(..., newParam []string) (string, []interface{}, error) {
    // ...
    
    // Validierung
    if len(newParam) > 0 {
        // Validation logic
        if invalid {
            return "", nil, fmt.Errorf("NewParam must be...")
        }
    }
    
    // Filter bauen (Multi-Value!)
    if len(newParam) > 0 {
        placeholders := make([]string, len(newParam))
        for i, val := range newParam {
            placeholders[i] = "?"
            args = append(args, val)
        }
        conditions = append(conditions, 
            fmt.Sprintf("ColumnName IN (%s)", strings.Join(placeholders, ",")))
    }
    
    // ...
}
```

**3. Testen:**
```bash
curl "http://localhost:8000/logs?NewParam=value1&NewParam=value2"
```

### Neuer Endpoint

**1. Handler erstellen:**
```go
func handleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Parse parameters
    
    // Query database
    
    // Return JSON
    respondJSON(w, http.StatusOK, data)
}
```

**2. Route registrieren:**
```go
// In main()
http.HandleFunc("/new-endpoint", 
    corsMiddleware(loggingMiddleware(authMiddleware(handleNewEndpoint))))
```

**3. Dokumentieren:**
```markdown
<!-- In docs/api-reference.md -->
## GET /new-endpoint

Description...
```

---

## 🐛 Debugging

### Lokales Debugging

**Run in Foreground:**
```bash
cd /opt/rsyslog-rest-api
export $(cat .env | xargs)
./rsyslog-rest-api

# Logs erscheinen direkt in Terminal
```

**Mit Delve Debugger:**
```bash
# Delve installieren
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug starten
dlv debug

# Breakpoints setzen
(dlv) break main.handleLogs
(dlv) continue
```

### SQL Query Logging

**Temporär in Code einfügen:**
```go
// In handleLogs() vor db.Query()
log.Printf("SQL: %s", sqlQuery)
log.Printf("Args: %v", args)
```

### HTTP Request Debugging

```bash
# Mit curl verbose
curl -v -H "X-API-Key: KEY" "http://localhost:8000/logs"

# Mit httpie
http -v localhost:8000/logs X-API-Key:KEY
```

---

## 🚢 Release Process

### Version Bump

```bash
# 1. Tag erstellen
git tag -a v0.2.3 -m "Release v0.2.3"

# 2. Tag pushen
git push origin v0.2.3

# 3. GitHub Actions baut automatisch
```

### Manueller Release

```bash
# 1. Build für Platforms
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags "-s -w -X main.Version=v0.2.3" \
  -o rsyslog-rest-api-linux-amd64 .

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
  go build -ldflags "-s -w -X main.Version=v0.2.3" \
  -o rsyslog-rest-api-linux-arm64 .

# 2. Checksums
sha256sum rsyslog-rest-api-* > SHA256SUMS

# 3. GitHub Release erstellen
# https://github.com/phil-bot/rsyslog-rest-api/releases/new
```

---

## 🤝 Contributing

### Workflow

1. **Fork Repository**
   ```bash
   # Auf GitHub: Fork Button
   git clone https://github.com/YOUR_USERNAME/rsyslog-rest-api.git
   ```

2. **Feature Branch erstellen**
   ```bash
   git checkout -b feature/awesome-feature
   ```

3. **Änderungen machen**
   ```bash
   # Code ändern
   vim main.go
   
   # Build & Test
   make build-static
   cd docker && docker-compose up -d
   ./test-v0.2.2.sh
   ```

4. **Commit**
   ```bash
   git add .
   git commit -m "Add awesome feature"
   
   # Commit Message Format:
   # - Add: Neue Features
   # - Fix: Bug Fixes
   # - Update: Updates
   # - Docs: Dokumentation
   ```

5. **Push & Pull Request**
   ```bash
   git push origin feature/awesome-feature
   
   # Auf GitHub: Create Pull Request
   ```

### Code Review Checkliste

Bevor PR erstellt wird:

- [ ] Code formatiert (`go fmt ./...`)
- [ ] Keine Linter-Warnings (`go vet ./...`)
- [ ] Gebaut ohne Errors (`make build-static`)
- [ ] Docker Tests bestanden (`./test-v0.2.2.sh`)
- [ ] Dokumentation aktualisiert
- [ ] Changelog aktualisiert (CHANGELOG.md)

---

## 📚 Dokumentation

### Markdown Files

Alle Dokumentation in `docs/`:

```bash
# Neue Seite erstellen
nano docs/new-page.md

# Struktur:
# [← Zurück zur Übersicht](index.md)
# 
# # Title
# Content...
# 
# [← Zurück zur Übersicht](index.md) | [Weiter →](next.md)
```

### Aktualisieren

```bash
# README.md - Kurz halten!
# docs/api-reference.md - Bei API-Änderungen
# docs/changelog.md - Bei jedem Release
# docs/examples.md - Bei neuen Features
```

---

## 🔮 Roadmap

### v0.3.0 (Geplant)

- [ ] Negation filters (`?exclude=FromHost:value`)
- [ ] Complex query support
- [ ] Unit Tests
- [ ] GitHub Actions CI/CD

### v0.4.0 (Geplant)

- [ ] Statistics endpoint (`/stats`)
- [ ] Aggregations
- [ ] Timeline/Histogram
- [ ] WebSocket support (live logs)

### Future

- [ ] PostgreSQL support?
- [ ] GraphQL API?
- [ ] Web UI?

---

## 📞 Support

### Fragen?

- **GitHub Discussions:** https://github.com/phil-bot/rsyslog-rest-api/discussions
- **Issues:** https://github.com/phil-bot/rsyslog-rest-api/issues

### Bug Reports

**Template:**
```markdown
**Environment:**
- OS: Ubuntu 22.04
- Go Version: 1.21.5
- API Version: v0.2.2

**Expected Behavior:**
[Was sollte passieren]

**Actual Behavior:**
[Was passiert tatsächlich]

**Steps to Reproduce:**
1. ...
2. ...

**Logs:**
```
[Relevante Logs]
```
```

---

## 🙏 Credits

Danke an alle Contributors!

**Maintainer:**
- [@phil-bot](https://github.com/phil-bot)

**Built with:**
- Go Programming Language
- go-sql-driver/mysql
- rsyslog
- MariaDB/MySQL

---

[← Zurück zur Übersicht](index.md)

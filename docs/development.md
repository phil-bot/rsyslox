# Development Guide

[← Back to overview](index.md)

Development Guide and Contributing Guidelines.

## 🏗️ Project architecture

### Overview

```
rsyslog-rest-api
├── main.go # Main application (single file!)
├── go.mod # Go dependencies
├── go.sum # Dependency checksums
├─── Makefile # Build Automation
├── .env.example # Config template
├── rsyslog-rest-api.service # systemd Service
├── docker/ # Docker test environment
│ ├─── Dockerfile
│ ├─── docker-compose.yml
│ ├─── entrypoint.sh
│ ├─── log-generator.sh
│ └── test-v0.2.2.sh
└── docs/ # Documentation
    └── ...
```

### Technology stack

- **Language:** Go 1.21+
- **Database Driver:** go-sql-driver/mysql v1.7.1
- **Database:** MySQL 5.7+ / MariaDB 10.3+
- **Testing:** Docker + Bash Scripts

---

## 🚀 Setup Development Environment

### Prerequisites

```bash
# Install Go (1.21+)
# Ubuntu/Debian
sudo apt-get install golang-1.21

# Or: https://go.dev/dl/

# Verify
go version
# go version go1.21.x linux/amd64

# Git
sudo apt-get install git

# Make
sudo apt-get install make

# Docker (for testing)
sudo apt-get install docker.io docker-compose
```

### Clone repository

```bash
git clone https://github.com/phil-bot/rsyslog-rest-api.git
cd rsyslog-rest-api
```

### Install dependencies

```bash
# Go Modules Download
go mod download

# Verify
go mod verify
```

---

## 🔨 Build

### Development Build

```bash
# Standard Build
make build

# Binary is in: ./build/rsyslog-rest-api
./build/rsyslog-rest-api
```

### Static Build (Production)

```bash
# Static Binary (no libc dependency)
make build-static

# Verify (no dependencies)
ldd ./build/rsyslog-rest-api
# not a dynamic executable
```

### With version

```bash
# Version from Git tag
VERSION=v0.2.2 make build

# Or manually
go build -ldflags "-s -w -X main.version=v0.2.2" -o build/rsyslog-rest-api .
```

### Clean

```bash
make clean
# Removes build/ directory
```

---

## 🧪 Testing

### Docker test environment

**Setup:**
```bash
# Build binary
make build-static

# Start container
cd docker
docker-compose up -d

# Track logs
docker-compose logs -f
```

**Manual tests:**
```bash
# Health Check
curl http://localhost:8000/health

# Retrieve logs
curl "http://localhost:8000/logs?limit=5"

# Multi-Value Filter
curl "http://localhost:8000/logs?FromHost=web01&FromHost=web02"
```

**Test-Suite:**
```bash
cd docker
./test-v0.2.2.sh

# Expected output:
# Passed: 22
# Failed: 0
```

**Cleanup:**
```bash
docker-compose down
```

### Unit Tests (TODO)

Currently no unit tests available. Planned for v0.3.0.

```bash
# Future
go test ./...
go test -cover ./...
```

---

## 📝 code structure

### main.go structure

```go
// Global Variables
var (
    config *Configuration
    db *sql.DB
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

// middleware
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

### Important functions

**buildFilters()** - Most critical function!

```go
func buildFilters(
    startDate, endDate string,
    fromHosts, priorities, facilities, messages, sysLogTags []string
) (string, []interface{}, error)
```

- Creates WHERE clause for SQL
- Validates all parameters
- Supports multi-value filters
- Returns: (whereClause, args, error)

**handleLogs()** - main endpoint

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

**Code style:**
```bash
# gofmt (automatic)
go fmt ./...

# golint (recommended)
go install golang.org/x/lint/golint@latest
golint ./...

# go vet (find errors)
go vet ./...
```

**Naming:**
- Exported functions: `PascalCase`
- Private functions: `camelCase`
- Constants: `ALL_CAPS` or `PascalCase

**Error handling:**
```go
// ✅ CORRECT
if err != nil {
    return fmt.Errorf("failed to do X: %v", err)
}

// ❌ WRONG
if err != nil {
    panic(err) // Only in main() at startup
}
```

### SQL Best Practices

**Prepared Statements:**
```go
// ✅ CORRECT - Prevents SQL injection
query := "SELECT * FROM SystemEvents WHERE FromHost = ?"
db.Query(query, fromHost)

// ❌ WRONG - SQL injection possible!
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

### New filter parameter

**1. parse URL parameter:**
```go
// In handleLogs()
newParam := query["NewParam"] // Multi-Value Array
```

**2. extend buildFilters():**
```go
func buildFilters(..., newParam []string) (string, []interface{}, error) {
    // ...
    
    // validation
    if len(newParam) > 0 {
        // Validation logic
        if invalid {
            return "", nil, fmt.Errorf("NewParam must be...")
        }
    }
    
    // Build filter (Multi-Value!)
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

**3. testing:**
```bash
curl "http://localhost:8000/logs?NewParam=value1&NewParam=value2"
```

### New endpoint

**1. create handler:**
```go
func handleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Parse parameters
    
    // Query database
    
    // Return JSON
    respondJSON(w, http.StatusOK, data)
}
```

**2. register route:**
```go
// In main()
http.HandleFunc("/new-endpoint",
    corsMiddleware(loggingMiddleware(authMiddleware(handleNewEndpoint))))
```

**3. document:**
```markdown
<!-- In docs/api-reference.md -->
## GET /new-endpoint

Description...
```

---

## 🐛 Debugging

### Local debugging

**Run in Foreground:**
```bash
cd /opt/rsyslog-rest-api
export $(cat .env | xargs)
./rsyslog-rest-api

# Logs appear directly in the terminal
```

**With Delve Debugger:**
```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Start debug
dlv debug

# Set breakpoints
(dlv) break main.handleLogs
(dlv) continue
```

### SQL Query Logging

**Insert temporary in code:**
```go
// In handleLogs() before db.Query()
log.Printf("SQL: %s", sqlQuery)
log.Printf("Args: %v", args)
```

### HTTP Request Debugging

```bash
# With curl verbose
curl -v -H "X-API-Key: KEY" "http://localhost:8000/logs"

# With httpie
http -v localhost:8000/logs X-API-Key:KEY
```

---

## 🚢 Release Process

### Version Bump

```bash
# Create 1st tag
git tag -a v0.2.3 -m "Release v0.2.3"

# push 2nd tag
git push origin v0.2.3

# 3. GitHub Actions builds automatically
```

### Manual release

```bash
# 1st build for platforms
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags "-s -w -X main.version=v0.2.3" \
  -o rsyslog-rest-api-linux-amd64 .

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
  go build -ldflags "-s -w -X main.version=v0.2.3" \
  -o rsyslog-rest-api-linux-arm64 .

# 2. checksums
sha256sum rsyslog-rest-api-* > SHA256SUMS

# 3. create GitHub release
# https://github.com/phil-bot/rsyslog-rest-api/releases/new
```

---

## 🤝 Contributing

### Workflow

1. **Fork Repository**
   ```bash
   # On GitHub: Fork button
   git clone https://github.com/YOUR_USERNAME/rsyslog-rest-api.git
   ```

2. **Create feature branch**
   ```bash
   git checkout -b feature/awesome-feature
   ```

3. **Make changes**
   ```bash
   # Change code
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
   # - Add: New features
   # - Fix: Bug fixes
   # - Update: Updates
   # - Docs: Documentation
   ```

5. **Push & Pull Request**
   ```bash
   git push origin feature/awesome-feature
   
   # On GitHub: Create Pull Request
   ```

### Code Review Checklist

Before PR is created:

- [ ] Code formatted (`go fmt ./...`)
- [ ] No linter warnings (`go vet ./...`)
- [ ] Built without errors (`make build-static`)
- [ ] Docker tests passed (`./test-v0.2.2.sh`)
- [ ] Documentation updated
- [ ] Changelog updated (CHANGELOG.md)

---

## 📚 Documentation

### Markdown Files

All documentation in `docs/`:

```bash
# Create new page
nano docs/new-page.md

# Structure:
# [← Back to overview](index.md)
#
# # Title
# Content...
#
# [← Back to overview](index.md) | [Next →](next.md)
```

### Refresh

```bash
# README.md - Keep it short!
# docs/api-reference.md - For API changes
# docs/changelog.md - For each release
# docs/examples.md - For new features
```

---

## 🔮 Roadmap

### v0.3.0 (Planned)

- [ ] Negation filters (`?exclude=FromHost:value`)
- [ ] Complex query support
- [ ] Unit tests
- [ ] GitHub Actions CI/CD

### v0.4.0 (Planned)

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

### Questions?

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
[What should happen]

**Actual Behavior:**
[What actually happens]

**Steps to Reproduce:**
1. ...
2. ...

**Logs:**
```
[Relevant logs]
```
```

---

## 🙏 Credits

Thanks to all contributors!

**Maintainer:**
- [@phil-bot](https://github.com/phil-bot)

**Built with:**
- Go Programming Language
- go-sql-driver/mysql
- rsyslog
- MariaDB/MySQL

---

[← Back to overview](index.md)

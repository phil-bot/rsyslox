# Changelog

All notable changes to rsyslog REST API.

## [Unreleased]

### Planned for v0.4.0
- Negation filters
- Advanced filter combinations
- Unit tests

### Planned for v0.5.0
- Statistics endpoint
- Aggregations
- Prometheus metrics

---

## [v0.3.0] - 2025-02-19

### Added
- **`Severity` field** in all `/logs` responses — RFC-5424 compliant severity value (0–7)
- **`Severity_Label` field** in all `/logs` responses — human-readable label (e.g. `"Error"`)
- **`Priority` field** now contains the true RFC PRI value (`Facility × 8 + Severity`)
- **Automatic rsyslog version detection** at API startup — determines whether the `Priority` DB column stores legacy Severity (< 8.2204.0) or modern RFC PRI (≥ 8.2204.0)
- **Per-row format detection** — mixed datasets (produced by a rsyslog upgrade without data migration) are handled correctly on a row-by-row basis
- **`?Severity=` filter parameter** for `/logs` and `/meta/{column}` — RFC-correct name
- **`/meta/Severity` virtual column** — returns distinct severity values with labels, computed via `Priority MOD 8`
- **`internal/database/priority_detection.go`** — isolated, testable detection logic
- **Cleanup / Housekeeping service** — automatically deletes the oldest `SystemEvents` records when disk usage exceeds a configurable threshold
- **Disk-based retention** — percentage-based configuration instead of fixed age or size limits
- **`internal/cleanup/cleaner.go`** — isolated, modular cleanup service with its own lifecycle (`Start` / `Stop`)

### Changed
- `Priority_Label` field **removed** from responses — RFC PRI has no standardised label; use `Severity_Label` instead
- `?Priority=` filter parameter is now a **deprecated alias** for `?Severity=` — existing clients continue to work
- `ErrCodeInvalidPriority` aliased to `ErrCodeInvalidSeverity` — error responses now use `INVALID_SEVERITY`
- `GetPriorityLabel` aliased to `GetSeverityLabel` internally
- `IsValidPriority` aliased to `IsValidSeverity` internally
- Facility meta endpoint (`/meta/Facility`) now returns values with RFC labels (same format as `/meta/Severity`)
- `/meta/{column}` no longer applies a default time filter — without explicit `start_date`/`end_date` all distinct values across the entire dataset are returned
- `internal/config/config.go` — extended with 5 new cleanup configuration fields and helper functions `getEnvFloat` / `getEnvInt`
- `main.go` — cleanup service initialized and started alongside the API server
- `.env.example` — documented all new cleanup variables
- Root `README.md` — streamlined to a single documentation link, added logo and program description
- Documentation — new [Cleanup Guide](guides/cleanup.md), updated sidebar, configuration reference and changelog

### Fixed
- Severity filter now uses `Priority MOD 8` in SQL — works correctly for legacy, modern, and mixed datasets without any configuration

### Cleanup Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `CLEANUP_ENABLED` | `false` | Enable the cleanup service |
| `CLEANUP_DISK_PATH` | `/var/lib/mysql` | Filesystem path to monitor |
| `CLEANUP_THRESHOLD_PERCENT` | `85` | Disk usage threshold in percent |
| `CLEANUP_BATCH_SIZE` | `1000` | Records deleted per cleanup run |
| `CLEANUP_INTERVAL` | `15m` | Check interval (Go duration format) |

---

## [v0.2.3] - 2025-02-15

### Added
- Enhanced multi-value filter performance
- Better error validation messages
- Improved meta endpoint filtering

### Fixed
- Multi-value filter edge cases
- Error message clarity
- Performance improvements

### Changed
- Updated dependencies
- Improved documentation
- Better code organization

---

## [v0.2.2] - 2025-02-09

### Added
- **Multi-value filters** for all parameters
- **Extended columns** - All 25+ SystemEvents fields
- **Enhanced meta endpoint** - Supports filtering
- **Live log generator** for Docker

### Changed
- API response includes all available columns
- Meta endpoint returns labels
- Improved filter validation

### Fixed
- Null handling for extended columns
- Meta endpoint performance
- Database index creation

---

## [v0.2.1] - 2025-01-15

### Fixed
- Database connection timeout
- Memory leak in queries
- CORS preflight handling

### Changed
- Improved error messages
- Better logging format

---

## [v0.2.0] - 2024-12-20

### Added
- RFC-5424 labels
- Meta endpoint
- SSL/TLS support
- CORS configuration

### Changed
- Response format includes labels
- Configuration via .env

---

## [v0.1.0] - 2024-10-01

Initial release.

### Features
- Basic REST API
- Authentication
- Pagination
- Docker testing

---

## Version History

| Version | Date | Summary |
|---------|------|---------|
| v0.3.0 | 2025-02-19 | RFC Severity/Priority, cleanup service, docs overhaul |
| v0.2.3 | 2025-02-15 | Performance, structured errors |
| v0.2.2 | 2025-02-09 | Multi-value filters, extended columns |
| v0.2.1 | 2025-01-15 | Bug fixes |
| v0.2.0 | 2024-12-20 | Labels, SSL |
| v0.1.0 | 2024-10-01 | Initial release |

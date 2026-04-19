# Development & Contributing

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- MySQL / MariaDB (or the [Docker test environment](docker.md))
- make

### Running Locally

**Backend:**
```bash
cp config.dev.toml.example config.dev.toml
# Edit config.dev.toml — set database credentials and generate an admin hash:
go run . hash-password yourpassword

export RSYSLOX_CONFIG=./config.dev.toml
go run .
# or: make dev
```

The backend serves on the port defined in the config (default 8000). If the config file does not exist, it starts in setup wizard mode.

**Frontend** (separate terminal):
```bash
cd frontend
npm install
npm run dev   # starts on http://localhost:5173
```

Vite proxies `/api/*` and `/health` to `http://localhost:8000`.

**Build:**
```bash
make all           # frontend + Redoc download + Go binary
make build         # Go binary only (requires frontend/dist to exist)
make build-static  # Static binary for production/Docker
```

---

## Project Structure

```
rsyslox/
├── main.go                 # Entry point; CLI commands (hash-password, …)
├── embed.go                # go:embed directives for frontend/dist and docs/api-ui
├── internal/
│   ├── auth/               # Session tokens, bcrypt, API key verification
│   ├── cleanup/            # Disk-based log retention goroutine
│   ├── config/             # TOML config: load, save, validate, AES-GCM encryption
│   ├── database/           # MySQL connection, query layer, TTL cache
│   └── server/             # HTTP server, routing, handlers, setup wizard
├── frontend/
│   ├── src/
│   │   ├── api/            # API client (fetch wrapper, auth headers)
│   │   ├── assets/         # Global CSS, CSS variables
│   │   ├── components/     # AppHeader, FilterPanel, LogTable, LogDetail
│   │   ├── composables/    # useLocale (i18n), useClickOutside
│   │   ├── i18n/           # en.json, de.json — all UI strings
│   │   ├── router/         # Vue Router — auth guard, setup detection
│   │   ├── stores/         # logs.js, auth.js, preferences.js
│   │   └── views/          # LogsView, AdminView, LoginView, SetupView
│   └── dist/               # Built output (embedded into binary via go:embed)
├── docs/                   # Docsify documentation (GitHub Pages)
├── docker/                 # Docker test environment with live log generator
├── scripts/
│   └── install.sh          # Installer / uninstaller
├── Makefile
└── rsyslox.service         # systemd unit file template
```

---

## Frontend Architecture

The frontend is a Vue 3 + Vite single-page application using the Composition API throughout. State is managed in plain reactive modules (`stores/`), not Pinia.

### State Stores

**`stores/logs.js`** — central log state (entries, filters, pagination, selection, auto-refresh). All filter changes trigger `resetPage()` + `fetchLogs()` via a single `watch`. Page changes trigger `fetchLogs()` via an arrow-wrapped watcher to prevent the page number being passed as the `fromRefresh` argument.

**`stores/auth.js`** — session token and role in `sessionStorage`. The API client reads `rsyslox_token` directly from `sessionStorage` to avoid circular imports.

**`stores/preferences.js`** — language, time format, font size, auto-refresh interval in `localStorage`. Exports reactive refs directly. Font size is applied immediately to `document.documentElement.style.fontSize` on load and on every change.

### i18n

Translation keys live in `src/i18n/en.json` and `src/i18n/de.json`. The `useLocale` composable provides:

- `t(key, vars?)` — returns the translated string; falls back to English; interpolates `{name}`-style variables
- `fmtNumber(n)` — formats numbers with locale-appropriate thousands separators

```javascript
import { useLocale } from '@/composables/useLocale'
const { t, fmtNumber } = useLocale()

t('filter.severity')           // → "Severity"
t('logs.showing', { n: 42 })   // → "Showing 42 entries"
fmtNumber(1234567)              // → "1,234,567" (EN) or "1.234.567" (DE)
```

To add a translation key: add it to both `i18n/en.json` and `i18n/de.json`, then use `t('your.key')`. Keys use dot notation with a section prefix (`nav.`, `filter.`, `table.`, `admin.`, `prefs.`).

### Dynamic Table Sizing

`LogsView.computePageSize()` is called by a `ResizeObserver` on `.logs-main` and on font size changes. It measures toolbar/thead/pagination heights from the live DOM, derives natural row height from the current `font-size` of `<html>`, calculates `n = floor(available / naturalRowH) - 1`, sets `exactRowH = floor(available / n)` as `--row-h` on `.table-scroll`, then calls `setPageSize(n)` + `fetchLogs()` if `n` changed. Rows use `height: var(--row-h)` to fill the container with no gap.

### Flash Animation

`fetchLogs(fromRefresh = false)` compares incoming row IDs against the previous set. Only IDs absent from the previous set receive the `row-new` class, triggering a `row-flash` keyframe animation. The `fromRefresh` flag is `true` only from the auto-refresh timer — never from filter changes or manual reloads.

---

## Backend Architecture

### Config

`internal/config` handles loading (`config.Load()`), saving (`config.Save()`), and validation. If the config file does not exist, `Load()` returns `setupMode = true` — the server then mounts only the setup routes.

The database password is encrypted with AES-GCM (`internal/config/crypto.go`). The encryption key is derived from `/etc/machine-id` using SHA-256. Passwords with an `enc:` prefix are decrypted when building the DSN; plain passwords (during initial setup) are encrypted before being saved.

### Auth

Admin sessions use a random 32-byte token stored in an in-memory map with expiry. Tokens are transmitted via `X-Session-Token`. Read-only API keys are verified by computing SHA-256 of the submitted value and comparing against stored hashes in the config.

### Cleanup

`internal/cleanup` runs as a goroutine. When enabled, it periodically checks disk usage at the configured path using `statvfs`. If usage exceeds the threshold, it deletes the oldest `batch_size` rows from the syslog table and repeats until usage drops below the threshold or no rows remain. Config changes are applied at runtime via `UpdateConfig(cfg Config)` without a restart.

---

## Contributing

### Coding Standards

- Go: run `go fmt ./...` and `go vet ./...` before committing
- Vue: keep components focused; use `<script setup>` syntax
- Commit messages: [Conventional Commits](https://www.conventionalcommits.org/) format

### Commit Message Format

```
<type>: <subject>

<body>

<footer>
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

**Examples:**
```
feat: add Facility filter to log viewer
```
```
fix: handle null SysLogTag in meta query

NULL values caused a panic in the distinct-value aggregation.
Filtered out with WHERE SysLogTag IS NOT NULL.

Fixes #42
```

### Branch Naming

```
feature/your-feature-name
fix/short-description
docs/what-you-are-documenting
```

### Workflow

1. Fork the repository on GitHub
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/rsyslox.git`
3. Create a feature branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Submit a pull request

### Pull Request Checklist

- [ ] `go fmt ./...` and `go vet ./...` pass
- [ ] All Go tests pass: `go test ./...`
- [ ] Frontend builds without errors: `cd frontend && npm run build`
- [ ] New translation keys added to both `en.json` and `de.json`
- [ ] Changelog entry added to `docs/development/changelog.md` under `[Unreleased]`
- [ ] Documentation updated if behaviour changed

### Tests

```bash
go test ./...          # all tests
go test -cover ./...   # with coverage
go test ./internal/config/  # specific package
```

---

## Release Process

For maintainers:

1. Move `[Unreleased]` entries to a new version section in `changelog.md`
2. Create and push a tag: `git tag -a v0.X.Y -m "Release v0.X.Y" && git push origin v0.X.Y`
3. GitHub Actions builds amd64/arm64 binaries, creates an offline package, and publishes the release with SHA-256 checksums automatically

Pre-releases are detected automatically from the tag name (e.g. `v0.5.0-beta`).

---

## Reporting Bugs & Feature Requests

- **Bugs:** [GitHub Issues](https://github.com/phil-bot/rsyslox/issues) — include the version (`/health` → `version`), OS/arch, steps to reproduce, and `sudo journalctl -u rsyslox -n 100` output
- **Feature requests:** Open a [GitHub Discussion](https://github.com/phil-bot/rsyslox/discussions) first for larger changes

# Changelog

All notable changes to rsyslox.

## [Unreleased]

---

## [v0.5.0] - 2026-03-31

This release overhauls the user interface and preference system. A new About
dialog shows version information, the navigation bar is redesigned without
a dropdown, and the filter panel gains collapsible sections. Server-side
defaults for all user preferences are now configurable. Three bugs are fixed:
the live-mode duration highlight, the cleanup threshold mismatch, and cleanup
configuration changes now take effect without a server restart.

### Added

**About dialog**
- New "About rsyslox" modal accessible from the navigation bar: shows version,
  license, author, source and documentation links

**Server-side preference defaults**
- `config.toml` now supports `default_language`, `default_font_size` and
  `default_time_format` under `[server]`, alongside the existing
  `default_time_range` and `auto_refresh_interval`
- `/health` response includes all five defaults in a `defaults` object
- The frontend applies server defaults to any preference key that has not yet
  been explicitly set by the user in `localStorage` — existing user settings
  are never overwritten
- All five defaults are editable in **Admin → Server → Default Values**

**Navigation bar**
- Redesigned as a flat icon bar — no dropdown menu
- Theme toggle replaced with an inline slide switch (sun / moon icons)
- Docs link has an external-tab indicator (↗ arrow badge)
- Logout, Settings, About, Docs are direct icon buttons; Settings and Docs
  are only shown to admin users

**Filter panel**
- Each filter section (Time Range, Severity, Facility, Tag, Host, Message
  Search) is individually collapsible via its header; collapsed state
  persists in `localStorage` (`rsyslox_filter_collapsed`)
- "Close filters" button inside the panel header (replaces the header toggle)
- When the sidebar is collapsed, a "Filters" button appears in the log table
  toolbar to re-open it

**Admin panel**
- Single **Save** button per tab, shown in a bottom bar (Server tab, Database
  tab); no intermediate save buttons inside form sections
- Database settings and Cleanup settings are saved together with one request
- Preferences tab: auto-saves on every change; no Save button
- API Keys tab: no Save button (create / revoke are immediate actions)

**Font sizes**
- Increased to 14 / 16 / 18 px (was 13 / 14 / 15 px)

**Default time range**
- Changed from `1h` to `24h` for new sessions without a stored preference

### Changed

- `server.New()` accepts a `*cleanup.Cleaner` parameter so cleanup config
  changes propagate at runtime
- `handlers.NewHealthHandler()` accepts `*config.Config` to supply defaults
- `admin.NewConfigHandler()` accepts `*cleanup.Cleaner`

### Fixed

- **Live-mode duration highlight** — the active duration button in the filter
  panel lost its highlight after a short time in live mode because `activeDur`
  compared the rolling `endDate` against the fixed `startDate`. Fix: return
  `relativeDur` directly when `autoRefresh` is true
- **Cleanup threshold mismatch** — `diskUsagePercent` used `stat.Bfree`
  (includes root-reserved blocks) while the disk widget used `stat.Bavail`.
  Both now use `stat.Bavail` for a consistent reading
- **Cleanup changes required restart** — the `Cleaner` goroutine held a copy
  of the initial config. A new `UpdateConfig(cfg Config)` method updates the
  running goroutine via a mutex + signal channel; changes from the Admin panel
  take effect immediately
- **Login flash on first load** — navigating to `/` briefly rendered the log
  view before redirecting to `/login`. Fix: the `/` route redirect now checks
  `auth.isAuthenticated` synchronously and goes directly to `/login` or
  `/logs` without an intermediate navigation

---

## [v0.4.3] - 2026-03-02

This release overhauls the Admin panel with fully editable server and database
configuration, browser-triggered server restart, SSL certificate management, and
a live disk usage widget. The filter panel receives visual fixes and message
search highlighting.

### Added

**Admin Panel — Server Settings**
- Host and port are now editable fields (previously read-only after setup)
- `use_ssl` toggle is now saved as part of the main Save action (no separate auto-save)
- Restart-required banner: appears at the top of the admin panel after saving any
  setting that requires a restart; stays visible across tab navigation until a
  restart is performed or dismissed
- Browser-triggered server restart (`POST /api/admin/restart`) using `syscall.Exec`
  — replaces the current process in-place without requiring a process manager; the
  frontend polls `/health` and reloads automatically once the server is back
- Context-sensitive hint texts on **Host** and **Allowed Origins** fields

**Admin Panel — SSL / TLS**
- New SSL section appears under Server settings when SSL is enabled
- `POST /api/admin/ssl/generate` — generates a self-signed ECDSA P-256 certificate
  (10-year validity) and writes it to the configured cert/key paths
- `POST /api/admin/ssl/upload` — multipart upload of a custom certificate and
  private key; rolls back the cert file if key upload fails
- Auto-generation on startup: `config.EnsureSSLCerts()` is called before
  `ListenAndServeTLS`; if either cert or key is missing, a self-signed certificate
  is generated automatically — no manual step required
- `internal/config/ssl.go` — shared cert generation logic used by both startup
  and the admin HTTP endpoint

**Admin Panel — Database**
- Database settings are now a fully editable form: host, port, name, user, password
- Password field accepts a new value to change it (blank = keep current);
  value is AES-GCM encrypted before writing to `config.toml`

**Admin Panel — Cleanup (merged into Database tab)**
- Cleanup configuration moved from its own tab into the Database tab as a sub-section
- Info callout warns that disk-usage monitoring reads the **local** filesystem and
  only works correctly when the database runs on the same host as rsyslox
- Disk usage widget: live progress bar (green < 75 %, amber < 90 %, red ≥ 90 %),
  used / free / total in human-readable units; refreshable on demand
- `GET /api/admin/disk` — returns `used_percent`, `used_bytes`, `free_bytes`,
  `total_bytes` for the configured `disk_path` via `syscall.Statfs`

**Filter Panel**
- Message search results highlighted inline in the Message column with `<mark>` tags

### Changed

- Facility pills use a dedicated `.fac-badge-btn` class; no longer inherit the white
  foreground colour from severity badges
- Tag filter section converted from a searchable list to a pill layout (consistent
  with facility); moved above the host section
- Panel header height aligned with the main toolbar (`min-height: 40px`)
- AppHeader: filter toggle button (funnel icon) placed to the right of the logo;
  Settings link moved from standalone header icon into the account dropdown
- Statistics nav item added as a placeholder (grayed out, "Coming soon" tooltip)
- `api/client.js` — `request()` reads the response body as text first, then attempts
  `JSON.parse`; connection-reset / non-JSON errors now produce readable messages
- `internal/server/server.go` — SSL, restart, and disk routes registered centrally

### Fixed

- SSL generation returned `JSON.parse: unexpected character` when the cert directory
  did not exist — caused by `defer file.Close()` keeping the handle open while the
  HTTP response was written; replaced with explicit `Close()` before each return
- Server failed to start (`open /etc/rsyslox/certs/cert.pem: no such file`) when
  `use_ssl = true` and no certificate was present; resolved by `EnsureSSLCerts`
  auto-generation on startup
- Database tab showed a read-only info grid with a note to edit `config.toml` manually
- Admin panel Vite build failure: `Element is missing end tag` at line 36 — unclosed
  `<template v-else>` in the Server section
- Stray `-->` comment fragment rendered as visible text throughout the admin panel

---

## [v0.4.0] - 2026-02-23

This release introduces the complete web frontend and the browser-based preferences
system, replacing the previous API-only interface. All changes previously tracked
under the unreleased `v1.0.0` label are included here.

### Added

**Web Frontend (Vue 3 + Vite)**
- Full-featured log viewer embedded in the binary via `go:embed`
- Dark/light theme with system preference detection and manual toggle
- Responsive layout: sidebar panel on desktop, slide-over modal on mobile
- Skeleton loading states for first render

**Log Viewer**
- Filter panel: time range (relative quick-select or absolute date/time), severity,
  facility, host, tag, and free-text message search
- Time range selector redesigned as a compact segmented control
- Date fields pre-filled on first render using the default duration (`1h`)
- Time shift buttons (`‹ Earlier` / `Later ›`) to step through log windows
- Log table with severity colour coding, monospace data columns, multi-row selection
- Detail panel: full message, all fields, expandable raw JSON, copy-to-clipboard
- Export selected or all visible rows as CSV or JSON (client-side, no server round-trip)
- Auto-refresh with countdown display; interval configurable from browser preferences
- Dynamic page size computed from viewport height; rows stretch to fill the container

**Admin Panel**
- Server settings editor (CORS origins, SSL toggle)
- Database info view (read-only; password always masked)
- Log cleanup configuration (disk path, threshold %, batch size, check interval)
- Read-only API key management: create named keys, list, revoke
- One-time key reveal after creation with copy button (plaintext never stored or logged)
- Preferences tab (default landing page)

**Internationalisation (i18n)**
- Translation files `src/i18n/en.json` and `src/i18n/de.json`; all UI strings externalised
- `useLocale` composable — reactive `t(key, vars?)` with variable interpolation
- All views and components fully translated (log viewer, admin panel, filter panel)

**User Preferences (browser-persisted)**
- Language, time format (12h/24h), font size, auto-refresh interval
- Applied immediately, no restart needed; stored in `localStorage`

**Configuration (TOML)**
- `/etc/rsyslox/config.toml` replaces `.env`
- AES-GCM encrypted database password, bcrypt admin password (cost 12), SHA-256 API keys
- First-run setup wizard (localhost-only until configured)

**Install & Operations**
- `scripts/install.sh` — guided installer with systemd hardening
- Offline API documentation via Redoc (embedded at `/docs`)

**CI/CD**
- GitHub Actions CI and release pipeline; multi-arch binaries (amd64/arm64)

### Changed

- Single binary embeds `frontend/dist/` and `docs/api-ui/` via `go:embed`
- `database.Connect()` uses `cfg.DSN()` from TOML config instead of env vars

### Fixed

- Flash animation on page change
- Impossible entry count display (e.g. `3,660 von 2,919`)
- Broken theme injection in `App.vue`
- Late imports in `FilterPanel.vue` and `LogsView.vue`

### Removed

- `.env` / `.env.example` — replaced by TOML config
- `API_KEY` env var — replaced by named, revocable read-only API keys

---

## [v0.3.0] - 2025-02-19

### Added
- `Severity` and `Severity_Label` fields in all `/logs` responses (RFC-5424)
- `Priority` field now contains true RFC PRI value (`Facility × 8 + Severity`)
- Automatic rsyslog version detection at startup
- `?Severity=` filter parameter (`?Priority=` kept as deprecated alias)
- `/meta/Severity` virtual column
- Cleanup service — disk-based log retention

### Changed
- `Priority_Label` removed from responses
- `/meta/{column}` no longer applies a default time filter

---

## [v0.2.3] - 2025-02-15

### Added
- Structured error responses (`code`, `message`, `details`, `field`)
- Enhanced multi-value filter performance

---

## [v0.2.2] - 2025-02-09

### Added
- Multi-value filters for all parameters
- All 25+ SystemEvents columns in responses
- Live log generator for Docker testing

---

## [v0.2.1] - 2025-01-15

### Fixed
- Database connection timeout
- Memory leak in queries
- CORS preflight handling

---

## [v0.2.0] - 2024-12-20

### Added
- RFC-5424 labels, meta endpoint, SSL/TLS support, CORS configuration

---

## [v0.1.0] - 2024-10-01

Initial release.

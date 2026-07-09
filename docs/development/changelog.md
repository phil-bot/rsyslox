# Changelog

All notable changes to rsyslox.

## [v0.5.5] - 2026-07-07

### Fixed

- **`EnsureFuturePartitions` failed permanently after older partitions were
  dropped** — the function always computed candidate weekly partitions as
  `now + 1..6 weeks`, independent of which partitions actually existed. Once
  cleanup had dropped the oldest weeks (e.g. `p2026_w27`–`w29`), a later run
  could try to insert one of those same week boundaries again via
  `REORGANIZE PARTITION p_future`, which MySQL rejects with `Error 1493:
  VALUES LESS THAN value must be strictly increasing for each partition`,
  since the newer partitions immediately before `p_future` already had a
  higher boundary. `EnsureFuturePartitions` now determines the actual highest
  existing partition boundary first (`highestPartitionBoundary`) and only
  ever appends weeks after it, making the strictly-increasing violation
  impossible by construction.
- **Retry-loop spammed the log every cleanup tick on failure** — because the
  above bug caused `EnsureFuturePartitions` to fail deterministically on
  every call, and `weeklyMaintenance` only updated its "last run" marker on
  *success*, the failing operation was retried on every single cleanup
  interval tick (e.g. every 60 s in tightly configured setups) instead of
  roughly once a week. `weeklyMaintenance` now rate-limits attempts to once
  per 24 hours regardless of outcome.

---

## [v0.5.4] - 2026-07-07

### Fixed

- **Unreclaimed disk space from unused FULLTEXT index** — `createIndexes()`
  created a `FULLTEXT` index on `SystemEvents.Message` at every startup, but
  the application has never queried it (message search uses `LIKE`, not
  `MATCH ... AGAINST` — see `internal/filters/builder.go`). MySQL/MariaDB's
  internal `FTS_*` auxiliary tables for this index are not reclaimed
  automatically and can grow to many gigabytes over time, independent of
  partition-based cleanup. `createIndexes()` no longer creates the index and
  now detects and drops any existing `FULLTEXT` index on `Message` (matched
  by type and column, not by name, so indexes created under any prior
  MySQL-assigned name are found) on every startup, freeing the associated
  `FTS_*` files immediately.

---

## [v0.5.3] - 2026-06-05

This release replaces the `DELETE`-based log cleanup with MySQL table
partitioning. Disk space is now reclaimed immediately when a partition is
dropped — no `OPTIMIZE TABLE` required. The Admin panel gains a dedicated
Cleanup tab with a live partition status widget.

### Added

**Partition-based log cleanup**
- `SystemEvents` is automatically migrated to weekly `RANGE COLUMNS`
  partitioning on first startup with cleanup enabled
- Partitions are named by ISO calendar week (`p2026_w23`)
- Up to 3 oldest partitions are dropped per cleanup run when disk usage
  exceeds the threshold; disk space is reclaimed immediately by the OS
- Weekly maintenance job pre-creates 6 weeks of future partitions
- `p_future` catch-all partition is always kept and never dropped
- Cleanup fails hard if the DB user lacks `ALTER` privilege — no silent
  `DELETE` fallback; a clear error with the required SQL grant is logged

**Admin panel — Cleanup tab**
- Dedicated **Cleanup** tab (previously a sub-section of the Database tab)
- Partition status badge: `Partition mode active` (green), `Migration in
  progress` (amber), `Partitioning failed` (red)
- Partition table showing name, row count, and size in MB for every partition
- Error block with SQL grant hint when `ALTER` privilege is missing
- `Batch size` field removed (not applicable in partition mode)

**New API endpoint**
- `GET /api/admin/cleanup/status` — returns partition mode, health flag,
  error string, and full partition list with row counts and sizes

### Changed

- `internal/cleanup/cleaner.go` — rewritten around partition mode; exposes
  `Status()` and `UpdateConfig()` methods; `Mode` type tracks
  `unknown / migrating / partition / failed`
- `internal/cleanup/partitions.go` — new module for all partition management:
  `IsPartitioned`, `MigrateToPartitions`, `EnsureFuturePartitions`,
  `DropOldestPartitions`, `HasAlterPrivilege`, `listPartitions`
- `server.New()` — `*cleanup.Cleaner` added as parameter so the status
  handler can query it directly
- `BatchSize` removed from `CleanupView` and `CleanupUpdateRequest` in the
  admin config API
- Makefile: `lint` removed from the `all` target (still available as
  `make lint`); `-buildvcs=false` added to `build` and `build-static`
- Docker `entrypoint.sh`: MariaDB startup timeout increased from 30 s to
  60 s with an initial 3 s pre-wait

### Fixed

- Disk space not reclaimed after cleanup — `DELETE` leaves InnoDB pages
  allocated internally; `DROP PARTITION` frees the `.ibd` file immediately

---

## [v0.5.2] - 2026-04-03

### Fixed

- **Active preset lost when time window crosses a DST boundary in Chrome** —
  Chrome parses `"YYYY-MM-DDTHH:mm"` strings (produced by `toISOString()`)
  as local time when no trailing `Z` is present; Firefox treats the same
  format as UTC. Two places were affected:
  - `LogsView` — `shiftTimeWindow()` computed the window size as
    `end - start` using local-time parsing. When the window straddled a
    DST transition (e.g. CET→CEST on 2026-03-29 or CEST→CET on
    2025-10-26) the UTC offset differed between the two strings, making
    the window 1 h short. Every subsequent shift moved by the wrong
    duration, permanently corrupting the stored dates.
  - `FilterPanel` — `activeDur` computed the diff the same way. For the
    single window position that spanned a DST boundary the diff was 1 h
    off, the 60 s tolerance check failed, and the active-preset highlight
    disappeared for exactly that one step.
  - Fix: append `'Z'` before parsing in both places, enforcing UTC
    interpretation consistently across all browsers.

---

## [v0.5.1] - 2026-03-31

### Fixed

- **`/health` response missing server-side defaults** — `NewHealthHandler` was
  called without the config parameter, so the `defaults` object introduced in
  v0.5.0 was never populated in the health response. Fixed by passing `cfg` to
  `NewHealthHandler` in `server.go`.

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
search highlighting. Several backend performance improvements are included.

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

**API / Backend**
- **`db_total` field in `/api/logs` response** — total entry count in `SystemEvents`
  (no filter applied), returned alongside the existing filtered `total` on every request
- **`db_total` and `oldest_entry` in `GET /api/meta` response** — database-level stats
  (total row count, timestamp of oldest entry) without a separate API call
- **`internal/database/cache.go`** — TTL cache (60 s) for `QueryDistinctValues` results;
  reduces redundant `SELECT DISTINCT` queries on poll-heavy setups
- **Toolbar DB total display** — log viewer toolbar shows
  `{filtered} entries · {db_total} total in DB` when a filter is active and the counts differ
- **Date range limit removed** — the 90-day hard cap on `start_date`/`end_date`
  has been dropped; arbitrarily large time windows are now accepted

### Changed

- **`/api/logs` runs three DB queries in parallel** (`CountLogs`, `QueryLogs`,
  `TotalCount`) via goroutines + `sync.WaitGroup`, reducing per-request latency
- **`QueryDistinctValues` results are cached** for 60 s per unique
  column + filter combination; subsequent identical meta requests are served from memory
- **Date range limit removed** — the 90-day hard cap on `start_date`/`end_date`
  has been dropped; arbitrarily large time windows are now accepted

### Fixed

- **Args slice mutation** — `QueryLogs` previously appended `LIMIT`/`OFFSET` directly
  to the caller's `args` slice; replaced with an internal copy to prevent data races
  when running queries concurrently in `QueryLogsWithTotal`
- **German strings in LogTable toolbar** — `ctrl-btn` title attributes were in German;
  replaced with English strings (consistent with UI language policy)
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

## [v0.4.2] - 2026-02-24

### Fixed

- **Frontend build failure after v0.4.0** — `frontend/package.json` and
  `package-lock.json` were missing entries added during the v0.4.0 feature
  work, causing `npm install` to fail on a clean checkout. Lockfile
  regenerated and committed.

---

## [v0.4.1] - 2026-02-24

### Fixed

- **Setup wizard unreachable on non-default port** — `config.Load()` ignored
  `RSYSLOX_PORT` in setup mode and always bound to the hardcoded default
  port 8000, making the wizard inaccessible when the installer had prompted
  for a different port. Fixed by reading `RSYSLOX_PORT` from the environment
  when no `config.toml` exists. The systemd unit now sets
  `Environment=RSYSLOX_PORT=<chosen-port>` so the value survives service
  restarts before setup is complete.
- **Installer showed wrong setup URL** — the post-install summary and
  browser-open logic used `DEFAULT_PORT` instead of the user-chosen
  `SETUP_PORT`, displaying `http://localhost:8000` even when a different port
  was selected. Fixed in `scripts/install.sh`.

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

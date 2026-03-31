# Configuration

All rsyslox settings are managed through the **Admin panel** at `/admin`. No manual config file editing is required.

## Admin Panel

Navigate to `http://<host>:8000/admin` and log in with your admin password.

### Server

| Setting | Description | Default |
|---|---|---|
| Host | Network interface to listen on. `0.0.0.0` binds all interfaces. | `0.0.0.0` |
| Port | TCP port rsyslox listens on | `8000` |
| Allowed Origins (CORS) | Comma-separated browser origins allowed to call the API. Use `*` to allow all. | `*` |
| Enable SSL / HTTPS | Serve HTTPS instead of HTTP — see [SSL section](#ssl--tls) below | off |

?> **Host and Port** changes require a server restart to take effect. Use the **Restart** button that appears in the banner after saving.

#### Default Values

These settings define the starting values applied to **new browser sessions** that have not yet configured their own preferences. Once a user changes a preference in **Admin → Preferences**, their personal setting takes over and the server default no longer applies for that user.

| Setting | Description | Default |
|---|---|---|
| Default Time Range | Initial time range shown on first visit | `24h` |
| Default Auto-Refresh | Auto-refresh interval in seconds | `30` |
| Default Language | UI language | `en` |
| Default Font Size | Table font size | `medium` |
| Default Time Format | 12-hour or 24-hour clock | `24h` |

#### SSL / TLS

When SSL is enabled, a certificate management section appears below the main form.

**Self-signed certificate** — generates an ECDSA P-256 certificate valid for 10 years. Suitable for internal use or testing. If no certificate exists when the server starts with `use_ssl = true`, one is generated automatically.

**Custom certificate** — upload your own `.pem` / `.crt` certificate and private key.

### Database

Editable form for the database connection. All changes require a server restart.

| Setting | Description |
|---|---|
| Host | Database server hostname or IP |
| Port | Database TCP port |
| Database name | The MySQL/MariaDB database (e.g. `Syslog`) |
| User | Database user |
| Password | Leave blank to keep the current password |

The password is AES-GCM encrypted before being written to `config.toml`.

#### Log Cleanup

Cleanup settings are part of the Database tab. The cleanup service monitors disk usage and deletes the oldest log entries when the threshold is exceeded. **Changes apply immediately — no restart needed.**

!> The cleanup service checks disk usage on the **local filesystem**. It only works correctly if the database runs on the same host as rsyslox.

| Setting | Description | Default |
|---|---|---|
| Enabled | Toggle the cleanup service | off |
| Disk path | Mount point to monitor (usually the MySQL data directory) | `/var/lib/mysql` |
| Threshold % | Delete entries when disk usage exceeds this | 85 % |
| Batch size | Rows deleted per cleanup run | 1 000 |
| Interval | Seconds between disk checks | 900 |

A live **disk usage bar** shows the current utilisation of the configured path.

### API Keys

Named, revocable read-only API keys for external tools. Keys are shown in plaintext **once** at creation time. Pass a key via:

```
X-API-Key: <plaintext key>
```

### Preferences

Browser-persisted settings stored in `localStorage`. Apply instantly without restart and are independent per browser. The server provides defaults (see [Default Values](#default-values) above) for keys that have not yet been set by the user.

| Setting | Options | Default (server) |
|---|---|---|
| Language | English, Deutsch | `en` |
| Time format | 24-hour, 12-hour | `24h` |
| Font size | Small (14 px), Medium (16 px), Large (18 px) | `medium` |
| Auto-refresh interval | 5–300 s | `30` |
| Default time range | 15m–30d | `24h` |

### Server Restart

Some settings (host, port, SSL, database connection) require a server restart to take effect. After saving, a yellow banner appears at the top of the Admin panel. Click **Restart Now** to restart the server in-place via `syscall.Exec`. The browser polls `/health` and reloads automatically once the server is back online.

---

## Configuration File Reference

`/etc/rsyslox/config.toml` is written by the setup wizard and updated by the Admin panel. Manual editing is not required.

```toml
[server]
host                  = "0.0.0.0"
port                  = 8000
use_ssl               = false
ssl_cert              = "/etc/rsyslox/certs/cert.pem"
ssl_key               = "/etc/rsyslox/certs/key.pem"
allowed_origins       = ["*"]

# Server-side defaults for new browser sessions
auto_refresh_interval = 30     # seconds
default_time_range    = "24h"  # "15m"|"1h"|"6h"|"24h"|"7d"|"30d"
default_language      = "en"   # "en"|"de"
default_font_size     = "medium"  # "small"|"medium"|"large"
default_time_format   = "24h"  # "24h"|"12h"

[database]
host     = "localhost"
port     = 3306
name     = "Syslog"
user     = "rsyslox"
password = "enc:<base64>"   # AES-GCM encrypted by setup wizard

[auth]
admin_password_hash = "$2a$12$..."   # bcrypt hash

[[auth.read_only_keys]]
name     = "monitoring"
key_hash = "<sha256 hex>"

[cleanup]
enabled           = false
disk_path         = "/var/lib/mysql"
threshold_percent = 85.0
batch_size        = 1000
interval          = "15m"
```

### Security Model

| Value | Storage |
|---|---|
| Database password | AES-GCM encrypted; key derived from `/etc/machine-id` — not portable between machines |
| Admin password | bcrypt hash (cost 12) |
| API key plaintext | Never stored; only SHA-256 hex hash written to disk |
| Config file | Mode `0640` — readable by `root` and group `rsyslox` only |

---

## Next Steps

- [Quick Start Guide](quick-start.md)
- [Deployment Guide](../guides/deployment.md)
- [Security Guide](../guides/security.md)
- [Cleanup Guide](../guides/cleanup.md)

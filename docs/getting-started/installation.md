# Installation

## Prerequisites

### System Requirements

- Linux (Ubuntu, Debian, CentOS, RHEL) — x86_64 or ARM64
- MySQL 5.7+ or MariaDB 10.3+ with rsyslog data
- systemd
- Root access for the installer

### rsyslog MySQL Setup

rsyslox reads from the `SystemEvents` table that rsyslog-mysql populates. If rsyslog is not yet configured with MySQL:

```bash
# Ubuntu / Debian
sudo apt-get install rsyslog-mysql

# CentOS / RHEL
sudo yum install rsyslog-mysql
```

`/etc/rsyslog.d/mysql.conf`:
```
module(load="ommysql")
action(type="ommysql" server="localhost" db="Syslog" uid="rsyslog" pwd="yourpassword")
```

```bash
sudo systemctl restart rsyslog
```

## Install

Download the latest release binary and run the installer:

```bash
wget https://github.com/phil-bot/rsyslox/releases/latest/download/rsyslox-linux-amd64
chmod +x rsyslox-linux-amd64
sudo ./install.sh
```

The installer:

1. Creates a dedicated system user and group `rsyslox`
2. Copies the binary to `/opt/rsyslox/rsyslox`
3. Installs and enables a hardened systemd service
4. Starts the service

At the end the installer prints the setup wizard URL. Open it in your browser to finish configuration.

## Setup Wizard

On first start, rsyslox has no configuration and serves a setup wizard on **`http://<server-ip>:8000`** — reachable from any machine on the network until setup is complete.

Fill in:

- **Database** — host, port, database name, user, password
- **Admin password** — minimum 12 characters (stored as bcrypt hash)
- **Server** — bind host, port, optional CORS origins

Click **Save** — rsyslox writes `/etc/rsyslox/config.toml` and immediately starts serving the log viewer. No restart is required.

## Verify

```bash
# Service status
sudo systemctl status rsyslox

# Health check
curl http://localhost:8000/health
```

Expected health response:
```json
{"status": "healthy", "database": "connected", "version": "v0.5.2", "timestamp": "..."}
```

Open `http://<your-host>:8000` in your browser — the log viewer should load.

## Uninstall

```bash
sudo ./install.sh --uninstall
```

Stops and removes the service and binary. Configuration at `/etc/rsyslox/` is kept intentionally — remove it manually if no longer needed.

## Install from Source

```bash
# Prerequisites: Go 1.21+, Node.js 18+, make

git clone https://github.com/phil-bot/rsyslox.git
cd rsyslox

make all           # builds frontend + downloads Redoc + builds binary
sudo ./install.sh
```

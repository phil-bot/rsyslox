<div align="center">
  <img src="https://rsyslox.grothu.net/rsyslox_light.svg" alt="rsyslox"/>
</div>

# rsyslox

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/phil-bot/rsyslox)](https://github.com/phil-bot/rsyslox/releases)

rsyslox is a self-hosted syslog viewer for rsyslog data stored in MySQL/MariaDB. It provides a full-featured web UI for browsing, filtering, and exporting log entries — alongside a clean REST API for integrating syslog data into dashboards, monitoring tools, or custom scripts.

A single binary embeds the frontend, API documentation, and all required assets. No config files need to be edited manually.

## 🚀 Quick Start

```bash
# Download latest release
wget https://github.com/phil-bot/rsyslox/releases/latest/download/rsyslox-linux-amd64
chmod +x rsyslox-linux-amd64

# Install (creates system user, registers systemd service)
sudo ./install.sh

# Open the setup wizard in your browser
# http://localhost:8000
```

The setup wizard walks you through database credentials, admin password, and server settings. rsyslox writes its own configuration — no manual config file editing required.

## ✨ Features

**Web UI**
- Log viewer with real-time filtering: time range, severity, facility, host, tag, message search
- Auto-refresh with configurable interval
- Dark / light theme, English / Deutsch, adjustable font size, 12h / 24h clock
- Multi-row selection, CSV / JSON export
- Detail panel with full field view and raw JSON

**Admin Panel**
- Server settings, database info, log cleanup configuration
- Named, revocable read-only API keys
- Browser-persisted preferences (no restart required)

**API**
- REST API with filtering, pagination, and metadata queries
- Two authentication methods: admin session token or read-only API key
- Interactive API documentation at `/docs`

## 📸 Screenshots

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="rsyslox-logs-dark.png">
  <img src="rsyslox-logs.png" alt="rsyslox log viewer" width="100%">
</picture>

## 📖 Documentation

**[https://rsyslox.grothu.net](https://rsyslox.grothu.net)**

## 📄 License

MIT License — see [LICENSE](https://github.com/phil-bot/rsyslox/blob/main/LICENSE) for details.

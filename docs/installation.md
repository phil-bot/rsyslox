# Installation Guide

[← Back to overview](index.md)

This guide describes all installation methods for rsyslog REST API.

## 📋 Prerequisites

### System requirements

- **Operating system:** Linux (Ubuntu, Debian, CentOS, RHEL, etc.)
- **Architecture:** x86_64 or ARM64
- **rsyslog:** Installed with MySQL/MariaDB support
- **Database:** MySQL 5.7+ or MariaDB 10.3+
- **Memory:** Minimum 256 MB RAM
- **Network:** Port 8000 (or configurable)

### rsyslog MySQL Setup

Before installation, rsyslog must be configured correctly:

```bash
# install rsyslog-mysql (if not already installed)
sudo apt-get install rsyslog-mysql # Ubuntu/Debian
sudo yum install rsyslog-mysql # CentOS/RHEL

# create rsyslog MySQL config
sudo nano /etc/rsyslog.d/mysql.conf
```

Contents of `/etc/rsyslog.d/mysql.conf`:
```
module(load="ommysql")
action(type="ommysql" server="localhost" db="Syslog" uid="rsyslog" pwd="password")
```

```bash
# restart rsyslog
sudo systemctl restart rsyslog
```

## 🚀 Installation methods

### Option 1: Binary Installation (Recommended)

The easiest and fastest method.

#### Download Latest Release

```bash
# Linux x86_64
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-amd64

# Linux ARM64
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-arm64

# Make executable
chmod +x rsyslog-rest-api-linux-amd64

# Move to /usr/local/bin
sudo mv rsyslog-rest-api-linux-amd64 /usr/local/bin/rsyslog-rest-api
```

#### Verify checksums (optional)

```bash
# Download checksums
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/SHA256SUMS

# Verify
sha256sum -c SHA256SUMS
```

#### Create configuration

```bash
# Create directory
sudo mkdir -p /opt/rsyslog-rest-api

# Create configuration
sudo nano /opt/rsyslog-rest-api/.env
```

Minimal configuration:
```bash
# API Key (IMPORTANT!)
API_KEY=your-secret-key-here

# Database
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=your-database-password

# Server
SERVER_PORT=8000
```

**generate API key:**
```bash
openssl rand -hex 32
```

#### Testing

```bash
# Start API (foreground for testing)
cd /opt/rsyslog-rest-api
rsyslog-rest-api

# Test in another terminal
curl http://localhost:8000/health
```

Continue to: [Production Setup](deployment.md)

---

### Option 2: Offline Package Installation

For systems without Internet access.

#### Download Package

```bash
# Download package (on systems with Internet)
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-v0.2.2-offline-linux-amd64.tar.gz

# Transfer to target system (USB, SCP, etc.)
scp rsyslog-rest-api-v0.2.2-offline-linux-amd64.tar.gz user@server:/tmp/
```

#### Installation

```bash
# Unpack
tar xzf rsyslog-rest-api-v0.2.2-offline-linux-amd64.tar.gz
cd package

# Execute install script
sudo ./install.sh
```

The script:
- ✅ Copies binary to `/opt/rsyslog-rest-api/`
- ✅ Creates `.env` template
- installs systemd service
- ✅ Reloads systemd daemon

#### Configuration

```bash
# Edit config (IMPORTANT!)
sudo nano /opt/rsyslog-rest-api/.env

# Set at least:
# - API_KEY
# - DB_HOST, DB_NAME, DB_USER, DB_PASS

# Set permissions
sudo chmod 600 /opt/rsyslog-rest-api/.env
```

#### Start service

```bash
sudo systemctl enable --now rsyslog-rest-api
sudo systemctl status rsyslog-rest-api
```

---

### Option 3: Installation from Source

For developers or special requirements.

#### Requirements

- Go 1.21 or higher
- git
- make

#### Clone Repository

```bash
git clone https://github.com/phil-bot/rsyslog-rest-api.git
cd rsyslog-rest-api
```

#### Build & Install

```bash
# Build static binary
make build-static

# Install (requires sudo)
sudo make install
```

The Makefile:
- Compiles static binary
- Installs to `/opt/rsyslog-rest-api/`
- Copies `.env.example` → `.env` (if not present)
- Installs systemd service
- Reloads systemd

#### Configuration

```bash
# Edit config
sudo nano /opt/rsyslog-rest-api/.env

# Set at least:
# - API_KEY
# - DB_HOST, DB_NAME, DB_USER, DB_PASS

# Secure permissions
sudo chmod 600 /opt/rsyslog-rest-api/.env
```

#### Start service

```bash
sudo systemctl enable --now rsyslog-rest-api
sudo systemctl status rsyslog-rest-api
```

---

## ✅ Verification

After each installation method:

### 1. check service status

```bash
sudo systemctl status rsyslog-rest-api
```

Expected output:
```
● rsyslog-rest-api.service - rsyslog REST API
     Loaded: loaded (/etc/systemd/system/rsyslog-rest-api.service; enabled)
     Active: active (running) since ...
```

### 2. health check

```bash
curl http://localhost:8000/health
```

Expected response:
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-02-09T10:30:00Z"
}
```

### 3rd API test with authentication

```bash
# Get API key from .env
API_KEY=$(sudo grep "^API_KEY=" /opt/rsyslog-rest-api/.env | cut -d'=' -f2)

# Retrieve logs
curl -H "X-API-Key: $API_KEY" "http://localhost:8000/logs?limit=5"
```

### 4. check logs

```bash
# API logs
sudo journalctl -u rsyslog-rest-api -n 50

# If error:
sudo journalctl -u rsyslog-rest-api -f
```

---

## 🔧 Post-installation

After successful installation:

1. **Optimize configuration:** → [Configuration Guide](configuration.md)
2. **Set up SSL/TLS:** → [Security Guide](security.md#ssltls)
3. **Deploy productively:** → [Deployment Guide](deployment.md)
4. **Set up monitoring:** → [Deployment: Monitoring](deployment.md#monitoring)

## 🆘 Troubleshooting

In case of problems, see [Troubleshooting Guide](troubleshooting.md).

Frequent installation problems:

- **Binary not found:** Check path, check permissions
- **Database connection failed:** Check credentials in `.env`
- **Permission denied:** `.env` permissions (`chmod 600`)
- **Service won't start:** Check logs with `journalctl`

---

## 🔄 Upgrade

### Upgrade to new version

```bash
# Download new version
wget https://github.com/phil-bot/rsyslog-rest-api/releases/download/v0.X.X/rsyslog-rest-api-linux-amd64

# Stop service
sudo systemctl stop rsyslog-rest-api

# Replace binary
sudo mv rsyslog-rest-api-linux-amd64 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api

# Start service
sudo systemctl start rsyslog-rest-api

# Check status
sudo systemctl status rsyslog-rest-api
```

**Important:** `.env` configuration is preserved!

See [Changelog](changelog.md) for breaking changes.

---

## 🗑️ Uninstallation

```bash
# Stop and disable service
sudo systemctl stop rsyslog-rest-api
sudo systemctl disable rsyslog-rest-api

# Remove service file
sudo rm /etc/systemd/system/rsyslog-rest-api.service
sudo systemctl daemon-reload

# Remove installation
sudo rm -rf /opt/rsyslog-rest-api

# Remove binary (if in /usr/local/bin)
sudo rm /usr/local/bin/rsyslog-rest-api
```

---

[← Back to overview](index.md) | [Next to Configuration →](configuration.md)

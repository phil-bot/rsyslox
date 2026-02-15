# Configuration Guide

[← Back to overview](index.md) | [← Installation](installation.md)

Complete configuration reference for rsyslog REST API.

## 📁 Configuration file

The API uses an `.env` file for configuration:

**location:** `/opt/rsyslog-rest-api/.env`

**Permissions:**
```bash
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
```

## 🔧 Complete configuration

### Template (.env.example)

```bash
# ===========================================
# rsyslog REST API Configuration
# ===========================================

# -------------------------------------------
# API Security (REQUIRED for production!)
# -------------------------------------------
# Generate with: openssl rand -hex 32
API_KEY=

# -------------------------------------------
# Server Configuration
# -------------------------------------------
SERVER_HOST=0.0.0.0
SERVER_PORT=8000

# -------------------------------------------
# SSL/TLS (Optional but recommended)
# -------------------------------------------
USE_SSL=false
SSL_CERTFILE=/opt/rsyslog-rest-api/certs/cert.pem
SSL_KEYFILE=/opt/rsyslog-rest-api/certs/key.pem

# -------------------------------------------
# CORS (Cross-Origin Resource Sharing)
# -------------------------------------------
# Comma-separated origins, or * for all
ALLOWED_ORIGINS=*

# -------------------------------------------
# Database Connection (RECOMMENDED)
# -------------------------------------------
# Direct database configuration (preferred method)
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=

# -------------------------------------------
# rsyslog Config Fallback (Legacy)
# -------------------------------------------
# Only used if DB_* variables above are not set
# WARNING: File must be readable by API user!
RSYSLOG_CONFIG_PATH=/etc/rsyslog.d/mysql.conf
```

---

## 🔐 API Security

### API_KEY

**Description:** API key for authentication
**Required:** Yes (Production) / No (Development)
**Default:** Empty (no authentication)

#### Generate API key:

```bash
# Generate strong 32-byte key
openssl rand -hex 32
```

#### Examples:

```bash
# Production (ALWAYS with key!)
API_KEY=a3d7f8c9e2b4a6d8f9c3e7b1a5d9f4c8e2b7a6d3f9c8e1b4a7d2f6c9e3b8a5d1

# Development (optional without key)
API_KEY=

# Test (simple key)
API_KEY=test123456789
```

#### Usage:

```bash
# Requests WITH API Key
curl -H "X-API-Key: YOUR_API_KEY" "http://localhost:8000/logs"

# Requests WITHOUT API key (if API_KEY is empty)
curl "http://localhost:8000/logs"
```

**⚠️ IMPORTANT:** ALWAYS use an API key in Production!

---

## 🌐 Server Configuration

### SERVER_HOST

**Description:** IP address on which the server is listening
**Default:** `0.0.0.0` (all interfaces)
**Type:** IP Address

#### Examples:

```bash
# All interfaces (default)
SERVER_HOST=0.0.0.0

# Only localhost (secure, local access only)
SERVER_HOST=127.0.0.1

# Specific interface
SERVER_HOST=192.168.1.100
```

### SERVER_PORT

**Description:** Port on which the server is listening
**Default:** `8000`
**Type:** Integer (1-65535)

#### Examples:

```bash
# Standard
SERVER_PORT=8000

# Alternative port
SERVER_PORT=8080

# Privileged port (requires root)
SERVER_PORT=80
```

---

## 🔒 SSL/TLS Configuration

### USE_SSL

**Description:** Activate SSL/TLS
**Default:** `false`
**Type:** Boolean (`true`/`false`)

### SSL_CERTFILE

**Description:** Path to the SSL certificate
**Default:** `/opt/rsyslog-rest-api/certs/cert.pem`
**Type:** File Path

### SSL_KEYFILE

**Description:** Path to the SSL private key
**Default:** `/opt/rsyslog-rest-api/certs/key.pem`
**Type:** File Path

#### Set up SSL:

```bash
#1 Create directory
sudo mkdir -p /opt/rsyslog-rest-api/certs

# 2. self-signed certificate (development)
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /opt/rsyslog-rest-api/certs/key.pem \
  -out /opt/rsyslog-rest-api/certs/cert.pem

# 3. permissions
sudo chmod 600 /opt/rsyslog-rest-api/certs/*.pem

# 4. Customize .env
sudo nano /opt/rsyslog-rest-api/.env
```

#### SSL configuration:

```bash
USE_SSL=true
SSL_CERTFILE=/opt/rsyslog-rest-api/certs/cert.pem
SSL_KEYFILE=/opt/rsyslog-rest-api/certs/key.pem
```

#### Production SSL:

For Production use **Let's Encrypt** or a commercial certificate:

```bash
USE_SSL=true
SSL_CERTFILE=/etc/letsencrypt/live/api.example.com/fullchain.pem
SSL_KEYFILE=/etc/letsencrypt/live/api.example.com/privkey.pem
```

→ More details: [Security Guide: SSL/TLS](security.md#ssltls)

---

## 🌍 CORS Configuration

### ALLOWED_ORIGINS

**Description:** Allowed origins for cross-origin requests
**Default:** `*` (all origins)
**Type:** Comma-separated list

#### Examples:

```bash
# Allow all origins (Development)
ALLOWED_ORIGINS=*

# Specific domains (Production)
ALLOWED_ORIGINS=https://dashboard.example.com,https://monitoring.example.com

# Only one domain
ALLOWED_ORIGINS=https://app.example.com

# Several with different schemas
ALLOWED_ORIGINS=http://localhost:3000,https://app.example.com,https://api.example.com
```

**⚠️ Security:** Never use `*` in Production!

---

## 💾 Database Configuration

There are two methods for configuring the database connection:

### Method 1: Environment Variables (Recommended)

Direct configuration via environment variables.

#### DB_HOST

**Description:** MySQL/MariaDB host name
**Default:** `localhost`
**Type:** Hostname or IP

#### DB_NAME

**Description:** Database name
**Default:** `Syslog`
**Type:** String

#### DB_USER

**Description:** Database user
**Default:** `rsyslog`
**Type:** String

#### DB_PASS

**Description:** Database password
**Required:** Yes
**Type:** String

#### Example configuration:

```bash
# Default (localhost)
DB_HOST=localhost
DB_NAME=syslog
DB_USER=rsyslog
DB_PASS=secure-password-here

# Remote Database
DB_HOST=192.168.1.50
DB_NAME=Syslog
DB_USER=api_user
DB_PASS=very-secure-password

# With port
DB_HOST=db.example.com:3307
DB_NAME=LogDatabase
DB_USER=readonly_api
DB_PASS=another-secure-password
```

**✅ Advantages:**
- Simple and direct
- Secure (only in `.env`, not in rsyslog config)
- Recommended for all setups

---

### Method 2: rsyslog config file (fallback)

Read the credentials from the rsyslog configuration.

#### RSYSLOG_CONFIG_PATH

**Description:** Path to the rsyslog MySQL configuration
**Default:** `/etc/rsyslog.d/mysql.conf`
**Type:** File Path

**Is only used if DB_* variables are NOT set!

#### Example:

```bash
# Only if DB_HOST, DB_NAME, DB_USER, DB_PASS are NOT set
RSYSLOG_CONFIG_PATH=/etc/rsyslog.d/mysql.conf
```

**⚠️ Disadvantages:**
- File must be readable for API users
- Security risk (other programs could read)
- Not recommended

**Recommendation:** Always use `DB_*` variables!

---

## 🎯 Configuration examples

### Minimal Setup (Development)

```bash
# Minimal for local testing
API_KEY=
DB_HOST=localhost
DB_NAME=syslog
DB_USER=rsyslog
DB_PASS=password
```

### Production Setup (Recommended)

```bash
# Production-ready configuration
API_KEY=a3d7f8c9e2b4a6d8f9c3e7b1a5d9f4c8e2b7a6d3f9c8e1b4a7d2f6c9e3b8a5d1
SERVER_HOST=0.0.0.0
SERVER_PORT=8000
USE_SSL=true
SSL_CERTFILE=/etc/letsencrypt/live/api.example.com/fullchain.pem
SSL_KEYFILE=/etc/letsencrypt/live/api.example.com/privkey.pem
ALLOWED_ORIGINS=https://dashboard.example.com
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog_api
DB_PASS=very-secure-database-password-here
```

### Behind Reverse Proxy

```bash
# API behind nginx/Apache
API_KEY=your-api-key
SERVER_HOST=127.0.0.1
SERVER_PORT=8000
USE_SSL=false # SSL terminated at the proxy
ALLOWED_ORIGINS=https://yourdomain.com
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=db-password
```

### High Security Setup

```bash
# Maximum Security
API_KEY=$(openssl rand -hex 32) # Regenerates at every startup
SERVER_HOST=127.0.0.1 # Only localhost
SERVER_PORT=8000
USE_SSL=true
SSL_CERTFILE=/etc/ssl/private/api.crt
SSL_KEYFILE=/etc/ssl/private/api.key
ALLOWED_ORIGINS=https://secure-dashboard.example.com
DB_HOST=127.0.0.1 # Local DB
DB_NAME=Syslog
DB_USER=rsyslog_readonly # READ-ONLY user!
DB_PASS=strong-complex-password-with-special-chars
```

---

## ✅ Validate configuration

After changes:

```bash
# 1. check config syntax
sudo cat /opt/rsyslog-rest-api/.env

# 2. restart service
sudo systemctl restart rsyslog-rest-api

# 3. check status
sudo systemctl status rsyslog-rest-api

# 4. check logs
sudo journalctl -u rsyslog-rest-api -n 50

# 5. health check
curl http://localhost:8000/health

# 6. test with API key
curl -H "X-API-Key: YOUR_KEY" "http://localhost:8000/logs?limit=1"
```

---

## 🔄 Apply changes

After configuration changes:

```bash
# Restart service
sudo systemctl restart rsyslog-rest-api

# Alternative: Reload (if supported)
sudo systemctl reload rsyslog-rest-api
```

---

## 🆘 Troubleshooting

### Config is not loaded

```bash
# Check permissions
ls -la /opt/rsyslog-rest-api/.env

# Should be: -rw------- root root
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
```

### Database Connection Failed

```bash
#1 Test credentials
mysql -h DB_HOST -u DB_USER -pDB_PASS DB_NAME

# 2. check network
ping DB_HOST

# 3. mysql port (default: 3306)
telnet DB_HOST 3306
```

### API key does not work

```bash
# Read key from .env
grep "^API_KEY=" /opt/rsyslog-rest-api/.env

# Test with correct key
curl -H "X-API-Key: EXACT_KEY_FROM_ENV" "http://localhost:8000/logs"
```

More: [Troubleshooting Guide](troubleshooting.md)

---

[← Back to overview](index.md) | [Next to API Reference →](api-reference.md)

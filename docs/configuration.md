# Configuration Guide

[← Zurück zur Übersicht](index.md) | [← Installation](installation.md)

Vollständige Konfigurationsreferenz für rsyslog REST API.

## 📁 Konfigurationsdatei

Die API verwendet eine `.env` Datei zur Konfiguration:

**Speicherort:** `/opt/rsyslog-rest-api/.env`

**Permissions:**
```bash
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
```

## 🔧 Vollständige Konfiguration

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

**Beschreibung:** API-Key für Authentifizierung  
**Required:** Ja (Production) / Nein (Development)  
**Default:** Leer (keine Authentifizierung)

#### API-Key generieren:

```bash
# Starken 32-Byte Key generieren
openssl rand -hex 32
```

#### Beispiele:

```bash
# Production (IMMER mit Key!)
API_KEY=a3d7f8c9e2b4a6d8f9c3e7b1a5d9f4c8e2b7a6d3f9c8e1b4a7d2f6c9e3b8a5d1

# Development (optional ohne Key)
API_KEY=

# Test (einfacher Key)
API_KEY=test123456789
```

#### Verwendung:

```bash
# Requests MIT API Key
curl -H "X-API-Key: YOUR_API_KEY" "http://localhost:8000/logs"

# Requests OHNE API Key (wenn API_KEY leer)
curl "http://localhost:8000/logs"
```

**⚠️ WICHTIG:** In Production IMMER einen API-Key verwenden!

---

## 🌐 Server Configuration

### SERVER_HOST

**Beschreibung:** IP-Adresse auf der der Server lauscht  
**Default:** `0.0.0.0` (alle Interfaces)  
**Type:** IP Address

#### Beispiele:

```bash
# Alle Interfaces (Standard)
SERVER_HOST=0.0.0.0

# Nur localhost (sicherer, nur lokaler Zugriff)
SERVER_HOST=127.0.0.1

# Spezifisches Interface
SERVER_HOST=192.168.1.100
```

### SERVER_PORT

**Beschreibung:** Port auf dem der Server lauscht  
**Default:** `8000`  
**Type:** Integer (1-65535)

#### Beispiele:

```bash
# Standard
SERVER_PORT=8000

# Alternativer Port
SERVER_PORT=8080

# Privilegierter Port (benötigt root)
SERVER_PORT=80
```

---

## 🔒 SSL/TLS Configuration

### USE_SSL

**Beschreibung:** SSL/TLS aktivieren  
**Default:** `false`  
**Type:** Boolean (`true`/`false`)

### SSL_CERTFILE

**Beschreibung:** Pfad zum SSL-Zertifikat  
**Default:** `/opt/rsyslog-rest-api/certs/cert.pem`  
**Type:** File Path

### SSL_KEYFILE

**Beschreibung:** Pfad zum SSL-Private-Key  
**Default:** `/opt/rsyslog-rest-api/certs/key.pem`  
**Type:** File Path

#### SSL einrichten:

```bash
# 1. Verzeichnis erstellen
sudo mkdir -p /opt/rsyslog-rest-api/certs

# 2. Self-Signed Certificate (Development)
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /opt/rsyslog-rest-api/certs/key.pem \
  -out /opt/rsyslog-rest-api/certs/cert.pem

# 3. Permissions
sudo chmod 600 /opt/rsyslog-rest-api/certs/*.pem

# 4. .env anpassen
sudo nano /opt/rsyslog-rest-api/.env
```

#### SSL Konfiguration:

```bash
USE_SSL=true
SSL_CERTFILE=/opt/rsyslog-rest-api/certs/cert.pem
SSL_KEYFILE=/opt/rsyslog-rest-api/certs/key.pem
```

#### Production SSL:

Für Production verwende **Let's Encrypt** oder ein kommerzielles Zertifikat:

```bash
USE_SSL=true
SSL_CERTFILE=/etc/letsencrypt/live/api.example.com/fullchain.pem
SSL_KEYFILE=/etc/letsencrypt/live/api.example.com/privkey.pem
```

→ Mehr Details: [Security Guide: SSL/TLS](security.md#ssltls)

---

## 🌍 CORS Configuration

### ALLOWED_ORIGINS

**Beschreibung:** Erlaubte Origins für Cross-Origin Requests  
**Default:** `*` (alle Origins)  
**Type:** Comma-separated list

#### Beispiele:

```bash
# Alle Origins erlauben (Development)
ALLOWED_ORIGINS=*

# Spezifische Domains (Production)
ALLOWED_ORIGINS=https://dashboard.example.com,https://monitoring.example.com

# Nur eine Domain
ALLOWED_ORIGINS=https://app.example.com

# Mehrere mit verschiedenen Schemas
ALLOWED_ORIGINS=http://localhost:3000,https://app.example.com,https://api.example.com
```

**⚠️ Sicherheit:** In Production niemals `*` verwenden!

---

## 💾 Database Configuration

Es gibt zwei Methoden, die Datenbankverbindung zu konfigurieren:

### Methode 1: Environment Variables (Empfohlen)

Direkte Konfiguration über Umgebungsvariablen.

#### DB_HOST

**Beschreibung:** MySQL/MariaDB Hostname  
**Default:** `localhost`  
**Type:** Hostname oder IP

#### DB_NAME

**Beschreibung:** Datenbank-Name  
**Default:** `Syslog`  
**Type:** String

#### DB_USER

**Beschreibung:** Datenbank-Benutzer  
**Default:** `rsyslog`  
**Type:** String

#### DB_PASS

**Beschreibung:** Datenbank-Passwort  
**Required:** Ja  
**Type:** String

#### Beispiel-Konfiguration:

```bash
# Standard (localhost)
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=secure-password-here

# Remote Database
DB_HOST=192.168.1.50
DB_NAME=Syslog
DB_USER=api_user
DB_PASS=very-secure-password

# Mit Port
DB_HOST=db.example.com:3307
DB_NAME=LogDatabase
DB_USER=readonly_api
DB_PASS=another-secure-password
```

**✅ Vorteile:**
- Einfach und direkt
- Sicher (nur in `.env`, nicht in rsyslog config)
- Empfohlen für alle Setups

---

### Methode 2: rsyslog Config File (Fallback)

Auslesen der Credentials aus rsyslog-Konfiguration.

#### RSYSLOG_CONFIG_PATH

**Beschreibung:** Pfad zur rsyslog MySQL-Konfiguration  
**Default:** `/etc/rsyslog.d/mysql.conf`  
**Type:** File Path

**Wird nur verwendet wenn DB_* Variablen NICHT gesetzt sind!**

#### Beispiel:

```bash
# Nur wenn DB_HOST, DB_NAME, DB_USER, DB_PASS NICHT gesetzt
RSYSLOG_CONFIG_PATH=/etc/rsyslog.d/mysql.conf
```

**⚠️ Nachteile:**
- Datei muss lesbar sein für API-User
- Sicherheitsrisiko (andere Programme könnten lesen)
- Nicht empfohlen

**Empfehlung:** Verwende immer `DB_*` Variablen!

---

## 🎯 Konfigurationsbeispiele

### Minimal Setup (Development)

```bash
# Minimal für lokales Testen
API_KEY=
DB_HOST=localhost
DB_NAME=Syslog
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
# API hinter nginx/Apache
API_KEY=your-api-key
SERVER_HOST=127.0.0.1
SERVER_PORT=8000
USE_SSL=false  # SSL terminiert am Proxy
ALLOWED_ORIGINS=https://yourdomain.com
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=db-password
```

### High Security Setup

```bash
# Maximum Security
API_KEY=$(openssl rand -hex 32)  # Generiert bei jedem Start neu
SERVER_HOST=127.0.0.1  # Nur localhost
SERVER_PORT=8000
USE_SSL=true
SSL_CERTFILE=/etc/ssl/private/api.crt
SSL_KEYFILE=/etc/ssl/private/api.key
ALLOWED_ORIGINS=https://secure-dashboard.example.com
DB_HOST=127.0.0.1  # Lokale DB
DB_NAME=Syslog
DB_USER=rsyslog_readonly  # READ-ONLY User!
DB_PASS=strong-complex-password-with-special-chars
```

---

## ✅ Konfiguration validieren

Nach Änderungen:

```bash
# 1. Config-Syntax prüfen
sudo cat /opt/rsyslog-rest-api/.env

# 2. Service neustarten
sudo systemctl restart rsyslog-rest-api

# 3. Status prüfen
sudo systemctl status rsyslog-rest-api

# 4. Logs checken
sudo journalctl -u rsyslog-rest-api -n 50

# 5. Health Check
curl http://localhost:8000/health

# 6. Mit API Key testen
curl -H "X-API-Key: YOUR_KEY" "http://localhost:8000/logs?limit=1"
```

---

## 🔄 Änderungen übernehmen

Nach Konfigurationsänderungen:

```bash
# Service neustarten
sudo systemctl restart rsyslog-rest-api

# Alternative: Reload (falls unterstützt)
sudo systemctl reload rsyslog-rest-api
```

---

## 🆘 Troubleshooting

### Config wird nicht geladen

```bash
# Permissions prüfen
ls -la /opt/rsyslog-rest-api/.env

# Sollte sein: -rw------- root root
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
```

### Database Connection Failed

```bash
# 1. Credentials testen
mysql -h DB_HOST -u DB_USER -pDB_PASS DB_NAME

# 2. Netzwerk prüfen
ping DB_HOST

# 3. MySQL Port (Standard: 3306)
telnet DB_HOST 3306
```

### API Key funktioniert nicht

```bash
# Key aus .env auslesen
grep "^API_KEY=" /opt/rsyslog-rest-api/.env

# Mit richtigem Key testen
curl -H "X-API-Key: EXACT_KEY_FROM_ENV" "http://localhost:8000/logs"
```

Mehr: [Troubleshooting Guide](troubleshooting.md)

---

[← Zurück zur Übersicht](index.md) | [Weiter zu API Reference →](api-reference.md)

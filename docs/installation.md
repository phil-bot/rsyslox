# Installation Guide

[← Zurück zur Übersicht](index.md)

Dieser Guide beschreibt alle Installationsmethoden für rsyslog REST API.

## 📋 Voraussetzungen

### System Requirements

- **Betriebssystem:** Linux (Ubuntu, Debian, CentOS, RHEL, etc.)
- **Architektur:** x86_64 oder ARM64
- **rsyslog:** Mit MySQL/MariaDB-Support installiert
- **Datenbank:** MySQL 5.7+ oder MariaDB 10.3+
- **Speicher:** Minimal 256 MB RAM
- **Netzwerk:** Port 8000 (oder konfigurierbar)

### rsyslog MySQL Setup

Vor der Installation muss rsyslog korrekt konfiguriert sein:

```bash
# rsyslog-mysql installieren (falls noch nicht vorhanden)
sudo apt-get install rsyslog-mysql  # Ubuntu/Debian
sudo yum install rsyslog-mysql       # CentOS/RHEL

# rsyslog MySQL Config erstellen
sudo nano /etc/rsyslog.d/mysql.conf
```

Inhalt von `/etc/rsyslog.d/mysql.conf`:
```
module(load="ommysql")
action(type="ommysql" server="localhost" db="Syslog" uid="rsyslog" pwd="password")
```

```bash
# rsyslog neustarten
sudo systemctl restart rsyslog
```

## 🚀 Installationsmethoden

### Option 1: Binary Installation (Empfohlen)

Die einfachste und schnellste Methode.

#### Download Latest Release

```bash
# Linux x86_64
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-amd64

# Linux ARM64
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-arm64

# Executable machen
chmod +x rsyslog-rest-api-linux-amd64

# Nach /usr/local/bin verschieben
sudo mv rsyslog-rest-api-linux-amd64 /usr/local/bin/rsyslog-rest-api
```

#### Checksums verifizieren (Optional)

```bash
# Checksums herunterladen
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/SHA256SUMS

# Verifizieren
sha256sum -c SHA256SUMS
```

#### Konfiguration erstellen

```bash
# Verzeichnis erstellen
sudo mkdir -p /opt/rsyslog-rest-api

# Konfiguration erstellen
sudo nano /opt/rsyslog-rest-api/.env
```

Minimal-Konfiguration:
```bash
# API Key (WICHTIG!)
API_KEY=your-secret-key-here

# Datenbank
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=your-database-password

# Server
SERVER_PORT=8000
```

**API-Key generieren:**
```bash
openssl rand -hex 32
```

#### Testen

```bash
# API starten (Vordergrund zum Testen)
cd /opt/rsyslog-rest-api
rsyslog-rest-api

# In anderem Terminal testen
curl http://localhost:8000/health
```

Weiter zu: [Production Setup](deployment.md)

---

### Option 2: Offline Package Installation

Für Systeme ohne Internet-Zugang.

#### Download Package

```bash
# Package herunterladen (auf System mit Internet)
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-v0.2.2-offline-linux-amd64.tar.gz

# Auf Zielsystem transferieren (USB, SCP, etc.)
scp rsyslog-rest-api-v0.2.2-offline-linux-amd64.tar.gz user@server:/tmp/
```

#### Installation

```bash
# Entpacken
tar xzf rsyslog-rest-api-v0.2.2-offline-linux-amd64.tar.gz
cd package

# Install-Script ausführen
sudo ./install.sh
```

Das Script:
- ✅ Kopiert Binary nach `/opt/rsyslog-rest-api/`
- ✅ Erstellt `.env` Template
- ✅ Installiert systemd Service
- ✅ Lädt systemd daemon neu

#### Konfiguration

```bash
# Config bearbeiten (WICHTIG!)
sudo nano /opt/rsyslog-rest-api/.env

# Mindestens setzen:
# - API_KEY
# - DB_HOST, DB_NAME, DB_USER, DB_PASS

# Permissions setzen
sudo chmod 600 /opt/rsyslog-rest-api/.env
```

#### Service starten

```bash
sudo systemctl enable --now rsyslog-rest-api
sudo systemctl status rsyslog-rest-api
```

---

### Option 3: Installation from Source

Für Entwickler oder spezielle Anforderungen.

#### Voraussetzungen

- Go 1.21 oder höher
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

# Install (benötigt sudo)
sudo make install
```

Das Makefile:
- Kompiliert statisches Binary
- Installiert nach `/opt/rsyslog-rest-api/`
- Kopiert `.env.example` → `.env` (falls nicht vorhanden)
- Installiert systemd Service
- Lädt systemd neu

#### Konfiguration

```bash
# Config bearbeiten
sudo nano /opt/rsyslog-rest-api/.env

# Mindestens setzen:
# - API_KEY
# - DB_HOST, DB_NAME, DB_USER, DB_PASS

# Secure permissions
sudo chmod 600 /opt/rsyslog-rest-api/.env
```

#### Service starten

```bash
sudo systemctl enable --now rsyslog-rest-api
sudo systemctl status rsyslog-rest-api
```

---

## ✅ Verifikation

Nach jeder Installationsmethode:

### 1. Service-Status prüfen

```bash
sudo systemctl status rsyslog-rest-api
```

Erwartete Ausgabe:
```
● rsyslog-rest-api.service - rsyslog REST API
     Loaded: loaded (/etc/systemd/system/rsyslog-rest-api.service; enabled)
     Active: active (running) since ...
```

### 2. Health Check

```bash
curl http://localhost:8000/health
```

Erwartete Antwort:
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-02-09T10:30:00Z"
}
```

### 3. API Test mit Authentication

```bash
# API Key aus .env holen
API_KEY=$(sudo grep "^API_KEY=" /opt/rsyslog-rest-api/.env | cut -d'=' -f2)

# Logs abrufen
curl -H "X-API-Key: $API_KEY" "http://localhost:8000/logs?limit=5"
```

### 4. Logs prüfen

```bash
# API Logs
sudo journalctl -u rsyslog-rest-api -n 50

# Falls Fehler:
sudo journalctl -u rsyslog-rest-api -f
```

---

## 🔧 Post-Installation

Nach erfolgreicher Installation:

1. **Konfiguration optimieren:** → [Configuration Guide](configuration.md)
2. **SSL/TLS einrichten:** → [Security Guide](security.md#ssltls)
3. **Produktiv deployen:** → [Deployment Guide](deployment.md)
4. **Monitoring einrichten:** → [Deployment: Monitoring](deployment.md#monitoring)

## 🆘 Troubleshooting

Bei Problemen siehe [Troubleshooting Guide](troubleshooting.md).

Häufige Installations-Probleme:

- **Binary not found:** Pfad prüfen, Permissions prüfen
- **Database connection failed:** Credentials in `.env` prüfen
- **Permission denied:** `.env` Permissions (`chmod 600`)
- **Service won't start:** Logs prüfen mit `journalctl`

---

## 🔄 Upgrade

### Upgrade auf neue Version

```bash
# Download neue Version
wget https://github.com/phil-bot/rsyslog-rest-api/releases/download/v0.X.X/rsyslog-rest-api-linux-amd64

# Service stoppen
sudo systemctl stop rsyslog-rest-api

# Binary ersetzen
sudo mv rsyslog-rest-api-linux-amd64 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api

# Service starten
sudo systemctl start rsyslog-rest-api

# Status prüfen
sudo systemctl status rsyslog-rest-api
```

**Wichtig:** `.env` Konfiguration bleibt erhalten!

Siehe [Changelog](changelog.md) für Breaking Changes.

---

## 🗑️ Deinstallation

```bash
# Service stoppen und deaktivieren
sudo systemctl stop rsyslog-rest-api
sudo systemctl disable rsyslog-rest-api

# Service-Datei entfernen
sudo rm /etc/systemd/system/rsyslog-rest-api.service
sudo systemctl daemon-reload

# Installation entfernen
sudo rm -rf /opt/rsyslog-rest-api

# Binary entfernen (falls in /usr/local/bin)
sudo rm /usr/local/bin/rsyslog-rest-api
```

---

[← Zurück zur Übersicht](index.md) | [Weiter zu Configuration →](configuration.md)

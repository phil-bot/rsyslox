# Troubleshooting Guide

[← Zurück zur Übersicht](index.md)

Lösungen für häufige Probleme und Fehlersuche.

## 🔍 Schnelldiagnose

### Service-Status prüfen

```bash
# Status
sudo systemctl status rsyslog-rest-api

# Aktive Logs
sudo journalctl -u rsyslog-rest-api -f

# Letzte 50 Zeilen
sudo journalctl -u rsyslog-rest-api -n 50
```

### Health Check

```bash
# Einfacher Health Check
curl http://localhost:8000/health

# Mit Details
curl -v http://localhost:8000/health
```

---

## ❌ Häufige Probleme

### 1. Service startet nicht

#### Symptom
```bash
sudo systemctl start rsyslog-rest-api
# Job for rsyslog-rest-api.service failed
```

#### Ursachen & Lösungen

**A) Binary nicht gefunden**

```bash
# Prüfen ob Binary existiert
ls -la /opt/rsyslog-rest-api/rsyslog-rest-api

# Falls nicht vorhanden: Neu installieren
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-amd64
sudo cp rsyslog-rest-api-linux-amd64 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api
```

**B) .env Datei fehlt**

```bash
# Prüfen
ls -la /opt/rsyslog-rest-api/.env

# Falls nicht vorhanden
sudo nano /opt/rsyslog-rest-api/.env
# Mindest-Config eintragen
```

**C) Permissions Problem**

```bash
# Korrekte Permissions setzen
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api
```

**D) Port bereits belegt**

```bash
# Prüfen welcher Prozess Port 8000 nutzt
sudo lsof -i :8000
sudo netstat -tlnp | grep 8000

# Anderen Port in .env setzen
sudo nano /opt/rsyslog-rest-api/.env
# SERVER_PORT=8080
```

**E) Logs prüfen**

```bash
# Detaillierte Fehler
sudo journalctl -u rsyslog-rest-api -n 100 --no-pager
```

---

### 2. Database Connection Failed

#### Symptom
```json
{
  "status": "unhealthy",
  "database": "disconnected"
}
```

#### Lösungen

**A) Credentials prüfen**

```bash
# .env Datei prüfen
sudo grep "^DB_" /opt/rsyslog-rest-api/.env

# Manuell testen
mysql -h DB_HOST -u DB_USER -pDB_PASS DB_NAME
```

**B) MySQL läuft nicht**

```bash
# Status prüfen
sudo systemctl status mysql
# oder
sudo systemctl status mariadb

# Starten falls gestoppt
sudo systemctl start mysql
```

**C) Benutzer/Rechte fehlen**

```bash
# MySQL als root
sudo mysql

# User erstellen
CREATE USER 'rsyslog'@'localhost' IDENTIFIED BY 'password';
GRANT SELECT ON Syslog.* TO 'rsyslog'@'localhost';
FLUSH PRIVILEGES;
```

**D) Datenbank existiert nicht**

```bash
sudo mysql

# Datenbank erstellen
CREATE DATABASE IF NOT EXISTS Syslog;

# Tabelle erstellen (falls nicht vorhanden)
USE Syslog;
CREATE TABLE IF NOT EXISTS SystemEvents (
    ID int unsigned not null auto_increment primary key,
    ReceivedAt datetime NULL,
    FromHost varchar(60) NULL,
    Priority smallint NULL,
    Facility smallint NULL,
    Message text,
    SysLogTag varchar(60)
);
```

**E) Netzwerk-Problem (Remote DB)**

```bash
# Verbindung testen
ping DB_HOST
telnet DB_HOST 3306

# Firewall prüfen
sudo ufw status
sudo iptables -L
```

---

### 3. API Key funktioniert nicht

#### Symptom
```json
{
  "error": "Invalid or missing API key"
}
```

#### Lösungen

**A) Key korrekt verwenden**

```bash
# Key aus .env lesen
API_KEY=$(sudo grep "^API_KEY=" /opt/rsyslog-rest-api/.env | cut -d'=' -f2)

# Mit exaktem Key testen
curl -H "X-API-Key: $API_KEY" "http://localhost:8000/logs?limit=1"
```

**B) Whitespace/Formatierung**

```bash
# .env prüfen (keine Leerzeichen!)
sudo cat /opt/rsyslog-rest-api/.env | grep API_KEY

# Korrekt: API_KEY=abc123
# Falsch:  API_KEY = abc123
# Falsch:  API_KEY=abc123 (Leerzeichen am Ende)
```

**C) Service neustarten nach Änderung**

```bash
sudo systemctl restart rsyslog-rest-api
```

---

### 4. Keine Logs / Leeres Ergebnis

#### Symptom
```json
{
  "total": 0,
  "offset": 0,
  "limit": 10,
  "rows": []
}
```

#### Lösungen

**A) rsyslog schreibt nicht in DB**

```bash
# rsyslog Status
sudo systemctl status rsyslog

# rsyslog Config prüfen
cat /etc/rsyslog.d/mysql.conf

# rsyslog neustarten
sudo systemctl restart rsyslog

# Test-Log senden
logger -t test "Test message from logger"

# In DB prüfen
sudo mysql Syslog -e "SELECT COUNT(*) FROM SystemEvents"
```

**B) Filter zu streng**

```bash
# Ohne Filter testen
curl -H "X-API-Key: YOUR_KEY" "http://localhost:8000/logs?limit=10"

# Zeitfenster erweitern
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?start_date=2025-01-01T00:00:00Z&limit=10"
```

**C) Falsche Tabelle**

```bash
# Richtige Tabelle prüfen
sudo mysql Syslog -e "SHOW TABLES"
# Sollte: SystemEvents

# Einträge zählen
sudo mysql Syslog -e "SELECT COUNT(*) FROM SystemEvents"
```

---

### 5. SSL/TLS Probleme

#### Symptom
```
curl: (60) SSL certificate problem
```

#### Lösungen

**A) Self-Signed Certificate (Development)**

```bash
# Warnung ignorieren (nur Development!)
curl -k https://localhost:8000/health
```

**B) Certificate Path falsch**

```bash
# Pfade in .env prüfen
sudo grep "^SSL_" /opt/rsyslog-rest-api/.env

# Files existieren?
ls -la /opt/rsyslog-rest-api/certs/

# Permissions
sudo chmod 600 /opt/rsyslog-rest-api/certs/*.pem
```

**C) Let's Encrypt Pfad**

```bash
# Korrekter Pfad für Let's Encrypt
SSL_CERTFILE=/etc/letsencrypt/live/yourdomain.com/fullchain.pem
SSL_KEYFILE=/etc/letsencrypt/live/yourdomain.com/privkey.pem

# API muss Zugriff haben
sudo chmod +r /etc/letsencrypt/live/yourdomain.com/*.pem
# ODER
sudo usermod -a -G ssl-cert root
```

---

### 6. CORS Errors (Browser)

#### Symptom
```
Access to fetch at 'http://api.example.com' from origin 
'http://dashboard.example.com' has been blocked by CORS policy
```

#### Lösungen

**A) ALLOWED_ORIGINS korrekt setzen**

```bash
sudo nano /opt/rsyslog-rest-api/.env

# Development (alle erlauben)
ALLOWED_ORIGINS=*

# Production (spezifisch)
ALLOWED_ORIGINS=https://dashboard.example.com,https://app.example.com

# WICHTIG: https:// prefix!
```

**B) Service neustarten**

```bash
sudo systemctl restart rsyslog-rest-api
```

**C) Browser-Cache leeren**

```bash
# Chrome: Strg+Shift+Delete
# Firefox: Strg+Shift+Del
# Oder: Inkognito-Modus testen
```

---

### 7. Performance Probleme

#### Symptom
Langsame Antwortzeiten, Timeouts

#### Lösungen

**A) Zeitfenster einschränken**

```bash
# Statt 90 Tage...
?start_date=2025-02-09T00:00:00Z&end_date=2025-02-09T23:59:59Z

# ...nur 1 Tag verwenden
```

**B) Limit verwenden**

```bash
# Nicht: ?limit=10000
# Besser: ?limit=100 mit Pagination
```

**C) Indexes prüfen**

```bash
sudo mysql Syslog -e "SHOW INDEX FROM SystemEvents"

# Falls fehlend, neu erstellen
sudo mysql Syslog <<EOF
CREATE INDEX idx_receivedat ON SystemEvents (ReceivedAt);
CREATE INDEX idx_host_time ON SystemEvents (FromHost, ReceivedAt);
CREATE INDEX idx_priority ON SystemEvents (Priority);
EOF
```

**D) Database Stats**

```bash
# Tabellengröße
sudo mysql Syslog -e "
SELECT 
  COUNT(*) as rows,
  ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
FROM information_schema.TABLES
WHERE table_schema = 'Syslog' AND table_name = 'SystemEvents'
"
```

---

## 🐞 Debug Mode

### Verbose Logging aktivieren

Für detaillierte Logs (Development):

```bash
# Service stoppen
sudo systemctl stop rsyslog-rest-api

# Manuell im Vordergrund starten
cd /opt/rsyslog-rest-api
sudo -E ./rsyslog-rest-api

# In anderem Terminal testen
curl http://localhost:8000/health

# Ausgabe ansehen
```

### SQL Queries loggen

In `main.go` (für Entwickler):

```go
// Vor db.Query()
log.Printf("SQL: %s", sqlQuery)
log.Printf("Args: %v", args)
```

---

## 🔬 Erweiterte Diagnose

### Network Debugging

```bash
# Port Listen prüfen
sudo ss -tlnp | grep 8000
sudo netstat -tlnp | grep 8000

# Firewall
sudo ufw status verbose
sudo iptables -L -n -v

# Test von anderem Host
telnet API_HOST 8000
```

### MySQL Connection Debugging

```bash
# MySQL Connections
sudo mysql -e "SHOW PROCESSLIST"

# MySQL Errors
sudo tail -f /var/log/mysql/error.log

# MySQL Slow Queries
sudo tail -f /var/log/mysql/slow.log
```

### System Resources

```bash
# CPU/Memory
top
htop

# Disk Space
df -h

# Disk I/O
iostat -x 1
```

---

## 📋 FAQ

### Kann ich die API ohne API Key verwenden?

**Ja**, setze `API_KEY=` (leer) in `.env`. Nur für Development empfohlen!

### Wie groß kann limit sein?

**Maximum 1000**. Bei größeren Datenmengen Pagination verwenden.

### Welcher Zeitraum ist erlaubt?

**Maximum 90 Tage** zwischen `start_date` und `end_date`.

### Funktioniert die API mit PostgreSQL?

**Nein**, nur MySQL/MariaDB. PostgreSQL Support ist nicht geplant.

### Kann ich mehrere API Keys verwenden?

**Nein**, aktuell nur ein globaler API Key. Für mehrere Clients den gleichen Key verwenden (oder Feature Request erstellen).

### Wie viele Requests/Sekunde sind möglich?

Aktuell **kein** Rate Limiting. Performance hängt von Hardware und Datenbankgröße ab. Für Production: Rate Limiting über Reverse Proxy.

### Kann ich eigene Felder hinzufügen?

**Ja**, alle Spalten in `SystemEvents` sind automatisch verfügbar. Einfach Spalte in DB hinzufügen, API neustarten.

### Wird HTTP/2 unterstützt?

**Nein**, aktuell nur HTTP/1.1. HTTP/2 über Reverse Proxy (nginx) möglich.

### Kann ich die API in Docker laufen lassen?

**Ja**, siehe [Docker Guide](docker.md). Empfohlen für Testing, nicht für Production.

---

## 🆘 Weitere Hilfe

### Logs sammeln

Für Bug Reports:

```bash
# System Info
uname -a
cat /etc/os-release

# API Version
/opt/rsyslog-rest-api/rsyslog-rest-api --version

# Service Status
sudo systemctl status rsyslog-rest-api

# Logs
sudo journalctl -u rsyslog-rest-api -n 100 --no-pager > api-logs.txt

# Config (Passwörter schwärzen!)
sudo cat /opt/rsyslog-rest-api/.env | sed 's/PASS=.*/PASS=REDACTED/' > config.txt
```

### GitHub Issue erstellen

→ [GitHub Issues](https://github.com/phil-bot/rsyslog-rest-api/issues)

**Template:**
```markdown
**Environment:**
- OS: Ubuntu 22.04
- API Version: v0.2.2
- Installation: Binary/Source/Docker

**Problem:**
[Beschreibung]

**Steps to Reproduce:**
1. ...
2. ...

**Expected:**
[Was sollte passieren]

**Actual:**
[Was passiert tatsächlich]

**Logs:**
```
[Logs hier einfügen]
```

**Config:**
[Config (ohne Passwörter!)]
```

---

[← Zurück zur Übersicht](index.md) | [Weiter zu Docker →](docker.md)

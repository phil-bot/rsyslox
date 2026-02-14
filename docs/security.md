# Security Guide

[← Zurück zur Übersicht](index.md)

Sicherheits-Best-Practices für rsyslog REST API.

## 🔒 Security Checkliste

### Production Deployment

- ✅ **API-Key aktiviert** und stark (32+ Bytes)
- ✅ **SSL/TLS aktiviert** (Let's Encrypt oder kommerziell)
- ✅ **CORS konfiguriert** (nicht `*`)
- ✅ **Reverse Proxy** mit Rate Limiting
- ✅ **Firewall aktiv** (nur 80/443 offen)
- ✅ **Read-Only DB-User** für API
- ✅ **File Permissions** (`.env` = 600)
- ✅ **Regular Updates** geplant
- ✅ **Monitoring & Logging** aktiv
- ✅ **Backup-Strategie** vorhanden

---

## 🔐 API Key Security

### Generierung

**Empfohlen: 32 Bytes (64 Hex-Zeichen)**

```bash
# Starker API Key
openssl rand -hex 32

# Beispiel-Output:
# a3d7f8c9e2b4a6d8f9c3e7b1a5d9f4c8e2b7a6d3f9c8e1b4a7d2f6c9e3b8a5d1
```

**Niemals verwenden:**
```bash
# ❌ Zu kurz
API_KEY=test123

# ❌ Wörterbuch-Wort
API_KEY=password

# ❌ Vorhersehbar
API_KEY=12345678901234567890
```

### Speicherung

**.env Permissions:**
```bash
# KRITISCH: Nur root darf lesen!
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env

# Verifizieren
ls -la /opt/rsyslog-rest-api/.env
# Sollte: -rw------- root root
```

**Niemals:**
- ❌ In Git committen
- ❌ Per E-Mail versenden (unverschlüsselt)
- ❌ In Logs ausgeben
- ❌ In URL-Parameter packen

### Rotation

**API-Key regelmäßig wechseln (z.B. vierteljährlich):**

```bash
# 1. Neuen Key generieren
NEW_KEY=$(openssl rand -hex 32)

# 2. In .env aktualisieren
sudo sed -i "s/^API_KEY=.*/API_KEY=$NEW_KEY/" /opt/rsyslog-rest-api/.env

# 3. Service neustarten
sudo systemctl restart rsyslog-rest-api

# 4. Clients updaten
# Sende neuen Key sicher an autorisierte Nutzer
```

**Bei Kompromittierung SOFORT wechseln!**

---

## 🔒 SSL/TLS Configuration

### Production: Let's Encrypt

**Installation:**
```bash
# Ubuntu/Debian
sudo apt-get install certbot python3-certbot-nginx

# CentOS/RHEL
sudo yum install certbot python3-certbot-nginx
```

**Zertifikat beantragen:**
```bash
# Für nginx
sudo certbot --nginx -d api.yourdomain.com

# Manuell (DNS/Webroot)
sudo certbot certonly --standalone -d api.yourdomain.com
```

**Auto-Renewal:**
```bash
# Test
sudo certbot renew --dry-run

# Cron (automatisch von certbot erstellt)
# Prüfen: sudo crontab -l
```

**API Konfiguration:**

```bash
# .env - API OHNE SSL (terminiert am Reverse Proxy!)
USE_SSL=false
```

**nginx Konfiguration:**

```nginx
server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    
    # SSL Best Practices (Mozilla Intermediate)
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_prefer_server_ciphers off;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # OCSP Stapling
    ssl_stapling on;
    ssl_stapling_verify on;
    ssl_trusted_certificate /etc/letsencrypt/live/api.yourdomain.com/chain.pem;
    
    # ...
}
```

### Development: Self-Signed

**Nur für Testing! Nicht für Production!**

```bash
# Zertifikat erstellen
sudo mkdir -p /opt/rsyslog-rest-api/certs
cd /opt/rsyslog-rest-api/certs

sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout key.pem \
  -out cert.pem \
  -subj "/C=DE/ST=State/L=City/O=Organization/CN=localhost"

# Permissions
sudo chmod 600 *.pem
```

**API Config:**
```bash
USE_SSL=true
SSL_CERTFILE=/opt/rsyslog-rest-api/certs/cert.pem
SSL_KEYFILE=/opt/rsyslog-rest-api/certs/key.pem
```

### SSL Testing

```bash
# SSL Labs (online)
# https://www.ssllabs.com/ssltest/analyze.html?d=api.yourdomain.com

# testssl.sh (lokal)
git clone --depth 1 https://github.com/drwetter/testssl.sh.git
cd testssl.sh
./testssl.sh https://api.yourdomain.com
```

---

## 🌐 CORS Security

### Production Setup

**Niemals `*` in Production!**

```bash
# ❌ FALSCH (Development only!)
ALLOWED_ORIGINS=*

# ✅ RICHTIG (spezifische Domains)
ALLOWED_ORIGINS=https://dashboard.yourdomain.com,https://monitoring.yourdomain.com
```

### Mehrere Subdomains

```bash
# Alle Subdomains erlauben (wenn sicher)
ALLOWED_ORIGINS=https://app.yourdomain.com,https://dashboard.yourdomain.com,https://admin.yourdomain.com
```

### Testing

```bash
# Preflight Request testen
curl -X OPTIONS \
  -H "Origin: https://dashboard.yourdomain.com" \
  -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: X-API-Key" \
  https://api.yourdomain.com/logs -v
```

---

## 💾 Database Security

### Read-Only User (Empfohlen!)

Die API **schreibt nicht** in die Datenbank - nur SELECT nötig!

```sql
-- Neuen Read-Only User erstellen
CREATE USER 'rsyslog_readonly'@'localhost' 
  IDENTIFIED BY 'very-secure-password-here';

-- Nur SELECT auf SystemEvents
GRANT SELECT ON Syslog.SystemEvents TO 'rsyslog_readonly'@'localhost';

FLUSH PRIVILEGES;

-- Verifizieren
SHOW GRANTS FOR 'rsyslog_readonly'@'localhost';
```

**API Config:**
```bash
DB_USER=rsyslog_readonly
DB_PASS=very-secure-password-here
```

### Starke Passwörter

```bash
# Datenbank-Passwort generieren (32 Zeichen)
openssl rand -base64 32

# Beispiel:
# jK8$mN2pQ4rS6tU8vW0xY2zA4bC6dE8fG
```

### Remote Database

**Wenn DB auf anderem Server:**

```bash
# Firewall: Nur API-Server erlauben
# UFW (Ubuntu)
sudo ufw allow from API_SERVER_IP to any port 3306

# firewalld (CentOS)
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="API_SERVER_IP" port protocol="tcp" port="3306" accept'
sudo firewall-cmd --reload
```

**MySQL/MariaDB:**
```sql
-- User nur von API-Server
CREATE USER 'rsyslog_api'@'API_SERVER_IP' 
  IDENTIFIED BY 'password';
  
GRANT SELECT ON Syslog.SystemEvents TO 'rsyslog_api'@'API_SERVER_IP';

-- bind-address anpassen
-- /etc/mysql/mariadb.conf.d/50-server.cnf
-- bind-address = 0.0.0.0  # oder spezifische IP
```

### SSL für Database Connection (optional)

```bash
# MySQL mit SSL
DB_HOST=db.example.com
# Zusätzlich SSL-Parameter in Connection String

# Erfordert Code-Änderung in main.go:
# dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&tls=true", ...)
```

---

## 🔥 Firewall Configuration

### ufw (Ubuntu/Debian)

```bash
# Default Deny
sudo ufw default deny incoming
sudo ufw default allow outgoing

# SSH (WICHTIG - sonst aussperren!)
sudo ufw allow ssh

# HTTP/HTTPS (für nginx)
sudo ufw allow http
sudo ufw allow https

# MySQL (nur wenn Remote DB)
# sudo ufw allow from DB_SERVER_IP to any port 3306

# Aktivieren
sudo ufw enable

# Status
sudo ufw status verbose
```

### firewalld (CentOS/RHEL)

```bash
# Services erlauben
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --permanent --add-service=ssh

# Reload
sudo firewall-cmd --reload

# Status
sudo firewall-cmd --list-all
```

### iptables (Advanced)

```bash
# Neue Chain für API
sudo iptables -N API_RATE_LIMIT

# Rate Limiting (10 req/sec)
sudo iptables -A API_RATE_LIMIT -m limit --limit 10/sec -j ACCEPT
sudo iptables -A API_RATE_LIMIT -j DROP

# API Port
sudo iptables -A INPUT -p tcp --dport 8000 -j API_RATE_LIMIT

# Persistent machen
sudo apt-get install iptables-persistent
sudo netfilter-persistent save
```

---

## 🛡️ Rate Limiting

### nginx Rate Limiting

```nginx
# In http {} Block
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_status 429;

# In server/location {} Block
location / {
    limit_req zone=api_limit burst=20 nodelay;
    # ...
}
```

**Parameter:**
- `rate=10r/s` - 10 Requests pro Sekunde
- `burst=20` - Burst bis 20 Requests
- `nodelay` - Sofort 429 wenn Limit überschritten

### Per-IP Blocking

```nginx
# Geografisches Blocking (benötigt GeoIP)
geo $blocked_country {
    default 0;
    CN 1;  # China
    RU 1;  # Russland
}

server {
    if ($blocked_country) {
        return 403;
    }
    # ...
}
```

---

## 🔍 Security Headers

### nginx Headers

```nginx
server {
    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;
    
    # Content Security
    add_header Content-Security-Policy "default-src 'self'" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    # Referrer Policy
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    # Permissions Policy
    add_header Permissions-Policy "geolocation=(), microphone=(), camera=()" always;
    
    # ...
}
```

### Testen

```bash
# securityheaders.com
curl -I https://api.yourdomain.com

# Mozilla Observatory
# https://observatory.mozilla.org/
```

---

## 📝 Logging & Auditing

### Access Logging

**nginx:**
```nginx
# Custom Log Format mit Details
log_format api_access '$remote_addr - $remote_user [$time_local] '
                      '"$request" $status $body_bytes_sent '
                      '"$http_referer" "$http_user_agent" '
                      '$request_time $upstream_response_time';

access_log /var/log/nginx/rsyslog-api-access.log api_access;
```

### Failed Auth Attempts

**Monitoring script:**
```bash
#!/bin/bash
# /usr/local/bin/monitor-auth-failures.sh

LOG="/var/log/nginx/rsyslog-api-access.log"
THRESHOLD=10
EMAIL="admin@yourdomain.com"

# 401 Responses zählen (letzte 5 Minuten)
FAILURES=$(tail -5000 "$LOG" | grep " 401 " | wc -l)

if [ "$FAILURES" -gt "$THRESHOLD" ]; then
    echo "WARNING: $FAILURES failed auth attempts in last 5 minutes!" | \
      mail -s "API Auth Failures Alert" "$EMAIL"
fi
```

**Cron:**
```bash
*/5 * * * * /usr/local/bin/monitor-auth-failures.sh
```

---

## 🔐 System Hardening

### Minimale Permissions

```bash
# API Binary
sudo chown root:root /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod 500 /opt/rsyslog-rest-api/rsyslog-rest-api

# Config
sudo chown root:root /opt/rsyslog-rest-api/.env
sudo chmod 400 /opt/rsyslog-rest-api/.env
```

### systemd Sandboxing

```ini
# /etc/systemd/system/rsyslog-rest-api.service
[Service]
# ...

# Security Hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/rsyslog-rest-api
ProtectKernelTunables=true
ProtectControlGroups=true
RestrictRealtime=true
RestrictNamespaces=true
```

### SELinux (CentOS/RHEL)

```bash
# SELinux Status
getenforce

# Für API Port
sudo semanage port -a -t http_port_t -p tcp 8000

# Für File Access
sudo chcon -R -t httpd_sys_content_t /opt/rsyslog-rest-api/
```

---

## 🚨 Incident Response

### Bei Kompromittierung

**1. API sofort stoppen:**
```bash
sudo systemctl stop rsyslog-rest-api
```

**2. API-Key ändern:**
```bash
NEW_KEY=$(openssl rand -hex 32)
sudo sed -i "s/^API_KEY=.*/API_KEY=$NEW_KEY/" /opt/rsyslog-rest-api/.env
```

**3. Logs analysieren:**
```bash
# Verdächtige Requests
sudo grep " 401 " /var/log/nginx/rsyslog-api-access.log
sudo journalctl -u rsyslog-rest-api --since "1 hour ago"

# IPs sammeln
sudo awk '{print $1}' /var/log/nginx/rsyslog-api-access.log | \
  sort | uniq -c | sort -rn | head -20
```

**4. Verdächtige IPs blocken:**
```bash
# nginx
sudo nano /etc/nginx/sites-available/rsyslog-api
# Hinzufügen: deny SUSPICIOUS_IP;

# ufw
sudo ufw deny from SUSPICIOUS_IP

# iptables
sudo iptables -A INPUT -s SUSPICIOUS_IP -j DROP
```

**5. Clients informieren:**
```bash
# Neuer API-Key an autorisierte Nutzer
```

**6. System prüfen:**
```bash
# Rootkit-Check
sudo rkhunter --check
sudo chkrootkit

# Integrity Check
sudo debsums -c  # Debian/Ubuntu
sudo rpm -Va     # CentOS/RHEL
```

---

## 📋 Security Audit Checkliste

### Monatlich

- [ ] API-Key rotieren (optional)
- [ ] Failed auth attempts prüfen
- [ ] SSL-Zertifikat-Ablauf checken
- [ ] Updates verfügbar?
- [ ] Logs auf Anomalien prüfen

### Vierteljährlich

- [ ] API-Key rotieren (empfohlen)
- [ ] Datenbank-Passwort ändern
- [ ] Backup-Strategie testen
- [ ] Permissions verifizieren
- [ ] Security Headers testen

### Jährlich

- [ ] Vollständige Security-Audit
- [ ] Penetration Test erwägen
- [ ] Dokumentation aktualisieren
- [ ] Disaster Recovery Plan testen

---

## 🔗 Weitere Ressourcen

- **OWASP Top 10:** https://owasp.org/www-project-top-ten/
- **Let's Encrypt:** https://letsencrypt.org/
- **Mozilla SSL Config Generator:** https://ssl-config.mozilla.org/
- **Security Headers:** https://securityheaders.com/

---

[← Zurück zur Übersicht](index.md) | [Weiter zu Performance →](performance.md)

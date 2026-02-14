# Deployment Guide

[← Back to overview](index.md)

Production deployment for rsyslog REST API.

## 🎯 Production Setup

### Preparation

**Checklist:**
- ✅ Server with Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+)
- ✅ rsyslog with MySQL/MariaDB installed and configured
- ✅ Database is running and accessible
- ✅ Firewall rules planned
- sSL certificate available (Let's Encrypt recommended)
- aPI key generated

---

## 📦 Installation (Production)

### 1. install binary

```bash
# Download
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-amd64

# Install
sudo mkdir -p /opt/rsyslog-rest-api
sudo mv rsyslog-rest-api-linux-amd64 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api
```

### 2. create configuration

```bash
sudo nano /opt/rsyslog-rest-api/.env
```

**Production .env:**
```bash
# API Security (REQUIRED!)
API_KEY=<generated-64-character-key>

# Server
SERVER_HOST=127.0.0.1 # Only localhost (behind reverse proxy!)
SERVER_PORT=8000

# SSL/TLS
USE_SSL=false # Terminate SSL at reverse proxy

# CORS
ALLOWED_ORIGINS=https://yourdomain.com,https://dashboard.yourdomain.com

# Database
DB_HOST=localhost # or remote DB server
DB_NAME=Syslog
DB_USER=rsyslog_api
DB_PASS=<secure-database-password>
```

*generate *API key:**
```bash
openssl rand -hex 32
```

**Permissions:**
```bash
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
```

### 3. systemd service

**Service file:**
```bash
sudo nano /etc/systemd/system/rsyslog-rest-api.service
```

```ini
[Unit]
Description=rsyslog REST API
After=network.target mysql.service mariadb.service
Wants=mysql.service mariadb.service

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/opt/rsyslog-rest-api
EnvironmentFile=/opt/rsyslog-rest-api/.env
ExecStart=/opt/rsyslog-rest-api/rsyslog-rest-api
Restart=on-failure
RestartSec=5s

# Security Hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/rsyslog-rest-api

[Install]
WantedBy=multi-user.target
```

**Activate:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable rsyslog-rest-api
sudo systemctl start rsyslog-rest-api
sudo systemctl status rsyslog-rest-api
```

---

## 🔄 Reverse Proxy Setup

### nginx (Recommended)

**Installation:**
```bash
sudo apt-get install nginx
```

**Configuration:**
```bash
sudo nano /etc/nginx/sites-available/rsyslog-api
```

```nginx
# /etc/nginx/sites-available/rsyslog-api

# Rate limiting
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_status 429;

# upstream
upstream rsyslog_api {
    server 127.0.0.1:8000;
    keepalive 32;
}

server {
    listen 80;
    server_name api.yourdomain.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    
    # SSL Best Practices
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Logging
    access_log /var/log/nginx/rsyslog-api-access.log;
    error_log /var/log/nginx/rsyslog-api-error.log;

    # API endpoints
    location / {
        # Rate limiting
        limit_req zone=api_limit burst=20 nodelay;
        
        # Proxy Headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Proxy Settings
        proxy_pass http://rsyslog_api;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Health Check (no rate limit)
    location /health {
        proxy_pass http://rsyslog_api/health;
        access_log off;
    }
}
```

**Activate:**
```bash
# Create link
sudo ln -s /etc/nginx/sites-available/rsyslog-api /etc/nginx/sites-enabled/

# Test
sudo nginx -t

# Reload
sudo systemctl reload nginx
```

**Let's Encrypt SSL:**
```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Request certificate
sudo certbot --nginx -d api.yourdomain.com

# Check auto-renewal
sudo certbot renew --dry-run
```

### Apache

**Installation:**
```bash
sudo apt-get install apache2
sudo a2enmod proxy proxy_http ssl headers
```

**Configuration:**
```bash
sudo nano /etc/apache2/sites-available/rsyslog-api.conf
```

```apache
<VirtualHost *:80>
    ServerName api.yourdomain.com
    Redirect permanent / https://api.yourdomain.com/
</VirtualHost>

<VirtualHost *:443>
    ServerName api.yourdomain.com
    
    # SSL
    SSLEngine on
    SSLCertificateFile /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem
    SSLCertificateKeyFile /etc/letsencrypt/live/api.yourdomain.com/privkey.pem
    
    # Security Headers
    Header always set Strict-Transport-Security "max-age=31536000"
    Header always set X-Content-Type-Options "nosniff"
    Header always set X-Frame-Options "DENY"
    
    # Proxy
    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:8000/
    ProxyPassReverse / http://127.0.0.1:8000/
    
    # Logging
    ErrorLog ${APACHE_LOG_DIR}/rsyslog-api-error.log
    CustomLog ${APACHE_LOG_DIR}/rsyslog-api-access.log combined
</VirtualHost>
```

**Activate:**
```bash
sudo a2ensite rsyslog-api
sudo systemctl reload apache2
```

---

## 🔥 Firewall Setup

### ufw (Ubuntu/Debian)

```bash
# Allow SSH (IMPORTANT!)
sudo ufw allow ssh

# HTTP/HTTPS (for nginx/Apache)
sudo ufw allow http
sudo ufw allow https

# Block API port (only runs on localhost)
# No rule necessary - 127.0.0.1:8000 is not externally accessible

# Enable firewall
sudo ufw enable

# Status
sudo ufw status verbose
```

### firewalld (CentOS/RHEL)

```bash
# Allow HTTP/HTTPS
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https

# Reload
sudo firewall-cmd --reload

# Status
sudo firewall-cmd --list-all
```

---

## 📊 Monitoring

### Log monitoring

```bash
# API logs live
sudo journalctl -u rsyslog-rest-api -f

# Nginx access logs
sudo tail -f /var/log/nginx/rsyslog-api-access.log

# Nginx Error Logs
sudo tail -f /var/log/nginx/rsyslog-api-error.log
```

### Health Checks

**Local health check:**
```bash
#!/bin/bash
# /usr/local/bin/health-check.sh

if curl -f http://localhost:8000/health > /dev/null 2>&1; then
    echo "API healthy"
    exit 0
else
    echo "API unhealthy!"
    # Optional: Restart the service
    # systemctl restart rsyslog-rest-api
    exit 1
fi
```

**Cron Job:**
```bash
# Check every 5 minutes
*/5 * * * * * /usr/local/bin/health-check.sh
```

### Performance monitoring

**Simple Stats:**
```bash
# Requests per minute (nginx)
tail -1000 /var/log/nginx/rsyslog-api-access.log | \
  awk '{print $4}' | cut -d: -f2 | sort | uniq -c

# Average response time
# (requires custom nginx log format with $request_time)
```

**Advanced: Prometheus + Grafana**

For production monitoring see [Performance Guide](performance.md#monitoring).

---

## 🔄 Backup & Recovery

### Backup configuration

```bash
#!/bin/bash
# backup-config.sh

BACKUP_DIR="/var/backups/rsyslog-api"
DATE=$(date +%Y%m%d-%H%M%S)

mkdir -p "$BACKUP_DIR"

# Backup .env
sudo cp /opt/rsyslog-rest-api/.env \
  "$BACKUP_DIR/.env-$DATE"

# Back up the systemd service
sudo cp /etc/systemd/system/rsyslog-rest-api.service \
  "$BACKUP_DIR/rsyslog-rest-api.service-$DATE"

# Back up nginx config
sudo cp /etc/nginx/sites-available/rsyslog-api \
  "$BACKUP_DIR/nginx-rsyslog-api-$DATE"

echo "Backup created: $BACKUP_DIR/*-$DATE"

# Delete old backups (older than 30 days)
find "$BACKUP_DIR" -name "*-*" -mtime +30 -delete
```

**Cron:**
```bash
# Daily at 2 o'clock
0 2 * * * * /usr/local/bin/backup-config.sh
```

### Disaster Recovery

```bash
#1 Reinstall binary
wget https://github.com/.../rsyslog-rest-api-linux-amd64
sudo mv rsyslog-rest-api-linux-amd64 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api

# 2. restore backup
sudo cp /var/backups/rsyslog-api/.env-LATEST /opt/rsyslog-rest-api/.env
sudo chmod 600 /opt/rsyslog-rest-api/.env

# 3. start service
sudo systemctl start rsyslog-rest-api
```

---

## 📈 Scaling

### Vertical scaling

**More resources:**
- CPU: API is I/O-bound (database), more CPU helps limited
- RAM: Minimum 256MB, recommended 512MB+
- Disk: Only for logs, minimum

**Database Optimization:**
```sql
-- Check indexes
SHOW INDEX FROM SystemEvents;

-- Connection Pool (Go defaults are good)
-- Max 25 connections, 5 idle
```

### Horizontal Scaling

**Load Balancer Setup (nginx):**

```nginx
upstream rsyslog_api_cluster {
    least_conn;
    server 192.168.1.10:8000;
    server 192.168.1.11:8000;
    server 192.168.1.12:8000;
    
    keepalive 32;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;
    
    # ... SSL config ...
    
    location / {
        proxy_pass http://rsyslog_api_cluster;
        # ... proxy settings ...
    }
}
```

**Shared Database:**
- All API instances connect to the same DB
- READ-ONLY user recommended (API does not write)

---

## 🔐 Security Hardening

### Minimum permissions

```bash
# Separate user (optional, but recommended)
sudo useradd -r -s /bin/false rsyslog-api

# Service as rsyslog-api user
# /etc/systemd/system/rsyslog-rest-api.service:
# User=rsyslog-api
# Group=rsyslog-api

# Permissions
sudo chown rsyslog-api:rsyslog-api /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod 500 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chown rsyslog-api:rsyslog-api /opt/rsyslog-rest-api/.env
sudo chmod 400 /opt/rsyslog-rest-api/.env
```

### Database Security

```sql
-- Create READ-ONLY user
CREATE USER 'rsyslog_readonly'@'localhost'
  IDENTIFIED BY 'secure-password';

-- SELECT on SystemEvents only
GRANT SELECT ON Syslog.SystemEvents TO 'rsyslog_readonly'@'localhost';

FLUSH PRIVILEGES;
```

### Further measures

→ See [Security Guide](security.md)

---

## 🧪 Testing Production Setup

```bash
# 1st health check
curl https://api.yourdomain.com/health

# 2. with API key
curl -H "X-API-Key: YOUR_KEY" \
  "https://api.yourdomain.com/logs?limit=1"

# 3. SSL/TLS check
curl -v https://api.yourdomain.com/health 2>&1 | grep "SSL connection"

# 4. performance test
ab -n 1000 -c 10 "https://api.yourdomain.com/health"

# 5. rate limiting test (should return 429)
for i in {1..100}; do
  curl -s -o /dev/null -w "%{http_code}\n" \
    "https://api.yourdomain.com/health"
  sleep 0.01
done
```

---

## 🔄 Maintenance

### Updates

```bash
# 1. download new version
wget https://github.com/.../rsyslog-rest-api-VERSION-linux-amd64

# 2. stop service
sudo systemctl stop rsyslog-rest-api

# 3. replace binary
sudo mv rsyslog-rest-api-VERSION-linux-amd64 \
  /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api

# 4. start service
sudo systemctl start rsyslog-rest-api

# 5. health check
curl http://localhost:8000/health
```

### Log rotation

**API Logs (journald):**
```bash
sudo nano /etc/systemd/journald.conf
```

```ini
[Journal]
SystemMaxUse=100M
SystemMaxFileSize=10M
```

**nginx Logs:**
```bash
# logrotate (usually already configured)
cat /etc/logrotate.d/nginx
```

---

[← Back to overview](index.md) | [Forward to Security →](security.md)

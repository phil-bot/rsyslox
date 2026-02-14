# Security Guide

[← Back to overview](index.md)

Security best practices for rsyslog REST API.

## 🔒 Security checklist

### Production Deployment

- ✅ **API key activated** and strong (32+ bytes)
- ✅ **SSL/TLS activated** (Let's Encrypt or commercial)
- ✅ **CORS configured** (not `*`)
- ✅ **Reverse proxy** with rate limiting
- ✅ **Firewall active** (only 80/443 open)
- ✅ **Read-only DB user** for API
- ✅ **File permissions** (`.env` = 600)
- ✅ **Regular updates** planned
- ✅ **Monitoring & logging** active
- ✅ **Backup strategy** available

---

## 🔐 API Key Security

### Generation

**Recommended: 32 bytes (64 hex characters)**

```bash
# Strong API key
openssl rand -hex 32

# Example output:
# a3d7f8c9e2b4a6d8f9c3e7b1a5d9f4c8e2b7a6d3f9c8e1b4a7d2f6c9e3b8a5d1
```

**never use:**
```bash
# ❌ Too short
API_KEY=test123

# ❌ Dictionary word
API_KEY=password

# ❌ Predictable
API_KEY=12345678901234567890
```

### Storage

**.env Permissions:**
```bash
# CRITICAL: Only root may read!
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env

# Verify
ls -la /opt/rsyslog-rest-api/.env
# Should: -rw------- root root
```

**Never:**
- ❌ Commit to Git
- ❌ Send by e-mail (unencrypted)
- output to logs
- ❌ Pack into URL parameters

### Rotation

**Change API key regularly (e.g. quarterly):**

```bash
#1 Generate new key
NEW_KEY=$(openssl rand -hex 32)

# 2. update in .env
sudo sed -i "s/^API_KEY=.*/API_KEY=$NEW_KEY/" /opt/rsyslog-rest-api/.env

# 3. restart the service
sudo systemctl restart rsyslog-rest-api

# 4. update clients
# Send new key securely to authorized users
```

**If compromised, change IMMEDIATELY!

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

**apply for certificate:**
```bash
# For nginx
sudo certbot --nginx -d api.yourdomain.com

# Manually (DNS/webroot)
sudo certbot certonly --standalone -d api.yourdomain.com
```

**Auto-Renewal:**
```bash
# Test
sudo certbot renew --dry-run

# Cron (automatically created by certbot)
# Check: sudo crontab -l
```

**API configuration:**

```bash
# .env - API WITHOUT SSL (terminated at reverse proxy!)
USE_SSL=false
```

**nginx configuration:**

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

**For testing only! Not for production!

```bash
# Create certificate
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

# testssl.sh (local)
git clone --depth 1 https://github.com/drwetter/testssl.sh.git
cd testssl.sh
./testssl.sh https://api.yourdomain.com
```

---

## 🌐 CORS Security

### Production Setup

**Never `*` in Production!

```bash
# ❌ FALSE (Development only!)
ALLOWED_ORIGINS=*

# ✅ CORRECT (specific domains)
ALLOWED_ORIGINS=https://dashboard.yourdomain.com,https://monitoring.yourdomain.com
```

### Multiple subdomains

```bash
# Allow all subdomains (if safe)
ALLOWED_ORIGINS=https://app.yourdomain.com,https://dashboard.yourdomain.com,https://admin.yourdomain.com
```

### Testing

```bash
# Test preflight request
curl -X OPTIONS \
  -H "Origin: https://dashboard.yourdomain.com" \
  -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: X-API-Key" \
  https://api.yourdomain.com/logs -v
```

---

## 💾 Database Security

### Read-Only User (Recommended!)

The API **does not** write to the database - only SELECT required!

```sql
-- Create new read-only user
CREATE USER 'rsyslog_readonly'@'localhost'
  IDENTIFIED BY 'very-secure-password-here';

-- SELECT on SystemEvents only
GRANT SELECT ON Syslog.SystemEvents TO 'rsyslog_readonly'@'localhost';

FLUSH PRIVILEGES;

-- Verify
SHOW GRANTS FOR 'rsyslog_readonly'@'localhost';
```

**API Config:**
```bash
DB_USER=rsyslog_readonly
DB_PASS=very-secure-password-here
```

### Strong passwords

```bash
# Generate database password (32 characters)
openssl rand -base64 32

# Example:
# jK8$mN2pQ4rS6tU8vW0xY2zA4bC6dE8fG
```

### Remote Database

**If DB on other server:**

```bash
# Firewall: Only allow API server
# UFW (Ubuntu)
sudo ufw allow from API_SERVER_IP to any port 3306

# firewalld (CentOS)
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="API_SERVER_IP" port protocol="tcp" port="3306" accept'
sudo firewall-cmd --reload
```

**MySQL/MariaDB:**
```sql
-- User only from API server
CREATE USER 'rsyslog_api'@'API_SERVER_IP'
  IDENTIFIED BY 'password';
  
GRANT SELECT ON Syslog.SystemEvents TO 'rsyslog_api'@'API_SERVER_IP';

-- customize bind-address
-- /etc/mysql/mariadb.conf.d/50-server.cnf
-- bind-address = 0.0.0.0 # or specific IP
```

### SSL for Database Connection (optional)

```bash
# MySQL with SSL
DB_HOST=db.example.com
# Additional SSL parameter in connection string

# Requires code change in main.go:
# dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&tls=true", ...)
```

---

## 🔥 Firewall Configuration

### ufw (Ubuntu/Debian)

```bash
# Default Deny
sudo ufw default deny incoming
sudo ufw default allow outgoing

# SSH (IMPORTANT - otherwise lock out!)
sudo ufw allow ssh

# HTTP/HTTPS (for nginx)
sudo ufw allow http
sudo ufw allow https

# MySQL (only if remote DB)
# sudo ufw allow from DB_SERVER_IP to any port 3306

# Enable
sudo ufw enable

# Status
sudo ufw status verbose
```

### firewalld (CentOS/RHEL)

```bash
# Allow services
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
# New chain for API
sudo iptables -N API_RATE_LIMIT

# Rate limiting (10 req/sec)
sudo iptables -A API_RATE_LIMIT -m limit --limit 10/sec -j ACCEPT
sudo iptables -A API_RATE_LIMIT -j DROP

# API port
sudo iptables -A INPUT -p tcp --dport 8000 -j API_RATE_LIMIT

# Make persistent
sudo apt-get install iptables-persistent
sudo netfilter-persistent save
```

---

## 🛡️ Rate Limiting

### nginx Rate Limiting

```nginx
# In http {} block
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_status 429;

# In server/location {} block
location / {
    limit_req zone=api_limit burst=20 nodelay;
    # ...
}
```

**Parameter:**
- `rate=10r/s` - 10 requests per second
- `burst=20` - Burst up to 20 requests
- `nodelay` - Immediately 429 if limit exceeded

### Per-IP Blocking

```nginx
# Geographic blocking (requires GeoIP)
geo $blocked_country {
    default 0;
    CN 1; # China
    RU 1; # Russia
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

### Testing

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
# Custom log format with details
log_format api_access '$remote_addr - $remote_user [$time_local] '
                      '"$request" $status $body_bytes_sent '
                      ' '$http_referer' '$http_user_agent' '
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

# Count 401 responses (last 5 minutes)
FAILURES=$(tail -5000 "$LOG" | grep " 401 " | wc -l)

if [ "$FAILURES" -gt "$THRESHOLD" ]; then
    echo "WARNING: $FAILURES failed auth attempts in last 5 minutes!" | \
      mail -s "API Auth Failures Alert" "$EMAIL"
fi
```

**Cron:**
```bash
*/5 * * * * * /usr/local/bin/monitor-auth-failures.sh
```

---

## 🔐 System Hardening

### Minimum permissions

```bash
# API Binary
sudo chown root:root /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod 500 /opt/rsyslog-rest-api/rsyslog-rest-api

# Config
sudo chown root:root /opt/rsyslog-rest-api/.env
sudo chmod 400 /opt/rsyslog-rest-api/.env
```

### systemd sandboxing

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
# SELinux status
getenforce

# For API port
sudo semanage port -a -t http_port_t -p tcp 8000

# For File Access
sudo chcon -R -t httpd_sys_content_t /opt/rsyslog-rest-api/
```

---

## 🚨 Incident Response

### In case of compromise

**1. stop API immediately:**
```bash
sudo systemctl stop rsyslog-rest-api
```

**2. change API key:**
```bash
NEW_KEY=$(openssl rand -hex 32)
sudo sed -i "s/^API_KEY=.*/API_KEY=$NEW_KEY/" /opt/rsyslog-rest-api/.env
```

**3. analyze logs:**
```bash
# Suspicious requests
sudo grep " 401 " /var/log/nginx/rsyslog-api-access.log
sudo journalctl -u rsyslog-rest-api --since "1 hour ago"

# Collect IPs
sudo awk '{print $1}' /var/log/nginx/rsyslog-api-access.log | \
  sort | uniq -c | sort -rn | head -20
```

**4. block suspicious IPs:**
```bash
# nginx
sudo nano /etc/nginx/sites-available/rsyslog-api
# Add: deny SUSPICIOUS_IP;

# ufw
sudo ufw deny from SUSPICIOUS_IP

# iptables
sudo iptables -A INPUT -s SUSPICIOUS_IP -j DROP
```

**5. inform clients:**
```bash
# New API key to authorized users
```

**6. check system:**
```bash
# Rootkit check
sudo rkhunter --check
sudo chkrootkit

# Integrity Check
sudo debsums -c # Debian/Ubuntu
sudo rpm -Va # CentOS/RHEL
```

---

## 📋 Security Audit Checklist

### Monthly

- [ ] Rotate API key (optional)
- [ ] Check failed auth attempts
- [ ] Check SSL certificate expiration
- [ ] Updates available?
- [ ] Check logs for anomalies

### Quarterly

- [ ] Rotate API key (recommended)
- [ ] Change database password
- [ ] Test backup strategy
- [ ] Verify permissions
- [ ] Test security headers

### Annually

- [ ] Complete security audit
- [ ] Consider penetration test
- [ ] Update documentation
- [ ] Test disaster recovery plan

---

## 🔗 More resources

- **OWASP Top 10:** https://owasp.org/www-project-top-ten/
- **Let's Encrypt:** https://letsencrypt.org/
- **Mozilla SSL Config Generator:** https://ssl-config.mozilla.org/
- **Security Headers:** https://securityheaders.com/

---

[← Back to overview](index.md) | [Next to Performance →](performance.md)

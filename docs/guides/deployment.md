# Deployment Guide

Production deployment guide for rsyslog REST API.

## Production Checklist

Before deploying:

- ✅ API key configured
- ✅ SSL/TLS certificates ready
- ✅ Database credentials secured
- ✅ Firewall rules configured
- ✅ Reverse proxy setup
- ✅ Monitoring configured
- ✅ Backup strategy in place

## systemd Service

### Installation

```bash
# Copy service file
sudo cp rsyslog-rest-api.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload

# Enable and start
sudo systemctl enable --now rsyslog-rest-api

# Check status
sudo systemctl status rsyslog-rest-api
```

### Service File

`/etc/systemd/system/rsyslog-rest-api.service`:

```ini
[Unit]
Description=rsyslog REST API
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/opt/rsyslog-rest-api
EnvironmentFile=/opt/rsyslog-rest-api/.env
ExecStart=/opt/rsyslog-rest-api/rsyslog-rest-api
Restart=on-failure
RestartSec=5s

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/rsyslog-rest-api

[Install]
WantedBy=multi-user.target
```

## Reverse Proxy

### nginx (Recommended)

```nginx
# /etc/nginx/sites-available/rsyslog-api

upstream rsyslog_api {
    server 127.0.0.1:8000;
}

# Rate limiting
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

server {
    listen 443 ssl http2;
    server_name api.example.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;

    # Logging
    access_log /var/log/nginx/rsyslog-api-access.log;
    error_log /var/log/nginx/rsyslog-api-error.log;

    location / {
        # Rate limiting
        limit_req zone=api_limit burst=20 nodelay;

        proxy_pass http://rsyslog_api;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Timeouts
        proxy_connect_timeout 10s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
}

# HTTP redirect
server {
    listen 80;
    server_name api.example.com;
    return 301 https://$server_name$request_uri;
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/rsyslog-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### SSL with Let's Encrypt

```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d api.example.com

# Auto-renewal (already configured)
sudo systemctl status certbot.timer
```

## Firewall Configuration

### ufw (Ubuntu/Debian)

```bash
# Allow API port (if direct access)
sudo ufw allow 8000/tcp

# Or only allow nginx
sudo ufw allow 'Nginx Full'

# Allow SSH
sudo ufw allow 22/tcp

# Enable
sudo ufw enable
```

### firewalld (CentOS/RHEL)

```bash
# Allow HTTPS
sudo firewall-cmd --permanent --add-service=https

# Reload
sudo firewall-cmd --reload
```

## Monitoring

### Health Check

```bash
# Simple check
curl https://api.example.com/health

# Monitoring script
#!/bin/bash
if ! curl -sf https://api.example.com/health > /dev/null; then
    echo "API is down!" | mail -s "API Alert" admin@example.com
fi
```

### Logs

```bash
# API logs
sudo journalctl -u rsyslog-rest-api -f

# nginx logs
sudo tail -f /var/log/nginx/rsyslog-api-access.log
sudo tail -f /var/log/nginx/rsyslog-api-error.log
```

### Performance Metrics

```bash
# Query performance
time curl -H "X-API-Key: $KEY" "https://api.example.com/logs?limit=1000"

# Database connections
mysql -e "SHOW STATUS LIKE 'Threads_connected'"

# System resources
htop
iostat -x 1
```

## Backup & Recovery

### Database Backup

```bash
# Backup SystemEvents table
mysqldump -u rsyslog -p Syslog SystemEvents > backup_$(date +%Y%m%d).sql

# Automated daily backup
cat > /etc/cron.daily/backup-syslog << 'BACKUP'
#!/bin/bash
mysqldump -u rsyslog -p"password" Syslog SystemEvents | \
  gzip > /backup/syslog_$(date +%Y%m%d).sql.gz
find /backup -name "syslog_*.sql.gz" -mtime +30 -delete
BACKUP

chmod +x /etc/cron.daily/backup-syslog
```

### Configuration Backup

```bash
# Backup API config
sudo cp /opt/rsyslog-rest-api/.env /backup/api-env-$(date +%Y%m%d)
```

## Scaling

### Vertical Scaling

```go
// In main.go - increase connection pool
db.SetMaxOpenConns(50)   // Default: 25
db.SetMaxIdleConns(10)   // Default: 5
```

### Horizontal Scaling

```
┌─────────┐
│ nginx   │ Load Balancer
│ (LB)    │
└────┬────┘
     │
     ├──────────┬──────────┐
     │          │          │
┌────▼───┐ ┌────▼───┐ ┌────▼───┐
│ API 1  │ │ API 2  │ │ API 3  │
└────┬───┘ └────┬───┘ └────┬───┘
     │          │          │
     └──────────┴──────────┘
            │
       ┌────▼─────┐
       │ MySQL    │
       │ (Master) │
       └──────────┘
```

nginx config for load balancing:

```nginx
upstream rsyslog_api_cluster {
    least_conn;
    server 10.0.1.10:8000;
    server 10.0.1.11:8000;
    server 10.0.1.12:8000;
}
```

## Maintenance

### Update API

```bash
# Download new version
wget https://github.com/.../rsyslog-rest-api-linux-amd64

# Stop service
sudo systemctl stop rsyslog-rest-api

# Backup old binary
sudo cp /opt/rsyslog-rest-api/rsyslog-rest-api \
       /opt/rsyslog-rest-api/rsyslog-rest-api.backup

# Replace binary
sudo mv rsyslog-rest-api-linux-amd64 /opt/rsyslog-rest-api/rsyslog-rest-api
sudo chmod +x /opt/rsyslog-rest-api/rsyslog-rest-api

# Start service
sudo systemctl start rsyslog-rest-api

# Verify
curl https://api.example.com/health
```

### Log Rotation

```bash
# /etc/logrotate.d/rsyslog-api
/var/log/rsyslog-rest-api.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 root root
    postrotate
        systemctl reload rsyslog-rest-api > /dev/null 2>&1 || true
    endscript
}
```

## Testing Production Setup

```bash
# 1. Health check
curl https://api.example.com/health

# 2. Authentication
curl -H "X-API-Key: $KEY" https://api.example.com/logs?limit=1

# 3. HTTPS redirect
curl -I http://api.example.com  # Should redirect to HTTPS

# 4. Rate limiting
for i in {1..100}; do 
    curl https://api.example.com/health &
done
# Should see 429 errors

# 5. Load test
ab -n 1000 -c 10 https://api.example.com/health
```

## Security Hardening

See [Security Guide](security.md) for complete security configuration.

# Security Guide

Security best practices for rsyslog REST API.

## Security Checklist

- ✅ Strong API key (32+ bytes)
- ✅ SSL/TLS enabled
- ✅ CORS properly configured
- ✅ Database read-only user
- ✅ Firewall configured
- ✅ Rate limiting enabled
- ✅ Logs monitored
- ✅ Regular security audits

## API Key Security

### Generate Strong Keys

```bash
# Minimum 32 bytes (64 hex characters)
openssl rand -hex 32

# Save to .env
echo "API_KEY=$(openssl rand -hex 32)" >> .env
```

### Secure Storage

```bash
# Restrict permissions
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
```

### Key Rotation

```bash
# Generate new key
NEW_KEY=$(openssl rand -hex 32)

# Update .env
sudo sed -i "s/^API_KEY=.*/API_KEY=$NEW_KEY/" /opt/rsyslog-rest-api/.env

# Restart service
sudo systemctl restart rsyslog-rest-api
```

## SSL/TLS Configuration

### Production (Let's Encrypt)

```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d api.example.com

# Auto-renewal
sudo certbot renew --dry-run
```

### Development (Self-Signed)

```bash
# Generate certificate
openssl req -x509 -newkey rsa:4096 -nodes \
  -keyout key.pem -out cert.pem -days 365 \
  -subj "/CN=localhost"

# Configure
USE_SSL=true
SSL_CERTFILE=/path/to/cert.pem
SSL_KEYFILE=/path/to/key.pem
```

## CORS Security

### Production Configuration

```bash
# NEVER use * in production!
ALLOWED_ORIGINS=https://dashboard.example.com,https://monitoring.example.com
```

### Testing

```bash
# Test CORS
curl -H "Origin: https://dashboard.example.com" \
     -H "Access-Control-Request-Method: GET" \
     -X OPTIONS \
     https://api.example.com/logs
```

## Database Security

### Read-Only User (Recommended)

```sql
-- Create read-only user
CREATE USER 'rsyslog_api'@'localhost' IDENTIFIED BY 'strong-password';
GRANT SELECT ON Syslog.* TO 'rsyslog_api'@'localhost';
FLUSH PRIVILEGES;
```

Update `.env`:
```bash
DB_USER=rsyslog_api
DB_PASS=strong-password
```

### Remote Database

```bash
# Use SSL for remote connections
DB_HOST=db.example.com

# MySQL SSL config (my.cnf)
[client]
ssl-ca=/path/to/ca.pem
ssl-cert=/path/to/client-cert.pem
ssl-key=/path/to/client-key.pem
```

## Firewall Configuration

### ufw

```bash
# Allow only nginx
sudo ufw allow 'Nginx Full'
sudo ufw deny 8000/tcp  # Block direct API access

# Allow SSH (be careful!)
sudo ufw allow 22/tcp

# Enable
sudo ufw enable
```

### IP Whitelisting

```nginx
# nginx - Allow only specific IPs
location / {
    allow 10.0.1.0/24;     # Internal network
    allow 203.0.113.0/24;  # Office network
    deny all;
    
    proxy_pass http://rsyslog_api;
}
```

## Rate Limiting

### nginx

```nginx
# Define rate limit zone
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_zone $http_x_api_key zone=api_key_limit:10m rate=100r/s;

server {
    location / {
        # Per IP: 10 req/s
        limit_req zone=api_limit burst=20 nodelay;
        
        # Per API key: 100 req/s
        limit_req zone=api_key_limit burst=200 nodelay;
        
        # Return 429 on limit
        limit_req_status 429;
        
        proxy_pass http://rsyslog_api;
    }
}
```

### iptables

```bash
# Limit connections per IP
sudo iptables -A INPUT -p tcp --dport 8000 -m connlimit \
  --connlimit-above 10 -j REJECT

# Limit new connections per minute
sudo iptables -A INPUT -p tcp --dport 8000 -m recent --set
sudo iptables -A INPUT -p tcp --dport 8000 -m recent --update \
  --seconds 60 --hitcount 100 -j REJECT
```

## Security Headers

### nginx Configuration

```nginx
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
add_header X-Frame-Options "DENY" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Content-Security-Policy "default-src 'self'" always;
```

## Monitoring & Auditing

### Failed Authentication Attempts

```bash
# Monitor auth failures
sudo journalctl -u rsyslog-rest-api | grep "Invalid or missing API key"

# Alert on multiple failures
#!/bin/bash
FAILURES=$(sudo journalctl -u rsyslog-rest-api --since "1 hour ago" | \
           grep -c "Invalid or missing API key")

if [ "$FAILURES" -gt 100 ]; then
    echo "High number of auth failures: $FAILURES" | \
      mail -s "Security Alert" admin@example.com
fi
```

### Access Logging

```nginx
# Detailed access log
log_format api_log '$remote_addr - [$time_local] '
                   '"$request" $status $body_bytes_sent '
                   '"$http_referer" "$http_user_agent" '
                   '$request_time $upstream_response_time '
                   '$http_x_api_key';

access_log /var/log/nginx/api-access.log api_log;
```

## System Hardening

### systemd Sandboxing

```ini
[Service]
# Security enhancements
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/rsyslog-rest-api
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true
RestrictRealtime=true
RestrictNamespaces=true
LockPersonality=true
MemoryDenyWriteExecute=true
RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX
```

### SELinux (CentOS/RHEL)

```bash
# Check SELinux status
sestatus

# Allow API to connect to database
sudo setsebool -P httpd_can_network_connect_db 1

# Allow API to bind to port
sudo semanage port -a -t http_port_t -p tcp 8000
```

## Incident Response

### Compromised API Key

```bash
# 1. Generate new key immediately
NEW_KEY=$(openssl rand -hex 32)

# 2. Update configuration
sudo sed -i "s/^API_KEY=.*/API_KEY=$NEW_KEY/" /opt/rsyslog-rest-api/.env

# 3. Restart service
sudo systemctl restart rsyslog-rest-api

# 4. Update all clients
# Notify users to update their API keys

# 5. Audit access logs
sudo grep "old-key" /var/log/nginx/api-access.log
```

### Suspicious Activity

```bash
# Check recent access
sudo journalctl -u rsyslog-rest-api --since "1 hour ago"

# Check nginx logs
sudo tail -1000 /var/log/nginx/api-access.log | \
  awk '{print $1}' | sort | uniq -c | sort -rn

# Block suspicious IP
sudo ufw deny from 203.0.113.50
```

## Security Audit Checklist

### Monthly

- [ ] Review access logs for anomalies
- [ ] Check for failed authentication attempts
- [ ] Verify SSL certificate expiry
- [ ] Review firewall rules
- [ ] Check for unauthorized database access

### Quarterly

- [ ] Rotate API keys
- [ ] Update dependencies
- [ ] Review user permissions
- [ ] Penetration testing
- [ ] Security scan with tools

### Yearly

- [ ] Full security audit
- [ ] Review disaster recovery plan
- [ ] Update incident response procedures
- [ ] Security training for team

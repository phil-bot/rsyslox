# Deployment Guide

Production deployment guide for rsyslox.

## Production Checklist

- ✅ Setup wizard completed
- ✅ Admin password set (12+ characters)
- ✅ SSL/TLS certificates in place
- ✅ CORS origins restricted in Admin → Server
- ✅ Firewall rules configured
- ✅ Reverse proxy set up (recommended)
- ✅ Monitoring configured

## Install

Follow the [Installation Guide](../getting-started/installation.md) to install the binary and complete the setup wizard.

## systemd Service

The installer handles this automatically. For reference, the service file at `/etc/systemd/system/rsyslox.service`:

```ini
[Unit]
Description=rsyslox — syslog viewer
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=rsyslox
Group=rsyslox
ExecStart=/opt/rsyslox/rsyslox
Restart=on-failure
RestartSec=5s
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl status rsyslox
sudo systemctl restart rsyslox
sudo journalctl -u rsyslox -f
```

## Reverse Proxy

Run rsyslox behind a reverse proxy to terminate TLS and add rate limiting.

### nginx (Recommended)

```nginx
# /etc/nginx/sites-available/rsyslox

upstream rsyslox {
    server 127.0.0.1:8000;
}

limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

server {
    listen 443 ssl http2;
    server_name rsyslox.example.com;

    ssl_certificate     /etc/letsencrypt/live/rsyslox.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/rsyslox.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;

    add_header Strict-Transport-Security "max-age=31536000" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;

    location / {
        limit_req zone=api_limit burst=20 nodelay;
        proxy_pass http://rsyslox;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name rsyslox.example.com;
    return 301 https://$server_name$request_uri;
}
```

```bash
sudo ln -s /etc/nginx/sites-available/rsyslox /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

### Apache

```apache
<VirtualHost *:443>
    ServerName rsyslox.example.com

    SSLEngine on
    SSLCertificateFile /etc/letsencrypt/live/rsyslox.example.com/fullchain.pem
    SSLCertificateKeyFile /etc/letsencrypt/live/rsyslox.example.com/privkey.pem

    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:8000/
    ProxyPassReverse / http://127.0.0.1:8000/

    Header always set Strict-Transport-Security "max-age=31536000"
    Header always set X-Frame-Options "DENY"
</VirtualHost>
```

## SSL/TLS

### Let's Encrypt (Recommended for Production)

```bash
sudo apt-get install certbot python3-certbot-nginx
sudo certbot --nginx -d rsyslox.example.com
sudo certbot renew --dry-run   # test auto-renewal
```

### Built-in SSL (Alternative)

If you prefer not to use a reverse proxy, rsyslox can terminate TLS directly — see [Configuration → SSL](../getting-started/configuration.md#ssl--tls) for setup instructions.

## Firewall

```bash
# UFW (Ubuntu/Debian)
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw deny 8000/tcp    # block direct access when using reverse proxy
sudo ufw enable

# firewalld (CentOS/RHEL)
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

## Database Security

Create a dedicated database user with minimal permissions. See [Security Guide → Database](security.md#database) for the required SQL.

## Monitoring

```bash
# Health check
curl http://localhost:8000/health
```

**Cron-based alert:**
```bash
#!/bin/bash
# /usr/local/bin/rsyslox-healthcheck.sh
if ! curl -sf "http://localhost:8000/health" > /dev/null; then
    echo "rsyslox health check failed" | mail -s "Alert" admin@example.com
    systemctl restart rsyslox
fi
```

```bash
# Check every 5 minutes
*/5 * * * * /usr/local/bin/rsyslox-healthcheck.sh
```

## Updates

```bash
# 1. Download new release
wget https://github.com/phil-bot/rsyslox/releases/download/vX.Y.Z/rsyslox-linux-amd64

# 2. Stop service
sudo systemctl stop rsyslox

# 3. Backup current binary
sudo cp /opt/rsyslox/rsyslox /opt/rsyslox/rsyslox.bak

# 4. Replace binary
sudo mv rsyslox-linux-amd64 /opt/rsyslox/rsyslox
sudo chmod +x /opt/rsyslox/rsyslox
sudo chown rsyslox:rsyslox /opt/rsyslox/rsyslox

# 5. Start service
sudo systemctl start rsyslox
sudo systemctl status rsyslox
curl http://localhost:8000/health
```

Configuration is preserved across updates.

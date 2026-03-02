# Security Guide

Security best practices for rsyslox.

## Security Checklist

- ✅ Strong admin password (12+ characters)
- ✅ Read-only API keys for external tools
- ✅ SSL/TLS enabled (via reverse proxy or built-in)
- ✅ CORS origins restricted to specific domains
- ✅ Database user has read-only access
- ✅ Firewall blocks direct port 8000 access (if using reverse proxy)
- ✅ Rate limiting in reverse proxy

## Authentication

rsyslox uses two separate authentication mechanisms:

**Admin session token** — full access. Obtained via `POST /api/admin/login`. Used for the web UI and admin API endpoints. Stored in the browser's `sessionStorage`.

**Read-only API key** — restricted access. Created in **Admin → API Keys**. Can only access `/api/logs` and `/api/meta`. Keys are stored as SHA-256 hashes; plaintext is shown only once at creation.

### API Key Best Practices

- Create one key per consumer (monitoring system, dashboard, script) so each can be revoked independently
- Revoke compromised keys immediately in Admin → API Keys
- Never commit key values to version control

### Admin Password

- Minimum 12 characters — use a passphrase or password manager
- The password is stored as a bcrypt hash (cost 12) — it cannot be recovered, only reset via the `hash-password` CLI command. See [Troubleshooting → Authentication Issues](troubleshooting.md#authentication-issues) for the exact steps.

## SSL/TLS

Use HTTPS in production. Options:

**Reverse proxy (recommended):** nginx or Apache handle TLS termination. See [Deployment Guide](deployment.md).

**Built-in SSL:** rsyslox can terminate TLS directly. Enable SSL in **Admin → Server**, then either:
- Click **Generate Self-Signed Certificate** for a self-signed ECDSA P-256 cert (development/internal use), or
- Upload your own certificate and key via **Upload Custom Certificate**.

If `use_ssl = true` is set in `config.toml` and no certificate files exist, rsyslox generates a self-signed certificate automatically on startup.

## CORS

Restrict origins in **Admin → Server → CORS origins**. Never leave `*` in production:

```
# Single origin
https://dashboard.example.com

# Multiple (comma-separated)
https://app1.example.com,https://monitoring.example.com
```

## Database

Use a read-only database user:

```sql
CREATE USER 'rsyslox'@'localhost' IDENTIFIED BY 'strong-password';
GRANT SELECT ON Syslog.SystemEvents TO 'rsyslox'@'localhost';
FLUSH PRIVILEGES;
```

Set this user in the setup wizard. If you later need the cleanup service (which requires `DELETE`), grant that additionally:

```sql
GRANT DELETE ON Syslog.SystemEvents TO 'rsyslox'@'localhost';
FLUSH PRIVILEGES;
```

## Firewall

If running behind a reverse proxy, block direct access to port 8000:

```bash
sudo ufw deny 8000/tcp
sudo ufw allow 443/tcp
sudo ufw allow 80/tcp
```

## Rate Limiting (nginx)

```nginx
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

location / {
    limit_req zone=api_limit burst=20 nodelay;
    limit_req_status 429;
    proxy_pass http://rsyslox;
}
```

## systemd Sandboxing

The installer applies these restrictions automatically via the service file:

```ini
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
```

## Configuration File

`/etc/rsyslox/config.toml` contains sensitive values and is protected:

- File mode `0640` — owner `root`, group `rsyslox`
- Database password is AES-GCM encrypted (key derived from `/etc/machine-id`)
- Admin password is a bcrypt hash
- API key plaintext is never stored — only SHA-256 hashes

## Incident Response

### Compromised API Key

1. Open **Admin → API Keys**
2. Click **Revoke** next to the compromised key
3. Create a new key and distribute it to the affected consumer

### Suspected Admin Password Compromise

```bash
# 1. Generate new bcrypt hash
/opt/rsyslox/rsyslox hash-password "new-strong-password"

# 2. Update config
sudo nano /etc/rsyslox/config.toml
# Replace admin_password_hash value

# 3. Restart
sudo systemctl restart rsyslox

# 4. Immediately revoke all API keys and reissue them
# (existing session tokens expire on restart)
```

## Security Audit Checklist

**Monthly**
- [ ] Review access logs for anomalies: `sudo journalctl -u rsyslox -n 1000`
- [ ] Check SSL certificate expiry
- [ ] Review active API keys in Admin panel — revoke unused ones

**Quarterly**
- [ ] Rotate read-only API keys
- [ ] Review CORS origin list
- [ ] Check for available rsyslox updates

**Yearly**
- [ ] Rotate admin password
- [ ] Full security review

## More Resources

- [Deployment Guide](deployment.md) — nginx/Apache config, TLS setup
- [Troubleshooting](troubleshooting.md) — Auth error diagnosis

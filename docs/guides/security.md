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
- Stored as a bcrypt hash (cost 12) — cannot be recovered, only reset via the `hash-password` CLI command (see [Troubleshooting → Authentication Issues](troubleshooting.md#authentication-issues))

## SSL/TLS

Use HTTPS in production. Options:

**Reverse proxy (recommended):** nginx or Apache handle TLS termination — see [Deployment Guide](deployment.md).

**Built-in SSL:** rsyslox can terminate TLS directly. Enable SSL in **Admin → Server**, then either:
- Click **Generate Self-Signed Certificate** for a self-signed ECDSA P-256 cert (development/internal use), or
- Upload your own certificate and key via **Upload Custom Certificate**

If `use_ssl = true` is set and no certificate files exist, rsyslox generates a self-signed certificate automatically on startup.

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

If the cleanup service is enabled, it additionally requires `DELETE`:

```sql
GRANT DELETE ON Syslog.SystemEvents TO 'rsyslox'@'localhost';
FLUSH PRIVILEGES;
```

## Firewall

If running behind a reverse proxy, block direct access to port 8000 — see [Deployment Guide → Firewall](deployment.md#firewall).

## Rate Limiting

Configure rate limiting in your reverse proxy. The nginx example in the [Deployment Guide](deployment.md#nginx-recommended) includes a ready-to-use `limit_req` setup.

## systemd Sandboxing

The installer applies these restrictions automatically:

```ini
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
```

## Configuration File

`/etc/rsyslox/config.toml` is protected:

| Value | Storage |
|---|---|
| Database password | AES-GCM encrypted; key derived from `/etc/machine-id` — not portable between machines |
| Admin password | bcrypt hash (cost 12) |
| API key plaintext | Never stored; only SHA-256 hex hash written to disk |
| Config file | Mode `0640` — owner `root`, group `rsyslox` |

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

# 3. Restart (invalidates all existing session tokens)
sudo systemctl restart rsyslox

# 4. Revoke all API keys and reissue them
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

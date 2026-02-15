# Configuration

Complete configuration reference for rsyslog REST API.

## Configuration File

The API uses environment variables loaded from a `.env` file.

**Location:** `/opt/rsyslog-rest-api/.env`

## Configuration Parameters

### API Security

#### API_KEY

API key for authentication.

```bash
API_KEY=a3d7f8c9e2b4...
```

**Generate secure key:**
```bash
openssl rand -hex 32
```

**Important:** Always use a strong API key in production!

### Server Configuration

#### SERVER_HOST

Host address to bind to.

```bash
SERVER_HOST=0.0.0.0  # Listen on all interfaces (default)
```

Options:
- `0.0.0.0` - All interfaces (default)
- `127.0.0.1` - Localhost only
- `192.168.1.10` - Specific IP

#### SERVER_PORT

Port to listen on.

```bash
SERVER_PORT=8000  # Default: 8000
```

### Database Configuration

#### Recommended: Environment Variables

```bash
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=your-secure-password
```

**Benefits:**
- ✅ Secure (no need to expose rsyslog config)
- ✅ Flexible (different DB for API)
- ✅ Clear separation of concerns

#### Alternative: rsyslog Config

If DB_* variables are not set, reads from rsyslog config:

```bash
RSYSLOG_CONFIG_PATH=/etc/rsyslog.d/mysql.conf
```

### SSL/TLS Configuration

#### USE_SSL

Enable SSL/TLS.

```bash
USE_SSL=false  # Default: false
```

Set to `true` for production with HTTPS.

#### SSL_CERTFILE

Path to SSL certificate.

```bash
SSL_CERTFILE=/opt/rsyslog-rest-api/certs/cert.pem
```

#### SSL_KEYFILE

Path to SSL private key.

```bash
SSL_KEYFILE=/opt/rsyslog-rest-api/certs/key.pem
```

**Generate self-signed certificate (development):**
```bash
openssl req -x509 -newkey rsa:4096 -nodes \
  -keyout key.pem -out cert.pem -days 365 \
  -subj "/CN=localhost"
```

**Production:** Use Let's Encrypt certificates!

### CORS Configuration

#### ALLOWED_ORIGINS

Allowed origins for CORS.

```bash
ALLOWED_ORIGINS=*  # Allow all (development only!)
```

**Production examples:**
```bash
# Single origin
ALLOWED_ORIGINS=https://dashboard.example.com

# Multiple origins (comma-separated)
ALLOWED_ORIGINS=https://app1.com,https://app2.com
```

**Security:** Never use `*` in production!

## Complete Example

### Development Configuration

```bash
# .env (development)
API_KEY=dev-key-12345
SERVER_HOST=127.0.0.1
SERVER_PORT=8000
USE_SSL=false
ALLOWED_ORIGINS=*

DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=devpassword
```

### Production Configuration

```bash
# .env (production)
API_KEY=a3d7f8c9e2b4a1c8f7e9d3b6a5c4e8f2
SERVER_HOST=0.0.0.0
SERVER_PORT=8000
USE_SSL=true
SSL_CERTFILE=/etc/letsencrypt/live/api.example.com/fullchain.pem
SSL_KEYFILE=/etc/letsencrypt/live/api.example.com/privkey.pem
ALLOWED_ORIGINS=https://dashboard.example.com,https://monitoring.example.com

DB_HOST=db.internal.example.com
DB_NAME=Syslog
DB_USER=rsyslog_api
DB_PASS=strong-secure-password-here
```

## File Permissions

Secure your configuration:

```bash
# Restrict access to .env
sudo chmod 600 /opt/rsyslog-rest-api/.env
sudo chown root:root /opt/rsyslog-rest-api/.env
```

## Environment Variables Priority

Configuration is loaded in this order:

1. **Environment variables** (highest priority)
2. **.env file**
3. **Default values**

Example:
```bash
# Override via environment
export SERVER_PORT=9000
rsyslog-rest-api  # Will use port 9000
```

## Validation

Test your configuration:

```bash
# Start in foreground
rsyslog-rest-api

# Check logs
sudo journalctl -u rsyslog-rest-api -n 50
```

Look for:
- ✅ "Database connection established"
- ✅ "Starting HTTP server on..."
- ❌ Any error messages

## Next Steps

- [Quick Start Tutorial](quick-start.md)
- [Deploy to Production](../guides/deployment.md)
- [Security Best Practices](../guides/security.md)

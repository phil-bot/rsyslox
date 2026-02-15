# Troubleshooting Guide

Common issues and solutions.

## Quick Diagnosis

```bash
# Check API status
sudo systemctl status rsyslog-rest-api

# View logs
sudo journalctl -u rsyslog-rest-api -n 50

# Test database
mysql -u rsyslog -p Syslog -e "SELECT COUNT(*) FROM SystemEvents"

# Test API
curl http://localhost:8000/health
```

## Common Issues

### API Won't Start

**Symptoms:**
- Service fails to start
- `sudo systemctl status` shows "failed"

**Solutions:**

```bash
# Check logs
sudo journalctl -u rsyslog-rest-api -n 100

# Common causes:
# 1. Database connection failed
mysql -u rsyslog -p Syslog  # Test manually

# 2. Port already in use
sudo lsof -i :8000

# 3. Permission denied on .env
sudo chmod 600 /opt/rsyslog-rest-api/.env

# 4. Binary not found
ls -la /opt/rsyslog-rest-api/rsyslog-rest-api
```

### Database Connection Failed

**Error:** `failed to ping database`

**Solutions:**

```bash
# 1. Check database is running
sudo systemctl status mysql

# 2. Test credentials
mysql -u rsyslog -p Syslog

# 3. Check .env file
cat /opt/rsyslog-rest-api/.env | grep DB_

# 4. Verify database exists
mysql -e "SHOW DATABASES"

# 5. Check user permissions
mysql -e "SHOW GRANTS FOR 'rsyslog'@'localhost'"
```

### API Key Issues

**Error:** `Invalid or missing API key`

**Solutions:**

```bash
# 1. Check API key in .env
grep "^API_KEY=" /opt/rsyslog-rest-api/.env

# 2. Verify header format
curl -H "X-API-Key: your-key" ...  # Correct
curl -H "API-Key: your-key" ...    # Wrong!

# 3. Check for spaces
# Wrong: API_KEY= abc123 (has space)
# Right: API_KEY=abc123
```

### No Logs Returned

**Symptoms:**
- API works but returns 0 logs
- `"total": 0`

**Solutions:**

```bash
# 1. Check database has data
mysql Syslog -e "SELECT COUNT(*) FROM SystemEvents"

# 2. Check time range
# Default is last 24 hours
# Try: ?start_date=2020-01-01T00:00:00Z

# 3. Remove all filters
curl ... /logs?limit=10  # No filters

# 4. Check column names
# Wrong: ?Host=web01
# Right: ?FromHost=web01
```

### SSL/TLS Problems

**Error:** `SSL certificate not found`

**Solutions:**

```bash
# 1. Check certificate files exist
ls -la /opt/rsyslog-rest-api/certs/

# 2. Check paths in .env
SSL_CERTFILE=/correct/path/to/cert.pem
SSL_KEYFILE=/correct/path/to/key.pem

# 3. Test certificate
openssl x509 -in cert.pem -text -noout

# 4. For development, disable SSL
USE_SSL=false
```

### CORS Errors

**Error:** Browser shows CORS policy error

**Solutions:**

```bash
# 1. Check ALLOWED_ORIGINS
ALLOWED_ORIGINS=https://your-domain.com

# 2. For development
ALLOWED_ORIGINS=*

# 3. Multiple origins (comma-separated)
ALLOWED_ORIGINS=https://app1.com,https://app2.com

# 4. Restart after changes
sudo systemctl restart rsyslog-rest-api
```

### Performance Issues

**Symptoms:**
- Slow responses (>1s)
- High CPU/memory

**Solutions:**

```bash
# 1. Check database performance
mysql -e "SHOW PROCESSLIST"

# 2. Add indexes
# See: performance.md

# 3. Reduce query size
?limit=100  # Instead of limit=10000

# 4. Optimize time range
# Last hour instead of last 30 days

# 5. Check system resources
htop
iostat -x 1
```

## FAQ

### How do I reset the API key?

```bash
# Generate new key
NEW_KEY=$(openssl rand -hex 32)

# Update .env
sudo sed -i "s/^API_KEY=.*/API_KEY=$NEW_KEY/" /opt/rsyslog-rest-api/.env

# Restart
sudo systemctl restart rsyslog-rest-api
```

### How do I enable debug logging?

```bash
# Check current logs
sudo journalctl -u rsyslog-rest-api -f

# For more details, start in foreground
sudo systemctl stop rsyslog-rest-api
cd /opt/rsyslog-rest-api
./rsyslog-rest-api  # Shows all logs
```

### Why are my filters not working?

```bash
# Common mistakes:

# ❌ Wrong: Comma-separated
?FromHost=web01,web02

# ✅ Right: Repeat parameter
?FromHost=web01&FromHost=web02

# ❌ Wrong: Lowercase
?fromhost=web01

# ✅ Right: CamelCase
?FromHost=web01
```

### How do I upgrade to a new version?

See [Deployment: Maintenance](deployment.md#maintenance)

### Can I use this with PostgreSQL?

Currently only MySQL/MariaDB is supported. PostgreSQL support is planned for future versions.

## Getting Help

- [GitHub Issues](https://github.com/phil-bot/rsyslog-rest-api/issues)
- [GitHub Discussions](https://github.com/phil-bot/rsyslog-rest-api/discussions)
- Check logs: `sudo journalctl -u rsyslog-rest-api -n 100`

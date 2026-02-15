# rsyslog REST API

> High-performance REST API for rsyslog/MySQL written in Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/phil-bot/rsyslog-rest-api/blob/main/LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![Version](https://img.shields.io/badge/version-v0.2.3-success.svg)](https://github.com/phil-bot/rsyslog-rest-api/releases)

## Features

- 🚀 **High Performance** - Compiled in Go for maximum speed
- 🔍 **Advanced Filtering** - Multi-value filters for complex queries
- 📊 **All Fields** - Access to all 25+ SystemEvents columns
- 🔐 **Secure** - API key authentication, SSL/TLS support
- 🐳 **Docker Ready** - Complete test environment with live data
- 📝 **REST API** - Clean JSON responses
- 🎯 **RFC-5424 Compliant** - Proper syslog severity and facility labels

## What's New in v0.2.3

- ✅ **Enhanced Multi-Value Filters** - Improved performance
- ✅ **Better Error Messages** - Clear validation feedback
- ✅ **Extended Meta Endpoint** - More filtering options
- ✅ **Bug Fixes** - Various stability improvements

[View Full Changelog](development/changelog.md)

## Quick Start

### Installation

```bash
# Download latest release
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-amd64

# Make executable
chmod +x rsyslog-rest-api-linux-amd64
sudo mv rsyslog-rest-api-linux-amd64 /usr/local/bin/rsyslog-rest-api
```

### Configuration

```bash
# Create .env file
cat > .env << EOF
API_KEY=$(openssl rand -hex 32)
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=your-password
EOF
```

### Run

```bash
rsyslog-rest-api
```

[Full Installation Guide →](getting-started/installation.md)

## Quick Examples

### Retrieve Logs

```bash
# Latest 10 logs
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?limit=10"

# Errors from multiple hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&Priority=3"
```

### Query Metadata

```bash
# All available hosts
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost"

# Hosts that logged errors
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3&Priority=4"
```

[More Examples →](api/examples.md)

## Documentation

### Getting Started
- [Installation Guide](getting-started/installation.md)
- [Configuration Reference](getting-started/configuration.md)
- [Quick Start Tutorial](getting-started/quick-start.md)

### API Documentation
- [API Reference](api/reference.md) - Complete endpoint documentation
- [Examples](api/examples.md) - Practical usage examples

### Guides
- [Production Deployment](guides/deployment.md)
- [Security Best Practices](guides/security.md)
- [Performance Tuning](guides/performance.md)
- [Troubleshooting](guides/troubleshooting.md)

### Development
- [Docker Testing Environment](development/docker.md)
- [Contributing Guidelines](development/contributing.md)
- [Changelog](development/changelog.md)

## Support

- **Issues:** [GitHub Issues](https://github.com/phil-bot/rsyslog-rest-api/issues)
- **Discussions:** [GitHub Discussions](https://github.com/phil-bot/rsyslog-rest-api/discussions)

## License

MIT License - see [LICENSE](https://github.com/phil-bot/rsyslog-rest-api/blob/main/LICENSE) for details.

---

**Built with ❤️ for the syslog community**

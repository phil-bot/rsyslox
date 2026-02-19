<div align="center">
  <img src="https://rsyslox.grothu.net/rsyslox_light.svg" alt="rsyslox"/>
</div>

# rsyslox

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/phil-bot/rsyslox)](https://github.com/phil-bot/rsyslox/releases)

rsyslox is a lightweight, high-performance REST API for accessing rsyslog data stored in MySQL/MariaDB. It exposes the `SystemEvents` table via a clean HTTP/JSON interface with filtering, pagination, and metadata queries — making it easy to integrate syslog data into dashboards, monitoring tools, or custom scripts.

Written in Go, it is designed to run as a standalone binary with minimal dependencies.

## 📖 Documentation

**[https://rsyslox.grothu.net](https://rsyslox.grothu.net)**

## 🚀 Quick Start

```bash
# Download latest release
wget https://github.com/phil-bot/rsyslox/releases/latest/download/rsyslox-linux-amd64

# Install
chmod +x rsyslox-linux-amd64
sudo mv rsyslox-linux-amd64 /usr/local/bin/rsyslox

# Create configuration
cat > .env << EOF
API_KEY=$(openssl rand -hex 32)
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=your-password
EOF

# Run
rsyslox
```

## 🤝 Support & Community

- **Documentation:** [https://rsyslox.grothu.net](https://rsyslox.grothu.net)
- **Issues:** [GitHub Issues](https://github.com/phil-bot/rsyslox/issues)
- **Discussions:** [GitHub Discussions](https://github.com/phil-bot/rsyslox/discussions)

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

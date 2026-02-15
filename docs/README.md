# rsyslog REST API

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/phil-bot/rsyslog-rest-api)](https://github.com/phil-bot/rsyslog-rest-api/releases)

High-performance REST API for rsyslog/MySQL written in Go. It queries rsyslog data from a MySQL/MariaDB database and makes it accessible via HTTP/JSON. Perfect for monitoring dashboards, log analysis, and system integration.

 → **[Main Repository](https://github.com/phil-bot/rsyslog-rest-api)**

# Documentation

Welcome to the complete documentation of the rsyslog REST API project.

## 📖 Overview

This documentation is divided into different sections depending on your role and requirements.

### Getting Started

| Document | Description |
|----------|-------------|
| [**Installation**](installation.md) | All installation methods (Binary, Source, Package) |
| [**Configuration**](configuration.md) | Complete configuration reference |
| [**Quick Examples**](examples.md) | Practical examples for common use cases |

### Usage & API

| Document | Description |
|----------|-------------|
| [**API Reference**](api-reference.md) | Complete API documentation with all endpoints |
| [**Deployment**](deployment.md) | Production setup, systemd, reverse proxy |
| [**Security**](security.md) | Best practices for secure operation |
| [**Performance**](performance.md) | Optimization and benchmarks |
| [**Troubleshooting**](troubleshooting.md) | Common issues, solutions, and FAQ |
| [**Changelog**](changelog.md) | Version history and breaking changes |

### Test environment

| Document | Description |
|----------|-------------|
| [**Docker Testing**](docker.md) | Test environment with live data |
| [**Development**](development.md) | Architecture, build, contributing |

## 🔍 Quick Access

### Common Tasks

- **Start installation:** → [Installation Guide](installation.md#quick-install)
- **Generate API key:** → [Configuration](configuration.md#api-key)
- **Setup SSL:** → [Security](security.md#ssltls)
- **Deploy to production:** → [Deployment](deployment.md#production-setup)
- **Troubleshoot:** → [Troubleshooting](troubleshooting.md)
- **Test with Docker:** → [Docker Guide](docker.md#quick-start)

### API Endpoints

- **Health check:** → [GET /health](api-reference.md#get-health)
- **Retrieve logs:** → [GET /logs](api-reference.md#get-logs)
- **Metadata:** → [GET /meta](api-reference.md#get-meta)

## 📚 Documentation Structure

```
docs/
├── index.md                 # This file - Overview
│
├── installation.md          # Installation (Binary, Source, Package)
├── configuration.md         # Complete configuration
├── api-reference.md         # API endpoints and parameters
├── examples.md              # Practical examples
├── troubleshooting.md       # Troubleshooting and FAQ
│
├── deployment.md            # Production deployment
├── security.md              # Security best practices
├── performance.md           # Performance tuning
│
├── docker.md                # Docker test environment
├── development.md           # Development and contributing
│
└── changelog.md             # Version history
```

## 🆘 Need Help?

- **GitHub Repository** [Main Repository](https://github.com/phil-bot/rsyslog-rest-api)
- **GitHub Issues:** [Report bugs](https://github.com/phil-bot/rsyslog-rest-api/issues)
- **GitHub Discussions:** [Ask questions](https://github.com/phil-bot/rsyslog-rest-api/discussions)
- **Troubleshooting:** [Browse FAQ](troubleshooting.md#faq)


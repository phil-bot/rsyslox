# rsyslog REST API

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/phil-bot/rsyslog-rest-api)](https://github.com/phil-bot/rsyslog-rest-api/releases)

High-performance REST API for rsyslog/MySQL written in Go.

## 📖 Über das Projekt

Ein moderner REST API Server, der rsyslog-Daten aus einer MySQL/MariaDB-Datenbank abfragt und über HTTP/JSON zugänglich macht. Perfekt für Monitoring-Dashboards, Log-Analysen und System-Integration.

### 🌟 Features

- 🚀 **High Performance** - In Go kompiliert für maximale Geschwindigkeit
- 🔍 **Erweiterte Filter** - Multi-Value-Filter für komplexe Abfragen
- 📊 **Alle Felder** - Zugriff auf alle 25+ SystemEvents-Spalten
- 🔐 **Sicher** - API-Key-Authentifizierung, SSL/TLS-Support
- 🐳 **Docker Ready** - Komplette Testumgebung mit Live-Daten
- 📝 **REST API** - Saubere JSON-Antworten
- 🎯 **RFC-5424** - Korrekte Syslog-Severity und Facility-Labels

## 🚀 Quick Start

### Binary Installation (empfohlen)

```bash
# Download latest release
wget https://github.com/phil-bot/rsyslog-rest-api/releases/latest/download/rsyslog-rest-api-linux-amd64

# Installieren
chmod +x rsyslog-rest-api-linux-amd64
sudo mv rsyslog-rest-api-linux-amd64 /usr/local/bin/rsyslog-rest-api

# Konfiguration erstellen
cat > .env << EOF
API_KEY=$(openssl rand -hex 32)
DB_HOST=localhost
DB_NAME=Syslog
DB_USER=rsyslog
DB_PASS=your-password
EOF

# Starten
rsyslog-rest-api
```

**API testen:**
```bash
curl http://localhost:8000/health
curl -H "X-API-Key: YOUR_KEY" "http://localhost:8000/logs?limit=5"
```

### Docker Testumgebung

Perfekt zum Testen mit live generierten Logs:

```bash
# Build binary
make build-static

# Start container
cd docker && docker-compose up -d

# Test
curl "http://localhost:8000/logs?limit=5"
```

→ [Ausführliche Installation](docs/installation.md)

## 📚 Dokumentation

### Erste Schritte
- [**Installation Guide**](docs/installation.md) - Alle Installationsmethoden
- [**Configuration**](docs/configuration.md) - Vollständige Konfiguration
- [**Quick Examples**](docs/examples.md) - Praktische Beispiele

### API & Nutzung
- [**API Reference**](docs/api-reference.md) - Alle Endpunkte und Parameter
- [**Troubleshooting**](docs/troubleshooting.md) - Fehlersuche und FAQ

### Administration
- [**Deployment**](docs/deployment.md) - Production Setup
- [**Security**](docs/security.md) - Sicherheits-Best-Practices

### Entwicklung
- [**Docker Testing**](docs/docker.md) - Testumgebung
- [**Development**](docs/development.md) - Architektur und Contributing

→ [**Vollständige Dokumentation**](docs/index.md)

## 💡 Beispiele

### Logs mit Filtern abrufen

```bash
# Alle Errors der letzten Stunde
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?Priority=3&start_date=2025-02-09T09:00:00Z"

# Logs von mehreren Hosts (Multi-Value!)
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&FromHost=web02&FromHost=db01"

# Kombinierte Filter
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/logs?FromHost=web01&Priority=3&Priority=4&limit=20"
```

### Metadaten abfragen

```bash
# Alle verfügbaren Hosts
curl -H "X-API-Key: YOUR_KEY" "http://localhost:8000/meta/FromHost"

# Alle Priorities mit Labels
curl -H "X-API-Key: YOUR_KEY" "http://localhost:8000/meta/Priority"

# Hosts die Errors hatten
curl -H "X-API-Key: YOUR_KEY" \
  "http://localhost:8000/meta/FromHost?Priority=3&Priority=4"
```

→ [Weitere Beispiele](docs/examples.md)

## 🆕 Was ist neu in v0.2.2?

- ✅ **Multi-Value-Filter** - Mehrere Werte pro Parameter
- ✅ **Erweiterte Spalten** - Alle 25+ SystemEvents-Felder
- ✅ **Live Log Generator** - Realistische Test-Logs (Docker)
- ✅ **Meta-Endpoint** - Filtert nun auch mit Multi-Value

→ [Changelog](docs/changelog.md)

## 🗺️ Roadmap

### v0.3.0 (Geplant)
- Negation filters (`exclude`, `not`)
- Erweiterte Filter-Kombinationen
- Complex Query Support

### v0.4.0 (Geplant)
- Statistics Endpoint (`/stats`)
- Aggregationen
- Timeline/Histogram

→ [GitHub Issues](https://github.com/phil-bot/rsyslog-rest-api/issues)

## 🤝 Support & Community

- **Issues:** [GitHub Issues](https://github.com/phil-bot/rsyslog-rest-api/issues)
- **Discussions:** [GitHub Discussions](https://github.com/phil-bot/rsyslog-rest-api/discussions)
- **Documentation:** [docs/](docs/index.md)

## 🙏 Contributing

Beiträge sind willkommen! Bitte lies [Contributing Guidelines](docs/development.md#contributing).

1. Fork the repository
2. Create feature branch
3. Make changes & add tests
4. Submit pull request

## 📄 License

MIT License - siehe [LICENSE](LICENSE) für Details.

## ✨ Credits

Erstellt mit ❤️ für die Syslog-Community.

**Built with:**
- [Go](https://go.dev/)
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- [rsyslog](https://www.rsyslog.com/)
- [MariaDB](https://mariadb.org/)

---

⭐ **Star dieses Projekt** wenn es dir hilft!

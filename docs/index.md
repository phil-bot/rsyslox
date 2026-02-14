# rsyslog REST API - Dokumentation

Willkommen zur vollständigen Dokumentation des rsyslog REST API Projekts.

## 📖 Übersicht

Diese Dokumentation ist in verschiedene Bereiche unterteilt, je nach Ihrer Rolle und Ihren Anforderungen.

## 🚀 Für Endnutzer

### Erste Schritte

| Dokument | Beschreibung |
|----------|--------------|
| [**Installation**](installation.md) | Alle Installationsmethoden (Binary, Source, Package) |
| [**Configuration**](configuration.md) | Vollständige Konfigurationsreferenz |
| [**Quick Examples**](examples.md) | Praktische Beispiele für häufige Anwendungsfälle |

### Nutzung & API

| Dokument | Beschreibung |
|----------|--------------|
| [**API Reference**](api-reference.md) | Vollständige API-Dokumentation mit allen Endpunkten |
| [**Troubleshooting**](troubleshooting.md) | Häufige Probleme, Lösungen und FAQ |
| [**Changelog**](changelog.md) | Versionshistorie und Breaking Changes |

## 🔧 Für Administratoren

| Dokument | Beschreibung |
|----------|--------------|
| [**Deployment**](deployment.md) | Production Setup, Systemd, Reverse Proxy |
| [**Security**](security.md) | Best Practices für sicheren Betrieb |
| [**Performance**](performance.md) | Optimierung und Benchmarks |

## 💻 Für Entwickler

| Dokument | Beschreibung |
|----------|--------------|
| [**Docker Testing**](docker.md) | Testumgebung mit Live-Daten |
| [**Development**](development.md) | Architektur, Build, Contributing |

## 🔍 Schnellzugriff

### Häufige Aufgaben

- **Installation starten:** → [Installation Guide](installation.md#quick-install)
- **API-Key generieren:** → [Configuration](configuration.md#api-key)
- **SSL einrichten:** → [Security](security.md#ssltls)
- **Produktiv deployen:** → [Deployment](deployment.md#production-setup)
- **Fehlersuche:** → [Troubleshooting](troubleshooting.md)
- **Docker testen:** → [Docker Guide](docker.md#quick-start)

### API-Endpunkte

- **Health Check:** → [GET /health](api-reference.md#get-health)
- **Logs abrufen:** → [GET /logs](api-reference.md#get-logs)
- **Metadaten:** → [GET /meta](api-reference.md#get-meta)

## 📚 Dokumentationsstruktur

```
docs/
├── index.md                 # Diese Datei - Übersicht
│
├── installation.md          # Installation (Binary, Source, Package)
├── configuration.md         # Vollständige Konfiguration
├── api-reference.md         # API-Endpunkte und Parameter
├── examples.md              # Praktische Beispiele
├── troubleshooting.md       # Fehlersuche und FAQ
│
├── deployment.md            # Production Deployment
├── security.md              # Sicherheits-Best-Practices
├── performance.md           # Performance-Tuning
│
├── docker.md                # Docker Testumgebung
├── development.md           # Entwicklung und Contributing
│
└── changelog.md             # Versionshistorie
```

## 🆘 Hilfe benötigt?

- **GitHub Issues:** [Fehler melden](https://github.com/phil-bot/rsyslog-rest-api/issues)
- **GitHub Discussions:** [Fragen stellen](https://github.com/phil-bot/rsyslog-rest-api/discussions)
- **Troubleshooting:** [FAQ durchsuchen](troubleshooting.md#faq)

## 🔄 Versionen

Diese Dokumentation gilt für:
- **Aktuelle Version:** v0.2.2
- **Mindest-Version:** v0.2.0

Für ältere Versionen siehe [Changelog](changelog.md).

---

[← Zurück zur README](../README.md)

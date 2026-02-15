# Contributing Guidelines

Thank you for contributing to rsyslog REST API!

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21+
- MySQL/MariaDB  
- rsyslog with MySQL support
- Docker (for testing)

### Build

```bash
make build-static
```

### Testing

```bash
cd docker
docker-compose up -d
./test.sh
```

## Coding Guidelines

- Follow Go best practices
- Use prepared statements for SQL
- Handle errors explicitly
- Add tests for new features
- Update documentation

## Pull Request Process

1. Update documentation
2. Add tests
3. Test locally
4. Create clear PR description
5. Respond to code review

## Questions?

- [GitHub Discussions](https://github.com/yourusername/rsyslog-rest-api/discussions)
- [Open an Issue](https://github.com/yourusername/rsyslog-rest-api/issues)

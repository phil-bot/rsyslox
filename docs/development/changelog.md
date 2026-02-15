# Changelog

All notable changes to rsyslog REST API.

## [Unreleased]

### Planned for v0.3.0
- Negation filters
- Advanced filter combinations
- Unit tests

### Planned for v0.4.0
- Statistics endpoint
- Aggregations
- Prometheus metrics

## [v0.2.3] - 2025-02-15

### Added
- Enhanced multi-value filter performance
- Better error validation messages
- Improved meta endpoint filtering

### Fixed
- Multi-value filter edge cases
- Error message clarity
- Performance improvements

### Changed
- Updated dependencies
- Improved documentation
- Better code organization

## [v0.2.2] - 2025-02-09

### Added
- **Multi-value filters** for all parameters
- **Extended columns** - All 25+ SystemEvents fields
- **Enhanced meta endpoint** - Supports filtering
- **Live log generator** for Docker

### Changed
- API response includes all available columns
- Meta endpoint returns labels
- Improved filter validation

### Fixed
- Null handling for extended columns
- Meta endpoint performance
- Database index creation

## [v0.2.1] - 2025-01-15

### Fixed
- Database connection timeout
- Memory leak in queries
- CORS preflight handling

### Changed
- Improved error messages
- Better logging format

## [v0.2.0] - 2024-12-20

### Added
- RFC-5424 labels
- Meta endpoint
- SSL/TLS support
- CORS configuration

### Changed
- Response format includes labels
- Configuration via .env

## [v0.1.0] - 2024-10-01

Initial release

### Features
- Basic REST API
- Authentication
- Pagination
- Docker testing

## Version History

| Version | Date | Changes |
|---------|------|---------|
| v0.2.3 | 2025-02-15 | Performance, errors |
| v0.2.2 | 2025-02-09 | Multi-value, columns |
| v0.2.1 | 2025-01-15 | Bug fixes |
| v0.2.0 | 2024-12-20 | Labels, SSL |
| v0.1.0 | 2024-10-01 | Initial release |

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.3] - 2026-02-15

### 🎉 Major: Complete Modularization

**Complete architectural refactor** - the codebase is now organized into logical packages for better maintainability and extensibility.

### ⚠️ Breaking Changes

#### Error Response Format (CHANGED)
- **Before:** `{"error": "message"}`
- **After:** `{"code": "ERROR_CODE", "message": "message", "details": "...", "field": "..."}`
- **Reason:** More structured and informative error responses
- **Migration:** Update error handling in API clients - see MIGRATION.md

### ✨ Added

#### New Configuration Options
- **DB_CONNECTION_STRING**: Full MySQL DSN support for advanced configurations
  - Example: `DB_CONNECTION_STRING=user:pass@tcp(host)/db?parseTime=true`
  - Provides more flexibility for complex database setups
  - Existing `DB_HOST/DB_NAME/DB_USER/DB_PASS` still supported

#### Improved Error Handling
- **Structured errors** with error codes (`INVALID_PARAMETER`, `DATABASE_ERROR`, etc.)
- **Field-specific validation errors** showing which parameter failed
- **Detailed error messages** with helpful hints (e.g., RFC-5424 references)
- **Consistent error format** across all endpoints

#### Enhanced Response Headers
- **Content-Type** now includes charset: `application/json; charset=utf-8`
- **X-Content-Type-Options: nosniff** for security
- Better HTTP compliance

#### Improved Logging
- **Timestamped logs**: `[2026-02-15 10:30:45] GET /logs - 200 - 15.3ms - IP`
- **Status code tracking**: All requests now show HTTP status
- **Client IP logging**: Better debugging and monitoring
- **Duration tracking**: Performance metrics for every request

#### Version Information
- `/health` endpoint now includes API version
- Helps with debugging and support

### 🔧 Changed

#### Code Organization (Internal - No API Impact)
- **Modular architecture**: Codebase split into logical packages
  - `internal/models`: Data structures and RFC mappings
  - `internal/config`: Configuration management
  - `internal/database`: Database connection and queries
  - `internal/filters`: Query builder and validation
  - `internal/middleware`: HTTP middleware (auth, CORS, logging)
  - `internal/handlers`: API endpoint handlers
  - `internal/server`: Server setup and routing
- **main.go**: Reduced from 850 lines to ~40 lines
- **Better separation of concerns**: Each package has a clear responsibility
- **Easier testing**: Components can be tested independently

#### Improved Validation
- **Earlier validation**: Parameters validated immediately upon parsing
- **Better error messages**: More specific validation feedback
- **Consistent behavior**: All endpoints use the same validation logic

#### Database Performance
- **Improved query helpers**: Better prepared statement support
- **Optimized column scanning**: More efficient data mapping
- **Connection pooling refinements**: Better resource management

### 📖 Documentation

#### New Files
- **MIGRATION.md**: Comprehensive migration guide from v0.2.2
- **CHANGELOG.md**: This file!

#### Updated Files
- **README.md**: Updated with v0.2.3 features and breaking changes
- **.env.example**: Added `DB_CONNECTION_STRING` example

### 🐛 Fixed
- More consistent NULL handling in responses
- Improved error context in database operations
- Better edge case handling in filter validation

### 🏗️ Technical Improvements
- **Go version**: Updated to 1.22 for better performance
- **Code quality**: Better code organization and readability
- **Maintainability**: Easier to add new features
- **Extensibility**: Prepared for v0.3.0 features (negation, advanced filters)

### 📊 Metrics
- **Lines of code**: More code total, but better organized
- **main.go**: 850 → 40 lines (95% reduction!)
- **Test coverage**: Foundation laid for comprehensive testing
- **Build time**: Slightly longer due to more files (marginal)
- **Runtime performance**: Same or slightly better

---

## [0.2.2] - 2025-02-09

### ✨ Added
- **Multi-value filters**: All filter parameters now support multiple values
  - Example: `?FromHost=web01&FromHost=web02&Priority=3&Priority=4`
- **Extended columns**: Response includes all 25 SystemEvents columns
- **Live log generator** (Docker): Generates realistic logs every 10 seconds

### 🔧 Changed
- Improved filter handling with OR logic for multi-value parameters
- Better NULL handling with `omitempty` tags

---

## [0.2.1] - 2025-02-01

### 🐛 Fixed
- Pagination edge cases
- Date range validation

---

## [0.2.0] - 2025-01-15

### ✨ Added
- Initial release with REST API
- Multi-value filter support
- Docker test environment

---

## Migration Guides

- **v0.2.2 → v0.2.3**: See [MIGRATION.md](MIGRATION.md)

---

## Versioning Strategy

- **Major (x.0.0)**: Breaking API changes
- **Minor (0.x.0)**: New features, may include breaking changes if necessary
- **Patch (0.0.x)**: Bug fixes and minor improvements

---

## Links

- [Repository](https://github.com/phil-bot/rsyslog-rest-api)
- [Issues](https://github.com/phil-bot/rsyslog-rest-api/issues)
- [Releases](https://github.com/phil-bot/rsyslog-rest-api/releases)

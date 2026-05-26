# Go Toolkit

[![Go Reference](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit.svg)](https://pkg.go.dev/github.com/alpardfm/go-toolkit)
[![CI](https://github.com/alpardfm/go-toolkit/actions/workflows/ci.yml/badge.svg)](https://github.com/alpardfm/go-toolkit/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpardfm/go-toolkit)](https://goreportcard.com/report/github.com/alpardfm/go-toolkit)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev/)

A collection of production-grade utility packages for building REST API services in Go. Go Toolkit provides reusable, well-tested building blocks for common backend tasks — from error handling and configuration management to database connections and JWT authentication — so you can focus on your application logic instead of reinventing infrastructure.

## Installation

```bash
# Core module (lightweight utilities)
go get github.com/alpardfm/go-toolkit

# Heavy module (database, AWS, logging, JWT, etc.)
go get github.com/alpardfm/go-toolkit/heavy
```

Requires **Go 1.21** or later.

## Module Structure

The toolkit is split into two Go modules to minimize dependency bloat:

| Module | Import Path | Description |
|--------|-------------|-------------|
| **Core** | `github.com/alpardfm/go-toolkit` | Lightweight utilities with minimal dependencies (operator, convert, errors, codes, concurrency, hash, encrypt, random, format, files, validation, distance, appcontext, language, header, sorter, parser) |
| **Heavy** | `github.com/alpardfm/go-toolkit/heavy` | Packages requiring large external dependencies (configbuilder, configreader, nosql, sql, smtp, log, tokens, query) |

If you only need lightweight utilities like type conversion, error handling, or hashing, import the core module — your project won't pull in AWS SDK, MongoDB driver, or database drivers.

For details on which dependencies each package requires, see [DEPENDENCIES.md](DEPENDENCIES.md). If you're migrating from the single-module layout, see [MIGRATION.md](MIGRATION.md) for updated import paths.

## Package Overview

| Package | Description | GoDoc |
|---------|-------------|-------|
| [appcontext](./appcontext) | Request-scoped context value management | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/appcontext)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/appcontext) |
| [codes](./codes) | Error code registry with bilingual messages (EN/ID) | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/codes)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/codes) |
| [concurrency](./concurrency) | Goroutine worker pool with configurable max workers | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/concurrency)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/concurrency) |
| [configbuilder](./configbuilder) | AWS SSM parameter store config generator | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/configbuilder)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/configbuilder) |
| [configreader](./configreader) | Viper-based configuration file reader | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/configreader)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/configreader) |
| [convert](./convert) | Type conversion utilities (int64, float64, string, struct copy) | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/convert)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/convert) |
| [distance](./distance) | Geographic distance calculation | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/distance)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/distance) |
| [encrypt](./encrypt) | AES-256-GCM encryption/decryption | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/encrypt)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/encrypt) |
| [errors](./errors) | Stacktrace-enriched error handling with error codes | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/errors)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/errors) |
| [files](./files) | File system utilities | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/files)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/files) |
| [format](./format) | Date and string formatting helpers | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/format)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/format) |
| [hash](./hash) | Password hashing (argon2, bcrypt, sha256) | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/hash)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/hash) |
| [header](./header) | HTTP header extraction utilities | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/header)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/header) |
| [language](./language) | Bilingual text support (EN/ID) | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/language)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/language) |
| [log](./log) | Structured logging via zerolog | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/log)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/log) |
| [nosql](./nosql) | MongoDB client wrapper | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/nosql)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/nosql) |
| [operator](./operator) | Ternary operator and generic utilities | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/operator)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/operator) |
| [parser](./parser) | JSON marshaling/unmarshaling interface | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/parser)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/parser) |
| [query](./query) | SQL/GQL query builder from struct tags | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/query)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/query) |
| [random](./random) | Random value generation (int, UUID-based) | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/random)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/random) |
| [smtp](./smtp) | Email sending via gomail | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/smtp)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/smtp) |
| [sorter](./sorter) | Generic sorting with custom comparators | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/sorter)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/sorter) |
| [sql](./sql) | Database connection management (leader/follower) | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/sql)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/sql) |
| [tokens](./tokens) | JWT token creation and validation | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/tokens)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/tokens) |
| [validation](./validation) | Input validation utilities (email, etc.) | [![GoDoc](https://pkg.go.dev/badge/github.com/alpardfm/go-toolkit/validation)](https://pkg.go.dev/github.com/alpardfm/go-toolkit/validation) |

## Breaking Changes (v1.0.0)

The following breaking changes were introduced in v1.0.0 to improve safety and correctness:

- **configbuilder**: `Init()` now returns `(Interface, error)` instead of `Interface`
- **configreader**: `Init()` now returns `(Interface, error)` instead of `Interface`
- **sql**: `Init()` now returns `(Interface, error)` instead of `Interface`
- **concurrency**: `Do()` now returns a multi-error (via `errors.Join`) containing all goroutine errors
- **tokens**: JWT library migrated from `dgrijalva/jwt-go/v4` to `golang-jwt/jwt/v5`

See [CHANGELOG.md](CHANGELOG.md) for full details.

## Versioning

This project follows [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html). Releases are tagged using annotated git tags in the format `vMAJOR.MINOR.PATCH`.

- **MAJOR** version increments indicate breaking API changes
- **MINOR** version increments add functionality in a backward-compatible manner
- **PATCH** version increments include backward-compatible bug fixes

See [CHANGELOG.md](CHANGELOG.md) for a history of changes by version.

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, coding standards, and the pull request process.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

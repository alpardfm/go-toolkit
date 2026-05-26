# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-05-26

### Added

- `appcontext` — Request-scoped context value management
- `codes` — Error code registry with bilingual messages (EN/ID)
- `concurrency` — Goroutine worker pool with configurable max workers
- `configbuilder` — AWS SSM parameter store config generator
- `configreader` — Viper-based configuration file reader
- `convert` — Type conversion utilities (int64, float64, string, struct copy)
- `distance` — Geographic distance calculation
- `encrypt` — AES-256-GCM encryption/decryption
- `errors` — Stacktrace-enriched error handling with error codes
- `files` — File system utilities
- `format` — Date and string formatting helpers
- `hash` — Password hashing (argon2, bcrypt, sha256)
- `header` — HTTP header extraction utilities
- `language` — Bilingual text support (EN/ID)
- `log` — Structured logging via zerolog
- `nosql` — MongoDB client wrapper
- `operator` — Ternary operator and generic utilities
- `parser` — JSON marshaling/unmarshaling interface
- `query` — SQL/GQL query builder from struct tags
- `random` — Random value generation (int, UUID-based)
- `smtp` — Email sending via gomail
- `sorter` — Generic sorting with custom comparators
- `sql` — Database connection management (leader/follower)
- `tokens` — JWT token creation and validation
- `validation` — Input validation utilities (email, etc.)
- `errors.Unwrap()` method enabling `errors.Is()` and `errors.As()` chain traversal
- CI pipeline with GitHub Actions (Go 1.21 + latest, golangci-lint)
- Comprehensive documentation (README, CONTRIBUTING, CHANGELOG)

### Changed

- **BREAKING:** `configbuilder.Init()` now returns `(Interface, error)` instead of `Interface`
- **BREAKING:** `configreader.Init()` now returns `(Interface, error)` instead of `Interface`
- **BREAKING:** `sql.Init()` now returns `(Interface, error)` instead of `Interface`
- **BREAKING:** `concurrency.Do()` now returns a multi-error (via `errors.Join`) containing all goroutine errors
- **BREAKING:** `tokens` package migrated from `dgrijalva/jwt-go/v4` to `golang-jwt/jwt/v5`

### Fixed

- Inverted MockDB conditional logic in `sql` package `connect()` method
- Race condition in `concurrency.Do()` when reading error list after `wg.Wait()`
- Deprecated API usage (`ioutil.ReadFile`, global `rand.Seed`)

### Security

- Migrated JWT library from unmaintained `dgrijalva/jwt-go` to actively maintained `golang-jwt/jwt/v5`

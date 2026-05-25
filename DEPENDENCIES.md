# Dependencies

This document lists the external dependencies required by each package in the Go Toolkit, organized by module.

## Module Structure

The toolkit is split into two Go modules to keep the dependency tree minimal for users who only need lightweight utilities:

- **Core module** (`github.com/alpardfm/go-toolkit`) — lightweight packages with minimal dependencies
- **Heavy module** (`github.com/alpardfm/go-toolkit/heavy`) — packages that require large external dependencies

## Core Module

**Import:** `github.com/alpardfm/go-toolkit`

The core module contains packages that depend only on the Go standard library and a few small, focused libraries:

| Package | External Dependencies |
|---------|----------------------|
| operator | _(none — stdlib only)_ |
| convert | `github.com/cstockton/go-conv` |
| errors | _(none — stdlib only)_ |
| codes | _(none — stdlib only)_ |
| concurrency | _(none — stdlib only)_ |
| hash | `golang.org/x/crypto` (argon2, bcrypt) |
| encrypt | _(none — stdlib only)_ |
| random | `github.com/google/uuid` |
| format | _(none — stdlib only)_ |
| files | _(none — stdlib only)_ |
| validation | _(none — stdlib only)_ |
| distance | _(none — stdlib only)_ |
| appcontext | _(none — stdlib only)_ |
| language | _(none — stdlib only)_ |
| header | _(none — stdlib only)_ |
| sorter | _(none — stdlib only)_ |
| parser | `github.com/json-iterator/go` |

**Total direct dependencies:** 4 lightweight libraries.

## Heavy Module

**Import:** `github.com/alpardfm/go-toolkit/heavy`

The heavy module contains packages that pull in large dependency trees. Each package and its primary external dependencies are listed below:

| Package | Primary Dependencies | Approximate Impact |
|---------|---------------------|-------------------|
| configbuilder | `github.com/aws/aws-sdk-go`, `github.com/cbroglie/mustache`, `github.com/spf13/viper` | AWS SDK (~30MB compiled) |
| configreader | `github.com/spf13/viper` | Viper + transitive deps |
| nosql | `go.mongodb.org/mongo-driver` | MongoDB driver (~15MB compiled) |
| sql | `github.com/jmoiron/sqlx`, `github.com/go-sql-driver/mysql`, `github.com/lib/pq` | SQL drivers |
| smtp | `github.com/go-mail/gomail` | Gomail |
| log | `github.com/rs/zerolog` | Zerolog structured logger |
| tokens | `github.com/golang-jwt/jwt/v5` | JWT library |
| query | _(uses only core toolkit + stdlib)_ | Minimal, but co-located for architectural reasons |

### Detailed Dependency Breakdown

#### configbuilder

```
github.com/aws/aws-sdk-go       — AWS SDK for SSM Parameter Store access
github.com/cbroglie/mustache    — Mustache template rendering
github.com/spf13/viper          — Configuration parsing (used for parameter mapping)
```

#### configreader

```
github.com/spf13/viper          — Configuration file reading (JSON, YAML, TOML, etc.)
```

#### nosql

```
go.mongodb.org/mongo-driver     — Official MongoDB Go driver
```

#### sql

```
github.com/jmoiron/sqlx         — Extensions to database/sql (named queries, struct scanning)
github.com/go-sql-driver/mysql  — MySQL driver
github.com/lib/pq               — PostgreSQL driver
```

#### smtp

```
github.com/go-mail/gomail       — Email sending via SMTP
```

#### log

```
github.com/rs/zerolog            — Zero-allocation structured logging
```

#### tokens

```
github.com/golang-jwt/jwt/v5    — JSON Web Token creation and validation
```

## Choosing the Right Module

| If you need... | Import |
|----------------|--------|
| Type conversions, error handling, hashing, encryption, formatting | `github.com/alpardfm/go-toolkit` |
| Database connections, AWS config, MongoDB, logging, JWT, email | `github.com/alpardfm/go-toolkit/heavy` |

If you only use packages from the core module, your project will **not** pull in AWS SDK, MongoDB driver, or any database drivers.

# Migration Guide

This guide documents import path changes introduced by the module split in Go Toolkit. If you were using packages that have moved to the `heavy` module, update your import paths as shown below.

## Import Path Changes

The following packages have moved from the root module to the `heavy` submodule:

| Old Import Path | New Import Path |
|----------------|-----------------|
| `github.com/alpardfm/go-toolkit/configbuilder` | `github.com/alpardfm/go-toolkit/heavy/configbuilder` |
| `github.com/alpardfm/go-toolkit/configreader` | `github.com/alpardfm/go-toolkit/heavy/configreader` |
| `github.com/alpardfm/go-toolkit/nosql` | `github.com/alpardfm/go-toolkit/heavy/nosql` |
| `github.com/alpardfm/go-toolkit/sql` | `github.com/alpardfm/go-toolkit/heavy/sql` |
| `github.com/alpardfm/go-toolkit/smtp` | `github.com/alpardfm/go-toolkit/heavy/smtp` |
| `github.com/alpardfm/go-toolkit/log` | `github.com/alpardfm/go-toolkit/heavy/log` |
| `github.com/alpardfm/go-toolkit/tokens` | `github.com/alpardfm/go-toolkit/heavy/tokens` |
| `github.com/alpardfm/go-toolkit/query` | `github.com/alpardfm/go-toolkit/heavy/query` |

## Unchanged Import Paths

The following packages remain in the core module with no import path changes:

- `github.com/alpardfm/go-toolkit/operator`
- `github.com/alpardfm/go-toolkit/convert`
- `github.com/alpardfm/go-toolkit/errors`
- `github.com/alpardfm/go-toolkit/codes`
- `github.com/alpardfm/go-toolkit/concurrency`
- `github.com/alpardfm/go-toolkit/hash`
- `github.com/alpardfm/go-toolkit/encrypt`
- `github.com/alpardfm/go-toolkit/random`
- `github.com/alpardfm/go-toolkit/format`
- `github.com/alpardfm/go-toolkit/files`
- `github.com/alpardfm/go-toolkit/validation`
- `github.com/alpardfm/go-toolkit/distance`
- `github.com/alpardfm/go-toolkit/appcontext`
- `github.com/alpardfm/go-toolkit/language`
- `github.com/alpardfm/go-toolkit/header`
- `github.com/alpardfm/go-toolkit/sorter`
- `github.com/alpardfm/go-toolkit/parser`

## Step-by-Step Migration

### 1. Update your `go.mod`

If you use any heavy packages, add the heavy module dependency:

```bash
go get github.com/alpardfm/go-toolkit/heavy@latest
```

### 2. Update import paths

Find and replace old import paths in your source files. For example:

```go
// Before
import (
    "github.com/alpardfm/go-toolkit/sql"
    "github.com/alpardfm/go-toolkit/log"
    "github.com/alpardfm/go-toolkit/tokens"
)

// After
import (
    "github.com/alpardfm/go-toolkit/heavy/sql"
    "github.com/alpardfm/go-toolkit/heavy/log"
    "github.com/alpardfm/go-toolkit/heavy/tokens"
)
```

You can use `sed` or your editor's find-and-replace to do this across your project:

```bash
# Example using sed (macOS)
find . -name "*.go" -exec sed -i '' \
  -e 's|"github.com/alpardfm/go-toolkit/configbuilder"|"github.com/alpardfm/go-toolkit/heavy/configbuilder"|g' \
  -e 's|"github.com/alpardfm/go-toolkit/configreader"|"github.com/alpardfm/go-toolkit/heavy/configreader"|g' \
  -e 's|"github.com/alpardfm/go-toolkit/nosql"|"github.com/alpardfm/go-toolkit/heavy/nosql"|g' \
  -e 's|"github.com/alpardfm/go-toolkit/sql"|"github.com/alpardfm/go-toolkit/heavy/sql"|g' \
  -e 's|"github.com/alpardfm/go-toolkit/smtp"|"github.com/alpardfm/go-toolkit/heavy/smtp"|g' \
  -e 's|"github.com/alpardfm/go-toolkit/log"|"github.com/alpardfm/go-toolkit/heavy/log"|g' \
  -e 's|"github.com/alpardfm/go-toolkit/tokens"|"github.com/alpardfm/go-toolkit/heavy/tokens"|g' \
  -e 's|"github.com/alpardfm/go-toolkit/query"|"github.com/alpardfm/go-toolkit/heavy/query"|g' \
  {} \;
```

### 3. Run `go mod tidy`

```bash
go mod tidy
```

This will clean up your `go.mod` and `go.sum` files, removing any dependencies that are no longer needed directly.

### 4. Verify

```bash
go build ./...
go test ./...
```

## API Changes

In addition to import path changes, the following function signatures have changed:

| Package | Function | Old Signature | New Signature |
|---------|----------|---------------|---------------|
| configbuilder | Init | `Init(opt Options) Interface` | `Init(opt Options) (Interface, error)` |
| configreader | Init | `Init(opt Options) Interface` | `Init(opt Options) (Interface, error)` |
| sql | Init | `Init(cfg Config, log log.Interface) Interface` | `Init(cfg Config, log log.Interface) (Interface, error)` |

These functions no longer panic on failure — they return errors instead. Update your initialization code:

```go
// Before
cb := configbuilder.Init(opts)

// After
cb, err := configbuilder.Init(opts)
if err != nil {
    // handle error
}
```

## Questions?

If you encounter issues during migration, please open an issue on the repository.

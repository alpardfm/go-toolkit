# Contributing to Go Toolkit

Thank you for your interest in contributing to Go Toolkit! This guide covers everything you need to get started.

## Development Setup

### Prerequisites

- **Go 1.21** or later ([download](https://go.dev/dl/))
- **Git** for version control
- **golangci-lint** for linting ([install](https://golangci-lint.run/welcome/install/))

### Getting Started

1. Fork the repository on GitHub.

2. Clone your fork:

   ```bash
   git clone https://github.com/<your-username>/go-toolkit.git
   cd go-toolkit
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Verify everything builds:

   ```bash
   go build ./...
   ```

5. Run the test suite:

   ```bash
   go test -race ./...
   ```

6. Run the linter:

   ```bash
   golangci-lint run
   ```

## Coding Standards

### Formatting

- All code must be formatted with `gofmt` (or `goimports`).
- Run `gofmt -s -w .` before committing.

### Naming Conventions

- Follow standard Go naming conventions: https://go.dev/doc/effective_go#names
- Exported identifiers use PascalCase.
- Unexported identifiers use camelCase.
- Acronyms are all-caps (e.g., `HTTP`, `URL`, `ID`).
- Interface names should describe behavior (e.g., `Reader`, `Interface`).

### Comments and Documentation

- All exported types, functions, and methods must have godoc-compatible comments.
- Comments must begin with the identifier name (e.g., `// Init initializes the...`).
- Each comment must be at least one complete sentence ending with a period.
- Package-level comments go in a `doc.go` file or at the top of the primary source file.

### Error Handling

- Use `errors.NewWithCode()` from the toolkit's `errors` package for all errors.
- Never use `panic()` in library code — return errors instead.
- Wrap errors with context using `fmt.Errorf("operation failed: %w", err)` or the toolkit's error wrapping.

### Package Structure

Each package should follow this layout:

```
package_name/
├── package.go              # Main implementation
├── package_test.go         # Unit tests (table-driven)
├── package_property_test.go # Property-based tests (if applicable)
└── doc.go                  # Package documentation (optional)
```

## Pull Request Process

### Branch Naming

Use descriptive branch names with a prefix:

- `feat/` — new features (e.g., `feat/add-redis-client`)
- `fix/` — bug fixes (e.g., `fix/concurrency-race-condition`)
- `docs/` — documentation changes (e.g., `docs/update-readme`)
- `refactor/` — code refactoring (e.g., `refactor/remove-panics`)
- `test/` — test additions or fixes (e.g., `test/add-hash-coverage`)

### Before Submitting

1. Ensure all tests pass: `go test -race ./...`
2. Ensure the linter passes: `golangci-lint run`
3. Ensure your code is formatted: `gofmt -s -w .`
4. Add or update tests for any changed functionality.
5. Update documentation if you changed public APIs.

### Review Process

1. Open a pull request against the `main` branch.
2. Fill in the PR template with a description of changes and testing done.
3. At least one maintainer must approve the PR before merge.
4. CI checks (tests + lint) must pass on all supported Go versions.
5. Address review feedback with additional commits (do not force-push during review).

### Merge Strategy

- PRs are merged via squash merge to keep the main branch history clean.
- The PR title becomes the commit message, so keep it concise and descriptive.

## Testing Requirements

### Coverage Expectations

- All packages should target **≥80% line coverage**.
- Critical packages (errors, convert, hash, encrypt, query, tokens) should target **≥90%**.

### Running Tests

```bash
# Run all tests with race detection
go test -race ./...

# Run tests for a specific package
go test -race ./errors/...

# Run tests with coverage report
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

- Use table-driven tests with `t.Run()` for functions with multiple input variations.
- Use `github.com/stretchr/testify` for assertions.
- Use `github.com/flyingmutant/rapid` for property-based tests.
- Mock external dependencies via interfaces — no live DB or API connections in unit tests.
- Include both positive (happy path) and negative (error) test cases.

### Test Naming

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {name: "valid input returns expected result", ...},
        {name: "empty input returns error", ...},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

## Questions?

If you have questions about contributing, feel free to open an issue for discussion.

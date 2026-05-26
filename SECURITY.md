# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| v1.x    | ✅ Active support  |
| < v1.0  | ❌ No support      |

## Reporting a Vulnerability

If you discover a security vulnerability in Go Toolkit, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please send an email to the maintainer with:

1. A description of the vulnerability
2. Steps to reproduce the issue
3. Potential impact assessment
4. Suggested fix (if any)

### What to expect

- **Acknowledgment**: Within 48 hours of your report
- **Assessment**: Within 7 days, we will assess the severity and confirm the vulnerability
- **Fix**: Critical vulnerabilities will be patched within 14 days
- **Disclosure**: We will coordinate disclosure timing with you

### Scope

The following are in scope for security reports:

- Cryptographic weaknesses in `encrypt`, `hash`, or `tokens` packages
- Authentication/authorization bypasses
- Injection vulnerabilities (SQL injection in `query` package, etc.)
- Denial of service through resource exhaustion
- Information disclosure (secrets in logs, error messages, etc.)

### Out of Scope

- Vulnerabilities in dependencies (report these upstream)
- Issues requiring physical access
- Social engineering

## Security Best Practices

When using this toolkit:

- Always use the latest version
- Keep dependencies updated (`go get -u`)
- Use `crypto/rand` for security-sensitive randomness (the `random` package already does this)
- Never hardcode secrets — use environment variables or secret managers
- Review the `encrypt` package key requirements (AES-256 needs exactly 32-byte keys)

## Dependencies

We monitor dependencies for known vulnerabilities using:

- `go vuln check` in CI
- GitHub Dependabot alerts

All direct dependencies are pinned to specific versions in `go.mod`.

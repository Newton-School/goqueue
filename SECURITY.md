# Security Policy

goqueue is intended for production services, so security issues should be
reported privately.

## Reporting a Vulnerability

Use GitHub Security Advisories for this repository. Do not open a public issue
for vulnerabilities, credential leaks, denial-of-service vectors, or unsafe
deserialization behavior.

## Secure Defaults

- Library code does not read `.env` files directly.
- Redis credentials should not be logged. Use `Config.RedactedRedisURL` when a
  connection URL must appear in diagnostics.
- Queue names and namespaces are validated before app construction.
- Future task payload handling must avoid unsafe deserialization and must not
  log raw payloads by default.

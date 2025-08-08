# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2024-08-07

### Added
- Initial release of GoUtils
- **Retry package**: Flexible retry mechanism with configurable strategies
  - Fixed delay, exponential backoff, and linear backoff strategies
  - Context support for cancellation
  - Custom retry conditions
  - Jitter support to avoid thundering herd
- **HTTP client package**: HTTP utilities with built-in retry
  - Configurable timeout and base URL
  - Automatic JSON encoding/decoding
  - Retry logic for transient failures
  - Context support
- **Cache package**: Multiple caching implementations
  - In-memory cache with TTL support
  - LRU (Least Recently Used) cache
  - Thread-safe operations
- **String utilities package**: Comprehensive string manipulation
  - Case conversions (camelCase, PascalCase, snake_case, kebab-case)
  - String validation (email, numeric, alpha)
  - Random string generation
  - String padding and truncation
- **Convert package**: Type conversion utilities
  - Safe type conversions with error handling
  - Slice conversions
  - JSON marshaling/unmarshaling helpers
  - Struct to map conversion
  - Time parsing with multiple formats
- Comprehensive test coverage (>90%)
- Examples and documentation
- CI/CD pipeline with GitHub Actions
- Makefile for development workflow

### Security
- All random string generation uses crypto/rand for security
- Input validation for all public functions

[Unreleased]: https://github.com/jelech/goutils/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/jelech/goutils/releases/tag/v1.0.0

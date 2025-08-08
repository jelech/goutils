# Contributing to GoUtils

Thank you for your interest in contributing to GoUtils! This document provides guidelines and information for contributors.

## How to Contribute

### Reporting Issues

- Use the GitHub issue tracker to report bugs or request features
- Search existing issues before creating a new one
- Provide as much detail as possible, including:
  - Go version
  - Operating system
  - Steps to reproduce the issue
  - Expected vs actual behavior

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Write tests** for any new functionality
3. **Ensure tests pass** by running `make test`
4. **Follow code style** by running `make fmt` and `make lint`
5. **Update documentation** if needed
6. **Write clear commit messages**

### Development Setup

1. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/goutils.git
   cd goutils
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Install development tools:
   ```bash
   make install-tools
   ```

4. Run tests to ensure everything works:
   ```bash
   make test
   ```

### Code Guidelines

#### Code Style

- Follow standard Go formatting (gofmt)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small

#### Testing

- Write unit tests for all new functionality
- Maintain test coverage above 90%
- Use table-driven tests when appropriate
- Include benchmark tests for performance-critical code

#### Documentation

- Document all exported functions and types
- Include usage examples in documentation
- Update README.md for new features
- Add entries to CHANGELOG.md

### Development Workflow

1. Create a feature branch:
   ```bash
   git checkout -b feature/new-feature
   ```

2. Make your changes and write tests

3. Run quality checks:
   ```bash
   make quality
   ```

4. Commit your changes:
   ```bash
   git commit -m "Add new feature: description"
   ```

5. Push to your fork and create a pull request

### Package Structure

When adding new packages:

```
pkg/
├── retry/           # Retry mechanisms
├── http/            # HTTP utilities
├── cache/           # Caching implementations
├── stringutil/      # String utilities
├── convert/         # Type conversions
└── newpackage/      # Your new package
    ├── package.go   # Main implementation
    └── package_test.go # Tests
```

### Commit Message Format

Use clear, descriptive commit messages:

```
type(scope): description

Examples:
- feat(retry): add exponential backoff with jitter
- fix(cache): resolve race condition in LRU cache
- docs(readme): update installation instructions
- test(http): add integration tests for retry logic
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `test`: Tests
- `refactor`: Code refactoring
- `perf`: Performance improvement
- `chore`: Maintenance

### Performance Considerations

- Use benchmarks to measure performance
- Avoid unnecessary allocations
- Consider memory usage in long-running operations
- Use appropriate data structures

### Security Guidelines

- Use `crypto/rand` for random generation
- Validate all inputs
- Avoid exposing sensitive information in logs
- Follow Go security best practices

### Release Process

Releases are handled by maintainers:

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create and push a version tag
4. GitHub Actions will handle the release

### Getting Help

- Check existing documentation and examples
- Search through existing issues
- Ask questions in discussions
- Contact maintainers if needed

## Code of Conduct

This project follows the Go Community Code of Conduct. Be respectful and inclusive in all interactions.

Thank you for contributing to GoUtils!

# Contributing to Howl

Thank you for your interest in contributing to Howl! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with the following information:

1. A clear, descriptive title
2. A detailed description of the issue
3. Steps to reproduce the bug
4. Expected behavior
5. Actual behavior
6. Environment information (Go version, OS, etc.)
7. Any additional context or screenshots

### Suggesting Enhancements

For feature requests or enhancements:

1. Create an issue with a clear title and detailed description
2. Explain why this enhancement would be useful
3. Suggest an implementation approach if possible

### Pull Requests

1. Fork the repository
2. Create a new branch from `develop`
3. Make your changes
4. Add or update tests as necessary
5. Ensure all tests pass with `go test ./...`
6. Update documentation if needed
7. Submit a pull request to the `develop` branch

## Development Guidelines

### Code Style

- All code should follow the standard Go formatting guidelines
- Run `go fmt` before committing
- Use `golint` and `go vet` to check for issues

### Documentation

- All exported functions, types, and methods must have proper documentation
- Comments should be in English
- Keep comments concise and focused on explaining "why" rather than "what"

### Testing

- Add tests for new functionality
- Ensure all tests pass before submitting a pull request
- Aim for high test coverage

### Commit Messages

- Use clear, descriptive commit messages
- Start with a short summary line (50 chars or less)
- Optionally followed by a blank line and a more detailed explanation

## Project Structure

```
howl/
├── examples/         # Example applications
├── *.go              # Core package files
└── *_test.go         # Test files
```

## License

By contributing to Howl, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).

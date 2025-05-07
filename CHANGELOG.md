# Changelog

All notable changes to the Howl logging package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-05-07

### Added
- Initial release of Howl logging package
- Support for multiple log levels (Debug, Info, Warn, Error, Fatal)
- Structured logging with fields using WithField and WithFields methods
- Error logging with WithError method
- JSON formatting with indentation
- Colorized output for improved readability
- Colorized JSON for better visualization of structured data
- Context integration with extractors
- Thread-safe logging for concurrent environments
- File output support
- Comprehensive test suite
- Example applications in the examples directory

### Changed
- Simplified API design for ease of use
- Improved error handling

### Fixed
- Automatic detection and disabling of ANSI color codes when writing to files

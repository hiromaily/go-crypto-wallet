## Testing Guidelines

This document describes the testing strategy and requirements for the go-crypto-wallet project.

### Testing Principles

- Use `//go:build integration` tag for integration tests
- Separate unit tests and integration tests
- Use [testify](https://github.com/stretchr/testify) package for assertions (`assert` and `require`)
- Measure and improve test coverage
- Write tests for all exported functions and methods
- Keep tests maintainable and readable

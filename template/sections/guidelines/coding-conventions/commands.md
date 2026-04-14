### Common Commands

After making code changes, use these commands to verify code correctness:

- `make go-lint`: Fix linting issues automatically
- `make check-build`: Verify that the code builds successfully
- `make go-test`: Run Go tests to verify functionality
- `make tidy`: Organize dependencies and clean up `go.mod`

**Important**: After modifying Go code, run these commands to ensure code quality and correctness.

**Command Constraints**:

- **DO NOT** use `go build -v` directly; use `make check-build` instead
- **DO NOT** use `go tool golangci-lint` directly; use `make go-lint` instead

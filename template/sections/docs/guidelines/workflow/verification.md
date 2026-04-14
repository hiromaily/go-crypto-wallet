### Verification Commands

After making code changes, always run these commands in order:

1. `make go-lint` - Fix linting issues automatically
2. `make tidy` - Organize dependencies and clean up `go.mod`
3. `make check-build` - Verify that the code builds successfully
4. `make go-test` - Run Go tests to verify functionality

**Optional but Recommended:**

- `make go-check-vuln` - Run security vulnerability scan (for security-related changes)
- `make go-test-integration` - Run integration tests (if applicable)

**Important**:

- Ensure no errors occur
- Ensure no files are modified (all changes should be committed)
- Ensure all commands pass successfully

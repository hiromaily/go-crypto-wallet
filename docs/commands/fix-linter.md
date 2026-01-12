# Fix Linter Command

## Purpose

Fix lint errors reported by `make go-lint` command.

## Parameters

None required. Works on current working directory.

## Process

### 1. Analyze Errors

1. Run `make go-lint` to get current errors
2. Categorize errors by type and severity:
   - Syntax errors (highest priority)
   - Security vulnerabilities
   - Type errors
   - Style/formatting issues (lowest priority)

### 2. Fix Errors

1. Start with critical errors
2. Group similar errors for batch fixing
3. Verify each fix doesn't break functionality
4. Check for auto-generated files (DO NOT EDIT)

### 3. Verify

Run verification commands after fixes.

## Verification

```bash
make go-lint
make tidy
make check-build
make gotest
```

## Guidelines

- **Preserve functionality**: Fixes must maintain original behavior
- **Use auto-fix**: Leverage linter's automatic fixing where appropriate
- **Never edit auto-generated files**: Check for `DO NOT EDIT` comments
- **Quality over speed**: Fix correctly in stages

## Output Format

For each fix:
1. Error category being addressed
2. Specific changes made
3. Reasoning behind the fix
4. Remaining errors (if any)

## Related Documents

- [Coding Standards](../standards/coding-conventions.md)
- [Code Generation](../ai-agents/guidelines/code-generation.md)

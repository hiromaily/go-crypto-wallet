# Fix Linter Command

Fix linter errors reported by `make go-lint`.

## Quick Reference

```bash
# 1. Run linter
make go-lint

# 2. Fix errors (priority: syntax > security > type > style)

# 3. Verify
make go-lint && make tidy && make check-build && make gotest
```

## Guidelines

- Never edit files with `DO NOT EDIT` comments
- Preserve original functionality
- Use linter auto-fix where appropriate

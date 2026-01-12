# Fix Linter Prompt

Refer to [docs/commands/fix-linter.md](../../docs/commands/fix-linter.md) for full documentation.

## Quick Reference

```
fix-linter
```

## Process Summary

1. Run `make go-lint` to identify errors
2. Categorize by severity (syntax > security > type > style)
3. Fix critical errors first
4. Verify fixes don't break functionality
5. Check for auto-generated files (DO NOT EDIT)

## Verification

```bash
make go-lint && make tidy && make check-build && make gotest
```

## Related

- [Coding Conventions](../../docs/standards/coding-conventions.md)
- [Code Generation](../../docs/ai-agents/guidelines/code-generation.md)

# Fix Linter Command

Refer to [docs/commands/fix-linter.md](../../docs/commands/fix-linter.md) for full documentation.

## Quick Reference

```
/fix-linter
```

## Process Summary

1. Run `make go-lint` to identify errors
2. Categorize by severity (syntax > security > type > style)
3. Fix critical errors first
4. Verify fixes don't break functionality

## Verification

```bash
make go-lint && make tidy && make check-build && make gotest
```

# Fix Linter

Fix linter errors reported by `make go-lint`.

## Skill Reference

**Use the `go-development` Skill** for verification commands.

## Process

1. **Run**: `make go-lint`
2. **Prioritize**: syntax > security > type > style
3. **Fix**: Address errors by priority
4. **Verify**: Run verification commands from Skill
5. **Commit**: Use `fix: resolve linter errors` format

## Guidelines

- Never edit files with `DO NOT EDIT` comments
- Preserve original functionality
- Use linter auto-fix where appropriate

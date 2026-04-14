### Pull Request Workflow

**Before Creating a PR:**

1. Run all verification commands
2. Ensure all tests pass
3. Review your own changes for:
   - Code quality and correctness
   - Adherence to Clean Architecture principles
   - Compliance with coding standards
   - Proper error handling
   - Security considerations

**PR Description Should Include:**

- Brief description of the changes
- Related issue number(s)
- Testing performed
- Any breaking changes
- Screenshots (if UI changes)

**PR Template:**

```markdown
## Description
[Brief description of the fix/feature]

## Changes
- [Change 1]
- [Change 2]

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing completed

## Verification
- [ ] `make go-lint` passes
- [ ] `make check-build` passes
- [ ] `make go-test` passes
- [ ] Security scan completed (if applicable)

Closes #[issue_number]
```

## Pull Request Guidelines

### PR Title Format

```
<type>: <description>
```

Example: `feat: add taproot address support`

### PR Description Template

```markdown
## Summary
Brief description of changes

## Changes
- Change 1
- Change 2

## Test Plan
- [ ] Unit tests pass
- [ ] Integration tests pass (if applicable)
- [ ] Manual testing completed

Closes #XXX
```

### Before Submitting

- [ ] Run `make go-lint` - no errors
- [ ] Run `make check-build` - build succeeds
- [ ] Run `make go-test` - all tests pass
- [ ] Commit messages follow Conventional Commits format
- [ ] PR title follows the format above

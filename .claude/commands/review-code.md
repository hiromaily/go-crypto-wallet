# Self-Review Changed Code

Review your own changed or added code and fix issues if found.

## Process

1. **Get changed files**: Run `git diff --name-only` and `git diff --cached --name-only` to identify modified files
2. **Get diff content**: Run `git diff` and `git diff --cached` to see actual changes
3. **Review checklist**: For each changed file, check:

### Code Quality

- [ ] No hardcoded values that should be configurable
- [ ] No commented-out code left behind
- [ ] No debug print statements (fmt.Println, console.log, etc.)
- [ ] Proper error messages with context
- [ ] Consistent naming conventions

### Architecture (for Go files)

- [ ] Domain layer has ZERO infrastructure dependencies
- [ ] Use cases depend only on interfaces they need (ISP)
- [ ] Error handling uses `fmt.Errorf("context: %w", err)`
- [ ] No circular dependencies

### Security (for sensitive areas)

- [ ] No private keys or sensitive data logged
- [ ] No secrets hardcoded
- [ ] Input validation at system boundaries

### Documentation

- [ ] Public functions have doc comments
- [ ] Complex logic has explanatory comments
- [ ] Interface definitions describe purpose

1. **Fix issues**: If issues are found, fix them immediately

## Output

Report:

- Files reviewed
- Issues found (if any)
- Fixes applied (if any)
- Verification results

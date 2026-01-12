# Fix Issue #{issue_number}

Fix a GitHub issue using appropriate skills based on labels.

## Skills to Use

1. **Always**: `git-workflow` - Branch, commit, PR
2. **Based on label**:
   - `lang:go` → `go-development`
   - `lang:typescript` → `typescript-development`
   - `lang:solidity` → `solidity-development`
   - `scope:*` → `git-workflow` only (+ scope-specific commands)

## Process

1. **Fetch issue**: `gh issue view {issue_number}`
2. **Identify labels**: Check `lang:*` or `scope:*` labels
3. **Create branch**: Follow `git-workflow` Skill
4. **Implement**: Follow language Skill (if applicable)
5. **Verify**: Run Skill-specific commands
6. **Commit & PR**: Follow `git-workflow` Skill

## Non-Code Tasks (scope:* labels)

| Scope | Verification |
|-------|--------------|
| `scope:docs` | Check markdown formatting |
| `scope:devops` | Test workflow locally |
| `scope:scripts` | `make shfmt` |
| `scope:makefile` | `make mk-lint` |
| `scope:config` | Validate syntax |
| `scope:db` | `make atlas-fmt && make atlas-lint` |

## Example

```
/fix-issue #123
```

Issue #123 with labels `enhancement`, `lang:go`, `chain:btc`:
→ Use `git-workflow` + `go-development` with Bitcoin context
